package mapcstd

type Caster interface {
	Caster()
}

type NopCaster struct{}

func (c NopCaster) Caster() {
}

type Caller struct {
	PkgPath    string
	Name       string
	CallerType CallerType
	ReturnType ReturnType
}

func (c Caller) Caster() {
}

type CallerType int

const (
	Unary CallerType = iota
	Type
	Func
)
