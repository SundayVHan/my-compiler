package lexer

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
	"strconv"
	"sync"
)

type Lexer struct {
	file       *bufio.Reader
	buf1       []byte
	buf1Filled bool
	mu1        sync.Mutex
	buf2       []byte
	buf2Filled bool
	mu2        sync.Mutex

	activeBuf *[]byte
	bufSize   int

	position int

	ch     byte
	line   int
	column int
	table  *Table
}

func NewLexer(file *os.File) *Lexer {
	reader := bufio.NewReader(file)
	l := &Lexer{
		file:     reader,
		bufSize:  1024,
		position: 0,
		line:     1,
		column:   0,
		table:    NewTable(),
	}

	l.fillBuffer()
	l.activeBuf = &l.buf1
	l.readChar()
	return l
}

func (lex *Lexer) fillBuffer() {
	if !lex.buf1Filled {
		lex.buf1 = make([]byte, lex.bufSize)
		n, err := lex.file.Read(lex.buf1)
		if err != nil && err != io.EOF {
			lex.reportError(err, lex.line, lex.column)
		}

		lex.mu1.Lock()
		lex.buf1 = lex.buf1[:n]
		lex.buf1Filled = true
		lex.mu1.Unlock()
	}

	if !lex.buf2Filled {
		lex.buf2 = make([]byte, lex.bufSize)
		n, err := lex.file.Read(lex.buf2)
		if err != nil && err != io.EOF {
			lex.reportError(err, lex.line, lex.column)
		}

		lex.mu2.Lock()
		lex.buf2 = lex.buf2[:n]
		lex.buf2Filled = true
		lex.mu2.Unlock()
	}
}

// 读取当前position所指向的字符，并让position指向下一个应读的字符
func (lex *Lexer) readChar() {
	if lex.position >= len(*lex.activeBuf) {
		lex.mu1.Lock()
		lex.mu2.Lock()
		if lex.activeBuf == &lex.buf1 && lex.buf2Filled {
			lex.activeBuf = &lex.buf2
			lex.buf1Filled = false
			go lex.fillBuffer()
		} else if lex.activeBuf == &lex.buf2 && lex.buf1Filled {
			lex.activeBuf = &lex.buf1
			lex.buf2Filled = false
			go lex.fillBuffer()
		}
		lex.mu1.Unlock()
		lex.mu2.Unlock()
		lex.position = 0
	}

	if len(*lex.activeBuf) > 0 {
		buf := *lex.activeBuf
		lex.ch = buf[lex.position]
		lex.position++
	} else {
		lex.ch = 0
	}
	lex.column++
}

func (lex *Lexer) NextToken() Token {
	var tok Token
	defer func() {
		fmt.Printf("%+v\n", tok)
		if tok.Type == IDENTIFIER {
			fmt.Printf("%v\n", *lex.table.Symbols.Table[tok.Value.(int)])
		}
		fmt.Printf("\n")
	}()
	lex.skipWhitespace()
	switch lex.ch {
	case 0:
		tok = newToken(EOF, -1, lex.line, lex.column)
	case '+':
		tok = newToken(PLUS, -1, lex.line, lex.column)
	case '-':
		tok = newToken(MINUS, -1, lex.line, lex.column)
	case '*':
		tok = newToken(MULTIPLY, -1, lex.line, lex.column)
	case '/':
		if lex.peekChar() == '/' {
			lex.skipSingleLineComment()
			tok = lex.NextToken() // 递归调用来获取下一个有效Token
			return tok
		} else if lex.peekChar() == '*' {
			lex.skipMultiLineComment()
			return lex.NextToken() // 递归调用来获取下一个有效Token
		} else {
			tok = newToken(DIVIDE, -1, lex.line, lex.column)
		}
	case '=':
		if lex.peekChar() == '=' {
			lex.readChar()
			tok = newToken(EQUALS, -1, lex.line, lex.column)
		} else {
			tok = newToken(ASSIGN, -1, lex.line, lex.column)
		}
	case '!':
		if lex.peekChar() == '=' {
			lex.readChar()
			tok = newToken(NOT_EQUALS, -1, lex.line, lex.column)
		} else {
			tok = newToken(ILLEGAL, -1, lex.line, lex.column)
		}
	case '<':
		if lex.peekChar() == '=' {
			lex.readChar()
			tok = newToken(LESS_OR_EQUALS, -1, lex.line, lex.column)
		} else {
			tok = newToken(LESS_THAN, -1, lex.line, lex.column)
		}
	case '>':
		if lex.peekChar() == '=' {
			lex.readChar()
			tok = newToken(GREATER_OR_EQUALS, -1, lex.line, lex.column)
		} else {
			tok = newToken(GREATER_THAN, -1, lex.line, lex.column)
		}
	case '&':
		if lex.peekChar() == '&' {
			lex.readChar() // 读取第二个 '&'
			tok = newToken(AND, -1, lex.line, lex.column)
		} else {
			tok = newToken(ADDR, -1, lex.line, lex.column)
		}
	case '|':
		if lex.peekChar() == '|' {
			lex.readChar() // 读取第二个 '|'
			tok = newToken(OR, -1, lex.line, lex.column)
		} else {
			tok = newToken(PIPE, -1, lex.line, lex.column)
		}
	case ',':
		tok = newToken(COMMA, -1, lex.line, lex.column)
	case ';':
		tok = newToken(SEMICOLON, -1, lex.line, lex.column)
	case ':':
		if lex.peekChar() == '=' {
			lex.readChar()
			tok = newToken(CREATE, -1, lex.line, lex.column)
		} else {
			tok = newToken(COLON, -1, lex.line, lex.column)
		}
	case '\n':
		tok = newToken(NEWLINE, -1, lex.line, lex.column)
		lex.line++
		lex.column = 0
	case '(':
		tok = newToken(LPAREN, -1, lex.line, lex.column)
	case ')':
		tok = newToken(RPAREN, -1, lex.line, lex.column)
	case '[':
		tok = newToken(LBRACKET, -1, lex.line, lex.column)
	case ']':
		tok = newToken(RBRACKET, -1, lex.line, lex.column)
	case '{':
		tok = newToken(LBRACE, -1, lex.line, lex.column)
	case '}':
		tok = newToken(RBRACE, -1, lex.line, lex.column)
	case '\'':
		if lex.peekChar() == '\'' {
			lex.readChar()
			tok = newToken(ILLEGAL, -1, lex.line, lex.column)
		} else {
			tok = newToken(BYTE, lex.readCharLiteral(), lex.line, lex.column)

			lex.readChar()
			if lex.ch != '\'' {
				tok.Type = ILLEGAL
			}
		}
	case '.':
		tok = newToken(DOT, -1, lex.line, lex.column)
	case '"':
		tok = newToken(STRING, lex.readStringLiteral(), lex.line, lex.column)
	default:
		if isLetter(lex.ch) {
			literal := lex.readIdentifier()
			typ := lookupIdentifier(literal)
			index := -1
			if typ == IDENTIFIER {
				index = lex.table.AddIdentifier(literal)
			}
			tok = newToken(typ, index, lex.line, lex.column)
			return tok
		}

		if isDigit(lex.ch) {
			literal, typ := lex.readNumber()
			var val interface{}
			if typ == INT {
				val, _ = strconv.ParseInt(literal, 10, 64)
			}
			if typ == FLOAT {
				val, _ = strconv.ParseFloat(literal, 64)
			}

			tok = newToken(typ, val, lex.line, lex.column)
			return tok
		}

		err := fmt.Errorf("Unrecognized character '%c'", lex.ch)
		lex.reportError(err, lex.line, lex.column)
	}

	lex.readChar()
	return tok
}

