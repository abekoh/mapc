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
	f, err := LoadFile("../../testdata/sample/sample.go", "github.com/abekoh/mapc/testdata/sample")
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
	t.Run("only selectors", func(t *testing.T) {
		idx, fn, ok := f.FindFunc("MapSrcUserToDestUser")
		require.True(t, ok)
		assert.Greater(t, idx, 0)
		assert.Equal(t, &Func{
			name:    "MapSrcUserToDestUser",
			argName: "x",
			srcTyp: &Typ{
				name:    "SrcUser",
				pkgPath: "github.com/abekoh/mapc/testdata/sample",
			},
			destTyp: &Typ{
				name:    "DestUser",
				pkgPath: "github.com/abekoh/mapc/testdata/sample",
			},
			mapExprs: MapExprList{
				&SimpleMapExpr{
					src:     "ID",
					dest:    "ID",
					casters: nil,
				},
				&SimpleMapExpr{
					src:     "Name",
					dest:    "Name",
					casters: nil,
				},
				&SimpleMapExpr{
					src:     "Age",
					dest:    "Age",
					casters: nil,
				},
				&SimpleMapExpr{
					src:     "RegisteredAt",
					dest:    "RegisteredAt",
					casters: nil,
				},
			},
		}, fn)
	})
	t.Run("only comments", func(t *testing.T) {
		idx, fn, ok := f.FindFunc("OnlyCommentsMapper")
		require.True(t, ok)
		assert.Greater(t, idx, 0)
		assert.Equal(t, &Func{
			name:    "OnlyCommentsMapper",
			argName: "x",
			srcTyp: &Typ{
				name:    "SrcUser",
				pkgPath: "github.com/abekoh/mapc/testdata/sample",
			},
			destTyp: &Typ{
				name:    "DestUser",
				pkgPath: "github.com/abekoh/mapc/testdata/sample",
			},
			mapExprs: MapExprList{
				&CommentedMapExpr{
					dest:    "ID",
					comment: "//ID:",
				},
				&CommentedMapExpr{
					dest:    "Name",
					comment: "//Name:",
				},
				&CommentedMapExpr{
					dest:    "Age",
					comment: "// Age:",
				},
				&CommentedMapExpr{
					dest:    "RegisteredAt",
					comment: "// RegisteredAt:",
				},
			},
		}, fn)
	})
	t.Run("first line is comment", func(t *testing.T) {
		idx, fn, ok := f.FindFunc("FirstLineIsCommentMapper")
		require.True(t, ok)
		assert.Greater(t, idx, 0)
		assert.Equal(t, &Func{
			name:    "FirstLineIsCommentMapper",
			argName: "x",
			srcTyp: &Typ{
				name:    "SrcUser",
				pkgPath: "github.com/abekoh/mapc/testdata/sample",
			},
			destTyp: &Typ{
				name:    "DestUser",
				pkgPath: "github.com/abekoh/mapc/testdata/sample",
			},
			mapExprs: MapExprList{
				&CommentedMapExpr{
					dest:    "ID",
					comment: "//ID:",
				},
				&SimpleMapExpr{
					src:     "Name",
					dest:    "Name",
					casters: nil,
				},
				&SimpleMapExpr{
					src:     "Age",
					dest:    "Age",
					casters: nil,
				},
				&SimpleMapExpr{
					src:     "RegisteredAt",
					dest:    "RegisteredAt",
					casters: nil,
				},
			},
		}, fn)
	})
	t.Run("last line is comment", func(t *testing.T) {
		idx, fn, ok := f.FindFunc("LastLineIsCommentMapper")
		require.True(t, ok)
		assert.Greater(t, idx, 0)
		assert.Equal(t, &Func{
			name:    "LastLineIsCommentMapper",
			argName: "x",
			srcTyp: &Typ{
				name:    "SrcUser",
				pkgPath: "github.com/abekoh/mapc/testdata/sample",
			},
			destTyp: &Typ{
				name:    "DestUser",
				pkgPath: "github.com/abekoh/mapc/testdata/sample",
			},
			mapExprs: MapExprList{
				&SimpleMapExpr{
					src:     "ID",
					dest:    "ID",
					casters: nil,
				},
				&SimpleMapExpr{
					src:     "Name",
					dest:    "Name",
					casters: nil,
				},
				&SimpleMapExpr{
					src:     "Age",
					dest:    "Age",
					casters: nil,
				},
				&CommentedMapExpr{
					dest:    "RegisteredAt",
					comment: "//RegisteredAt:",
				},
			},
		}, fn)
	})
	t.Run("no selectors", func(t *testing.T) {
		idx, fn, ok := f.FindFunc("NoSelectorsMapper")
		require.True(t, ok)
		assert.Greater(t, idx, 0)
		assert.Len(t, fn.mapExprs, 4)
		for _, expr := range fn.mapExprs {
			assert.IsType(t, &UnknownMapExpr{}, expr)
		}
	})
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
		_, got, ok := f.FindFunc("MapSrcUserToDestUser")
		assert.True(t, ok)
		assert.Equal(t, fn, got)
	})
	t.Run("when Func is not found, append", func(t *testing.T) {
		f := NewFile(outPkgPath)
		err = f.Attach(fn, Deterministic)
		require.Nil(t, err)
		_, got, ok := f.FindFunc("MapSrcUserToDestUser")
		assert.True(t, ok)
		assert.Equal(t, fn, got)
	})
}
