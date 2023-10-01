package owl

var (
	defaultNS = NewNamespace()
)

// RegisterDirectiveExecutor registers a named executor globally (in the default
// namespace).
func RegisterDirectiveExecutor(name string, exe DirectiveExecutor) {
	defaultNS.RegisterDirectiveExecutor(name, exe)
}

// ReplaceDirectiveExecutor replaces a named executor globally (in the default
// namespace). It works like RegisterDirectiveExecutor without panic on
// duplicate names.
func ReplaceDirectiveExecutor(name string, exe DirectiveExecutor) {
	defaultNS.ReplaceDirectiveExecutor(name, exe)
}

// LookupExecutor returns the executor by name globally (from the default
// namespace).
func LookupExecutor(name string) DirectiveExecutor {
	return defaultNS.LookupExecutor(name)
}

// Namespace isolates the executors as a collection.
type Namespace struct {
	executors map[string]DirectiveExecutor
}

// NewNamespace creates a new namespace. Which is a collection of executors.
func NewNamespace() *Namespace {
	return &Namespace{
		executors: make(map[string]DirectiveExecutor),
	}
}

// RegisterDirectiveExecutor registers a named executor in the namespace.
// The executor should implement the DirectiveExecutor interface.
// Will panic if the name were taken or the executor is nil.
func (ns *Namespace) RegisterDirectiveExecutor(name string, exe DirectiveExecutor) {
	if _, ok := ns.executors[name]; ok {
		panic(duplicateExecutor(name))
	}
	ns.ReplaceDirectiveExecutor(name, exe)
}

// ReplaceDirectiveExecutor works like RegisterDirectiveExecutor without panicing
// on duplicate names. But it will panic if the executor is nil or the name is invalid.
func (ns *Namespace) ReplaceDirectiveExecutor(name string, exe DirectiveExecutor) {
	if exe == nil {
		panic(nilExecutor(name))
	}
	if !isValidDirectiveName(name) {
		panic(invalidDirectiveName(name))
	}
	ns.executors[name] = exe
}

// LookupExecutor returns the executor by name.
func (ns *Namespace) LookupExecutor(name string) DirectiveExecutor {
	return ns.executors[name]
}
