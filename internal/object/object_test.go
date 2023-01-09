package object

import (
	"reflect"
	"testing"

	"github.com/abekoh/mapc/internal/testdata/sample"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewStruct(t *testing.T) {
	o := sample.Object{}
	got, err := NewStruct(reflect.TypeOf(o))
	require.Nil(t, err)
	assert.Equal(t, "Object", got.Name)
	assert.Equal(t, "github.com/abekoh/mapc/internal/testdata/sample", got.PkgPath)
	var gotFieldNames []string
	for _, f := range got.Fields {
		switch f.Name() {
		case "ExternalType":
			assert.Equal(t, "github.com/google/uuid", f.PkgPath())
		default:
		}
		gotFieldNames = append(gotFieldNames, f.Name())
	}
	assert.ElementsMatch(t, sample.ObjectFieldNames, gotFieldNames)
}
