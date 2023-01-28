package str

import (
	"reflect"
	"testing"

	"github.com/abekoh/mapc/testdata/sample"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewStruct(t *testing.T) {
	got, err := NewStruct(reflect.TypeOf(sample.Object{}))
	require.Nil(t, err)
	assert.Equal(t, "Object", got.Name())
	assert.Equal(t, "github.com/abekoh/mapc/testdata/sample", got.PkgPath())
	var gotFieldNames []string
	for _, f := range got.Fields {
		switch f.Name() {
		case "ExternalType":
			assert.Equal(t, "github.com/google/uuid", f.Typ().PkgPath())
		default:
		}
		gotFieldNames = append(gotFieldNames, f.Name())
	}
	assert.ElementsMatch(t, sample.ObjectFieldNames, gotFieldNames)
}
