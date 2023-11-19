package owl

import "context"

// Option is an option for New.
type Option interface {
	Apply(context.Context) context.Context
}

// OptionFunc is a function that implements Option.
type OptionFunc func(context.Context) context.Context

func (f OptionFunc) Apply(ctx context.Context) context.Context {
	return f(ctx)
}

// WithNamespace binds a namespace to the resolver. The namespace is used to
// lookup directive executors. There's a default namespace, which is used when
// the namespace is not specified. The namespace set in New will be overridden
// by the namespace set in Resolve or Scan.
func WithNamespace(ns *Namespace) Option {
	return WithValue(ckNamespace, ns)
}

// WithValue binds a value to the context.
//
// When used in New, the value is bound to Resolver.Context.
//
// When used in Resolve or Scan, the value is bound to DirectiveRuntime.Context.
// See DirectiveRuntime.Context for more details.
func WithValue(key, value interface{}) Option {
	return OptionFunc(func(ctx context.Context) context.Context {
		return context.WithValue(ctx, key, value)
	})
}

// WithNestedDirectivesEnabled controls whether to resolve nested directives.
// The default value is true. When set to false, the nested directives will not
// be executed. The value set in New will be overridden by the value set in
// Resolve or Scan.
func WithNestedDirectivesEnabled(resolve bool) Option {
	return WithValue(ckResolveNestedDirectives, resolve)
}

type DirectiveRunOrder func(*Directive, *Directive) bool

// WithDirectiveRunOrder sets the order of execution for directives.
//
// When used in New, the directives will be sorted at the tree building stage.
//
// When Used in Resolve or Scan, a copy of the directives will be sorted and used.
func WithDirectiveRunOrder(runOrder DirectiveRunOrder) Option {
	return WithValue(ckDirectiveRunOrder, runOrder)
}
