package evaluation

import (
	"fmt"

	"github.com/joshua-zingale/indigo/indigo/interfaces"
)

type StandardEvaluator struct{}

func NewStandardEvaluator() interfaces.IndigoEvaluator {
	return &StandardEvaluator{}
}

func (se *StandardEvaluator) Eval(expression any, namespace interfaces.NameSpace) (any, error) {
	return se.evalInNamespace(expression, namespace)
}

func (se *StandardEvaluator) evalInNamespace(expression any, namespace interfaces.NameSpace) (any, error) {
	switch typedExpression := expression.(type) {
	case interfaces.Symbol:
		if value, ok := namespace.Get(typedExpression); ok {
			return value, nil
		}
		return nil, interfaces.UndefinedSymbolError(typedExpression)
	case interfaces.Cons:
		if list, ok := typedExpression.(interfaces.List); ok {
			return se.evalList(list, namespace)
		}
		return typedExpression, nil
	default:
		return typedExpression, nil
	}
}

func (se *StandardEvaluator) evalList(list interfaces.List, namespace interfaces.NameSpace) (any, error) {
	if list.Empty() {
		return nil, fmt.Errorf("cannot evaluate empty list")
	}
	symbol, ok := list.Car().(interfaces.Symbol)
	if !ok {
		return nil, fmt.Errorf("cannot use %v as a function: must be a symbol", list.Car())
	}

	value, ok := namespace.Get(symbol)
	if !ok {
		return nil, interfaces.UndefinedSymbolError(symbol)
	}

	function, ok := value.(interfaces.IndigoFunction)
	if !ok {
		return nil, interfaces.ExpectedButFoundTypeError("function", value)
	}

	expressionValue, err := function.Call(se, namespace, list.Cdr().(interfaces.List))
	if err != nil {
		return nil, err
	}

	return expressionValue, nil
}
