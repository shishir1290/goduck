package cli

import (
	_ "embed"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"

	"github.com/spf13/cobra"
)

//go:embed goduck-language-0.0.1.vsix
var vsixBytes []byte

var installCmd = &cobra.Command{
	Use:   "install",
	Short: "Install Goduck CLI and the VS Code extension on this PC",
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println("Installing Goduck CLI on your system...")

		// 1. Install CLI binary
		exePath, err := os.Executable()
		if err != nil {
			return fmt.Errorf("failed to get current executable path: %w", err)
		}

		gopath := os.Getenv("GOPATH")
		if gopath == "" {
			home, err := os.UserHomeDir()
			if err != nil {
				return fmt.Errorf("failed to get user home directory: %w", err)
			}
			gopath = filepath.Join(home, "go")
		}

		binDir := filepath.Join(gopath, "bin")
		if err := os.MkdirAll(binDir, 0755); err != nil {
			return fmt.Errorf("failed to create bin directory: %w", err)
		}

		binaryName := "goduck"
		if runtime.GOOS == "windows" {
			binaryName += ".exe"
		}
		targetPath := filepath.Join(binDir, binaryName)

		// Copy executable
		src, err := os.Open(exePath)
		if err != nil {
			return fmt.Errorf("failed to open source binary: %w", err)
		}
		defer src.Close()

		// Attempt to delete target if exists to avoid "text file busy" error
		_ = os.Remove(targetPath)

		dst, err := os.OpenFile(targetPath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0755)
		if err != nil {
			return fmt.Errorf("failed to open destination path %s: %w", targetPath, err)
		}
		defer dst.Close()

		if _, err := io.Copy(dst, src); err != nil {
			return fmt.Errorf("failed to copy binary: %w", err)
		}

		fmt.Printf("✓ Installed Goduck CLI to: %s\n", targetPath)

		// 2. Install Editor extensions
		if err := installExtensions(false); err != nil {
			return err
		}

		// Write marker file to avoid auto-installing on subsequent runs
		if home, err := os.UserHomeDir(); err == nil {
			goduckDir := filepath.Join(home, ".goduck")
			_ = os.MkdirAll(goduckDir, 0755)
			markerFile := filepath.Join(goduckDir, ".extension_installed_0.0.1")
			_ = os.WriteFile(markerFile, []byte("installed"), 0644)
		}

		return nil
	},
}

func findEditorCLIs() []string {
	var foundPaths []string

	// Define base list of CLI names
	editorNames := []string{"code", "cursor", "antigravity", "codium", "windsurf", "code-insiders"}

	// First check via PATH (using LookPath)
	for _, name := range editorNames {
		path, err := exec.LookPath(name)
		if err == nil {
			foundPaths = append(foundPaths, path)
			continue
		}

		// On Windows, check with .cmd extension
		if runtime.GOOS == "windows" {
			path, err = exec.LookPath(name + ".cmd")
			if err == nil {
				foundPaths = append(foundPaths, path)
				continue
			}
		}
	}

	// Check typical installation paths if not found in PATH
	if runtime.GOOS == "windows" {
		localAppData := os.Getenv("LOCALAPPDATA")
		programFiles := os.Getenv("ProgramFiles")

		// Standard paths on Windows relative to LocalAppData/Programs or ProgramFiles
		winRelativePaths := map[string][]string{
			"code": {
				filepath.Join(localAppData, "Programs", "Microsoft VS Code", "bin", "code.cmd"),
				filepath.Join(programFiles, "Microsoft VS Code", "bin", "code.cmd"),
			},
			"code-insiders": {
				filepath.Join(localAppData, "Programs", "Microsoft VS Code Insiders", "bin", "code-insiders.cmd"),
				filepath.Join(programFiles, "Microsoft VS Code Insiders", "bin", "code-insiders.cmd"),
			},
			"cursor": {
				filepath.Join(localAppData, "Programs", "cursor", "resources", "app", "bin", "cursor.cmd"),
				filepath.Join(programFiles, "cursor", "resources", "app", "bin", "cursor.cmd"),
			},
			"antigravity": {
				filepath.Join(localAppData, "Programs", "Antigravity", "bin", "antigravity.cmd"),
				filepath.Join(programFiles, "Antigravity", "bin", "antigravity.cmd"),
			},
			"codium": {
				filepath.Join(localAppData, "Programs", "VSCodium", "bin", "codium.cmd"),
				filepath.Join(programFiles, "VSCodium", "bin", "codium.cmd"),
			},
			"windsurf": {
				filepath.Join(localAppData, "Programs", "Windsurf", "bin", "windsurf.cmd"),
				filepath.Join(programFiles, "Windsurf", "bin", "windsurf.cmd"),
			},
		}

		for _, paths := range winRelativePaths {
			for _, p := range paths {
				if _, err := os.Stat(p); err == nil {
					if !contains(foundPaths, p) {
						foundPaths = append(foundPaths, p)
					}
				}
			}
		}
	} else if runtime.GOOS == "darwin" {
		// macOS default paths
		home := os.Getenv("HOME")
		macPaths := []string{
			"/Applications/Visual Studio Code.app/Contents/Resources/app/bin/code",
			"/Applications/Visual Studio Code - Insiders.app/Contents/Resources/app/bin/code",
			"/Applications/Cursor.app/Contents/Resources/app/bin/cursor",
			"/Applications/Antigravity.app/Contents/Resources/app/bin/antigravity",
			"/Applications/VSCodium.app/Contents/Resources/app/bin/codium",
			"/Applications/Windsurf.app/Contents/Resources/app/bin/windsurf",

			filepath.Join(home, "Applications", "Visual Studio Code.app", "Contents", "Resources", "app", "bin", "code"),
			filepath.Join(home, "Applications", "Cursor.app", "Contents", "Resources", "app", "bin", "cursor"),
			filepath.Join(home, "Applications", "Antigravity.app", "Contents", "Resources", "app", "bin", "antigravity"),
			filepath.Join(home, "Applications", "VSCodium.app", "Contents", "Resources", "app", "bin", "codium"),
			filepath.Join(home, "Applications", "Windsurf.app", "Contents", "Resources", "app", "bin", "windsurf"),
		}

		for _, p := range macPaths {
			if _, err := os.Stat(p); err == nil {
				if !contains(foundPaths, p) {
					foundPaths = append(foundPaths, p)
				}
			}
		}
	}

	return foundPaths
}

