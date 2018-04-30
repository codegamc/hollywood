package main

import (
	"bufio"
	"fmt"
	"os"

	// Note: `codegamc` is me, as this is code written for the project
	// also, this is actually the import path. In golang, if the import path
	// is a URL, then it can automatically "go get" the imports from that URL
	// and it makes tracking the origin of the code better, and its easier for
	// sharing to the EECS shared server, since it automatically pulls it in over
	// git (when ready - currently not uploaded)
	environment "github.com/codegamc/hollywood/environment"
	printer "github.com/codegamc/hollywood/printer"
	reader "github.com/codegamc/hollywood/reader"
	stdlib "github.com/codegamc/hollywood/stdlib"
	types "github.com/codegamc/hollywood/types"
)

/*

Some structural information:

The interpreter is written as a "Read Evaluate Print" loop

The reader (parsing) is in the "reader" package. It is a simple REGEX based LISP parser.  The regex breaks up the input string into tokens,
and then parses those tokens into an Abstract Syntax Tree

The "Eval" (interpreting) is in the main package, but could be moved in the future.  The scope portion of the interpreting is in the
Environmental package.  The stdlib functions are in the stdlib package (printing, equality, etc). The determining factor for is something
should be in the Eval or Stdlib is "Should this be considered a reserved word", and if it is a reserved word, it goes in the eval switch
statement instead of in the standard library of functions

The "Print" (printing) is in the "printer" package, but it is very small for now.  It would get more complicated should a "pretty printing"
be implemented, where formatting is applied to the string as it gets printed.
*/

// This is the entrypoint into the code
func main() {

	envi := environment.MakeEnvironment(nil, nil)
	envi = setBaseEnv(envi)

	commandLineReader := bufio.NewReader(os.Stdin)
	// execute infinitely...
	for {
		fmt.Print("hollywood>> ")
		text, _ := commandLineReader.ReadString('\n')
		fmt.Println(rep(text, envi))
	}
}

// this is the Read/Evaluate/Print command
func rep(code string, envi environment.Environment) string {
	return PRINT(EVAL(READ(code), &envi))
}

// READ is the read function that turns strings into a tree of HWLists
func READ(str string) types.HWType {
	result := reader.ReadStr(str)
	return result
}

