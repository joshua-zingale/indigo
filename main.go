package main

import (
	"fmt"

	"github.com/joshua-zingale/indigo/indigo"
)

func main() {
	v, _ := indigo.Read("(+ 2 (* 3 1))")
	fmt.Printf("%v", v)
}
