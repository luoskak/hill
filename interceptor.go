package mist

import "context"

type ServerInfo struct {
	Server         interface{}
	IsClientStream bool
	IsServerStream bool
	FullMethod     string
}

type Handler func(ctx context.Context, req interface{}) (interface{}, error)

type Interceptor func(ctx context.Context, req interface{}, info *ServerInfo, handler Handler) (interface{}, error)
