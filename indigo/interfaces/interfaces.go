package interfaces

import (
	"reflect"
)

type Cons interface {
	Car() any
	Cdr() any
	Empty() bool
}

type IndigoReader interface {
	// Returns the next-parsed object
	Read() (any, error)
}

type IndigoEvaluator interface {
	Eval(any, NameSpace) (any, error)
}

type IndigoFunction interface {
	// Calls the function with the passed in arguments.
	Call(args ...any) (any, error)
	ParameterTypes() []reflect.Type
}

type Symbol string

type NameSpace interface {

	// Create a Child NameSpace
	NewChild() NameSpace

	// Gets the value associated with a Symbol from the current namespace.
	// If the Symbol is undefined in this namespace, its ancestry is recursively
	// searched for a definition. The second value is false iff no definition
	// is found int he hierarchy
	Get(Symbol) (any, bool)

	// Sets a Symbol's value in this NameSpace
	Set(symbol Symbol, value any)
}
