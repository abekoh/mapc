package util

import (
	"os"
	"path/filepath"

	"golang.org/x/mod/modfile"
)

func UpperFirst(inp string) string {
	f := inp[0]
	if f < 0x61 || f > 0x7A {
		return inp
	}
	return string(f-0x20) + inp[1:]
}

type Queue[T any] struct {
	elements []T
}

func (q *Queue[T]) Push(el T) {
	q.elements = append(q.elements, el)
}

func (q *Queue[T]) Pop() (T, bool) {
	if len(q.elements) == 0 {
		var t T
		return t, false
	}
	r := q.elements[0]
	q.elements = q.elements[1:]
	return r, true
}

func Prepend[T any](x []T, y T) []T {
	var zero T
	x = append(x, zero)
	copy(x[1:], x)
	x[0] = y
	return x
}

func RootPkgPath() (rootDirPath string, rootPkgPath string, err error) {
	modFileName := "go.mod"
	wd, err := os.Getwd()
	if err != nil {
		return
	}
	dirPath := filepath.Dir(wd)
	for dirPath != "/" && dirPath != "." {
		info, fErr := os.Stat(filepath.Join(dirPath, modFileName))
		if fErr == nil && !info.IsDir() {
			rootDirPath = dirPath
			break
		}
		dirPath = filepath.Dir(dirPath)
	}
	f, err := os.ReadFile(filepath.Join(dirPath, modFileName))
	if err != nil {
		return
	}
	modFile, err := modfile.Parse(modFileName, f, nil)
	if err != nil {
		return
	}
	rootPkgPath = modFile.Module.Mod.Path
	return
}
