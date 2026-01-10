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

func NewNameSpaceFromMap(namespaceMap map[string]any) interfaces.NameSpace {
	namespace := NewNameSpace()
	for key, value := range namespaceMap {
		namespace.Set(interfaces.Symbol(key), value)
	}
	return namespace
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

func (ns *NameSpace) Symbols() []interfaces.Symbol {
	var symbols []interfaces.Symbol

	curr_namespace := ns
	for curr_namespace != nil {
		keys := make([]interfaces.Symbol, 0, len(ns.namespace))
		for s := range ns.namespace {
			keys = append(keys, s)
		}
		symbols = append(symbols, keys...)

		curr_namespace = ns.parent
	}

	return symbols
}
