package main

import (
	"fmt"
	"myCompiler/pkg/lexer"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: lexer <filename>")
		return
	}
	filename := os.Args[1]

	file, err := os.Open(filename)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	l := lexer.NewLexer(file)
	for tok := l.NextToken(); tok.Type != lexer.EOF; tok = l.NextToken() {

	}
}
