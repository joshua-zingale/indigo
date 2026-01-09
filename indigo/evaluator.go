package indigo

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
)

type IndigoFunction interface {
	Call(args ...any) (any, error)
	ParameterTypes() []reflect.Type
}

type goFunction struct {
	function       any
	parameterTypes []reflect.Type
}

func NewIndigoFunctionFromGoFunction(function any) IndigoFunction {
	if reflect.ValueOf(function).Kind() != reflect.Func {
		panic("must be function")
	}
	return &goFunction{
		function:       function,
		parameterTypes: getFuncParameters(function),
	}
}

func (gf *goFunction) ParameterTypes() []reflect.Type {
	return gf.parameterTypes
}

func (gf *goFunction) Call(args ...any) (any, error) {
	if err := validateFunctionArgs(gf, args); err != nil {
		return nil, err
	}
	return invoke(gf.function, args)
}

func validateFunctionArgs(function IndigoFunction, args []any) error {
	if len(args) != len(function.ParameterTypes()) {
		return fmt.Errorf("expected %d arguments but found %d", len(function.ParameterTypes()), len(args))
	}

	var typeErrors []string
	for i := range args {
		argType := reflect.TypeOf(args[i])
		paramType := function.ParameterTypes()[i]

		if argType != paramType {
			typeErrors = append(typeErrors, fmt.Sprintf("argument %d should be of type %v but found type %v", i, paramType, argType))
		}
	}

	if len(typeErrors) > 0 {
		return errors.New(strings.Join(typeErrors, "; "))
	}
	return nil
}

func Eval(any) {

}

type IndigoEvaluator interface {
	Eval(any) any
}

type StandardEvaluator struct {
	namespace nameSpace
}

func (se *StandardEvaluator) Eval(expression any) (any, error) {
	return se.evalInNamespace(expression, se.namespace)
}

func (se *StandardEvaluator) evalInNamespace(expression any, namespace nameSpace) (any, error) {
	switch typedExpression := expression.(type) {
	case Symbol:
		if value, ok := namespace.get(typedExpression); ok {
			return value, nil
		}
		return nil, undefinedSymbol(typedExpression)
	case *Cons:
		if _, ok := typedExpression.car.(Symbol); ok {
		}
	}
	panic("unreachable!")
}

func evalList(head *Cons, namespace nameSpace) (any, error) {
	symbol, ok := head.car.(Symbol)
	if !ok {
		return nil, fmt.Errorf("cannot evaluate %v as a function: must be a symbol", head.car)
	}
	function, ok := namespace.get(symbol)
	if !ok {
		return nil, undefinedSymbol(symbol)
	}
	if !isGoFunc(function) {
		return nil, expectedButFoundType("function", function)
	}
	panic("TODO")
}

type nameSpace struct {
	parent    *nameSpace
	namespace map[Symbol]any
}

func newnameSpace() *nameSpace {
	return &nameSpace{
		parent:    nil,
		namespace: map[Symbol]any{},
	}
}

func (ns *nameSpace) getChild() *nameSpace {
	return &nameSpace{parent: ns, namespace: make(map[Symbol]any)}
}

func (ns *nameSpace) get(symbol Symbol) (any, bool) {
	if v, ok := ns.namespace[symbol]; ok {
		return v, true
	}

	if ns.parent == nil {
		return nil, false
	}

	return ns.parent.get(symbol)
}

func expectedButFoundType(expectedType string, foundValue any) error {
	return fmt.Errorf("expected value of type '%s' but got '%v' of type '%v'", expectedType, foundValue, reflect.TypeOf(foundValue))
}

func undefinedSymbol(symbol Symbol) error {
	return fmt.Errorf("undefined symbol '%s'", symbol)
}
