package code

import (
	"bytes"
	_ "embed"
	"os"
	"testing"

	"github.com/abekoh/mapc/internal/mapping"
	"github.com/abekoh/mapc/testdata/sample"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func loadSampleRaw(t *testing.T) string {
	t.Helper()
	f, err := os.ReadFile("../../testdata/sample/sample.go")
	if err != nil {
		t.Fatal(err)
	}
	return string(f)
}

func loadSample(t *testing.T) *File {
	t.Helper()
	f, err := LoadFile("../../testdata/sample/sample.go", "github.com/abekoh/mapc/internal/code/testdata/sample")
	if err != nil {
		t.Fatal(err)
	}
	return f
}

func TestNewFile(t *testing.T) {
	f := NewFile("github.com/abekoh/mapc/main")

	// check writing
	var buf bytes.Buffer
	if err := f.Write(&buf); err != nil {
		t.Fatal(err)
	}
	expected := `package main
`
	got := buf.String()
	assert.Equal(t, expected, got)
}

func TestLoadFile(t *testing.T) {
	loadSample(t)
}

func TestFile_Write(t *testing.T) {
	f := loadSample(t)
	raw := loadSampleRaw(t)
	var buf bytes.Buffer
	if err := f.Write(&buf); err != nil {
		t.Fatal(err)
	}
	r := buf.String()
	assert.Equal(t, raw, r)
}

func TestFile_FindFunc(t *testing.T) {
	f := loadSample(t)
	fnName := "MapSrcUserToDestUser"
	fn, ok := f.FindFunc(fnName)
	if !ok {
		t.Errorf("not found %s", fnName)
	}
	assert.Equal(t, fn.name, "MapSrcUserToDestUser")
}

func TestFile_Attach(t *testing.T) {
	outPkgPath := "github.com/abekoh/mapc/testdata/sample"
	mapper := mapping.Mapper{}
	m, err := mapper.NewMapping(sample.SrcUser{}, sample.DestUser{}, outPkgPath)
	require.Nil(t, err)
	fn := NewFuncFromMapping(m, nil)
	t.Run("when Func is found, replace", func(t *testing.T) {
		f := loadSample(t)
		err = f.Attach(fn, Deterministic)
		require.Nil(t, err)
		got, ok := f.FindFunc("MapSrcUserToDestUser")
		assert.True(t, ok)
		assert.Equal(t, fn, got)
	})
	t.Run("when Func is not found, append", func(t *testing.T) {
		f := NewFile(outPkgPath)
		err = f.Attach(fn, Deterministic)
		require.Nil(t, err)
		got, ok := f.FindFunc("MapSrcUserToDestUser")
		assert.True(t, ok)
		assert.Equal(t, fn, got)
	})
}
