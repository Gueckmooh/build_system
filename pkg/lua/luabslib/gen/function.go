package main

import (
	"fmt"
	"go/ast"
	"reflect"
	"strings"
)

type Callable interface {
	CallExpr(params ...string) string
}

type Function struct {
	Name           string
	MappingName    string
	LuaMappingName string
	MayFail        bool
	Type           *TFunction
}

type Method struct {
	Function
	This        Type
	MappingName string
}

func (m *Method) GenerateLuaBinding() string {
	return MustExecuteTemplate("lua_method_binding.gotmpl", m)
}

func (f *Function) Returns() bool {
	_, ok := f.Type.ReturnType.(*TNil)
	return !ok
}

func (f *Function) Signature() string {
	var paramStrings []string
	for _, p := range f.Type.Parameters {
		paramStrings = append(paramStrings, p.String())
	}
	if _, ok := f.Type.ReturnType.(*TNil); ok {
		return fmt.Sprintf("func %s(%s)", f.Name, strings.Join(paramStrings, ", "))
	}
	return fmt.Sprintf("func %s(%s) %s", f.Name, strings.Join(paramStrings, ", "), f.Type.ReturnType.GoString())
}

func (f *Method) Signature() string {
	var paramStrings []string
	for _, p := range f.Type.Parameters {
		paramStrings = append(paramStrings, p.String())
	}
	if _, ok := f.Type.ReturnType.(*TNil); ok {
		return fmt.Sprintf("func (%s) %s(%s)", f.This.GoString(), f.Name, strings.Join(paramStrings, ", "))
	}
	return fmt.Sprintf("func (%s) %s(%s) %s", f.This.GoString(), f.Name, strings.Join(paramStrings, ", "),
		f.Type.ReturnType.GoString())
}

func (f *Method) LuaSignature() string {
	if _, ok := f.This.(*TClass); ok {
		var paramStrings []string
		for _, p := range f.Type.Parameters {
			paramStrings = append(paramStrings, p.String())
		}
		return fmt.Sprintf("func %s(L *lua.LState) int", f.MappingName)
	} else {
		panic("Could not produce lua signature for method " + f.Name)
	}
}

func (f *Function) String() string {
	return f.Signature()
}

func (f *Method) String() string {
	return f.Signature()
}

func (f *Function) CallExpr(params ...string) string {
	if len(params) != len(f.Type.Parameters) {
		panic(fmt.Errorf("Wrong number of argument, required %d, got %d", len(f.Type.Parameters), len(params)))
	}
	return fmt.Sprintf("%s(%s)", f.Name, strings.Join(params, ", "))
}

func (f *Method) CallExpr(params ...string) string {
	if len(params) != (len(f.Type.Parameters) + 1) {
		panic(fmt.Errorf("Wrong number of argument, required %d, got %d", len(f.Type.Parameters)+1, len(params)))
	}
	return fmt.Sprintf("%s.%s(%s)", params[0], f.Name, strings.Join(params[1:], ", "))
}

func getMethodsForClass(c *Class, node ast.Node) {
	mr := &methodReader{
		className:   c.Name,
		methods:     &[]*Function{},
		constructor: new(*Function),
	}
	ast.Walk(mr, node)
	var methods []*Method
	for _, f := range *mr.methods {
		methods = append(methods, &Method{
			Function:    *f,
			This:        &TClass{c},
			MappingName: fmt.Sprintf("__lua%s%s", c.Name, f.Name),
		})
	}
	c.Ctor = *mr.constructor
	c.Methods = methods
}

type methodReader struct {
	className   string
	methods     *[]*Function
	constructor **Function
}

func (m methodReader) Visit(node ast.Node) ast.Visitor {
	switch n := node.(type) {
	case *ast.FuncDecl:
		if n.Recv == nil {
			if n.Name.Name != makeConstructorName(m.className) {
				break
			}
			if len(n.Type.Results.List) == 0 || len(n.Type.Results.List[0].Names) > 1 {
				break
			}

			if v, ok := skipStar(n.Type.Results.List[0].Type).(*ast.Ident); ok {
				if v.Name == m.className {
					ty, e := newTypeFromNode(n)
					if ty, ok := ty.(*TFunction); ok {
						*m.constructor = &Function{
							Name:    n.Name.Name,
							Type:    ty,
							MayFail: e,
						}
					}
				}
			}
			break
		}
		if len(n.Recv.List) > 0 {
			if v, ok := skipStar(n.Recv.List[0].Type).(*ast.Ident); ok {
				if v.Name == m.className {
					ty, e := newTypeFromNode(n)
					if ty, ok := ty.(*TFunction); ok {
						method := &Function{
							Name:    n.Name.Name,
							Type:    ty,
							MayFail: e,
						}
						*m.methods = append(*m.methods, method)
					} else {
						panic("Function type should be found..")
					}
				}
			}
		}
	}
	return m
}

type FunctionGen struct {
	Function
	templateName string
	Class        *Class
}