// EVAL is the evaluation function. It takes the Abstract Syntax Tree,
// and it evaluates it against the input environment.  It returns another AST
// that is printed or evaluated later. It can either call evaluateAST if the
// object it is acting on is a symbol, in order to determine the meaning of that symbol,
// or it switches and evaluates it if the input is a list,
// making certain evaluations in the function
func EVAL(ast types.HWType, envi *environment.Environment) types.HWType {

	if ast.GetType() != types.LIST_TYPE {
		return evaluateAST(ast, envi)
	}
	// ast is a list:
	// if empty: return

	//else:
	if ast.(types.HWList).Val[0].GetType() == types.SYM_TYPE {
		sym := ast.(types.HWList).Val[0].(types.HWSymbol)
		switch sym.Val {
		//its a define statement that ads a value to the env.
		case "def!":
			key := ast.(types.HWList).Val[1].(types.HWSymbol)
			value := EVAL(ast.(types.HWList).Val[2], envi)
			val := envi.Bind(key, value)
			return val

		//it is a namespace creation statement
		case "let*":
			newEnvi := environment.MakeEnvironment(nil, nil)
			newEnvi = *(&newEnvi).BindParent(envi)
			bindingList := ast.(types.HWList).Val[1]
			if bindingList.GetType() == types.LIST_TYPE {
				list := bindingList.(types.HWList)
				// list is even length, so it is a set of pairs
				if len(list.Val)%2 == 0 {
					for i := 0; i < len(list.Val); i = i + 2 {
						//ensuring that the key is a symbol
						if list.Val[i].GetType() == types.SYM_TYPE {
							newEnvi.Bind(list.Val[i].(types.HWSymbol), EVAL(list.Val[i+1], &newEnvi))
						}
					}
				}
			}

			evalTerm := ast.(types.HWList).Val[2]
			val := EVAL(evalTerm, &newEnvi)
			return val
		// it is an iterative do statement
		case "do":
			// this should evaluate each element in the list, return final element
			length := len(ast.(types.HWList).Val)
			// deal with all of the statemnets in the do block
			for i := 1; i < length-1; i++ {
				//evaluateAST(ast.(types.HWList).Val[i], envi)
				EVAL(ast.(types.HWList).Val[i], envi)
			}
			// return the last one, for printing
			//return evaluateAST(ast.(types.HWList).Val[length-1], envi)
			return EVAL(ast.(types.HWList).Val[length-1], envi)

		// it is an if statement
		case "if":
			cond := EVAL(ast.(types.HWList).Val[1], envi)
			// CONDITION MET
			if types.NotFalsey(cond) {
				return EVAL(ast.(types.HWList).Val[2], envi)
			}
			// ELSE STATEMENT
			if len(ast.(types.HWList).Val) > 3 {
				return EVAL(ast.(types.HWList).Val[3], envi)
			}
			// no else statement
			return types.MakeNull()
		// it is a lambda function
		case "fn*":

			//(  (fn* (args_symbol) (eval_this) )  _args_  )
			argsSymList := ast.(types.HWList).Val[1].(types.HWList)
			argsSym := []types.HWSymbol{}
			for i := 0; i < len(argsSymList.Val); i++ {
				if argsSymList.Val[i].GetType() == types.SYM_TYPE {
					argsSym = append(argsSym, argsSymList.Val[i].(types.HWSymbol))
				}
			}
			evalThis := ast.(types.HWList).Val[2]

			f := types.MakeFunc(func(args []types.HWType) types.HWType {
				newEnvi := environment.MakeEnvironment(argsSym, args)
				newEnvi = *(&newEnvi).BindParent(envi)
				printer.PrintStr(evalThis)
				return evaluateAST(evalThis, &newEnvi)
			}, "__anon__")

			return f
			// (defun symbol (LIST_ARGS_SYMBOLS) (EVAL_THIS))
		// explicitly defining a function, that has an arg list
		case "defun":
			// ast.().Val[0] == "defun"
			symbol := ast.(types.HWList).Val[1].(types.HWSymbol)
			argsSymList := ast.(types.HWList).Val[2].(types.HWList)

			argsSym := []types.HWSymbol{}
			for i := 0; i < len(argsSymList.Val); i++ {
				if argsSymList.Val[i].GetType() == types.SYM_TYPE {
					argsSym = append(argsSym, argsSymList.Val[i].(types.HWSymbol))
				}
			}

			evalThis := ast.(types.HWList).Val[3]

			f := types.MakeFunc(func(args []types.HWType) types.HWType {
				newEnvi := environment.MakeEnvironment(argsSym, args)
				newEnvi = *(&newEnvi).BindParent(envi)

				//printer.PrintStr(evalThis)
				//fmt.Println("Type: " + evalThis.(types.HWList).Val[0].GetMeta())
				return EVAL(evalThis, &newEnvi)
			}, symbol.Val)

			envi.Bind(symbol, f)

			return f
		// maps a function to act on each item in a list
		// (map myFunction myList)
		// 		=> returns a list of results
		case "map":
			// its another symbol
			//map := ast.(types.HWList).Val[0] // the map symbol
			function := ast.(types.HWList).Val[1]
			sequence := ast.(types.HWList).Val[2]

			// This is more parallelism than concurrency... oh well
			length := len(sequence.(types.HWList).Val)
			returnDataChannels := make([](chan types.HWType), length)
			for i := 0; i < length; i++ {
				//
				c := make(chan types.HWType)
				returnDataChannels[i] = c
				go func(c chan types.HWType, sequenceElement types.HWType) {
					// take function, and call it on the correct part of the sequence
					f := envi.Get(function.(types.HWSymbol))
					args := make([]types.HWType, 1)
					args[0] = sequenceElement
					c <- f.(types.HWFunc).Val(args)
				}(c, sequence.(types.HWList).Val[i])
			}

			returnData := make([]types.HWType, length)
			for i := 0; i < length; i++ {
				returnData[i] = <-returnDataChannels[i]
			}

			return types.MakeList(returnData)

			// spin up a channel to store the return values
			// eval each item in a seperate gorouting
			// collect the results, and store them as a

		// its a while loop
		// (while (condition) (eval_this_every_loop))
		case "while":
			//while := ast.(types.HWList).Val[0] // the map symbol
			condition := ast.(types.HWList).Val[1]
			loopBody := ast.(types.HWList).Val[2]

			cond := EVAL(condition, envi)
			for types.NotFalsey(cond) {
				EVAL(loopBody, envi)
				cond = EVAL(condition, envi)
			}
			return types.MakeNullImplicit()

		// this is how true concurrency will be acheived here...
		// ACT is the keyword for triggering the "actor model" style of concurrency
		// an actor is an independant "thread" that has a "mailbox" that can receive messages,
		// actors can also send messages
		// (act myFunc args) => resolves to an int (thread-id)
		case "act":
			//act := ast.(types.HWList).Val[0] // the map symbol
			function := ast.(types.HWList).Val[1]
			args := ast.(types.HWList).Val[2]
			// this environment is unique to the actor
			actorGlobal := environment.MakeEnvironment(nil, nil)
			//actorGlobal = setBaseEnv(actorGlobal)
			actorEnvi := environment.MakeEnvironment(nil, nil)
			actorEnvi.BindParent(&actorGlobal)
			// it still has a common parent, since that actor might want to use globals (like the stdlib)
			// the only risk here is that some action on the global level changes what actors see,
			// but this a problem for the writers
			actorEnvi = *(&actorEnvi).BindParent(envi)
			// this is where the studio should generate an actor number, and add it to the actor's envi
			actor := envi.GetStudio().NewActor()
			actorEnvi.Bind(types.MakeSymbol("actor/id"), types.MakeInt(int64(actor.ActorID)))
			actorEnvi.Bind(types.MakeSymbol("actor/mail"), types.MakeFunc(func(args []types.HWType) types.HWType {
				// this is the get mail function
				nextMail := actor.GetMail()
				return nextMail
			}, "actor/mail"))

			// (actor/send (address message))
			actorEnvi.Bind(types.MakeSymbol("actor/send"), types.MakeFunc(func(args []types.HWType) types.HWType {
				address := args[0]
				//message := args[1]
				messageCouplet := args[0:2]
				actor.SendMail(messageCouplet)
				return address
			}, "actor/send"))

			// this is where the new AST to be evaluated is created
			actorASTlist := []types.HWType{}
			actorASTlist = append(actorASTlist, function)
			for i := 0; i < len(args.(types.HWList).Val); i++ {
				actorASTlist = append(actorASTlist, args.(types.HWList).Val[i])
			}

			// converting that into a hollywood type
			actorAST := types.MakeList(actorASTlist)

			// this is where the new corouting is created, and run
			go EVAL(actorAST, &actorEnvi)

			// returning some value that represents the actor, so it can be used later
			return types.MakeInt(actor.ActorID)

		// its just a symbol
		default:
			// or just evaluate the list, with the first term as a function
			// it could be that this term is a symbol, and it is a created function
			list := evaluateAST(ast, envi)
			var result types.HWType
			result = list.(types.HWList).Val[0]
			fmt.Println(printer.PrintStr(result))

			if result.GetType() == types.SYM_TYPE {
				// this problem again...
				fmt.Println(printer.PrintStr(envi.Get(result.(types.HWSymbol))))
			}

			listFunction := result.(types.HWFunc).Val
			//func_(list.(types.HWList).Val[1:])
			args := []types.HWType{}
			for i := 1; i < len(list.(types.HWList).Val); i++ {
				args = append(args, list.(types.HWList).Val[i])
			}

			results := listFunction(args)
			return results

		}
	}
	return ast
}

// PRINT is the print func
func PRINT(exp types.HWType) string {
	return printer.PrintStr(exp)
}

//  This evaluates symbols in accordance with the environment
func evaluateAST(ast types.HWType, envi *environment.Environment) types.HWType {
	//fmt.Println(ast.ToString())
	if ast.GetType() == types.SYM_TYPE {
		// if found in env, return, else raise error
		val := envi.Get(ast.(types.HWSymbol))
		if val.GetType() != types.NULL_TYPE {
			//fmt.Println("Evaluated symbol: " + ast.(types.HWSymbol).Val + " and returned: " + val.GetMeta())
			return val
		}
		return ast
	}

	if ast.GetType() == types.LIST_TYPE {
		// evaluate ast for each of the items in the list...
		length := len(ast.(types.HWList).Val)
		list := make([]types.HWType, length)
		//list := []types.HWType{}
		for i := 0; i < length; i++ {
			list[i] = EVAL(ast.(types.HWList).Val[i], envi)
		}
		return types.MakeList(list)
	}

	return ast
}

// adds global stuff like stdlib to the base environment
func setBaseEnv(envi environment.Environment) environment.Environment {
	stdlib.ImportStdLib(envi)

	return envi
}
