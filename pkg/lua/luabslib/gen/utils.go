package main

import (
	"fmt"
	"go/ast"
)

func skipStar(node ast.Node) ast.Node {
	if v, ok := node.(*ast.StarExpr); ok {
		return v.X
	}
	return node
}

func makeConstructorName(typeName string) string {
	return fmt.Sprintf("New%s", typeName)
}
