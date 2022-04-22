package main

import (
	"fmt"
	"strings"
)

const methodSetterTemplate = `func {{.MethodName}}(L *lua.LState) int {
	self := L.ToTable(1)
	value := L.Get(2)
{{genTypeCheck .ValueType "value"}}
	L.SetField(self, "{{.FieldName}}", value)
	return 0
}
`

const methodAppendTemplate = `func {{.MethodName}}(L *lua.LState) int {
	self := L.ToTable(1)
	value := L.Get(2)
{{genTypeCheck .ValueType "value"}}
	vtable := L.GetField(self, "{{.FieldName}}")
{{genTypeCheck .TargetType "vtable"}}
	table := vtable.(*lua.LTable)
{{genAppendForTypes .ValueType "value" "table"}}
	return 0
}
`

const totoTemplate = `var {{ .FuncMapName }} = map[string]lua.LGFunction{
{{- range .Mappings }}
	"{{ .LuaName }}": {{ .GoName }},
{{- end }}
}`

const libraryBodyTemplate = `// Code generated by go generate; DO NOT EDIT.
package {{.PackageName}}

import (
	"fmt"
	lua "github.com/yuin/gopher-lua"
)

{{ .TypeDefinition }}

{{ .PublicTypeDefinition }}

{{ .PublicGetters }}

{{ .ConvertLuaTable }}

{{ .PublicConvertLuaTable }}

{{ .FunctionMapping }}

{{ .IntegrityChecker }}

{{ .NewTable }}

{{ .LuaNewTable }}

{{ .PublicNewTable }}

{{ .PublicLuaNewTable }}

// go bindings for lua methods:
{{- range .GoMethods }}
{{ . }}
{{- end }}
`

const typeCheckTemplate = `if {{ genTypeCheckCond .VarType .VarName }} {
	fmt.Printf("Incorrect type %s\n", {{.VarName}}.Type().String())
	L.Panic(L)
}`

const tableTypeCheckTemplate = `if {{.VarName}}.Type() == lua.LTTable {
	L.ForEach({{.VarName}}.(*lua.LTable), func (_, v lua.LValue) {
		{{ genTypeCheck .VarType "v" }}
	})
}`

const typeCheckErrorTemplate = `if {{ genTypeCheckCond .VarType .VarName }} {
	{{.ErrName}} = fmt.Errorf("Unknown type %s", {{.VarName}}.Type().String())
}`

const tableTypeCheckErrorTemplate = `if {{.VarName}}.Type() == lua.LTTable {
	var err error
	L.ForEach({{.VarName}}.(*lua.LTable), func (_, v lua.LValue) {
		var locError error
		{{ genTypeCheckError .VarType "v" "locError" }}
		if err == nil {
			err = locError
		}
	})
	if err != nil {
		return err
	}
}`

const checkTableIntegrityTemplate = `func {{.FuncName}}(L *lua.LState, T *lua.LTable) error {
	{{- range .Fields}}
	{
		value := L.GetField(T, "{{.Name}}")
		var err error
		{{genTypeCheckError .Type "value" "err"}}
		if err != nil {
			return err
		}
	}
	{{- end}}
	return nil
}`

const newTableTemplate = `func {{.FuncName}}(L *lua.LState) *lua.LTable {
	table := L.SetFuncs(L.NewTable(), {{.FunctionMapping}})

	{{ range .Fields}}
		{{- genFieldInit . "table"}}
	{{ end}}

	return table
}`

const luaNewTableTemplate = `func {{.FuncName}}(L *lua.LState) int {
	table := {{.NewFuncName}}(L)

	L.Push(table)

	return 1
}`

const getLuaStringFromTableFieldTemplate = `var {{.VarName}} string
{
	__luaFieldValue := L.GetField({{.TableName}}, "{{.FieldName}}")
	{{.VarName}} = __luaFieldValue.String()
}`

const getLuaStringFromValueTemplate = `{{.VarName}} := {{.LuaVarName}}.String()`

const getLuaTableFromTableFieldTemplate = `var {{.VarName}} []{{.GoType}}
{
	__luaFieldValue := L.GetField({{.TableName}}, "{{.FieldName}}")
	__luaFieldTable := __luaFieldValue.(*lua.LTable)
	L.ForEach(__luaFieldTable, func(_, v lua.LValue) {
		{{genGetGoValueFromValue "__subField" "v" .LuaType }}
		{{.VarName}} = append({{.VarName}}, __subField)
	})
}`

const tableTypeDefTemplate = `type {{.TypeName}} struct {
	{{- range .Fields }}
		{{genTypeDefField .}}
	{{- end }}
}`

const tableConversionTemplate = `func {{.FuncName}}(L *lua.LState, T *lua.LTable) (*{{.TypeName}}, error) {
	err := {{.CheckIntegrity}}(L, T)
	if err != nil {
		return nil, err
	}
	{{- range .Fields }}
		{{genGetGoValueFromLuaField .GoName .Name "T" .Type}}
	{{- end }}

	return &{{.TypeName}}{
		{{- range .Fields }}
			{{.GoName}}: {{.GoName}},
		{{- end }}
	}, nil
}`

const publicTableGetterTemplate = `func (t *{{.TableName}}) {{.FuncName}}() {{.TypeName}} {
	return t.{{.FieldName}}
}`

const publicConvertTableTemplate = `func {{.FuncName}}(L *lua.LState, T *lua.LTable) (*{{.TableName}}, error) {
	privateTable, err := {{.ConvertName}}(L, T)
	if err != nil {
		return nil, err
	}
	publicTable := {{.TableName}}(*privateTable)
	return &publicTable, nil
}`

const publicNewTableTemplate = `func {{.FuncName}}(L *lua.LState) *lua.LTable {
	return {{.NewTable}}(L)
}`

const publicLuaNewTableTemplate = `func {{.FuncName}}(L *lua.LState) int {
	return {{.NewTable}}(L)
}`

// @todo move
func luaTypeToGoType(ty string) string {
	switch ty {
	case "String":
		return "string"
	default:
		if typeIsTable(ty) {
			return fmt.Sprintf("[]%s", luaTypeToGoType(getInnerType(ty)))
		}
		return ""
	}
}

func snakeCaseToCamelCase(s string) string {
	subs := strings.Split(s, "_")
	var nsubs []string
	for _, s := range subs {
		nsubs = append(nsubs, strings.ToUpper(string(s[0]))+s[1:])
	}
	return strings.Join(nsubs, "")
}
