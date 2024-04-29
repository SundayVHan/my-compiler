package lexer

type TokenType int

const (
	EOF     TokenType = iota // 文件结束
	ILLEGAL                  // 非法字符

	// 标识符与字面量
	IDENTIFIER // 标识符
	INT        // 整数
	FLOAT      // 浮点数
	STRING     // 字符串
	BOOL       // 布尔值
	BYTE       // 字节

	TYPE_INT    // int
	TYPE_FLOAT  //float // n
	TYPE_STRING // string
	TYPE_BOOL   // bool
	TYPE_BYTE   // byte

	// 运算符
	PLUS              // +
	MINUS             // -
	MULTIPLY          // *
	DIVIDE            // /
	ASSIGN            // =
	CREATE            // :=
	LESS_THAN         // <
	LESS_OR_EQUALS    // <=
	GREATER_THAN      // >
	GREATER_OR_EQUALS // >=
	EQUALS            // ==
	NOT_EQUALS        // !=
	ADDR              // &
	AND               // &&
	PIPE              // |
	OR                // ||

	// 分隔符
	COMMA     // ,
	SEMICOLON // ;
	COLON     // :
	NEWLINE   // \n
	DOT       // .

	// 括号
	LPAREN   // (
	RPAREN   // )
	LBRACE   // {
	RBRACE   // }
	LBRACKET // [
	RBRACKET // ]

	// 关键字
	FUNCTION // func
	VAR      // var
	TRUE     // true
	FALSE    // false
	IF       // if
	ELSE     // else
	RETURN   // return
	SWITCH   // switch
	CASE     // case
	TYPE     // type
	STRUCT   // struct
	PACKAGE  // package
	LEN
	IMPORT
	MAKE
)

func (t TokenType) String() string {
	switch t {
	case EOF:
		return "EOF"
	case ILLEGAL:
		return "ILLEGAL"

	case IDENTIFIER:
		return "IDENTIFIER"
	case INT:
		return "INT"
	case FLOAT:
		return "FLOAT"
	case STRING:
		return "STRING"
	case BOOL:
		return "BOOL"
	case BYTE:
		return "BYTE"

	case TYPE_INT:
		return "TYPE_INT"
	case TYPE_FLOAT:
		return "TYPE_FLOAT"
	case TYPE_STRING:
		return "TYPE_STRING"
	case TYPE_BOOL:
		return "TYPE_BOOL"
	case TYPE_BYTE:
		return "TYPE_BYTE"

	case PLUS:
		return "PLUS"
	case MINUS:
		return "MINUS"
	case MULTIPLY:
		return "MULTIPLY"
	case DIVIDE:
		return "DIVIDE"
	case ASSIGN:
		return "ASSIGN"
	case CREATE:
		return "create"

	case COMMA:
		return "COMMA"
	case SEMICOLON:
		return "SEMICOLON"
	case COLON:
		return "COLON"
	case NEWLINE:
		return "NEWLINE"
	case DOT:
		return "DOT"

	case LPAREN:
		return "LPAREN"
	case RPAREN:
		return "RPAREN"
	case LBRACE:
		return "LBRACE"
	case RBRACE:
		return "RBRACE"
	case LBRACKET:
		return "LBRACKET"
	case RBRACKET:
		return "RBRACKET"

	case FUNCTION:
		return "FUNCTION"
	case VAR:
		return "VAR"
	case TRUE:
		return "TRUE"
	case FALSE:
		return "FALSE"
	case IF:
		return "IF"
	case ELSE:
		return "ELSE"
	case RETURN:
		return "RETURN"
	case SWITCH:
		return "SWITCH"
	case CASE:
		return "CASE"
	case TYPE:
		return "TYPE"
	case STRUCT:
		return "STRUCT"
	case PACKAGE:
		return "PACKAGE"
	case LEN:
		return "LEN"
	case IMPORT:
		return "IMPORT"
	case MAKE:
		return "MAKE"

	case ADDR:
		return "ADDR"
	case AND:
		return "AND"
	case PIPE:
		return "PIPE"
	case OR:
		return "OR"
	case EQUALS:
		return "EQUALS"
	case NOT_EQUALS:
		return "NOT_EQUALS"
	case LESS_THAN:
		return "LESS_THAN"
	case GREATER_THAN:
		return "GREATER_THAN"
	case LESS_OR_EQUALS:
		return "LESS_OR_EQUALS"
	case GREATER_OR_EQUALS:
		return "GREATER_OR_EQUALS"
	default:
		return "UNKNOWN"
	}
}

type Token struct {
	Type   TokenType
	Value  interface{}
	Line   int
	Column int
}

func newToken(tokenType TokenType, value interface{}, line int, column int) Token {
	return Token{
		Type:   tokenType,
		Value:  value,
		Line:   line,
		Column: column,
	}
}

func lookupIdentifier(literal string) TokenType {
	switch literal {
	case "switch":
		return SWITCH
	case "case":
		return CASE

	case "int":
		return TYPE_BYTE
	case "float":
		return TYPE_FLOAT
	case "string":
		return TYPE_STRING
	case "bool":
		return TYPE_BOOL
	case "byte":
		return TYPE_BYTE

	case "func":
		return FUNCTION
	case "var":
		return VAR
	case "true":
		return TRUE
	case "false":
		return FALSE
	case "if":
		return IF
	case "else":
		return ELSE
	case "return":
		return RETURN
	case "type":
		return TYPE
	case "struct":
		return STRUCT
	case "package":
		return PACKAGE
	case ".":
		return DOT
	case "len":
		return LEN
	case "make":
		return MAKE
	case "import":
		return IMPORT

	default:
		return IDENTIFIER
	}
}
