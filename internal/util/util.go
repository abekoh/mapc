package util

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"golang.org/x/mod/modfile"
)

var (
	RootDirPath string
	RootPkgPath string
)

func init() {
	rootDirPath, rootPkgPath, err := rootPath()
	if err != nil {
		log.Fatal(err)
	}
	RootDirPath = rootDirPath
	RootPkgPath = rootPkgPath
}

func rootPath() (rootDirPath string, rootPkgPath string, err error) {
	modFileName := "go.mod"
	wd, err := os.Getwd()
	if err != nil {
		return
	}
	dirPath := wd
	for dirPath != "/" && dirPath != "." {
		info, fErr := os.Stat(filepath.Join(dirPath, modFileName))
		if fErr == nil && !info.IsDir() {
			rootDirPath, err = filepath.Abs(dirPath)
			if err != nil {
				return
			}
			break
		}
		dirPath = filepath.Dir(dirPath)
	}
	if rootDirPath == "" {
		err = fmt.Errorf("go.mod is not found")
		return
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

func PkgPathFromFilePath(filePath string) string {
	if !strings.HasPrefix(filePath, RootDirPath) {
		return ""
	}
	// TODO: not work when RootPkgPath like: `github.com/abekoh/mapc/foo/bar.test`
	return filepath.Join(RootPkgPath, filePath[len(RootDirPath):])
}

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

func PkgName(pkgPath string) string {
	sp := strings.Split(pkgPath, "/")
	if len(sp) == 0 {
		return ""
	}
	return sp[len(sp)-1]
}
