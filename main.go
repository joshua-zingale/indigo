package main

import (
	"fmt"

	"github.com/joshua-zingale/indigo/indigo"
	"github.com/joshua-zingale/indigo/indigo/evaluation"
	"github.com/joshua-zingale/indigo/indigo/interfaces"
	"github.com/joshua-zingale/indigo/indigo/reading"
)

func main() {
	evaluator := evaluation.NewStandardEvaluator()
	namespace := indigo.NewNameSpace()
	namespace.Set(interfaces.Symbol("+"), evaluation.NewTypeCheckedIndigoFunctionFromGo(func(a int, b int) (int, error) {
		return a + b, nil
	}))
	namespace.Set(interfaces.Symbol("*"), evaluation.NewTypeCheckedIndigoFunctionFromGo(func(a int, b int) (int, error) {
		return a * b, nil
	}))
	namespace.Set(interfaces.Symbol("Bob"), "WOW!")

	syntax, err := reading.NewStandardReader("(+ 1 3)").Read()

	found, err := evaluator.Eval(syntax, namespace)
	if err != nil {
		panic(err)
	}

	fmt.Println(found)
}
