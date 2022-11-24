package main

import (
	"fmt"
	"os"

	"github.com/pingcap/tidb/parser"
	"github.com/pingcap/tidb/parser/ast"
	_ "github.com/pingcap/tidb/types/parser_driver"
)

type colX struct {
	colNames []string
}

func (x *colX) Enter(in ast.Node) (ast.Node, bool) {
	if name, ok := in.(*ast.ColumnName); ok {
		x.colNames = append(x.colNames, name.Name.O)
	}
	return in, false
}

func (x *colX) Leave(in ast.Node) (ast.Node, bool) {
	return in, true
}

func parse(sql string) (ast.StmtNode, error) {
	parser := parser.New()

	stmtNodes, _, err := parser.Parse(sql, "", "")
	if err != nil {
		return nil, err
	}
	return stmtNodes[0], nil
}

func extract(stmtNode ast.StmtNode) []string {
	x := &colX{}
	stmtNode.Accept(x)
	return x.colNames
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: colx 'SQL statement'")
		return
	}

	sql := os.Args[1]
	node, err := parse(sql)
	if err != nil {
		fmt.Printf("parse error: %v\n", err.Error())
		return
	}

	fmt.Printf("%v\n", extract(node))
}
