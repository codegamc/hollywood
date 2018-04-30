#"Hollywood" By Michael Codega

  Hollywood is a simple lisp interpreter written in go. Its name comes from the "Actor Model" of concurrency, which hollywood is loosely trying to implement.  The Actor Model is based on a 1978 paper by Tony Hoare, where he described "Communicating Sequential Processes".  Hollywood does not implement a communication method between the "Actors" or coroutines.  

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
 
Since compiling a go binary requires installing a go compiler and setting up an environment, a linux compatable binary is included in this repository. It is called `hollywood-linux` and it works on the EECS Server. 
