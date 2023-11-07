package owl

import "context"

type Option interface {
	Apply(context.Context) context.Context
}

type OptionFunc func(context.Context) context.Context

func (f OptionFunc) Apply(ctx context.Context) context.Context {
	return f(ctx)
}

func WithValue(key, value interface{}) Option {
	return OptionFunc(func(ctx context.Context) context.Context {
		return context.WithValue(ctx, key, value)
	})
}

func WithNamespace(ns *Namespace) Option {
	return WithValue(ckNamespace, ns)
}
