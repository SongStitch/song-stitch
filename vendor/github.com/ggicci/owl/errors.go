package owl

import (
	"errors"
	"fmt"
)

var (
	ErrUnsupportedType      = errors.New("unsupported type")
	ErrNilNamespace         = errors.New("nil namespace")
	ErrInvalidDirectiveName = errors.New("invalid directive/executor name")
	ErrDuplicateExecutor    = errors.New("duplicate executor")
	ErrDuplicateDirective   = errors.New("duplicate directive")
	ErrNilExecutor          = errors.New("nil executor")
	ErrMissingExecutor      = errors.New("missing executor")
)

func invalidDirectiveName(name string) error {
	return fmt.Errorf("%w: %q (should comply with %s)", ErrInvalidDirectiveName, name, reDirectiveName.String())
}

func duplicateExecutor(name string) error {
	return fmt.Errorf("%w: %q (registered to the same namespace)", ErrDuplicateExecutor, name)
}

func duplicateDirective(name string) error {
	return fmt.Errorf("%w: %q (defined in the same struct tag)", ErrDuplicateDirective, name)
}

func nilExecutor(name string) error {
	return fmt.Errorf("%w: %q", ErrNilExecutor, name)
}

type ResolveError struct {
	Err      error
	Resolver *Resolver
}

func (e *ResolveError) Error() string {
	return fmt.Sprintf("resolve field %q failed: %s", e.Resolver.String(), e.Err)
}

func (e *ResolveError) Unwrap() error {
	return e.Err
}

func (e *ResolveError) AsDirectiveExecutionError() *DirectiveExecutionError {
	var de *DirectiveExecutionError
	errors.As(e.Err, &de)
	return de
}

type DirectiveExecutionError struct {
	Err error
	Directive
}

func (e *DirectiveExecutionError) Error() string {
	return fmt.Sprintf("execute directive %q with args %v failed: %s", e.Directive.Name, e.Directive.Argv, e.Err)
}

func (e *DirectiveExecutionError) Unwrap() error {
	return e.Err
}
