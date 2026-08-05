package gobuilder

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"

	"github.com/shishir1290/goduck/internal/generator"
)

type Builder struct {
	project *generator.Project
}

func New(project *generator.Project) *Builder {
	return &Builder{
		project: project,
	}
}

func (b *Builder) Build() error {

	projectDir := filepath.Join("build", b.project.Name)

	if _, err := os.Stat(projectDir); err != nil {
		return fmt.Errorf("project directory not found: %s", projectDir)
	}

	output := b.project.Name

	if runtime.GOOS == "windows" {
		output += ".exe"
	}

	// -------------------------
	// Run: go mod tidy
	// -------------------------

	fmt.Println("Running: go mod tidy")

	tidy := exec.Command("go", "mod", "tidy")
	tidy.Dir = projectDir
	tidy.Stdout = os.Stdout
	tidy.Stderr = os.Stderr

	if err := tidy.Run(); err != nil {
		return fmt.Errorf("go mod tidy failed: %w", err)
	}

	// -------------------------
	// Run: go build
	// -------------------------

	fmt.Println("Running: go build")

	cmd := exec.Command(
		"go",
		"build",
		"-o",
		output,
	)

	cmd.Dir = projectDir
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("go build failed: %w", err)
	}

	fmt.Printf(
		"Executable created: %s\n",
		filepath.Join(projectDir, output),
	)

	return nil
}