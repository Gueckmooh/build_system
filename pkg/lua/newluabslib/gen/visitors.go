package main

import (
	"fmt"
	"go/ast"
)

type classSearcher struct {
	classes *[]*Class
}

func (c classSearcher) Visit(node ast.Node) ast.Visitor {
	switch node.(type) {
	case *ast.TypeSpec:
		class := &Class{}
		*c.classes = append(*c.classes, class)
		reader := &classReader{
			class: class,
		}
		ast.Walk(reader, node)
		return nil
	}
	return c
}

type classReader struct {
	class *Class
}

func (c classReader) Visit(node ast.Node) ast.Visitor {
	switch n := node.(type) {
	case *ast.TypeSpec:
		c.class.Name = n.Name.String()
		c.class.MethodMapName = fmt.Sprintf("__%sMethods", c.class.Name)
		c.class.MetatableName = fmt.Sprintf("__%sMetatableName", c.class.Name)
		c.class.Ast = n
		return c
	case *ast.Field:
		field := newFieldFromNode(n)
		c.class.Fields = append(c.class.Fields, field)
	}
	return c
}

type fieldReader struct {
	name *string
	ty   *Type
}

func (f fieldReader) Visit(node ast.Node) ast.Visitor {
	switch n := node.(type) {
	case *ast.Field:
		*f.name = n.Names[0].String()
		*f.ty = newTypeFromNode(n)
	}
	return nil
}
