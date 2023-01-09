package typemapper

type Caster interface {
	PkgPath() string
	Func() string
}

type NopCaster struct{}

func (n NopCaster) PkgPath() string {
	return ""
}

func (n NopCaster) Func() string {
	return ""
}

type SimpleCaster struct {
	pkgPath string
	fn      string
}

func (s SimpleCaster) PkgPath() string {
	return s.pkgPath
}

func (s SimpleCaster) Func() string {
	return s.fn
}
