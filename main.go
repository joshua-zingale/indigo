package main

import "fmt"

func main() {
	var a int
	b := &a

	*b = 3

	fmt.Println(a)
}
