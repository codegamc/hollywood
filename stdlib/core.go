package stdlib

import (
	"fmt"

	environment "github.com/codegamc/hollywood/environment"
	printer "github.com/codegamc/hollywood/printer"
	types "github.com/codegamc/hollywood/types"
)

// imports the core functions into the environment (work in progress)
func importCore(envi *environment.Environment) *environment.Environment {

	//
	//list := types.MakeSymbol("core/list")
	// Checks if the next object is a list
	//listy := types.MakeSymbol("core/list?")
	// Checks if the next object is an empty list
	//empty := types.MakeSymbol("core/empty?")
	// Counts how many items are in the list
	//count := types.MakeSymbol("core/count")

	// Prints a symbol in "pretty text"
	envi.Bind(types.MakeSymbol("core/print"), types.MakeFunc(print, "print"))
	// Creates a list
	envi.Bind(types.MakeSymbol("core/list"), types.MakeFunc(list, "list"))
	// returns true if a list
	envi.Bind(types.MakeSymbol("core/list?"), types.MakeFunc(listq, "list?"))

	return envi
}

// Print is the core print function (core/print myAtom)
func print(args []types.HWType) types.HWType {
	fmt.Println(printer.PrintStr(args[0]))
	return types.MakeNullImplicit()
}

func list(args []types.HWType) types.HWType {
	return types.MakeList(args)
}

func listq(args []types.HWType) types.HWType {
	if args[0].GetType() == types.LIST_TYPE {
		return types.MakeBool(true)
	}
	return types.MakeBool(false)
}
