package object

func NewEnclosedEnvironment(outer *Environment) *Environment {
	env := NewEnvironment()
	env.outer = outer
	return env
}

func NewEnvironment() *Environment {
	s := make(map[string]Entity)
	return &Environment{store: s, outer: nil}
}

type Entity struct {
	Object  Object
	Mutable bool
}

type Environment struct {
	store map[string]Entity
	outer *Environment
}

const NULL_NAME = "_"

var NULL_ENTITY = &Entity{Object: NULL, Mutable: false}

func (e *Environment) Get(name string) (*Entity, bool) {
	if name == NULL_NAME {
		return NULL_ENTITY, true
	}
	pair, ok := e.store[name]
	if !ok && e.outer != nil {
		return e.outer.Get(name)
	}
	return &pair, ok
}

func (e *Environment) Set(name string, val Object, mutable bool, initialize bool) (*Entity, bool) {
	if name == NULL_NAME {
		return NULL_ENTITY, true
	}

	if _, ok := e.store[name]; ok || initialize {
		entity := Entity{Object: val, Mutable: mutable}
		e.store[name] = entity
		return &entity, true
	} else if e.outer != nil {
		return e.outer.Set(name, val, mutable, false)
	} else {
		return nil, false
	}
}
