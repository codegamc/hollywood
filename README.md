# "Hollywood" By Michael Codega

 Hollywood is a simple lisp interpreter written in go. It is currently only functional as a REPL. Its name comes from the "Actor Model" of concurrency, which hollywood is loosely trying to implement.  The Actor Model is based on a 1978 paper by Tony Hoare, where he described "Communicating Sequential Processes".  Hollywood does not implement a communication method between the "Actors" or coroutines.  

Hollywood has 9 reserved words. These words are 
* `var`, used for defining variables in the environment.
* `do`, used to describe sequentially executed statements.
* `if`, used to construct conditional statements.
* `fn*`, used to define anonymous functions.
* `func`, used to define named functions.
* `map`, used to map a function onto a list of values. This runs each evaluation in parallel.
* `s/map`, a serial version of map. 
* `while`, a looping mechanism that loops while a condition holds true
* `act`, which runs a function on a seperate coroutine
  
The two map functions also have built in timers to calculate how long each takes to execute.  There is more overhead for the parallel map to run, so for simple functions, or small lists, its often faster to run in serial. 
 
# Building and Running
### Building
Since compiling a go binary requires installing a go compiler and setting up an environment, a linux compatable binary is included in this repository. It is called `hollywood-linux` and it works on the EECS Server. 

A makefile is included, that can compile either for MacOS or linux with respeective commands: `make mac` or `make linux`, assuming that proper golang build toolchains are installed, which is not the case on the EECS Server.

### Running 
Assuming that the linux binary is being run on the EECS Server, simply run `./hollywood-linux` and the REPL prompt will appear right away. 

> Please note that there is almost no error checking, so something like a missed closing paren can crash the imterpreter.

# Examples to run

Lisp, the basis for the hollywood language looks different compared to most languages. It is structed very differently, and uses only parens to organize the code.

```lisp
(myFunction arg1 arg2 ... arg_n )

To build complex expressions, simply nest them.

(myFunction arg1 (myOtherFunction _arg1 _arg2) arg3)
                 |~~~~~~~~~~~~~~~~~~~~~~~~~~~|
                    This all evaluates, then becomes the second 
                    argument to myFunction
```

### Simple tests
###### Math 
Only simple integer math is defined, and there are no error-checks for mathematically undefined behavior. Something like Divide-By-Zero will crash the interpeter.

* ` (+ 1 1) ` => This returns  Int: 2
* ` (- 1 1) ` => This returns  Int: 0
* ` (* 12 12)` => This returns Int: 144
* ` (/ 24 2)` => This returns Int: 12
* ` (/ (* (+ 6 6) (- 24 12)) 12)` => This returns Int: 12

###### More complex expressions
More complex expressions are also valid. 
You can define a function:
* `(func square (a) (* a a))` Here we are defining the function square, which takes the argument `a` and returns the result of `a` multiplied by `a`.
*  `(square 3)` => this returns Int: 6
