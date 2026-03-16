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

func (e *Environment) Get(name string) (Entity, bool) {
	pair, ok := e.store[name]
	if !ok && e.outer != nil {
		return e.outer.Get(name)
	}
	return pair, ok
}
func (e *Environment) Set(name string, val Object, mutable bool) (Entity, bool) {
	e.store[name] = Entity{Object: val, Mutable: mutable}
	return e.Get(name)
}
