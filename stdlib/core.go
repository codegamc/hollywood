package stdlib

import (
	"fmt"

	environment "github.com/codegamc/hollywood/environment"
	printer "github.com/codegamc/hollywood/printer"
	types "github.com/codegamc/hollywood/types"
)

// imports the core functions into the environment (work in progress)
func importCore(envi *environment.Environment) *environment.Environment {
	// Prints a symbol in "pretty text"
	envi.Bind(types.MakeSymbol("core/print"), types.MakeFunc(print, "core/print"))
	// Creates a list
	envi.Bind(types.MakeSymbol("core/list"), types.MakeFunc(list, "core/list"))
	// returns true if a list
	envi.Bind(types.MakeSymbol("core/list?"), types.MakeFunc(listq, "core/list?"))
	// return the length of a list
	envi.Bind(types.MakeSymbol("core/count"), types.MakeFunc(count, "core/count"))

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

func count(args []types.HWType) types.HWType {
	if args[0].GetType() == types.LIST_TYPE {
		count := len(args[0].(types.HWList).Val)
		return types.MakeInt(int64(count))
	}
	return types.MakeInt(int64(-1))
}
