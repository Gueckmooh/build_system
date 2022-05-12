package main

import (
	"bytes"
	"fmt"
	"go/ast"
)

type Class struct {
	Name           string
	Fields         []*Field
	Methods        []*Method
	Ast            *ast.TypeSpec
	FunctionBundle FunctionNameBundle
	MethodMapName  string
	MetatableName  string
}

func (c *Class) Constructor() string {
	return fmt.Sprintf("&%s{}", c.Name)
}

func newFieldFromNode(node *ast.Field) *Field {
	fr := &fieldReader{
		name: new(string),
		ty:   new(Type),
	}
	ast.Walk(fr, node)
	return &Field{
		Name: *fr.name,
		Type: *fr.ty,
	}
}

func getClasses(node ast.Node) []*Class {
	cv := classSearcher{
		classes: &[]*Class{},
	}
	ast.Walk(cv, node)
	return *cv.classes
}

func (c *Class) String() string {
	var buff bytes.Buffer
	fmt.Fprintf(&buff, "class %s {\n", c.Name)
	fmt.Fprintf(&buff, "fields:\n")
	for _, field := range c.Fields {
		fmt.Fprintf(&buff, "  %s\n", field)
	}
	fmt.Fprintf(&buff, "methods:\n")
	for _, method := range c.Methods {
		fmt.Fprintf(&buff, "  %s\n", method)
	}
	fmt.Fprintf(&buff, "generated methods:\n")
	for _, method := range c.FunctionBundle.List() {
		fmt.Fprintf(&buff, "  %s\n", method)
	}
	fmt.Fprintf(&buff, "}")
	return buff.String()
}
