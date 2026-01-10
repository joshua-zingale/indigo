package evaluation

import (
	"fmt"

	"github.com/joshua-zingale/indigo/indigo/interfaces"
)

type StandardEvaluator struct{}

func NewStandardEvaluator() interfaces.IndigoEvaluator {
	return &StandardEvaluator{}
}

func (se *StandardEvaluator) Eval(object any, namespace interfaces.NameSpace) (any, error) {
	return se.evalInNamespace(object, namespace)
}

func NewNameSpacedEval(namespace interfaces.NameSpace) func(any) (any, error) {
	evaluator := NewStandardEvaluator()
	return func(a any) (any, error) {
		return evaluator.Eval(a, namespace)
	}
}

func (se *StandardEvaluator) evalInNamespace(object any, namespace interfaces.NameSpace) (any, error) {
	switch typedObject := object.(type) {
	case interfaces.Symbol:
		if value, ok := namespace.Get(typedObject); ok {
			return value, nil
		}
		return nil, interfaces.UndefinedSymbolError(typedObject)
	case interfaces.Cons:
		if list, ok := typedObject.(interfaces.List); ok {
			return se.evalList(list, namespace)
		}
		return typedObject, nil
	default:
		return typedObject, nil
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

	objectValue, err := function.Call(se, namespace, list.Cdr().(interfaces.List))
	if err != nil {
		return nil, err
	}

	return objectValue, nil
}
