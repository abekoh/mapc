package code

import (
	"bytes"
	_ "embed"
	"testing"

	"github.com/abekoh/mapc/internal/code/testdata/sample"
	"github.com/abekoh/mapc/internal/mapping"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

//go:embed testdata/sample/sample.go
var mapperRawFile string

func loadSample(t *testing.T) *File {
	t.Helper()
	f, err := LoadFile("testdata/sample/sample.go", "github.com/abekoh/mapc/internal/code/testdata/sample")
	if err != nil {
		t.Fatal(err)
	}
	return f
}

func TestLoadFile(t *testing.T) {
	loadSample(t)
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
	r := buf.String()
	if r != expected {
		t.Errorf("output is must be empty")
	}
}

func TestFile_Write(t *testing.T) {
	f := loadSample(t)
	var buf bytes.Buffer
	if err := f.Write(&buf); err != nil {
		t.Fatal(err)
	}
	r := buf.String()
	assert.Equal(t, mapperRawFile, r)
}

func TestFile_FindFunc(t *testing.T) {
	f := loadSample(t)
	funcName := "ToBUser"
	_, ok := f.FindFunc(funcName)
	if !ok {
		t.Errorf("not found %s", funcName)
	}
}

func TestFile_Apply(t *testing.T) {
	mapper := mapping.Mapper{}
	m, err := mapper.NewMapping(sample.AUser{}, sample.BUser{})
	require.Nil(t, err)
	fn := NewFromMapping(m)
	t.Run("when Func is found, replace", func(t *testing.T) {
		f := loadSample(t)
		err = f.Apply(fn)
		require.Nil(t, err)
		got, ok := f.FindFunc("ToBUser")
		assert.True(t, ok)
		assert.Equal(t, fn, got)
	})
	t.Run("when Func is not found, append", func(t *testing.T) {
		f := NewFile("github.com/abekoh/mapc/internal/code/testdata/sample")
		err = f.Apply(fn)
		require.Nil(t, err)
		got, ok := f.FindFunc("ToBUser")
		assert.True(t, ok)
		assert.Equal(t, fn, got)
	})
}
