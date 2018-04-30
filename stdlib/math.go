package stdlib

import (
	"github.com/codegamc/hollywood/environment"
	types "github.com/codegamc/hollywood/types"
)

// imports the math functions into the environment (work in progress)
func importMath(envi *environment.Environment) *environment.Environment {
	envi.Bind(types.MakeSymbol("*"), types.MakeFunc(MULTIPLY, "*"))
	envi.Bind(types.MakeSymbol("+"), types.MakeFunc(ADD, "+"))
	envi.Bind(types.MakeSymbol("-"), types.MakeFunc(SUBTRACT, "-"))
	envi.Bind(types.MakeSymbol("/"), types.MakeFunc(DIVIDE, "/"))

	return envi
}

//MULTIPLY is the multiply builtin
func MULTIPLY(args []types.HWType) types.HWType {
	result := int64(1)
	for i := 0; i < len(args); i++ {
		if args[i].GetType() != types.INT_TYPE {
			return types.MakeNull()
		}
		result = result * args[i].(types.HWInt).Val
	}
	return types.MakeInt(result)
}

// ADD is the add builtin
func ADD(args []types.HWType) types.HWType {
	result := int64(0)
	for i := 0; i < len(args); i++ {
		if args[i].GetType() != types.INT_TYPE {
			return types.MakeNull()
		}
		result = result + args[i].(types.HWInt).Val
	}
	return types.MakeInt(result)
}

// SUBTRACT is the subtract builtin
func SUBTRACT(args []types.HWType) types.HWType {
	if args[0].GetType() != types.INT_TYPE {
		return types.MakeNull()
	}
	result := int64(args[0].(types.HWInt).Val)
	for i := 1; i < len(args); i++ {
		if args[i].GetType() != types.INT_TYPE {
			return types.MakeNull()
		}
		result = result - args[i].(types.HWInt).Val
	}
	return types.MakeInt(result)
}

// DIVIDE is the divide builtin
func DIVIDE(args []types.HWType) types.HWType {
	if args[0].GetType() != types.INT_TYPE {
		return types.MakeNull()
	}
	result := int64(args[0].(types.HWInt).Val)
	for i := 1; i < len(args); i++ {
		if args[i].GetType() != types.INT_TYPE {
			return types.MakeNull()
		}
		result = result / args[i].(types.HWInt).Val
	}
	return types.MakeInt(result)
}