// 继续向后读取，直到读到一个非空白（换行）字符为止
func (lex *Lexer) skipWhitespace() {
	for lex.ch == ' ' || lex.ch == '\t' || lex.ch == '\r' {
		lex.readChar()
	}
}

func (lex *Lexer) skipSingleLineComment() {
	for lex.ch != '\n' && lex.ch != 0 {
		lex.readChar()
	}
}

func (lex *Lexer) skipMultiLineComment() {
	for {
		lex.readChar()
		if lex.ch == '*' && lex.peekChar() == '/' {
			lex.readChar()
			lex.readChar()
			break
		}
		if lex.ch == 0 {
			break
		}
	}
}

// 读取一个标识符，以字符开头，向后一直读到非字符且非数字
// 如果在结束后，发现终止的字符并非空格（换行）则报错并退出
func (lex *Lexer) readIdentifier() string {
	ident := make([]byte, 0, 10)
	for isLetter(lex.ch) || isDigit(lex.ch) {
		ident = append(ident, lex.ch)
		lex.readChar()
	}

	return string(ident)
}

// 读取一个标识符，以数字开头，向后一直读到非数字
// 如果在结束后，发现终止的字符并非空格（换行）则报错并退出
func (lex *Lexer) readNumber() (string, TokenType) {
	number := make([]byte, 0, 10)
	isFloat := false                       // 标记是否是浮点数
	for isDigit(lex.ch) || lex.ch == '.' { // 添加对小数点的检查
		if lex.ch == '.' {
			isFloat = true
		}
		number = append(number, lex.ch)
		lex.readChar()
	}
	if !(lex.ch == ' ' || lex.ch == '\t' || lex.ch == '\r' || lex.ch == '\n') {
		lex.reportError(errors.New("illegal number"), lex.line, lex.column)
	}
	if isFloat {
		return string(number), FLOAT
	} else {
		return string(number), INT
	}
}

// 打印错误并退出程序
func (lex *Lexer) reportError(err error, line, column int) {
	err = fmt.Errorf("Line %d, Column %d: %s", line, column, err.Error())
	fmt.Println(err)
	os.Exit(1)
}

func (lex *Lexer) peekChar() byte {
	if lex.position >= len(*lex.activeBuf) {
		if lex.activeBuf == &lex.buf1 && lex.buf2Filled {
			if len(lex.buf2) > 0 {
				return lex.buf2[0]
			}
		} else if lex.activeBuf == &lex.buf2 && lex.buf1Filled {
			if len(lex.buf1) > 0 {
				return lex.buf1[0]
			}
		}
		return 0
	}

	buf := *lex.activeBuf
	return buf[lex.position]
}

// 读取一个字符的字面量，会排除掉转义符号
func (lex *Lexer) readCharLiteral() string {
	// 处理转义字符
	if lex.peekChar() == '\\' {
		lex.readChar()
	}
	lex.readChar() // 更新l.ch
	return string(lex.ch)
}

// 读取一个字符串的字面量
func (lex *Lexer) readStringLiteral() string {
	str := make([]byte, 0, 10)
	for {
		lex.readChar()
		if lex.ch == '"' || lex.ch == 0 {
			break
		}

		str = append(str, lex.ch)
	}

	return string(str)
}
