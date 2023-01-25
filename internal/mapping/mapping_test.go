package mapping

import (
	"testing"

	"github.com/abekoh/mapc/mapcstd"
	"github.com/abekoh/mapc/testdata/sample"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMapper_NewMapping(t *testing.T) {
	outPkgPath := "github.com/abekoh/mapc/testdata/sample"
	t.Run("map same field", func(t *testing.T) {
		mapper := Mapper{
			FieldMappers: mapcstd.DefaultFieldMappers,
			TypeMappers:  mapcstd.DefaultTypeMappers,
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
			FieldMappers: []mapcstd.FieldMapper{},
			TypeMappers:  []mapcstd.TypeMapper{},
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
