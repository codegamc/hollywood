package hollywood

import (
	"regexp"
	"strconv"

	types "github.com/codegamc/hollywood/types"
)

// Maybe use const types?

// Reader is a stateful reader
type Reader struct {
	currentPosition int
	tokens          []string
}

// Next gets next token from reader, removes token from reader
func (reader *Reader) Next() string {
	if len(reader.tokens) == (reader.currentPosition + 1) {
		// cannot do
		// TODO fault
		return ""
	} else {
		reader.currentPosition++
	}

	return reader.tokens[reader.currentPosition]
}

// Peek peeks at the next token, wihtout removing it from reader
func (reader *Reader) Peek() string {
	if len(reader.tokens) == (reader.currentPosition + 1) {
		// cannot do
		// TODO fault
		return ""
	}
	token := reader.tokens[reader.currentPosition+1]
	return token
}

// Tokenizer takes a string and returns a list of tokens
func Tokenizer(code string) []string {
	tokens := make([]string, 0)
	// regex: [\s,]*(~@|[\[\]{}()'`~^@]|"(?:\\.|[^\\"])*"|;.*|[^\s\[\]{}('"`,;)]*)
	// This form is used to handle the regex without worrying about trying to escape
	// all the different characters in the original one since quotes are needed in the regex

	// The internet helped me figure out this regex. It breaks up strings based on the
	// LISP syntax of parens and whitespace while keeping strings inside quotes together
	r := regexp.MustCompile(`[\s,]*(~@|[\[\]{}()'` + "`" +
		`~^@]|"(?:\\.|[^\\"])*"|;.*|[^\s\[\]{}('"` + "`" +
		`,;)]*)`)

	for _, toks := range r.FindAllStringSubmatch(code, -1) {
		if (toks[1] == "") || (toks[1][0] == ';') {
			continue
		}
		tokens = append(tokens, toks[1])
	}

	return tokens
}

// ReadStr reads a string, and creates a reader, and returns an AST
func ReadStr(str string) types.HWType {
	tokens := Tokenizer(str)
	// test to ensure tokens isnt length 0
	if len(tokens) == 0 {
		return types.MakeNull()
	}
	for i := 0; i < len(tokens); i++ {
		//fmt.Println("Next token: ", tokens[i])
	}

	reader := Reader{-1, tokens}
	//fmt.Println("Readers tokens:", tokens)
	return ReadForm(&reader)
}

// ReadForm breaks down the tokens as either:
// 1: new list opening
// 2: atomic type
// it then returns this structure, as the Abstract syntax tree
func ReadForm(reader *Reader) types.HWType {
	// determine where the next token should handle
	token := reader.Peek()

	switch token[0] {
	case '(':
		return ReadList(reader)
	default:
		return ReadAtom(reader)
	}
}

// ReadList converts a list of items into a Type
// it converts the string into a list of HWTypes,
// recursively calling ReadForm on it
func ReadList(reader *Reader) types.HWType {
	token := reader.Next() //move past the open parens
	//fmt.Println(token)
	// The token is '('
	if token != "(" {
		// throw error
	}
	// This is the list of items that are inside that List
	// eg. (* 4 5 (+ 3 3) 7 8)
	// { *, 4, 5, LIST, 7, 8}
	ast := make([]types.HWType, 0)

	for {
		// assuming the next token isnt end list
		token = reader.Peek()
		if token != ")" {
			// get the next List Entry

			entity := ReadForm(reader)
			// Add it to the list
			ast = append(ast, entity)
		} else {
			break
		}
	}

	reader.Next() // to account for the ) peak
	entity := types.MakeList(ast)
	return entity
}

// ReadAtom reads an atom, and returns it as a data type
func ReadAtom(reader *Reader) types.HWType {
	token := reader.Next()
	// test if its a string, this is a gross thing to read but its checking equivilence to this character ->> "
	if token[0] == "\""[0] {
		// its a string
		return types.MakeString(token)
	}

	// testing if the value is an integer
	i, e := strconv.Atoi(token)
	if e == nil { // This means its actually an integer
		//fmt.Println(i)
		atom := types.MakeInt(int64(i))
		return atom
	}

	if token == "true" {
		return types.MakeBool(true)
	}

	if token == "false" {
		return types.MakeBool(false)
	}

	if token == "null" {
		return types.MakeNull()
	}

	// if not, assume it is a symbol for now
	atom := types.MakeSymbol(token)
	return atom
}
