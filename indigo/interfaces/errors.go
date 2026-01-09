package interfaces

import (
	"fmt"
	"reflect"
)

func ExpectedButFoundTypeError(expectedType string, foundValue any) error {
	return fmt.Errorf("expected value of type '%s' but got '%v' of type '%v'", expectedType, foundValue, reflect.TypeOf(foundValue))
}

func UndefinedSymbolError(symbol Symbol) error {
	return fmt.Errorf("undefined symbol '%s'", symbol)
}
