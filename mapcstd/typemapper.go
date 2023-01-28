package mapcstd

type TypeMapper interface {
	Map(src, dest *Typ) (Caster, bool)
}

var DefaultTypeMappers = []TypeMapper{
	&AssignMapper{},
	&ConvertMapper{},
	&RefMapper{},
	&DerefMapper{},
}

type AssignMapper struct {
}

func (a AssignMapper) Map(src, dest *Typ) (Caster, bool) {
	if src.AssignableTo(dest) {
		return &NopCaster{}, true
	}
	return nil, false
}

type ConvertMapper struct {
}

func (c ConvertMapper) Map(src, dest *Typ) (Caster, bool) {
	if src.ConvertibleTo(dest) {
		return &SimpleCaster{
			caller: &Caller{
				Name:       dest.Name(),
				PkgPath:    dest.PkgPath(),
				CallerType: Type,
			},
		}, true
	}
	return nil, false
}

type RefMapper struct {
}

func (p RefMapper) Map(src, dest *Typ) (Caster, bool) {
	if srcElm, ok := dest.Elem(); ok && src.AssignableTo(srcElm) {
		return &SimpleCaster{
			caller: &Caller{
				PkgPath:    "",
				Name:       "&",
				CallerType: Unary,
			},
		}, true
	}
	return nil, false
}

type DerefMapper struct {
}

func (p DerefMapper) Map(src, dest *Typ) (Caster, bool) {
	if destElm, ok := src.Elem(); ok && destElm.AssignableTo(dest) {
		return &SimpleCaster{
			caller: &Caller{
				PkgPath:    "",
				Name:       "*",
				CallerType: Unary,
			},
		}, true
	}
	return nil, false
}

type MapTypeMapper map[*Typ]map[*Typ]Caster

func (m MapTypeMapper) Map(src, dest *Typ) (Caster, bool) {
	m2, ok := m[src]
	if !ok {
		return nil, false
	}
	c, ok := m2[dest]
	if !ok {
		return nil, false
	}
	return c, true
}

type TypeMapperFunc func(src, dest *Typ) (Caster, bool)

func (m TypeMapperFunc) Map(src, dest *Typ) (Caster, bool) {
	return m(src, dest)
}
