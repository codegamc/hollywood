package stdlib

import (
	environment "github.com/codegamc/hollywood/environment"
)

// ImportStdLib binds the stdlib to the parent envi
func ImportStdLib(envi environment.Environment) environment.Environment {
	envir := &envi
	envir = importLogic(envir)
	envir = importMath(envir)
	envir = importCore(envir)
	return *envir
}
