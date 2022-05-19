package main

import (
	"fmt"
	"go/ast"
	"strings"
)

type Type interface {
	GoString() string
	LuaString() string
	CheckFunction() Callable
	ToLuaType(v string) string
	ToGoType(v string) string
	IsContainer() bool
	InsideType() Type
	NeedsEllipsis() bool
}

type (
	TNil    struct{}
	TString struct{}
	TInt    struct{}
	TCustom struct {
		Name string
	}
	TPointer struct {
		X Type
	}
	TFunction struct {
		ReturnType Type
		Parameters []*Field
	}
	TArray struct {
		X Type
	}
	TEllipsis struct {
		X Type
	}
	TClass struct {
		X *Class
	}

	TState     struct{}
	TUserData  struct{}
	TGFunction struct{}
)

func (t *TNil) GoString() string       { return "<nil>" }
func (t *TString) GoString() string    { return "string" }
func (t *TCustom) GoString() string    { return t.Name }
func (t *TPointer) GoString() string   { return "*" + t.X.GoString() }
func (t *TClass) GoString() string     { return t.X.Name }
func (t *TState) GoString() string     { return "lua.LState" }
func (t *TUserData) GoString() string  { return "lua.LUserData" }
func (t *TGFunction) GoString() string { return "lua.LGFunction" }
func (t *TInt) GoString() string       { return "int" }
func (t *TArray) GoString() string     { return "[]" + t.X.GoString() }
func (t *TEllipsis) GoString() string  { return "[]" + t.X.GoString() }

func (t *TFunction) GoString() string {
	var paramStrings []string
	for _, p := range t.Parameters {
		paramStrings = append(paramStrings, p.String())
	}
	return fmt.Sprintf("func(%s) %s", strings.Join(paramStrings, ", "), t.ReturnType.GoString())
}

func (t *TNil) LuaString() string       { return "<nil>" }
func (t *TString) LuaString() string    { return "lua.LTString" }
func (t *TCustom) LuaString() string    { return "<nil>" }
func (t *TPointer) LuaString() string   { return "<nil>" }
func (t *TClass) LuaString() string     { return "<nil>" }
func (t *TState) LuaString() string     { return "lua.LTState" }
func (t *TUserData) LuaString() string  { return "lua.LTUserData" }
func (t *TGFunction) LuaString() string { return "<nil>" }
func (t *TInt) LuaString() string       { return "lua.LTNumber" }
func (t *TArray) LuaString() string     { return "<nil>" }
func (t *TEllipsis) LuaString() string  { return "<nil>" }
func (t *TFunction) LuaString() string  { return "<nil>" }

func (t *TNil) InsideType() Type       { return nil }
func (t *TString) InsideType() Type    { return nil }
func (t *TCustom) InsideType() Type    { return nil }
func (t *TPointer) InsideType() Type   { return t.X }
func (t *TClass) InsideType() Type     { return nil }
func (t *TState) InsideType() Type     { return nil }
func (t *TUserData) InsideType() Type  { return nil }
func (t *TGFunction) InsideType() Type { return nil }
func (t *TInt) InsideType() Type       { return nil }
func (t *TArray) InsideType() Type     { return t.X }
func (t *TEllipsis) InsideType() Type  { return t.X }
func (t *TFunction) InsideType() Type  { return nil }

func (t *TNil) CheckFunction() Callable       { return nil }
func (t *TCustom) CheckFunction() Callable    { return nil }
func (t *TPointer) CheckFunction() Callable   { return nil }
func (t *TState) CheckFunction() Callable     { return nil }
func (t *TUserData) CheckFunction() Callable  { return nil }
func (t *TGFunction) CheckFunction() Callable { return nil }
func (t *TFunction) CheckFunction() Callable  { return nil }
func (t *TEllipsis) CheckFunction() Callable  { return nil }

func (t *TString) CheckFunction() Callable {
	return &Method{
		This: &TState{},
		Function: Function{
			Name: "CheckString",
			Type: &TFunction{
				ReturnType: t,
				Parameters: []*Field{
					{
						Name: "n",
						Type: &TInt{},
					},
				},
			},
		},
	}
}

func (t *TArray) CheckFunction() Callable {
	return &Method{
		This: &TState{},
		Function: Function{
			Name: "CheckArray",
			Type: &TFunction{
				ReturnType: t,
				Parameters: []*Field{
					{
						Name: "n",
						Type: &TInt{},
					},
				},
			},
		},
	}
}

func (t *TClass) CheckFunction() Callable {
	return &Function{
		Name: t.X.FunctionBundle.LuaCheckType.Name,
		Type: &TFunction{
			ReturnType: t,
			Parameters: []*Field{
				{
					Name: "L",
					Type: &TState{},
				},
				{
					Name: "n",
					Type: &TInt{},
				},
			},
		},
	}
}

func (t *TInt) CheckFunction() Callable {
	return &Method{
		This: &TState{},
		Function: Function{
			Name: "CheckInt",
			Type: &TFunction{
				ReturnType: t,
				Parameters: []*Field{
					{
						Name: "n",
						Type: &TInt{},
					},
				},
			},
		},
	}
}

