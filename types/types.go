package hollywood

import "strconv"

/*

NOTE: not all of the types defined here are implementd

Some where made with ambition intentions, or with the intentions of creating a clear scope for the
remaining types, so show the bounds of what is used

All of the types defined in the Const block are the types that are used in the code

*/

// Function is a definition of "Functions"
type Function func([]HWType) HWType

// This const block contains the constants defining types, think Enums
const (
	// this indicates that something without implementing everything correctly
	// since 0 is the first iota, and 0 is the integer "nil type" any HWTYpe struct created
	// and not properly implemented will have a 0 as the type
	undefined_type = iota
	// ERROR_TYPE represents a EXCEPTION type
	ERROR_TYPE = iota
	// LIST_TYPE represents a list type
	LIST_TYPE = iota
	// ATOM_TYPE represents an atom type
	ATOM_TYPE = iota
	// SYM_TYPE represents
	SYM_TYPE = iota
	// INT_TYPE represents
	INT_TYPE = iota
	// CHAR_TYPE
	//CHAR_TYPE = iota
	// FUNC_TYPE represents
	FUNC_TYPE = iota
	// NULL_TYPE represents
	NULL_TYPE = iota
	// BOOL_TYPE represents a boolean type
	BOOL_TYPE = iota
	// STR_TYPE represents a string type
	STR_TYPE = iota
)

// HWType is a HW type
type HWType interface {
	ToString() string
	GetMeta() string
	GetType() int
}

//GetType uses meta to find type
func GetType(object HWType) string {
	// todo fix this to use type field
	return object.GetMeta()
}

// NotFalsey determines the "falseness" of an atom
// this should be implemented at the Eval stage, and evalute the input first
func NotFalsey(atom HWType) bool {
	if atom.GetType() == BOOL_TYPE {
		if atom.(HWBool).Val == false {
			return false
		}
		return true
	}

	if atom.GetType() == NULL_TYPE {
		return false
	}

	if atom.GetType() == INT_TYPE {
		if atom.(HWInt).Val == 0 {
			return false
		}

		return true
	}

	return true
}

//////////////////////////////////////////////////////////////////// LIST

// MakeList makes a list
func MakeList(val []HWType) HWList {
	list := HWList{Val: val, Meta: "list", Type: LIST_TYPE}
	return list
}

// HWList is a list of HWTypes
type HWList struct {
	Val  []HWType
	Meta string
	Type int
}

// ToString converts the list to a string
func (list HWList) ToString() string {
	result := " HWList: "
	for _, item := range list.Val {
		str := item.ToString()
		result = result + str + " "
	}
	return result
}

// GetMeta returns a list's meta data
func (list HWList) GetMeta() string {
	return list.Meta
}

// GetType returns a lists type
func (list HWList) GetType() int {
	return list.Type
}

//////////////////////////////////////////////////////////////////// ATOM

// MakeAtom creates an atoms
func MakeAtom(val HWType) HWAtom {
	return HWAtom{Val: val, Meta: "atom", Type: ATOM_TYPE}
}

// HWAtom is a HW primative
type HWAtom struct {
	// will be a single 64 bit int, unsigned that stores any singleton
	// worry about types later...
	Val  HWType
	Meta string
	Type int
}

// ToString converts the atom to a string
func (atom HWAtom) ToString() string {
	str := "HWAtom: " + atom.Val.ToString()
	return str
}

// GetMeta returns a Atoms's meta data
func (atom HWAtom) GetMeta() string {
	return atom.Meta
}

// GetType returns
func (atom HWAtom) GetType() int {
	return atom.Type
}

//////////////////////////////////////////////////////////////////// SYMBOL

// MakeSymbol makes a sumbol
func MakeSymbol(sym string) HWSymbol {
	return HWSymbol{Val: sym, Meta: "symbol", Type: SYM_TYPE}
}

// HWSymbol is a symbol / string
type HWSymbol struct {
	Val  string
	Meta string
	Type int
}

// ToString Converts HWInto to string
func (sym HWSymbol) ToString() string {
	return "HWSymbol " + sym.Val
}

// GetMeta returns a symbols's meta data
func (sym HWSymbol) GetMeta() string {
	return sym.Meta
}

// GetType returns
func (sym HWSymbol) GetType() int {
	return sym.Type
}

//////////////////////////////////////////////////////////////////// INT

