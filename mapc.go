package mapc

type MapC struct {
	pairs     []pair
	globalOpt Option
}

type pair struct {
	from any
	to   any
	opt  Option
}

type Option struct {
	OutPath string
}

func (o Option) override(opts ...Option) Option {
	for _, opt := range opts {
		// TODO: to be more simple
		if len(opt.OutPath) != 0 {
			o.OutPath = opt.OutPath
		}
	}
	return o
}

func New() *MapC {
	return &MapC{}
}

func (mc *MapC) Global(opt Option) {
	mc.globalOpt = opt
}

func (mc *MapC) Register(from, to any, options ...Option) {
	mc.pairs = append(mc.pairs, pair{from: from, to: to})
}

func (mc MapC) Generate() (errs []error) {
	return
}
