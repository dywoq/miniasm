package token

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
