package main

import (
	"fmt"

	"github.com/joshua-zingale/indigo/indigo"
	"github.com/joshua-zingale/indigo/indigo/functools"
	"github.com/joshua-zingale/indigo/indigo/standard/library"
)

func main() {

	interpreter := indigo.NewStandardInterpreter()
	interpreter.LoadModule(library.IndigoCore)

	read := functools.Must(indigo.Read("(+ (+ 1 2) 3 4.5)"))
	fmt.Println(read)
	fmt.Println(interpreter.Eval(read))
}
