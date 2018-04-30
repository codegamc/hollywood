package printer

import (
	types "github.com/codegamc/hollywood/types"
)

// PrintStr takes a HWType and converts it to a String
func PrintStr(ast types.HWType) string {
	return ast.ToString()
}
