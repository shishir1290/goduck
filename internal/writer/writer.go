package writer

import (
	"os"
	"path/filepath"

	"github.com/shishir1290/goduck/internal/generator"
)

type Writer struct {
	project *generator.Project
}

func New(
	project *generator.Project,
) *Writer {

	return &Writer{
		project: project,
	}
}

func (w *Writer) Write() error {

	root := filepath.Join(
		"build",
		w.project.Name,
	)

	if err := os.RemoveAll(root); err != nil {
		return err
	}

	if err := os.MkdirAll(root, 0755); err != nil {
		return err
	}

	for _, dir := range w.project.Directories {

		err := os.MkdirAll(
			filepath.Join(root, dir),
			0755,
		)

		if err != nil {
			return err
		}
	}
	for _, file := range w.project.Files {

		path := filepath.Join(
			root,
			file.Path,
		)

		if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
			return err
		}

		err := os.WriteFile(
			path,
			[]byte(file.Content),
			0644,
		)

		if err != nil {
			return err
		}
	}

	return nil
}