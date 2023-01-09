//go:build !windows

package util

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_rootPath(t *testing.T) {
	rootDirPath, rootPkgPath, err := rootPath()
	require.Nil(t, err)
	assert.NotEmpty(t, rootDirPath)
	// TODO: Windows path
	assert.True(t, strings.HasPrefix(rootDirPath, "/"))
	assert.Equal(t, "github.com/abekoh/mapc", rootPkgPath)
}

func TestPkgPathFromRelativePath(t *testing.T) {
	RootDirPath = "/home/abekoh/repos/mapc"
	RootPkgPath = "github.com/abekoh/mapc"
	assert.Equal(t, "github.com/abekoh/mapc/pkg/mapping", PkgPathFromRelativePath("pkg/mapping/sample.go"))
	assert.Equal(t, "github.com/abekoh/mapc/pkg/mapping", PkgPathFromRelativePath("pkg/mapping"))
}
