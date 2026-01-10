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
		if list, err := consToSlice(typedExpression); err == nil {
			return se.evalList(list, namespace)
		}
		return typedExpression, nil
	default:
		return typedExpression, nil
	}
}

func (se *StandardEvaluator) evalList(list []any, namespace interfaces.NameSpace) (any, error) {
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

	expressionValue, err := function.Call(se, namespace, list[1:])
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
