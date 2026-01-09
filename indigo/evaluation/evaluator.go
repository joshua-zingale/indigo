package evaluation

import (
	"fmt"
	"reflect"

	"github.com/joshua-zingale/indigo/indigo/interfaces"
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

type StandardEvaluator struct{}

func NewStandardEvaluator() interfaces.IndigoEvaluator {
	return &StandardEvaluator{}
}

func (se *StandardEvaluator) Eval(expression any, namespace interfaces.NameSpace) (any, error) {
	return evalInNamespace(expression, namespace)
}

func evalInNamespace(expression any, namespace interfaces.NameSpace) (any, error) {
	switch typedExpression := expression.(type) {
	case interfaces.Symbol:
		if value, ok := namespace.Get(typedExpression); ok {
			return value, nil
		}
		return nil, interfaces.UndefinedSymbolError(typedExpression)
	case interfaces.Cons:
		if list, err := consToSlice(typedExpression); err == nil {
			return evalList(list, namespace)
		}
		return typedExpression, nil
	default:
		return typedExpression, nil
	}
}

func evalList(list []any, namespace interfaces.NameSpace) (any, error) {
	if len(list) == 0 {
		return nil, fmt.Errorf("cannot evaluate empty list")
	}
	symbol, ok := list[0].(interfaces.Symbol)
	if !ok {
		return nil, fmt.Errorf("cannot use %v as a function: must be a symbol", list[0])
	}

	value, ok := namespace.Get(symbol)
	if !ok {
		return nil, interfaces.UndefinedSymbolError(symbol)
	}

	function, ok := value.(interfaces.IndigoFunction)
	if !ok {
		return nil, interfaces.ExpectedButFoundTypeError("function", value)
	}

	var evaluatedArgList []any
	for _, element := range list[1:] {
		evaluatedElement, err := evalInNamespace(element, namespace)
		if err != nil {
			return nil, err
		}
		evaluatedArgList = append(evaluatedArgList, evaluatedElement)
	}

	expressionValue, err := function.Call(evaluatedArgList...)
	if err != nil {
		return nil, err
	}

	return expressionValue, nil
}

func consToSlice(maybeCons any) ([]any, error) {
	var result []any

	c, ok := maybeCons.(interfaces.Cons)
	if !ok {
		return nil, fmt.Errorf("invalid list")
	}

	for !c.Empty() {
		result = append(result, c.Car())

		next := c.Cdr()
		if next == nil {
			break
		}

		nextCons, ok := next.(interfaces.Cons)
		if !ok {
			return nil, fmt.Errorf("invalid list: must be nil terminated")
		}
		c = nextCons
	}

	return result, nil
}
