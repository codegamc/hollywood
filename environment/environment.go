package environment

import (
	types "github.com/codegamc/hollywood/types"
)

// MakeEnvironment makes an environment, and can take expressions and bindings as an input
func MakeEnvironment(bindings []types.HWSymbol, expressions []types.HWType) Environment {
	envMap := make(map[string]types.HWType)
	env := Environment{parent: nil, environment: envMap}
	env.GiveStudio(MakeStudio())

	if len(bindings) != len(expressions) {
		return env // this should throw an error, but error catching is a slow process
		// there are way too many errors to catch in a programming language environment
		// that i just wont do any if it that isnt absolutely critical
	}

	for i := 0; i < len(bindings); i++ {
		env.Bind(bindings[i], expressions[i])
	}

	return env
}

// Environment is an Environment implementation
type Environment struct {
	parent      *Environment
	environment map[string]types.HWType
	Studio      *Studio
}

//GiveStudio ads a studio to the environment
func (e *Environment) GiveStudio(studio *Studio) {
	e.Studio = studio
}

// GetStudio returns the Studio, looking for the parent's first // we want to use only 1 studio if
// possible
func (e *Environment) GetStudio() *Studio {
	if e.parent != nil {
		return e.parent.GetStudio()
	}
	return e.Studio
}

// BindParent does something
func (e *Environment) BindParent(parent *Environment) *Environment {
	e.parent = parent
	return e
}

// Find recursively searches the tree for an environment containing a symbol
func (e *Environment) Find(key types.HWSymbol) *Environment {
	if _, ok := e.environment[key.Val]; ok {
		return e
	}
	if e.parent != nil {
		//fmt.Println("The environment has a parent")
		return e.parent.Find(key)
	}
	//fmt.Println("Find failed for " + key.Val)
	return nil
}

// Get retrieves a symbol from an env
func (e *Environment) Get(key types.HWSymbol) types.HWType {
	//fmt.Println("Getting: " + key.Val)
	env := e.Find(key)
	if env == nil {
		return types.MakeNull()
	}
	val, _ := env.environment[key.Val]
	return val
}

// Bind adds a value to the environment
func (e *Environment) Bind(key types.HWSymbol, value types.HWType) types.HWType {
	e.environment[key.Val] = value
	return value
}
