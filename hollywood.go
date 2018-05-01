package main

import (
	"bufio"
	"fmt"
	"os"
	"time"

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

The reader (parsing) is in the "reader" package. It is a simple REGEX based LISP parser.
The regex breaks up the input string into tokens, and then parses those tokens into an
Abstract Syntax Tree

The "Eval" (interpreting) is in the main package, but could be moved in the future.  The
scope portion of the interpreting is in the Environmental package.  The stdlib functions
are in the stdlib package (printing, equality, etc). The determining factor for is something
should be in the Eval or Stdlib is "Should this be considered a reserved word", and if it is
a reserved word, it goes in the eval switch statement instead of in the standard library of
functions.

The "Print" (printing) is in the "printer" package, but it is very small for now.  It would
get more complicated should a "pretty printing" be implemented, where formatting is applied
to the string as it gets printed.  This is not an issue as long as strings are not fully
implemented into the type system (they are, but there is nothing special about them) eg. no
escape characters, etc
*/

// This is the entrypoint into the code
func main() {

	// create a new environment, its the root envi.
	envi := environment.MakeEnvironment(nil, nil)
	// set it as the base, this means that it will do things to it like
	// set the standard library into the namespace
	envi = setBaseEnv(envi)

	// create the CLI reader
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

	// if its not a list, then you cannot evaluate the expression as a function
	// so send it to EvaluateAST
	if ast.GetType() != types.LIST_TYPE {
		return evaluateAST(ast, envi)
	}
	// ast is a list: // this will never be done, its not worth the work
	// the parser just crashes at list-length = 0
	// if empty: return

	//else, evaluate the list assuming the first word is either
	// 	(1. a control word (see switch statement))
	//		 so handle accordingly
	// 	(2. a function (stores as a symbol))
	//		 so call it on the rest of the list as args
	if ast.(types.HWList).Val[0].GetType() == types.SYM_TYPE {
		sym := ast.(types.HWList).Val[0].(types.HWSymbol)
		switch sym.Val {
		//its a define statement that adds a value to the env.
		case "var":
			key := ast.(types.HWList).Val[1].(types.HWSymbol)
			value := EVAL(ast.(types.HWList).Val[2], envi)
			val := envi.Bind(key, value)
			return val

		// (do (eval_1) (eval_2) ... (eval_n) (eval_final)) => returns the result of eval_final
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
		// ( if (condition) (eval_this) (else_eval_this) )
		// this uses the types.NotFalsey method, which is probably not robust
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

			// no else statement, return explicit null
			return types.MakeNull()
		// it is a lambda function
		// (	(fn* (args_symbols) (eval_this) ) _args_ )
		// it generates an anonymous function, that has no symbol and does not bind to the environment
		// it does however, have a generated parent environment, so as it gets passed around, it holds
		// a constant environment parent
		case "fn*":
			argsSymList := ast.(types.HWList).Val[1].(types.HWList)
			argsSym := []types.HWSymbol{}
			// generating the list of symbols
			for i := 0; i < len(argsSymList.Val); i++ {
				if argsSymList.Val[i].GetType() == types.SYM_TYPE {
					argsSym = append(argsSym, argsSymList.Val[i].(types.HWSymbol))
				}
			}
			// evaluation AST
			evalThis := ast.(types.HWList).Val[2]

			// generated HWFunction
			f := types.MakeFunc(func(args []types.HWType) types.HWType {
				newEnvi := environment.MakeEnvironment(argsSym, args)
				newEnvi = *(&newEnvi).BindParent(envi)
				return evaluateAST(evalThis, &newEnvi)
			}, "__anon__")

			return f
		// (func symbol_to_represent_function (LIST_ARGS_SYMBOLS) (EVAL_THIS))
		// explicitly defining a function, that has an args list
		case "func":
			// ast.().Val[0] == "func"
			symbol := ast.(types.HWList).Val[1].(types.HWSymbol)
			argsSymList := ast.(types.HWList).Val[2].(types.HWList)

			// for each symbol in the args list, ensure that it is a symbol, and add it to
			// the args being passed into the function
			argsSym := []types.HWSymbol{}
			for i := 0; i < len(argsSymList.Val); i++ {
				if argsSymList.Val[i].GetType() == types.SYM_TYPE {
					argsSym = append(argsSym, argsSymList.Val[i].(types.HWSymbol))
				}
			}

			// the AST being evaluated when the function is run
			evalThis := ast.(types.HWList).Val[3]

			// this is the HWFunction that gets run when the function is run
			f := types.MakeFunc(func(args []types.HWType) types.HWType {
				newEnvi := environment.MakeEnvironment(argsSym, args)
				newEnvi = *(&newEnvi).BindParent(envi)
				return EVAL(evalThis, &newEnvi)
			}, symbol.Val)

			// binding the function to the environment
			envi.Bind(symbol, f)

			// returning the function
			return f
		// maps a function to act on each item in a list
		// (map myFunction myList)
		// 		=> returns a list of results, which is each item in myList,
		//			 being input to myFunction
		case "map":

			start := time.Now()
			// its another symbol
			//map := ast.(types.HWList).Val[0] // the map symbol
			function := ast.(types.HWList).Val[1]
			sequence := ast.(types.HWList).Val[2]

			// This is more parallelism than concurrency... oh well
			// take each item in the sequence, and make a goroutine to handle that item
			// make a channel to allow that item to return data when its done
			// then run the goroutines
			length := len(sequence.(types.HWList).Val)
			returnDataChannels := make([](chan types.HWType), length)
			for i := 0; i < length; i++ {
				//
				c := make(chan types.HWType)
				returnDataChannels[i] = c
				go func(c chan types.HWType, sequenceElement types.HWType) {
					//
					//fmt.Println("This is a mapped function running in its own thread")

					// take function, and call it on the correct part of the sequence
					f := envi.Get(function.(types.HWSymbol))
					args := make([]types.HWType, 1)
					args[0] = sequenceElement
					c <- f.(types.HWFunc).Val(args)
				}(c, sequence.(types.HWList).Val[i])
			}

			// create an array to store the return data, and collect it from the channels
			// note that it blocks on sequence(id: i), so if a channel has no data, it will wait.
			// maybe there should be a timeout or something, since a malformed input can lead to
			// no output ever, and it holds up everything
			returnData := make([]types.HWType, length)
			for i := 0; i < length; i++ {
				returnData[i] = <-returnDataChannels[i]
			}

			t := time.Now()
			elapsed := t.Sub(start)
			fmt.Println("The time that was taken to complete map in parallel is: " + elapsed.String())

			return types.MakeList(returnData)

			// spin up a channel to store the return values
			// eval each item in a seperate gorouting
			// collect the results, and store them

		// This is a version of map that runs in serial, and does not take advantage of threads
		case "s/map":
			// its another symbol
			//map := ast.(types.HWList).Val[0] // the map symbol
			start := time.Now()
			function := ast.(types.HWList).Val[1]
			sequence := ast.(types.HWList).Val[2]
			length := len(sequence.(types.HWList).Val)
			f := envi.Get(function.(types.HWSymbol))
			returnData := make([]types.HWType, length)

			for i := 0; i < length; i++ {
				sequenceElement := sequence.(types.HWList).Val[i]
				args := make([]types.HWType, 1)
				args[0] = sequenceElement
				returnData[i] = f.(types.HWFunc).Val(args)
			}

			t := time.Now()
			elapsed := t.Sub(start)
			fmt.Println("The time that was taken to complete map in series is: " + elapsed.String())
			return types.MakeList(returnData)

		// its a while loop
		// (while (condition) (eval_this_every_loop))
		case "while":
			//while := ast.(types.HWList).Val[0] // the map symbol
			condition := ast.(types.HWList).Val[1]
			loopBody := ast.(types.HWList).Val[2]

			// evaluate the condition, if its false from the start, dont run it
			cond := EVAL(condition, envi)
			// repeatedly evaluate the loopbody while the condition is not false
			for types.NotFalsey(cond) {
				EVAL(loopBody, envi)
				cond = EVAL(condition, envi)
			}
			// returns an implicit null (has no printed value)
			return types.MakeNullImplicit()

		// this is how true concurrency will be acheived here...
		// act is the keyword for triggering the "green thread" goroutine
		// an actor is an independant "thread" that has a "mailbox" that can receive messages,
		// actors can also send messages
		// (act myFunc args) => resolves to an int (thread-id)
		case "act":
			//act := ast.(types.HWList).Val[0] // the map symbol
			function := ast.(types.HWList).Val[1]
			args := ast.(types.HWList).Val[2]
			// this environment is unique to the actor
			// the actorGlobal is used to ensure that certain values can be set to the actor that
			// cannot be over written (nothing implemented though)
			actorGlobal := environment.MakeEnvironment(nil, nil)
			actorEnvi := environment.MakeEnvironment(nil, nil)
			actorEnvi.BindParent(&actorGlobal)
			// it still has a common parent, since that actor might want to use globals (like the stdlib)
			// the only risk here is that some action on the global level changes what actors see,
			// but this a problem for the writers of the code. It is a way to signal to actors to change
			// some sort of behavior? it would be a "read only" communication from the actors perspective
			actorEnvi = *(&actorEnvi).BindParent(envi)

			// this is where the studio should generate an actor
			// this does not actually do anything, but it could do stuff in the future
			actor := envi.GetStudio().NewActor()

			// this is where the new AST to be evaluated is created
			actorASTlist := []types.HWType{}
			actorASTlist = append(actorASTlist, function)
			for i := 0; i < len(args.(types.HWList).Val); i++ {
				actorASTlist = append(actorASTlist, args.(types.HWList).Val[i])
			}

			// converting that into a hollywood type
			actorAST := types.MakeList(actorASTlist)

			// this is where the new corouting is created, and run, its a non-blocking call
			// that starts the "greenthread" or goroutine (coroutine) that is managed by the
			// golang runtime.  This means that actors will be multiplexed across several of the
			// operating system's threads, to ensure that nothing is blocking
			go EVAL(actorAST, &actorEnvi)

			// returning some value that represents the actor, so it can be used later
			// currently, there is nothing that can be done with such an ID, but one day maybe
			return types.MakeInt(actor.ActorID)

		// its just a symbol
		default:
			// or just evaluate the list, with the first term as a function
			// it could be that this term is a symbol, and it is a created function
			list := evaluateAST(ast, envi)
			var result types.HWType
			result = list.(types.HWList).Val[0]

			// this breaks out the function from the rest of the AST
			listFunction := result.(types.HWFunc).Val

			// This breaks the args out of the AST and formats them as an array of inputs for
			// functions, since that is the standard interface
			args := []types.HWType{}
			for i := 1; i < len(list.(types.HWList).Val); i++ {
				args = append(args, list.(types.HWList).Val[i])
			}

			// calling the actual function
			results := listFunction(args)
			return results

		}
	}
	//
	return ast
}

// PRINT is the print func for the REPL
func PRINT(exp types.HWType) string {
	return printer.PrintStr(exp)
}

//  This evaluates symbols in accordance with the environment
func evaluateAST(ast types.HWType, envi *environment.Environment) types.HWType {
	//this takes a symbol and resolves it based on the environment, so things defined as vars
	// become their value at runtime
	if ast.GetType() == types.SYM_TYPE {
		// if found in env, return, else raise error
		val := envi.Get(ast.(types.HWSymbol))
		if val.GetType() != types.NULL_TYPE {
			//fmt.Println("Evaluated symbol: " + ast.(types.HWSymbol).Val + " and returned: " + val.GetMeta())
			return val
		}
		return ast
	}

	// this takes a list, and evaluates each item in the list, in order to resolve that
	// its how recursively defined expressions get handled
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

// adds global stuff like stdlib to the base environment (currently only stdlib)
func setBaseEnv(envi environment.Environment) environment.Environment {
	stdlib.ImportStdLib(envi)
	return envi
}
