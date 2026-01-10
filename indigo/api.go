package indigo

import (
	"github.com/joshua-zingale/indigo/indigo/interfaces"
	"github.com/joshua-zingale/indigo/indigo/internal"
	"github.com/joshua-zingale/indigo/indigo/reading"
)

func Read(source string) (any, error) {
	lexer := reading.NewStandardReader(source)
	return lexer.Read()
}

func NewCons(car any, cdr any) interfaces.Cons {
	return internal.NewCons(car, cdr)
}

func NewList(elements ...any) interfaces.Cons {
	return internal.NewList(elements...)
}

func Symbol(symbol string) interfaces.Symbol {
	return interfaces.Symbol(symbol)
}

func NewNameSpace() interfaces.NameSpace {
	return internal.NewNameSpace()
}
