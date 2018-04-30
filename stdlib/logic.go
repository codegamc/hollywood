package stdlib

import (
	environment "github.com/codegamc/hollywood/environment"
	types "github.com/codegamc/hollywood/types"
)

// imports the logic functions into the environment (work in progress)
func importLogic(envi *environment.Environment) *environment.Environment {
	//not := types.MakeSymbol("!")

	envi.Bind(types.MakeSymbol("="), types.MakeFunc(equal, "equals"))
	envi.Bind(types.MakeSymbol("<="), types.MakeFunc(lteq, "less than or equal to"))
	envi.Bind(types.MakeSymbol(">="), types.MakeFunc(gteq, "greater than or equal to"))
	envi.Bind(types.MakeSymbol("<"), types.MakeFunc(lt, "less than"))
	envi.Bind(types.MakeSymbol(">"), types.MakeFunc(gt, "greater than"))
	envi.Bind(types.MakeSymbol("!"), types.MakeFunc(not, "not"))

	return envi
}

func equal(args []types.HWType) types.HWType {
	// INTEGER EQUALITY
	if args[0].GetType() == types.INT_TYPE {
		if args[1].GetType() == types.INT_TYPE {
			if args[0].(types.HWInt).Val == args[1].(types.HWInt).Val {
				return types.MakeBool(true)
			}
		}
	}

	// BOOLEAN EQUALITY
	if args[0].GetType() == types.BOOL_TYPE {
		if args[1].GetType() == types.BOOL_TYPE {
			if args[0].(types.HWBool).Val == args[1].(types.HWBool).Val {
				return types.MakeBool(true)
			}
		}
	}

	return types.MakeBool(false)
}

func lteq(args []types.HWType) types.HWType {
	// INTEGER EQUALITY
	if args[0].GetType() == types.INT_TYPE {
		if args[1].GetType() == types.INT_TYPE {
			if args[0].(types.HWInt).Val <= args[1].(types.HWInt).Val {
				return types.MakeBool(true)
			}
		}
	}

	return types.MakeBool(false)
}

func gteq(args []types.HWType) types.HWType {
	// INTEGER EQUALITY
	if args[0].GetType() == types.INT_TYPE {
		if args[1].GetType() == types.INT_TYPE {
			if args[0].(types.HWInt).Val >= args[1].(types.HWInt).Val {
				return types.MakeBool(true)
			}
		}
	}

	return types.MakeBool(false)
}

func lt(args []types.HWType) types.HWType {
	// INTEGER EQUALITY
	if args[0].GetType() == types.INT_TYPE {
		if args[1].GetType() == types.INT_TYPE {
			if args[0].(types.HWInt).Val < args[1].(types.HWInt).Val {
				return types.MakeBool(true)
			}
		}
	}

	return types.MakeBool(false)
}

func gt(args []types.HWType) types.HWType {
	// INTEGER EQUALITY
	if args[0].GetType() == types.INT_TYPE {
		if args[1].GetType() == types.INT_TYPE {
			if args[0].(types.HWInt).Val > args[1].(types.HWInt).Val {
				return types.MakeBool(true)
			}
		}
	}

	return types.MakeBool(false)
}

func not(args []types.HWType) types.HWType {
	if equal(args).(types.HWBool).Val == true {
		return types.MakeBool(false)
	}

	return types.MakeBool(true)
}
