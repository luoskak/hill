package mist

import "context"

type Middleware interface {
	Inter(full bool) Interceptor
	Init(opt []Option)
}

type Runner interface {
	Run(ctx context.Context, info *ServerInfo) error
}
