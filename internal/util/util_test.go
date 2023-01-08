package util

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_rootPath(t *testing.T) {
	rootDirPath, rootPkgPath, err := rootPath()
	require.Nil(t, err)
	assert.NotEmpty(t, rootDirPath)
	assert.Equal(t, "github.com/abekoh/mapc", rootPkgPath)
}