func (t *TNil) ToLuaType(v string) string       { return "<nil>" }
func (t *TString) ToLuaType(v string) string    { return fmt.Sprintf("lua.LString(%s)", v) }
func (t *TCustom) ToLuaType(v string) string    { return fmt.Sprintf("__Convert%s(%s)", t.Name, v) }
func (t *TPointer) ToLuaType(v string) string   { return t.X.ToLuaType(v) }
func (t *TClass) ToLuaType(v string) string     { return "" }
func (t *TState) ToLuaType(v string) string     { return "" }
func (t *TUserData) ToLuaType(v string) string  { return "" }
func (t *TGFunction) ToLuaType(v string) string { return "" }
func (t *TInt) ToLuaType(v string) string       { return "lua.LNumber(" + v + ")" }
func (t *TFunction) ToLuaType(v string) string  { return "" }
func (t *TArray) ToLuaType(v string) string     { return "lua.LArray(" + v + ")" }
func (t *TEllipsis) ToLuaType(v string) string  { return "lua.LEllipsis(" + v + ")" }

func (t *TNil) ToGoType(v string) string       { return "<nil>" }
func (t *TString) ToGoType(v string) string    { return fmt.Sprintf("%s.String()", v) }
func (t *TCustom) ToGoType(v string) string    { return "<nil>" }
func (t *TPointer) ToGoType(v string) string   { return "<nil>" }
func (t *TClass) ToGoType(v string) string     { return "<nil>" }
func (t *TState) ToGoType(v string) string     { return "<nil>" }
func (t *TUserData) ToGoType(v string) string  { return "<nil>" }
func (t *TGFunction) ToGoType(v string) string { return "<nil>" }
func (t *TInt) ToGoType(v string) string       { return "<nil>" }
func (t *TFunction) ToGoType(v string) string  { return "<nil>" }
func (t *TArray) ToGoType(v string) string     { return "<nil>" }
func (t *TEllipsis) ToGoType(v string) string  { return "<nil>" }

func (t *TNil) IsContainer() bool       { return false }
func (t *TString) IsContainer() bool    { return false }
func (t *TCustom) IsContainer() bool    { return false }
func (t *TPointer) IsContainer() bool   { return false }
func (t *TClass) IsContainer() bool     { return false }
func (t *TState) IsContainer() bool     { return false }
func (t *TUserData) IsContainer() bool  { return false }
func (t *TGFunction) IsContainer() bool { return false }
func (t *TInt) IsContainer() bool       { return false }
func (t *TFunction) IsContainer() bool  { return false }
func (t *TArray) IsContainer() bool     { return true }
func (t *TEllipsis) IsContainer() bool  { return true }

func (t *TNil) NeedsEllipsis() bool       { return false }
func (t *TString) NeedsEllipsis() bool    { return false }
func (t *TCustom) NeedsEllipsis() bool    { return false }
func (t *TPointer) NeedsEllipsis() bool   { return false }
func (t *TClass) NeedsEllipsis() bool     { return false }
func (t *TState) NeedsEllipsis() bool     { return false }
func (t *TUserData) NeedsEllipsis() bool  { return false }
func (t *TGFunction) NeedsEllipsis() bool { return false }
func (t *TInt) NeedsEllipsis() bool       { return false }
func (t *TFunction) NeedsEllipsis() bool  { return false }
func (t *TArray) NeedsEllipsis() bool     { return false }
func (t *TEllipsis) NeedsEllipsis() bool  { return true }

func newTypeFromNode(node ast.Node) Type {
	tr := &typeReader{
		ty: new(Type),
	}
	ast.Walk(tr, node)
	return *tr.ty
}

func newTypeFromName(name string) Type {
	switch name {
	case "string":
		return &TString{}
	default:
		return &TCustom{
			Name: name,
		}
	}
}

func newTypeFromExpr(expr ast.Expr) Type {
	switch e := expr.(type) {
	case *ast.StarExpr:
		return &TPointer{
			X: newTypeFromExpr(e.X),
		}
	case *ast.Ident:
		return newTypeFromName(e.Name)
	case *ast.ArrayType:
		return &TArray{
			X: newTypeFromExpr(e.Elt),
		}
	case *ast.Ellipsis:
		return &TEllipsis{
			X: newTypeFromExpr(e.Elt),
		}
	default:
		panic(fmt.Errorf("Unhandled type %#v", e))
	}
}

type typeReader struct {
	ty *Type
}

func (t typeReader) Visit(node ast.Node) ast.Visitor {
	switch n := node.(type) {
	case *ast.Field:
		*t.ty = newTypeFromExpr(n.Type)
		return nil
	case *ast.FuncDecl:
		var ret Type = nil
		if n.Type.Results != nil {
			ret = newTypeFromNode(n.Type.Results.List[0])
		} else {
			ret = &TNil{}
		}
		var fields []*Field
		for _, field := range n.Type.Params.List {
			ty := newTypeFromExpr(field.Type)
			for _, name := range field.Names {
				fields = append(fields, &Field{
					Name: name.Name,
					Type: ty,
				})
			}
		}
		*t.ty = &TFunction{
			ReturnType: ret,
			Parameters: fields,
		}
		return nil
	default:
		fmt.Printf("--> %#v\n", n)
	}
	return t
}
