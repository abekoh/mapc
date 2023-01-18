package mapc

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

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

func pkgPathFromRelativePath(relativePath string) string {
	ext := filepath.Ext(relativePath)
	if ext == "" {
		return filepath.Join(RootPkgPath, relativePath)
	} else {
		return filepath.Join(RootPkgPath, filepath.Dir(relativePath))
	}
}
