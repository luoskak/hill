package mist

type Options interface{}

type Option interface {
	Apply(Options)
	Name() string
}

type funcOption struct {
	n string
	f func(Options)
}

func (fo funcOption) Apply(o Options) {
	fo.f(o)
}

func (fo funcOption) Name() string {
	return fo.n
}

func NewFuncMyOption(n string, f func(Options)) *funcOption {
	return &funcOption{
		n: n,
		f: f,
	}
}
