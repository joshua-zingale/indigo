package internal

import "github.com/joshua-zingale/indigo/indigo/interfaces"

type NameSpace struct {
	parent    *NameSpace
	namespace map[interfaces.Symbol]any
}

func NewNameSpace() interfaces.NameSpace {
	return &NameSpace{
		parent:    nil,
		namespace: map[interfaces.Symbol]any{},
	}
}

func (ns *NameSpace) NewChild() interfaces.NameSpace {
	return &NameSpace{parent: ns, namespace: make(map[interfaces.Symbol]any)}
}

func (ns *NameSpace) Get(symbol interfaces.Symbol) (any, bool) {
	if v, ok := ns.namespace[symbol]; ok {
		return v, true
	}

	if ns.parent == nil {
		return nil, false
	}

	return ns.parent.Get(symbol)
}

func (ns *NameSpace) Set(symbol interfaces.Symbol, value any) {
	ns.namespace[symbol] = value
}
