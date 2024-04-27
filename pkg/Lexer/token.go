package Lexer

const (
	TokenTypeEOF = iota
	TokenTypeIllegal
	TokenTypeIdentifier
	TokenTypeInteger
	TokenTypePlus
	TokenTypeMinus
	TokenTypeMultiply
	TokenTypeDivide
	TokenTypeAssign
	TokenTypeSemicolon
)

type Token struct {
	TokenType int
	Line      int
	Column    int
}
