package code

import (
	"fmt"
	"os"

	"github.com/dave/dst"
	"github.com/dave/dst/decorator"
)

type File struct {
	file *dst.File
}

func NewFile(path string) (*File, error) {
	f, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %w", err)
	}
	df, err := decorator.Parse(f)
	if err != nil {
		return nil, fmt.Errorf("failed to parse file: %w", err)
	}
	return &File{file: df}, nil
}
