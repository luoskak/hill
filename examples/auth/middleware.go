package auth

import (
	"context"

	"github.com/luoskak/mist"
)

const MistName = "auth"

var mwKey mist.ContextKey = MistName

func init() {
	mist.DefaultManager.Register(MistName, &middleware{})
}

type middleware struct {
	opts *mwOptions
}

func (m *middleware) Inter(full bool) mist.Interceptor {
	return func(ctx context.Context, req interface{}, info *mist.ServerInfo, handler mist.Handler) (interface{}, error) {
		return handler(ctx, req)
	}
}

func (m *middleware) Init(opt []mist.Option) {
	opts := defaultMwOptions
	for _, o := range opt {
		o.Apply(&opts)
	}
	m.opts = &opts
}
