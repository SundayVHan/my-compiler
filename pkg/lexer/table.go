package lexer

import (
	"strconv"
)

type Table struct {
	Symbols   SymbolTable
	Constants ConstantTable
}

type SymbolTable struct {
	Table []*SymbolInfo
}

type SymbolInfo struct {
	Name string
	Type string
	Val  interface{}
}

func NewTable() *Table {
	s := SymbolTable{
		Table: make([]*SymbolInfo, 0, 100),
	}
	c := ConstantTable{
		Table: make([]interface{}, 0, 100),
	}

	return &Table{
		Symbols:   s,
		Constants: c,
	}
}

// 向符号表中添加一个标识符，并返回该标识符所在的位置
func (t *Table) AddIdentifier(name string) int {
	t.Symbols.Table = append(t.Symbols.Table, &SymbolInfo{
		Name: name,
	})

	return len(t.Symbols.Table) - 1
}

type ConstantTable struct {
	Table []interface{}
}

func (t *Table) AddConstant(val string) (int, error) {
	intVal, err := strconv.Atoi(val)
	if err != nil {
		return 0, err
	}
	t.Constants.Table = append(t.Constants.Table, intVal)

	return len(t.Constants.Table) - 1, nil
}
