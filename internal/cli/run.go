package cli

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
	"syscall"
	"time"

	"github.com/shishir1290/goduck/internal/compiler"
	"github.com/spf13/cobra"
)

var (
	cmdMutex   sync.Mutex
	runningCmd *exec.Cmd
)

var runCmd = &cobra.Command{
	Use:   "run [directory]",
	Short: "Run the .duck application with hot-reloading (like air)",
	Args:  cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		watchDir := "src"
		if len(args) > 0 {
			watchDir = args[0]
		}

		if _, err := os.Stat(watchDir); err != nil {
			return fmt.Errorf("directory not found: %s", watchDir)
		}

		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		// Listen for termination signals (Ctrl+C)
		sigChan := make(chan os.Signal, 1)
		signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
		go func() {
			<-sigChan
			fmt.Println("\n[goduck run] Shutting down...")
			cancel()
			killRunningProcess()
			os.Exit(0)
		}()

		// Function to rebuild and run the project
		rebuildAndRun := func() {
			fmt.Println("\n[goduck run] File changes detected. Rebuilding...")
			killRunningProcess()

			appName, err := compiler.Build(watchDir)
			if err != nil {
				fmt.Printf("[goduck run] Build failed: %v\n", err)
				fmt.Println("[goduck run] Waiting for fixes before restarting...")
				return
			}

			startProcess(appName)
		}

		// Initial build and start
		fmt.Println("[goduck run] Starting initial build...")
		appName, err := compiler.Build(watchDir)
		if err != nil {
			fmt.Printf("[goduck run] Initial build failed: %v\n", err)
			fmt.Println("[goduck run] Waiting for fixes...")
		} else {
			startProcess(appName)
		}

		// Setup custom file watcher
		watcher := NewWatcher(watchDir, rebuildAndRun)
		fmt.Printf("[goduck run] Watching directory '%s' for changes...\n", watchDir)
		watcher.Start(ctx)

		return nil
	},
}

func startProcess(appName string) {
	cmdMutex.Lock()
	defer cmdMutex.Unlock()

	executable := appName
	if runtime.GOOS == "windows" {
		executable += ".exe"
	}

	projectDir := filepath.Join("build", appName)
	binaryPath := filepath.Join(projectDir, executable)
	
	absPath, err := filepath.Abs(binaryPath)
	if err != nil {
		absPath = binaryPath
	}

	fmt.Printf("[goduck run] Starting process: %s\n", absPath)

	cmd := exec.Command(absPath)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	runningCmd = cmd

	go func() {
		if err := cmd.Run(); err != nil {
			// Do not log exit error if we killed the process ourselves
			if cmd.ProcessState != nil && !cmd.ProcessState.Exited() {
				fmt.Printf("[goduck run] Process exited with error: %v\n", err)
			}
		} else {
			fmt.Println("[goduck run] Process exited cleanly")
		}
	}()
}

func killRunningProcess() {
	cmdMutex.Lock()
	defer cmdMutex.Unlock()

	if runningCmd != nil && runningCmd.Process != nil {
		fmt.Println("[goduck run] Stopping current process...")
		_ = runningCmd.Process.Kill()
		runningCmd = nil
	}
}

type Watcher struct {
	dir      string
	files    map[string]time.Time
	mu       sync.Mutex
	onChange func()
}

func NewWatcher(dir string, onChange func()) *Watcher {
	return &Watcher{
		dir:      dir,
		files:    make(map[string]time.Time),
		onChange: onChange,
	}
}

func (w *Watcher) Scan() (bool, error) {
	w.mu.Lock()
	w.mu.Unlock() // Use defer instead
	return w.scanUnlocked()
}

func (w *Watcher) scanUnlocked() (bool, error) {
	w.mu.Lock()
	defer w.mu.Unlock()

	currentFiles := make(map[string]time.Time)
	err := filepath.Walk(w.dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && strings.HasSuffix(info.Name(), ".duck") {
			currentFiles[path] = info.ModTime()
		}
		return nil
	})
	if err != nil {
		return false, err
	}

	changed := false

	// Check for additions or modifications
	for path, modTime := range currentFiles {
		oldTime, exists := w.files[path]
		if !exists || modTime.After(oldTime) {
			changed = true
		}
	}

	// Check for deletions
	for path := range w.files {
		if _, exists := currentFiles[path]; !exists {
			changed = true
		}
	}

	w.files = currentFiles
	return changed, nil
}

func (w *Watcher) Start(ctx context.Context) {
	// Initial scan to populate files map
	_, _ = w.scanUnlocked()

	ticker := time.NewTicker(500 * time.Millisecond)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			changed, err := w.scanUnlocked()
			if err == nil && changed {
				w.onChange()
			}
		}
	}
}
