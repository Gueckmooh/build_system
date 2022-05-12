package main

import "go/ast"

func skipStar(node ast.Node) ast.Node {
	if v, ok := node.(*ast.StarExpr); ok {
		return v.X
	}
	return node
}
