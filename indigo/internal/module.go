package internal

import "github.com/joshua-zingale/indigo/indigo/interfaces"

type Module struct {
	interfaces.NameSpace
	name string
}

func NewModuleWithName(name string) interfaces.IndigoModule {
	return &Module{
		NameSpace: NewNameSpace(),
		name:      name,
	}
}

func NewModule(name string, namespace interfaces.NameSpace) interfaces.IndigoModule {
	return &Module{
		NameSpace: namespace,
		name:      name,
	}
}

func (m *Module) Name() string {
	return m.name
}