func (f *FunctionGen) Generate() string {
	return MustExecuteTemplate(f.templateName, f)
}

type FunctionNameBundle struct {
	LuaCheckType    *FunctionGen
	LuaTypeCtor     *FunctionGen
	LuaCtor         *FunctionGen
	LuaTypeCvtor    *FunctionGen
	LuaRegisterType *FunctionGen
	LuaNewLoader    *FunctionGen
}

func (f FunctionNameBundle) List() []*FunctionGen {
	var functions []*FunctionGen
	valueR := reflect.ValueOf(f)
	for i := 0; i < valueR.NumField(); i++ {
		fieldR := valueR.Field(i)
		functions = append(functions, fieldR.Interface().(*FunctionGen))
	}
	return functions
}

func SetFunctionNameBundle(c *Class) {
	c.FunctionBundle = FunctionNameBundle{
		LuaCheckType:    newLuaCheckTypeFunc(c),
		LuaTypeCtor:     newLuaTypeCtorFunc(c),
		LuaCtor:         newLuaCtorFunc(c),
		LuaTypeCvtor:    newLuaTypeConvertorFunc(c),
		LuaRegisterType: newLuaRegisterTypeFunc(c),
		LuaNewLoader:    newLuaNewLoaderFunc(c),
	}
}

func newLuaCheckTypeFunc(c *Class) *FunctionGen {
	f := &FunctionGen{
		Function:     Function{Name: fmt.Sprintf("__Check%s", c.Name), Type: &TFunction{}},
		templateName: "check_type_function.gotmpl",
		Class:        c,
	}
	f.Type.ReturnType = &TPointer{&TClass{c}}
	f.Type.Parameters = append(f.Type.Parameters, &Field{
		Name: "L",
		Type: &TPointer{&TState{}},
	})
	f.Type.Parameters = append(f.Type.Parameters, &Field{
		Name: "n",
		Type: &TInt{},
	})
	return f
}

func newLuaTypeCtorFunc(c *Class) *FunctionGen {
	f := &FunctionGen{
		Function:     Function{Name: fmt.Sprintf("__New%s", c.Name), Type: &TFunction{}},
		templateName: "type_constructor_function.gotmpl",
		Class:        c,
	}
	if c.Ctor != nil {
		c.Ctor.MappingName = fmt.Sprintf("__New%s", c.Name)
	}
	f.Type.ReturnType = &TPointer{&TUserData{}}
	f.Type.Parameters = append(f.Type.Parameters, &Field{
		Name: "L",
		Type: &TPointer{&TState{}},
	})
	return f
}

func newLuaCtorFunc(c *Class) *FunctionGen {
	f := &FunctionGen{
		Function:     Function{Name: fmt.Sprintf("__LuaNew%s", c.Name), Type: &TFunction{}},
		templateName: "lua_type_constructor_function.gotmpl",
		Class:        c,
	}
	if c.Ctor != nil {
		c.Ctor.LuaMappingName = fmt.Sprintf("__LuaNew%s", c.Name)
	}
	f.Type.ReturnType = &TInt{}
	f.Type.Parameters = append(f.Type.Parameters, &Field{
		Name: "L",
		Type: &TPointer{&TState{}},
	})
	return f
}

func newLuaTypeConvertorFunc(c *Class) *FunctionGen {
	f := &FunctionGen{
		Function:     Function{Name: fmt.Sprintf("__Convert%s", c.Name), Type: &TFunction{}},
		templateName: "type_convertor_function.gotmpl",
		Class:        c,
	}
	f.Type.ReturnType = &TPointer{&TUserData{}}
	f.Type.Parameters = append(f.Type.Parameters, &Field{
		Name: "L",
		Type: &TPointer{&TState{}},
	})
	f.Type.Parameters = append(f.Type.Parameters, &Field{
		Name: "val",
		Type: &TPointer{&TClass{c}},
	})
	return f
}

func newLuaRegisterTypeFunc(c *Class) *FunctionGen {
	f := &FunctionGen{
		Function:     Function{Name: fmt.Sprintf("__Register%sType", c.Name), Type: &TFunction{}},
		templateName: "register_type.gotmpl",
		Class:        c,
	}
	f.Type.ReturnType = &TNil{}
	f.Type.Parameters = append(f.Type.Parameters, &Field{
		Name: "L",
		Type: &TPointer{&TState{}},
	})
	return f
}

func newLuaNewLoaderFunc(c *Class) *FunctionGen {
	f := &FunctionGen{
		Function:     Function{Name: fmt.Sprintf("__New%sLoader", c.Name), Type: &TFunction{}},
		templateName: "new_loader.gotmpl",
		Class:        c,
	}
	f.Type.ReturnType = &TGFunction{}
	f.Type.Parameters = append(f.Type.Parameters, &Field{
		Name: "ret",
		Type: &TPointer{&TPointer{&TClass{c}}},
	})
	return f
}
