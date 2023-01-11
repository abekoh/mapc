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
	t.Run("map same field", func(t *testing.T) {
		mapper := Mapper{
			FieldMappers: fieldmapper.Default,
			TypeMappers:  typemapper.Default,
		}
		got, err := mapper.NewMapping(sample.SrcUser{}, sample.DestUser{})
		require.Nil(t, err)
		assert.Equal(t, "SrcUser", got.From.Name)
		assert.Equal(t, "github.com/abekoh/mapc/testdata/sample", got.From.PkgPath)
		assert.Equal(t, "DestUser", got.To.Name)
		assert.Equal(t, "github.com/abekoh/mapc/testdata/sample", got.To.PkgPath)
		assert.Len(t, got.FieldPairs, 4)
	})
	t.Run("no fieldmappers,typemappers", func(t *testing.T) {
		mapper := Mapper{
			FieldMappers: []fieldmapper.FieldMapper{},
			TypeMappers:  []typemapper.TypeMapper{},
		}
		got, err := mapper.NewMapping(sample.SrcUser{}, sample.DestUser{})
		require.Nil(t, err)
		assert.Equal(t, "SrcUser", got.From.Name)
		assert.Equal(t, "github.com/abekoh/mapc/testdata/sample", got.From.PkgPath)
		assert.Equal(t, "DestUser", got.To.Name)
		assert.Equal(t, "github.com/abekoh/mapc/testdata/sample", got.To.PkgPath)
		assert.Len(t, got.FieldPairs, 0)
	})
}
