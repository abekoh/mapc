package e2e

import (
	"os"
	"path"
	"testing"

	"github.com/abekoh/mapc"
	"github.com/abekoh/mapc/e2e/testdata/ab"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_SaveToSameFile(t *testing.T) {
	m := mapc.New()
	m.Option.FuncComment = false
	m.Option.OutPath = "dest.go"
	m.Option.OutPkgPath = "github.com/abekoh/mapc/e2e/testdata/ab"
	m.Register(ab.AUser{}, ab.BUser{}, func(option *mapc.Option) {
		option.FuncName = "Func1"
	})
	m.Register(ab.AUser{}, ab.BUser{}, func(option *mapc.Option) {
		option.FuncName = "Func2"
	})
	gs, errs := m.Generate()
	requireNoErrors(t, errs)
	got, err := gs[0].Sprint()
	require.Nil(t, err)
	assert.Equal(t, `package ab

func Func1(x AUser) BUser {
	return BUser{
		ID:           x.ID,
		Name:         x.Name,
		Age:          x.Age,
		RegisteredAt: x.RegisteredAt,
	}
}
func Func2(x AUser) BUser {
	return BUser{
		ID:           x.ID,
		Name:         x.Name,
		Age:          x.Age,
		RegisteredAt: x.RegisteredAt,
	}
}
`, got)
}

func Test_SaveToFile(t *testing.T) {
	tempDirPath := t.TempDir()
	outPath := path.Join(tempDirPath, "dest.go")
	m := mapc.New()
	m.Option.FuncComment = false
	m.Option.OutPath = outPath
	m.Option.OutPkgPath = "github.com/abekoh/mapc/e2e/testdata/ab"
	m.Register(ab.AUser{}, ab.BUser{})
	gs, errs := m.Generate()
	requireNoErrors(t, errs)
	err := gs[0].Save()
	require.Nil(t, err)
	got, err := os.ReadFile(outPath)
	require.Nil(t, err)
	assert.Equal(t, `package ab

func MapAUserToBUser(x AUser) BUser {
	return BUser{
		ID:           x.ID,
		Name:         x.Name,
		Age:          x.Age,
		RegisteredAt: x.RegisteredAt,
	}
}
`, string(got))
}

func Test_OutPathIsNotSet(t *testing.T) {
	m := mapc.New()
	m.Option.OutPkgPath = "github.com/abekoh/mapc/e2e/testdata/ab"
	m.Register(ab.AUser{}, ab.BUser{})
	gs, errs := m.Generate()
	requireNoErrors(t, errs)
	err := gs[0].Save()
	assert.NotNil(t, err)
}
