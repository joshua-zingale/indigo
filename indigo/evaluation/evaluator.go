package evaluation

import (
	"fmt"
	"reflect"

	"github.com/joshua-zingale/indigo/indigo/interfaces"
	"github.com/joshua-zingale/indigo/indigo/internal"
)

type goFunction struct {
	function       any
	parameterTypes []reflect.Type
}

func NewIndigoFunctionFromGoFunction(function any) interfaces.IndigoFunction {
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
	if err := interfaces.ValidateFunctionArgs(gf, args); err != nil {
		return nil, err
	}
	return invoke(gf.function, args)
}

type StandardEvaluator struct {
	namespace internal.NameSpace
}

func (se *StandardEvaluator) Eval(expression any) (any, error) {
	return se.evalInNamespace(expression, se.namespace)
}

func (se *StandardEvaluator) evalInNamespace(expression any, namespace internal.NameSpace) (any, error) {
	switch typedExpression := expression.(type) {
	case interfaces.Symbol:
		if value, ok := namespace.Get(typedExpression); ok {
			return value, nil
		}
		return nil, interfaces.UndefinedSymbolError(typedExpression)
	case interfaces.Cons:
		if _, ok := typedExpression.Car().(interfaces.Symbol); ok {
		}
	}
	panic("unreachable!")
}

func evalList(head interfaces.Cons, namespace interfaces.NameSpace) (any, error) {
	symbol, ok := head.Car().(interfaces.Symbol)
	if !ok {
		return nil, fmt.Errorf("cannot evaluate %v as a function: must be a symbol", head.Car())
	}
	function, ok := namespace.Get(symbol)
	if !ok {
		return nil, interfaces.UndefinedSymbolError(symbol)
	}
	if !isGoFunc(function) {
		return nil, interfaces.ExpectedButFoundTypeError("function", function)
	}
	panic("TODO")
}
