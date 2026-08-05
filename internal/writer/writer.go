package writer

import (
	"os"
	"path/filepath"

	"github.com/shishir1290/goduck/internal/generator"
)

type Writer struct {
	project *generator.Project
}

func New(project *generator.Project) *Writer {
	return &Writer{
		project: project,
	}
}

func (w *Writer) Write() error {

	root := filepath.Join("build", w.project.Name)

	if err := os.MkdirAll(root, 0755); err != nil {
		return err
	}

	for _, file := range w.project.Files {

		full := filepath.Join(root, file.Path)

		dir := filepath.Dir(full)

		if err := os.MkdirAll(dir, 0755); err != nil {
			return err
		}

		if err := os.WriteFile(full, []byte(file.Content), 0644); err != nil {
			return err
		}
	}

	return nil
}