func contains(slice []string, val string) bool {
	for _, item := range slice {
		if item == val {
			return true
		}
	}
	return false
}

func installExtensions(silent bool) error {
	tempDir := os.TempDir()
	tempVsixPath := filepath.Join(tempDir, "goduck-language-0.0.1.vsix")

	if err := os.WriteFile(tempVsixPath, vsixBytes, 0644); err != nil {
		return fmt.Errorf("failed to write temporary vsix file: %w", err)
	}
	defer os.Remove(tempVsixPath)

	editorPaths := findEditorCLIs()
	if len(editorPaths) == 0 {
		if !silent {
			fmt.Println("No supported editors (VS Code, Cursor, Antigravity, VSCodium, Windsurf) found in PATH or standard installation directories.")
		}
		return nil
	}

	for _, editorPath := range editorPaths {
		if !silent {
			fmt.Printf("Installing goduck extension to %s...\n", filepath.Base(editorPath))
		}

		var cmd *exec.Cmd
		if runtime.GOOS == "windows" && (filepath.Ext(editorPath) == ".cmd" || filepath.Ext(editorPath) == ".bat") {
			cmd = exec.Command("cmd", "/c", editorPath, "--install-extension", tempVsixPath)
		} else {
			cmd = exec.Command(editorPath, "--install-extension", tempVsixPath)
		}

		if !silent {
			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr
		}

		if err := cmd.Run(); err != nil {
			if !silent {
				fmt.Printf("Warning: failed to install extension for %s: %v\n", editorPath, err)
			}
		} else if !silent {
			fmt.Printf("✓ Successfully installed extension for %s\n", filepath.Base(editorPath))
		}
	}

	return nil
}

func autoInstallExtensionsIfNeeded() {
	// Skip if we are running the explicit install command
	for _, arg := range os.Args {
		if arg == "install" {
			return
		}
	}

	home, err := os.UserHomeDir()
	if err != nil {
		return
	}

	goduckDir := filepath.Join(home, ".goduck")
	markerFile := filepath.Join(goduckDir, ".extension_installed_0.0.1")

	// If marker file exists, we already installed it for this version. Skip!
	if _, err := os.Stat(markerFile); err == nil {
		return
	}

	// Create .goduck directory if it doesn't exist
	if err := os.MkdirAll(goduckDir, 0755); err != nil {
		return
	}

	// Run installation silently
	fmt.Println("First-time setup: Installing Goduck editor extension for VS Code / Cursor / Antigravity...")
	_ = installExtensions(true)

	// Write marker file so we don't run it again
	_ = os.WriteFile(markerFile, []byte("installed"), 0644)
	fmt.Println("✓ Setup complete!")
}
