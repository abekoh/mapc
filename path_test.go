//go:build !windows

package mapc

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

func Test_pkgPathFromRelativePath(t *testing.T) {
	RootDirPath = "/home/abekoh/repos/mapc"
	RootPkgPath = "github.com/abekoh/mapc"
	assert.Equal(t, "github.com/abekoh/mapc/pkg/mapping", pkgPathFromRelativePath("pkg/mapping/sample.go"))
	assert.Equal(t, "github.com/abekoh/mapc/pkg/mapping", pkgPathFromRelativePath("pkg/mapping"))
}
