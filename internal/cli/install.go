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

		// 2. Install VS Code extension
		fmt.Println("Installing VS Code extension...")
		tempDir := os.TempDir()
		tempVsixPath := filepath.Join(tempDir, "goduck-language-0.0.1.vsix")

		if err := os.WriteFile(tempVsixPath, vsixBytes, 0644); err != nil {
			return fmt.Errorf("failed to write temporary vsix file: %w", err)
		}
		defer os.Remove(tempVsixPath)

		// Check if vscode is installed
		vscodeCmd := "code"
		if runtime.GOOS == "windows" {
			_, err = exec.LookPath("code")
			if err != nil {
				vscodeCmd = "code.cmd"
			}
		}

		installExtCmd := exec.Command(vscodeCmd, "--install-extension", tempVsixPath)
		installExtCmd.Stdout = os.Stdout
		installExtCmd.Stderr = os.Stderr

		if err := installExtCmd.Run(); err != nil {
			if runtime.GOOS == "windows" {
				fmt.Println("Retrying extension install using cmd...")
				retryCmd := exec.Command("cmd", "/c", "code", "--install-extension", tempVsixPath)
				retryCmd.Stdout = os.Stdout
				retryCmd.Stderr = os.Stderr
				err = retryCmd.Run()
			}
			if err != nil {
				return fmt.Errorf("failed to install VS Code extension: %w (make sure 'code' is in your PATH)", err)
			}
		}

		fmt.Println("✓ VS Code extension installed successfully!")
		return nil
	},
}
