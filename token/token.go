package token

import "unicode"

// Kind represents the token kind.
type Kind string

// Slice contains allowed literals of a specified token kind.
// To see if your literal is allowed, you can use:
//
//	 userLit := "bye"
//	 lits := token.Slice{"hi", "bye"}
//	 if !slices.Contains(lits, userLit) {
//			panic("no user literal found in lits")
//	 }
type Slice []string

// Positions represents the token position.
type Position struct {
	Line     int `json:"line"`
	Column   int `json:"column"`
	Position int `json:"position"`
}

// Token represents the literal in the code,
// which is used by parser to build AST tree.
type Token struct {
	Literal  string    `json:"literal"`
	Position *Position `json:"position"`
	Kind     Kind      `json:"kind"`
}

const (
	Identifier Kind = "identifier"
	Number     Kind = "number"
	Float      Kind = "float"
	String     Kind = "string"
	Char       Kind = "char"
	Separator  Kind = "separator"
)

var (
	Separators = Slice{
		";",
		",",
		"[",
		"]",
		"{",
		"}",
	}
)

// IsIdentifier checks whether str is a valid identifier:
//   - Must not start with digit;
//   - Must not contain any special symbols except _;
//   - Must contain whitespaces;
//   - The length must be not longer than 255 or empty.
func IsIdentifier(str string) bool {
	if len(str) == 0 || len(str) > 255 {
		return false
	}
	for idx, r := range str {
		if idx == 0 && unicode.IsDigit(r) {
			return false
		}
		if !unicode.IsLetter(r) && !unicode.IsDigit(r) && r != '_' {
			return false
		}
	}
	return true
}
