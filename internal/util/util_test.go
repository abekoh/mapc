package util

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRootPkgPath(t *testing.T) {
	rootDirPath, rootPkgPath, err := RootPath()
	require.Nil(t, err)
	assert.NotEmpty(t, rootDirPath)
	assert.Equal(t, "github.com/abekoh/mapc", rootPkgPath)
}
