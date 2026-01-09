package internal

import (
	"reflect"

	"github.com/joshua-zingale/indigo/indigo/interfaces"
)

func IndigoEqual(a any, b any) bool {
	if cons1, ok := a.(interfaces.Cons); ok {
		if cons2, ok := b.(interfaces.Cons); ok {
			return cons1.Empty() && cons2.Empty() || IndigoEqual(cons1.Car(), cons2.Car()) && IndigoEqual(cons1.Cdr(), cons2.Cdr())
		}
		return false
	}

	return reflect.DeepEqual(a, b)
}
