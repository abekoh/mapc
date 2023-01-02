package code

import (
	"fmt"
	"os"
)

type File struct {
	code    string
	imports map[string]struct{}
	mapping map[string]Mapping
}

func NewFile(path string) (*File, error) {
	f, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %w", err)
	}
	code := string(f)
	return &File{
		code: code,
	}, nil
}

type Mapping struct {
	name string
	code string
}