// MakeInt makes a list
func MakeInt(val int64) HWInt {
	integer := HWInt{Val: val, Meta: "int", Type: INT_TYPE}
	return integer
}

// HWInt is an integer
type HWInt struct {
	Val  int64
	Meta string
	Type int
}

// ToString Converts HWInto to string
func (atom HWInt) ToString() string {
	// figure out what this means
	return "HWInt: " + strconv.FormatInt(atom.Val, 10)
}

// GetMeta returns a ints's meta data
func (atom HWInt) GetMeta() string {
	return atom.Meta
}

// GetType returns the type
func (atom HWInt) GetType() int {
	return atom.Type
}

//////////////////////////////////////////////////////////////////// CHAR

// HWChar is a 8 bit Value
type HWChar struct {
	Val  int8
	Meta string
	Type int
}

// ToString Converts HWChar to string
func (atom HWChar) ToString() string {
	return string(atom.Val)
}

// GetMeta returns a chars's meta data
func (atom HWChar) GetMeta() string {
	return atom.Meta
}

//////////////////////////////////////////////////////////////////// DOUBLE

// HWDouble is a double presision floating point value
type HWDouble struct {
	Val  float64
	Meta string
}

// ToString Converts HWInto to string
func (atom HWDouble) ToString() string {
	return strconv.FormatFloat(atom.Val, 'f', 6, 64)
}

// GetMeta returns a doubles's meta data
func (atom HWDouble) GetMeta() string {
	return atom.Meta
}

//////////////////////////////////////////////////////////////////// FUNCTION

// MakeFunc creates a hw function
func MakeFunc(f Function, sym string) HWFunc {
	fun := HWFunc{Val: f, Meta: "function", Symbol: sym, Type: FUNC_TYPE}
	return fun
}

// HWFunc is a function that accepts a HWType (List, or Atom)
type HWFunc struct {
	Val    Function
	Meta   string
	Symbol string
	Type   int
}

// ToString Converts HWInto to string
func (function HWFunc) ToString() string {
	return "Function: <" + function.Symbol + "> "
}

// GetMeta returns a func's meta data
func (function HWFunc) GetMeta() string {
	return function.Meta
}

// GetType returns the type
func (function HWFunc) GetType() int {
	return function.Type
}

//////////////////////////////////////////////////////////////////// NULL

// MakeNull returns a hollywood null type that prints itself as null
func MakeNull() HWNull {
	return HWNull{Meta: "<NULL TYPE>", Val: 0, Type: NULL_TYPE}
}

// MakeNullImplicit returns a hollywood null type that does not print anything out
func MakeNullImplicit() HWNull {
	return HWNull{Meta: " ", Val: 0, Type: NULL_TYPE}
}

// HWNull is hollywood null type
type HWNull struct {
	Val  int64
	Meta string
	Type int
}

// ToString Converts HWNull into to string
func (n HWNull) ToString() string {
	return n.Meta
}

// GetMeta returns a nulls's meta data
func (n HWNull) GetMeta() string {
	return n.Meta
}

// GetType returns the type
func (n HWNull) GetType() int {
	return n.Type
}

//////////////////////////////////////////////////////////////////// BOOLEAN

// MakeBool returns a hollywood boolean type
func MakeBool(truthiness bool) HWBool {
	return HWBool{Meta: "bool", Val: truthiness, Type: BOOL_TYPE}
}

// HWBool is hollywood boolean type
type HWBool struct {
	Val  bool
	Meta string
	Type int
}

// ToString Converts HWNull into to string
func (b HWBool) ToString() string {
	if b.Val {
		return "true"
	}

	return "false"
}

// GetMeta returns a nulls's meta data
func (b HWBool) GetMeta() string {
	return b.Meta
}

// GetType returns the type
func (b HWBool) GetType() int {
	return b.Type
}

//////////////////////////////////////////////////////////////////// STRING

// MakeString returns a hollywood string type
func MakeString(str string) HWString {
	return HWString{Val: str, Meta: "string - 0", Type: STR_TYPE}
}

// HWString is hollywood string type
type HWString struct {
	Val  string
	Meta string
	Type int
}

// ToString Converts HWInto to string
func (str HWString) ToString() string {
	return str.Val
}

// GetMeta returns a symbols's meta data
func (str HWString) GetMeta() string {
	return str.Meta
}

// GetType returns
func (str HWString) GetType() int {
	return str.Type
}
