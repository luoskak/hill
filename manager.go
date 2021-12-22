package mist

import (
	"context"
	"fmt"
	"io"
	"sync"
)

// DefaultManager is the default middleware manager
var DefaultManager = &defaultManager

var defaultManager = Manager{
	mu: new(sync.Mutex),
	mm: make(map[string]int),
}

type Manager struct {
	// rootCtx the manager's context
	rootCtx context.Context
	mu      *sync.Mutex
	mm      map[string]int
	ms      []Middleware
	inter   Interceptor
}

func (m *Manager) Register(name string, mw Middleware) {
	m.mu.Lock()
	defer m.mu.Unlock()
	_, duplicated := m.mm[name]
	if duplicated {
		panic("duplicate register middleware")
	}
	m.mm[name] = len(m.ms)
	m.ms = append(m.ms, mw)
}

func (m *Manager) Init(opt ...Option) {
	m.mu.Lock()
	defer m.mu.Unlock()
	optMap := make(map[string][]Option, len(m.mm))
	for _, op := range opt {
		optMap[op.Name()] = append(optMap[op.Name()], op)
	}
	for n, mw := range m.ms {
		var name string
		for k, v := range m.mm {
			if v == n {
				name = k
			}
		}
		mw.Init(optMap[name])
	}

	var inters []Interceptor
	for _, mw := range m.ms {
		inters = append(inters, mw.Inter(true))
	}
	fullInter := chainInterceptors(inters)
	fullInter(context.Background(), nil, nil, func(ctx context.Context, req interface{}) (interface{}, error) {
		m.rootCtx = ctx
		return nil, nil
	})
}

func (m *Manager) Run() error {
	m.mu.Lock()
	defer m.mu.Unlock()
	info := ServerInfo{}
	for _, mw := range m.ms {
		if runner, ok := mw.(Runner); ok {
			if err := runner.Run(m.rootCtx, &info); err != nil {
				return err
			}
		}
	}
	return nil
}

func (m *Manager) Close() error {
	m.mu.Lock()
	defer m.mu.Unlock()
	var errs error
	for _, mw := range m.ms {
		if app, ok := mw.(io.Closer); ok {
			err := app.Close()
			if errs == nil {
				errs = err
			} else {
				errs = fmt.Errorf("%v; %w", errs, err)
			}
		}
	}
	return errs
}

func (m *Manager) Inter() Interceptor {
	if m.inter != nil {
		return m.inter
	}
	var inters []Interceptor
	for _, mw := range m.ms {
		inters = append(inters, mw.Inter(false))
	}
	if len(inters) == 0 {
		return nil
	}
	m.inter = chainInterceptors(inters)
	return m.inter
}

func chainInterceptors(inters []Interceptor) Interceptor {
	return func(ctx context.Context, req interface{}, info *ServerInfo, handler Handler) (interface{}, error) {
		var i int
		var next Handler
		next = func(ctx context.Context, req interface{}) (interface{}, error) {
			if i == len(inters)-1 {
				return inters[i](ctx, req, info, handler)
			}
			i++
			return inters[i-1](ctx, req, info, next)
		}
		return next(ctx, req)
	}
}
