# "Hollywood"

 Hollywood is a simple lisp interpreter written in go. It is currently only functional as a REPL. Its name comes from the "Actor Model" of concurrency, which hollywood is loosely trying to implement.  The Actor Model is based on a 1978 paper by Tony Hoare, where he described "Communicating Sequential Processes".  Hollywood does not implement a communication method between the "Actors" or coroutines (yet).  

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
* You can combine these two steps together with a  `do` statement: `(do (func square (a) (* a a)) (square 3))`.  `do` statements take a list of expressions, evaluate them, and return the result of the last one.  The above expression returns Int: 9.
* You can also use the `map` function to run this square function over several inputs:
```lisp
(do 
   (func square (a) 
     (* a a)) 
   (map square (1 2 3 4 5 6 7 8 9))
)
```
Since the interpreter cannot handle carriage returns:
```lisp
(do (func square (a) (* a a)) (map square (1 2 3 4 5 6 7 8 9)))
```
* You can also create loops:
```lisp
(do
  (var a 0)
  (while (< a 10)
    (do 
      (var a (+ a 1))
      (core/print "A:")
      (core/print a)
    )
  )
)
```
This loop prints the value of a 10 times, as a is incremented from 0->10
```lisp
(do (var a 0) (while (< a 10) (do (var a (+ a 1)) (core/print "A:") (core/print a))))
```
    
* You can write conditional if statements:
```
(do
  (var a 5)
  (if (< a 10) 
    (core/print a)
  )
)
(do (var a 5) (if (< a 10) (core/print a) ))
(do
  (var a 5)
  (var b 10)
  (if (< a 1) 
    (core/print a)
    (core/print b)
  )
)
(do (var a 5) (var b 10) (if (< a 1) (core/print a) (core/print b) ))
```
* Lastly, you can make functions run concurrently to the REPL. This is hard to visualize, but it is possible using the `act` keyword.  Below, functions are defined that waste CPU by adding adding i, untill a large number (blocking), while another calls this function then prints an input aferwards. This function is run on a seperate actor, which returns the input (10 in this case). It takes ~ 4 seconds on my 2015 Macbook Pro. 
```lisp
(do 
 (func sleep (_) (do (var a 0) (while (< a 1000000) (do (var a (+ a 1)))) ))
 (func print-after-wait (val) (do (sleep 0) (core/print val)))
 (act print-after-wait (10))
 (core/print 1)
)

(do (func sleep (_) (do (var a 0) (while (< a 1000000) (do (var a (+ a 1)))) )) (func print-after-wait (val) (do (sleep 0) (core/print val))) (act print-after-wait (10)) (core/print 1) )
```


