package mapping

import (
	"testing"

	"github.com/abekoh/mapc/fieldmapper"
	"github.com/abekoh/mapc/testdata/sample"
	"github.com/abekoh/mapc/typemapper"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMapper_NewMapping(t *testing.T) {
	outPkgPath := "github.com/abekoh/mapc/testdata/sample"
	t.Run("map same field", func(t *testing.T) {
		mapper := Mapper{
			FieldMappers: fieldmapper.DefaultMappers,
			TypeMappers:  typemapper.Defaults,
		}
		got, err := mapper.NewMapping(sample.SrcUser{}, sample.DestUser{}, outPkgPath)
		require.Nil(t, err)
		assert.Equal(t, "SrcUser", got.Src.Name())
		assert.Equal(t, outPkgPath, got.Src.PkgPath())
		assert.Equal(t, "DestUser", got.Dest.Name())
		assert.Equal(t, outPkgPath, got.Dest.PkgPath())
		assert.Len(t, got.FieldPairs, 4)
	})
	t.Run("no fieldmappers,typemappers", func(t *testing.T) {
		mapper := Mapper{
			FieldMappers: []fieldmapper.FieldMapper{},
			TypeMappers:  []typemapper.TypeMapper{},
		}
		got, err := mapper.NewMapping(sample.SrcUser{}, sample.DestUser{}, outPkgPath)
		require.Nil(t, err)
		assert.Equal(t, "SrcUser", got.Src.Name())
		assert.Equal(t, outPkgPath, got.Src.PkgPath())
		assert.Equal(t, "DestUser", got.Dest.Name())
		assert.Equal(t, outPkgPath, got.Dest.PkgPath())
		assert.Len(t, got.FieldPairs, 0)
	})
}
