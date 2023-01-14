package typemapper

type Caster interface {
	PkgPath() string
	Caller() string
}

type NopCaster struct{}

func (n NopCaster) PkgPath() string {
	return ""
}

func (n NopCaster) Caller() string {
	return ""
}

type SimpleCaster struct {
	pkgPath string
	caller  string
}

func (s SimpleCaster) PkgPath() string {
	return s.pkgPath
}

func (s SimpleCaster) Caller() string {
	return s.caller
}
