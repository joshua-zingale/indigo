package indigo

import (
	"fmt"

	"github.com/joshua-zingale/indigo/indigo/interfaces"
	"github.com/joshua-zingale/indigo/indigo/internal"
	"github.com/joshua-zingale/indigo/indigo/standard/evaluation"
	"github.com/joshua-zingale/indigo/indigo/standard/reading"
)

type IndigoInterpreter struct {
	Evaluator       interfaces.IndigoEvaluator
	GlobalNamespace interfaces.NameSpace
}

func NewStandardInterpreter() IndigoInterpreter {
	return IndigoInterpreter{
		Evaluator:       evaluation.NewStandardEvaluator(),
		GlobalNamespace: NewNameSpace(),
	}
}

func (ii *IndigoInterpreter) LoadModule(module interfaces.IndigoModule) {
	for _, symbol := range module.Symbols() {
		value, _ := module.Get(symbol)
		ii.GlobalNamespace.Set(symbol, value)
	}
}

func (ii *IndigoInterpreter) Eval(object any) (any, error) {
	return ii.Evaluator.Eval(object, ii.GlobalNamespace)
}

func Read(source string) (any, error) {
	lexer := reading.NewStandardReader(source)
	readObject, err := lexer.Read()
	if err != nil {
		return nil, err
	}

	if _, err := lexer.Read(); err == nil {
		return nil, fmt.Errorf("multiple objects read")
	}

	return readObject, nil
}

func NewCons(car any, cdr any) interfaces.Cons {
	return internal.NewCons(car, cdr)
}

func NewList(elements ...any) interfaces.List {
	return internal.NewList(elements...)
}

func Symbol(symbol string) interfaces.Symbol {
	return interfaces.Symbol(symbol)
}

func NewNameSpace() interfaces.NameSpace {
	return internal.NewNameSpace()
}

func NewNameSpaceFromMap(namespaceMap map[string]any) interfaces.NameSpace {
	return internal.NewNameSpaceFromMap(namespaceMap)
}
