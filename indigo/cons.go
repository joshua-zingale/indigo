package indigo

type Cons struct {
	car interface{}
	cdr interface{}
}

func NewCons(car interface{}, cdr interface{}) Cons {
	return Cons{car: car, cdr: cdr}
}

func EmptyCons() Cons {
	return Cons{}
}

func (c *Cons) Car() interface{} {
	return c.car
}

func (c *Cons) Cdr() interface{} {
	return c.cdr
}
