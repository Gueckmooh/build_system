package main

import (
	"bytes"
	"fmt"
	"regexp"
	"strings"
	"text/template"
)

type DefaultValue interface{}

type DefaultFromParam struct {
	ParamId int
}

type NoDefault struct{}

var Dependencies []*TableGenerator

type FieldDescriptor struct {
	Name         string
	GoName       string
	Types        TypeDescriptors
	GetterName   string
	DefaultValue DefaultValue
}

func (fd *FieldDescriptor) String() string {
	return fmt.Sprintf("Field: %s(%s), %s", fd.Name, fd.GoName, fd.Types)
}

func NewDefaultValue(f *Field) DefaultValue {
	if f.DefaultFromParam != 0 {
		return &DefaultFromParam{
			ParamId: f.DefaultFromParam,
		}
	}
	return &NoDefault{}
}

func NewFieldDescriptor(f *Field) *FieldDescriptor {
	var name string
	if f.Private {
		name = fmt.Sprintf("_%s_", f.Name)
	} else {
		name = f.Name
	}
	return &FieldDescriptor{
		Name:   name,
		GoName: fmt.Sprintf("__v%s", f.Name),
		Types:  NewTypeDescriptors(f.Type),
		// Type:         f.Type,
		// GoType:       luaTypeToGoType(f.Type),
		GetterName:   fmt.Sprintf("get%s", snakeCaseToCamelCase(f.Name)),
		DefaultValue: NewDefaultValue(f),
	}
}

type ParamDescriptor struct {
	// LuaType string
	// GoType  string
	Type TypeDescriptor
	Name string
}

func (pd *ParamDescriptor) String() string {
	return fmt.Sprintf("Param: %s, %s", pd.Name, pd.Type)
}

func NewParamDescriptor(ty, name string) *ParamDescriptor {
	// luaType := ty
	// goType := luaTypeToGoType(ty)
	return &ParamDescriptor{
		// LuaType: luaType,
		// GoType:  goType,
		Type: NewTypeDescriptor(ty),
		Name: name,
	}
}

type MethodDescriptor struct {
	Name    string
	LuaName string
	Kind    string
	Types   TypeDescriptors
	Target  string
}

func NewMethodDescriptor(m *Method, tableName string) *MethodDescriptor {
	name := fmt.Sprintf("__lua_%s_%s", tableName, m.Name)
	return &MethodDescriptor{
		Name:    name,
		LuaName: m.Name,
		Kind:    m.Kind,
		Types:   NewTypeDescriptors(m.Type),
		Target:  m.Target,
	}
}

type ConstructorDescription struct {
	NParams int
	Params  []*ParamDescriptor
}

func NewConstructorDescription(c *Constructor) *ConstructorDescription {
	if c == nil {
		return nil
	}
	types := strings.Split(c.Type, ",")
	nparams := len(types)
	var params []*ParamDescriptor
	for n, ty := range types {
		params = append(params, NewParamDescriptor(ty, fmt.Sprintf("__p%d", n)))
	}
	return &ConstructorDescription{
		NParams: nparams,
		Params:  params,
	}
}

type (
	FieldMap  map[string]*FieldDescriptor
	MethodMap map[string]*MethodDescriptor
)

type TableGenerator struct {
	TableName                     string
	LuaMappingName                string
	TableIntegrityCheckerName     string
	LuaNewTableName               string
	NewTableName                  string
	PublicLuaNewTableName         string
	PublicNewTableName            string
	TableTypeName                 string
	PublicTableTypeName           string
	ConvertTableFromLuaName       string
	PackageName                   string
	PublicInterfaceName           string
	PublicConvertTableFromLuaName string
	Fields                        FieldMap
	Methods                       MethodMap
	Constructor                   *ConstructorDescription
}

func (m FieldMap) List() []*FieldDescriptor {
	var fields []*FieldDescriptor
	for _, f := range m {
		fields = append(fields, f)
	}
	return fields
}

func (m MethodMap) List() []*MethodDescriptor {
	var methods []*MethodDescriptor
	for _, f := range m {
		methods = append(methods, f)
	}
	return methods
}

type TableGeneratorOption func(*TableGenerator)

func WithPackageName(name string) TableGeneratorOption {
	return func(t *TableGenerator) {
		t.PackageName = name
	}
}

func WithPublicInterface(name string) TableGeneratorOption {
	return func(t *TableGenerator) {
		t.PublicInterfaceName = name
	}
}

func NewTableGenerator(t *Table, opts ...TableGeneratorOption) *TableGenerator {
	tg := &TableGenerator{
		TableName:                     t.Name,
		LuaMappingName:                fmt.Sprintf("__%s_FunctionsMapping", t.Name),
		TableIntegrityCheckerName:     fmt.Sprintf("__%s_CheckTableIntegrity", t.Name),
		LuaNewTableName:               fmt.Sprintf("__%s_LuaNewTable", t.Name),
		NewTableName:                  fmt.Sprintf("__%s_NewTable", t.Name),
		PublicLuaNewTableName:         "",
		PublicNewTableName:            "",
		TableTypeName:                 fmt.Sprintf("__%s_Table", t.Name),
		PublicTableTypeName:           "",
		ConvertTableFromLuaName:       fmt.Sprintf("__%s_ConvertTableFromLuaTable", t.Name),
		PackageName:                   "",
		PublicInterfaceName:           "",
		PublicConvertTableFromLuaName: "",
		Fields:                        make(map[string]*FieldDescriptor),
		Methods:                       make(map[string]*MethodDescriptor),
		Constructor:                   NewConstructorDescription(t.Constructor),
	}
	for _, field := range t.Fields {
		tg.Fields[field.Name] = NewFieldDescriptor(&field)
	}
	for _, method := range t.Methods {
		tg.Methods[method.Name] = NewMethodDescriptor(&method, t.Name)
	}
	for _, opt := range opts {
		opt(tg)
	}
	if len(tg.PublicInterfaceName) > 0 {
		tg.PublicTableTypeName = tg.PublicInterfaceName
		tg.PublicConvertTableFromLuaName = fmt.Sprintf("Get%sFromLuaTable", tg.PublicInterfaceName)
		tg.PublicLuaNewTableName = fmt.Sprintf("LuaNew%s", tg.PublicInterfaceName)
		tg.PublicNewTableName = fmt.Sprintf("New%s", tg.PublicInterfaceName)
	}
	return tg
}

func (tg *TableGenerator) GetConstructorParams() []*ParamDescriptor {
	if tg.Constructor == nil {
		return nil
	}
	return tg.Constructor.Params
}

func (tg *TableGenerator) FieldExists(name string) bool {
	_, ok := tg.Fields[name]
	return ok
}

func (tg *TableGenerator) GetField(name string) *FieldDescriptor {
	return tg.Fields[name]
}

func genTypeCheckCond(tys TypeDescriptors, varName string) string {
	var checks []string

	for _, ty := range tys {
		checks = append(checks, ty.LuaTypeCheck(varName))
	}
	return fmt.Sprintf("!(%s)", strings.Join(checks, " || "))
}

var templateFuncMap = template.FuncMap{
	"genTypeCheck":      genTypeCheck,
	"genTypeCheckError": genTypeCheckError,
	"genAppendForTypes": genAppendForTypes,
	"genFieldInit":      genFieldInit,
	"genFieldsInit":     genFieldsInit,
	"print":             func(s ...string) string { return fmt.Sprintf("%s", strings.Join(s, " -- ")) },
}

func getDefaultValue(field *FieldDescriptor, params []*ParamDescriptor) string {
	switch def := field.DefaultValue.(type) {
	case *DefaultFromParam:
		if def.ParamId > len(params) {
			panic(fmt.Sprintf("Cannot index param %d for field %s", def.ParamId, field.Name))
		}
		param := params[def.ParamId-1]
		if !field.Types.Contains(param.Type) {
			panic(fmt.Sprintf("Incompatible types %s and %s for field %s", field.Types, param.Type.LuaType(),
				field.Name))
		}
		return param.Name
	case *NoDefault:
		return field.Types.GetDefaultValue()
	default:
		return field.Types.GetDefaultValue()
	}
}

func genFieldInit(field *FieldDescriptor, params []*ParamDescriptor, tableName string) string {
	linit := getDefaultValue(field, params)
	return fmt.Sprintf(`L.SetField(%s, "%s", %s)`, tableName, field.Name, linit)
}

func genFieldsInit(fields []*FieldDescriptor, params []*ParamDescriptor, tableName string) string {
	var inits []string
	for _, field := range fields {
		inits = append(inits, genFieldInit(field, params, tableName))
	}
	return strings.Join(inits, "\n")
}

func genAppendForTypes(tys TypeDescriptors, varName string, tabName string) string {
	var appenders []string
	for _, t := range tys {
		var temp string
		if t.IsTable() {
			temp = fmt.Sprintf(`if %s {
	L.ForEach({{.VarName}}.(*lua.LTable), func(_, v lua.LValue) {
		{{.TabName}}.Append(v)
	})
}`, t.LuaTypeCheck(varName))
		} else {
			temp = fmt.Sprintf(`if %s {
	{{.TabName}}.Append({{.VarName}})
}`, t.LuaTypeCheck(varName))
		}
		appenders = append(appenders,
			MustExecuteTemplate("genAppendForTypes", temp, template.FuncMap{},
				struct {
					TabName string
					VarName string
				}{
					TabName: tabName,
					VarName: varName,
				}))
	}
	return strings.Join(appenders, " else ")
}

func genTableTypeCheck(ty TypeDescriptor, varName string) string {
	templateFuncMap := template.FuncMap{
		"genTypeCheckCond": genTypeCheckCond,
		"genTypeCheck":     genTypeCheck,
	}
	return MustExecuteTemplate("genTableTypeCheck", tableTypeCheckTemplate, templateFuncMap,
		struct {
			VarType TypeDescriptors
			VarName string
		}{
			VarType: TypeDescriptors{ty.(*TableDescriptor).InnerType},
			VarName: varName,
		})
}

func genTypeCheck(tys TypeDescriptors, varName string) string {
	templateFuncMap := template.FuncMap{
		"genTypeCheckCond": genTypeCheckCond,
		"genTypeCheck":     genTypeCheck,
	}
	check := MustExecuteTemplate("genTypeCheck", typeCheckTemplate, templateFuncMap,
		struct {
			VarTypes []TypeDescriptor
			VarName  string
		}{
			VarTypes: tys,
			VarName:  varName,
		})
	var tableChecks []string
	for _, ty := range tys {
		if ty.IsTable() {
			genTableTypeCheck(ty, varName)
		}
	}

	if len(tableChecks) > 0 {
		return fmt.Sprintf("%s\n%s", check, strings.Join(tableChecks, "\n"))
	}
	return check
}

func genTableTypeCheckError(ty TypeDescriptor, varName string) string {
	templateFuncMap := template.FuncMap{
		"genTypeCheckCond":  genTypeCheckCond,
		"genTypeCheckError": genTypeCheckError,
	}
	return MustExecuteTemplate("getTableTypeCheckError", tableTypeCheckErrorTemplate, templateFuncMap,
		struct {
			VarTypes TypeDescriptors
			VarName  string
		}{
			VarTypes: TypeDescriptors{ty.(*TableDescriptor).InnerType},
			VarName:  varName,
		})
}

func genTypeCheckError(tys []TypeDescriptor, varName, errName string) string {
	templateFuncMap := template.FuncMap{
		"genTypeCheckCond": genTypeCheckCond,
		"genTypeCheck":     genTypeCheck,
	}
	check := MustExecuteTemplate("genTypeCheckError", typeCheckErrorTemplate, templateFuncMap,
		struct {
			VarTypes []TypeDescriptor
			ErrName  string
			VarName  string
		}{
			VarTypes: tys,
			ErrName:  errName,
			VarName:  varName,
		})
	var tableChecks []string
	for _, ty := range tys {
		if ty.IsTable() {
			genTableTypeCheckError(ty, varName)
		}
	}

	if len(tableChecks) > 0 {
		return fmt.Sprintf("%s\n%s", check, strings.Join(tableChecks, "\n"))
	}
	return check
}

func (tg *TableGenerator) GenSetter(m *MethodDescriptor) string {
	if !tg.FieldExists(m.Target) {
		panic("!tg.FieldExists(m.Target)")
	}
	f := tg.GetField(m.Target)

	ut, err := template.New("GenSetter").Funcs(templateFuncMap).Parse(methodSetterTemplate)
	if err != nil {
		panic(err.Error())
	}
	var buff bytes.Buffer
	err = ut.Execute(&buff, struct {
		MethodName string
		FieldName  string
		ValueType  TypeDescriptors
	}{
		MethodName: m.Name,
		FieldName:  f.Name,
		ValueType:  m.Types,
	})
	if err != nil {
		panic(err.Error())
	}
	return buff.String()
}

var tableTypeRe = regexp.MustCompile(`^Table\(([^)]*)\)$`)

func typeIsTable(t string) bool {
	return tableTypeRe.MatchString(t)
}

func typesHasTable(types []TypeDescriptor) bool {
	for _, ty := range types {
		if ty.IsTable() {
			return true
		}
	}
	return false
}

func getInnerType(types string) string {
	for _, t := range strings.Split(types, "|") {
		if typeIsTable(t) {
			return tableTypeRe.FindAllStringSubmatch(t, -1)[0][1]
		}
	}
	return ""
}

func (tg *TableGenerator) GenAppend(m *MethodDescriptor) string {
	if !tg.FieldExists(m.Target) {
		panic("!tg.FieldExists(m.Target)")
	}
	f := tg.GetField(m.Target)
	if !f.Types.HasTable() {
		panic(fmt.Sprintf("Cannot append to non table field: %s", f.Types))
	}
	return MustExecuteTemplate("GenAppend", methodAppendTemplate, templateFuncMap,
		struct {
			MethodName string
			FieldName  string
			ValueType  TypeDescriptors
			TargetType TypeDescriptors
		}{
			MethodName: m.Name,
			FieldName:  f.Name,
			ValueType:  m.Types,
			TargetType: f.Types,
		})
}

func (tg *TableGenerator) GenMethod(m *MethodDescriptor) string {
	switch m.Kind {
	case "Setter":
		return tg.GenSetter(m)
	case "Append":
		return tg.GenAppend(m)
	}
	panic("Unknown kind")
}

func (tg *TableGenerator) GenIntegrityChecker() string {
	return MustExecuteTemplate("GenIntegrityChecker", checkTableIntegrityTemplate, templateFuncMap,
		struct {
			FuncName string
			Fields   []*FieldDescriptor
		}{
			FuncName: tg.TableIntegrityCheckerName,
			Fields:   tg.Fields.List(),
		})
}

func (tg *TableGenerator) GenNewTable() string {
	paramDecls := []string{"L *lua.LState"}
	if tg.Constructor != nil {
		for _, param := range tg.Constructor.Params {
			paramDecls = append(paramDecls, fmt.Sprintf("%s *%s", param.Name,
				param.Type.LuaType()))
		}
	}

	return MustExecuteTemplate("GenNewTable", newTableTemplate, templateFuncMap,
		struct {
			FuncName        string
			Fields          []*FieldDescriptor
			FunctionMapping string
			ParamsDecl      string
			Params          []*ParamDescriptor
		}{
			FuncName:        tg.NewTableName,
			Fields:          tg.Fields.List(),
			FunctionMapping: tg.LuaMappingName,
			ParamsDecl:      strings.Join(paramDecls, ", "),
			Params:          tg.GetConstructorParams(),
		})
}

func genParamGets(params []*ParamDescriptor) []string {
	var gets []string
	for n, param := range params {
		gets = append(gets, fmt.Sprintf("%s := L.Get(%d)", param.Name, n+1))
	}
	return gets
}

func (tg *TableGenerator) GenLuaNewTable() string {
	var paramTypeChecks []string
	if tg.Constructor != nil {
		for _, param := range tg.Constructor.Params {
			paramTypeChecks = append(paramTypeChecks, genTypeCheck([]TypeDescriptor{param.Type}, param.Name))
		}
	}
	paramUse := []string{"L"}
	if tg.Constructor != nil {
		for _, param := range tg.Constructor.Params {
			paramUse = append(paramUse, fmt.Sprintf("%s.(*%s)", param.Name, param.Type.LuaType()))
		}
	}

	return MustExecuteTemplate("GenLuaNewTable", luaNewTableTemplate, templateFuncMap,
		struct {
			FuncName        string
			NewFuncName     string
			ParamGets       string
			ParamTypeChecks string
			ParamsUse       string
		}{
			FuncName:        tg.LuaNewTableName,
			NewFuncName:     tg.NewTableName,
			ParamGets:       strings.Join(genParamGets(tg.GetConstructorParams()), "\n"),
			ParamTypeChecks: strings.Join(paramTypeChecks, "\n"),
			ParamsUse:       strings.Join(paramUse, ", "),
		})
}

func (tg *TableGenerator) GenLuaMappingTable() string {
	var mappingList []struct {
		LuaName string
		GoName  string
	}
	for _, v := range tg.Methods {
		mappingList = append(mappingList, struct {
			LuaName string
			GoName  string
		}{
			LuaName: v.LuaName,
			GoName:  v.Name,
		})
	}
	var buff bytes.Buffer
	t, err := template.New("luaMapping").Parse(luaFunctionMapTemplate)
	if err != nil {
		panic(err.Error())
	}
	err = t.Execute(&buff, struct {
		FuncMapName string
		Mappings    []struct {
			LuaName string
			GoName  string
		}
	}{
		FuncMapName: tg.LuaMappingName,
		Mappings:    mappingList,
	})
	if err != nil {
		panic(err.Error())
	}
	return buff.String()
}

func genGetGoStringFromLuaField(varName, fieldName, tableName string) string {
	t, err := template.New("genGetGoStringFromLuaField").Parse(getLuaStringFromTableFieldTemplate)
	if err != nil {
		panic(err)
	}
	var buff bytes.Buffer
	err = t.Execute(&buff, struct {
		VarName   string
		TableName string
		FieldName string
	}{
		VarName:   varName,
		TableName: tableName,
		FieldName: fieldName,
	})
	if err != nil {
		panic(err)
	}
	return buff.String()
}

func genGetGoObjectFromLuaField(varName, fieldName, tableName string, ty *CustomTypeDescriptor) string {
	return MustExecuteTemplate("genGetGoObjectFromLuaField", getLuaObjectFromTableFieldTemplate,
		template.FuncMap{},
		struct {
			VarName       string
			TableName     string
			FieldName     string
			ConverterName string
			VarType       string
		}{
			VarName:       varName,
			TableName:     tableName,
			FieldName:     fieldName,
			ConverterName: ty.ConverterFunc,
			VarType:       ty.GoType(),
		})
}

func genGetGoTableFromLuaField(varName, fieldName, tableName string, ty TypeDescriptor) string {
	return MustExecuteTemplate("genGetGoTableFromLuaField", getLuaTableFromTableFieldTemplate,
		template.FuncMap{
			"genGetGoValueFromValue": genGetGoValueFromValue,
		},
		struct {
			VarName   string
			TableName string
			FieldName string
			Type      TypeDescriptor
			GoType    string
		}{
			VarName:   varName,
			TableName: tableName,
			FieldName: fieldName,
			Type:      ty,
			GoType:    ty.GoType(),
		})
}

// func genGetGoStringFromValue(varName, luaVarName string) string {
// 	t, err := template.New("genGetGoStringFromValue").Parse(getLuaStringFromValueTemplate)
// 	if err != nil {
// 		panic(err)
// 	}
// 	var buff bytes.Buffer
// 	err = t.Execute(&buff, struct {
// 		VarName    string
// 		LuaVarName string
// 	}{
// 		VarName:    varName,
// 		LuaVarName: luaVarName,
// 	})
// 	if err != nil {
// 		panic(err)
// 	}
// 	return buff.String()
// }

func genGetGoValueFromLuaField(varName, fieldName, tableName string, tys TypeDescriptors) string {
	switch ty := tys[0].(type) {
	case *StringDescriptor:
		return genGetGoStringFromLuaField(varName, fieldName, tableName)
	case *TableDescriptor:
		return genGetGoTableFromLuaField(varName, fieldName, tableName, ty.InnerType)
	case *CustomTypeDescriptor:
		return genGetGoObjectFromLuaField(varName, fieldName, tableName, ty)
	default:
		panic(fmt.Sprintf("genGetGoValueFromLuaField: Unhandled type %s", ty))
	}
}

func genGetGoValueFromValue(varName, luaVarName string, ty TypeDescriptor) string {
	switch ty.(type) {
	case *StringDescriptor:
		return fmt.Sprintf("%s := %s.String()", varName, luaVarName)
	default:
		panic(fmt.Sprintf("genGetGoValueFromValue: Unhandled type %s", ty))
	}
	// switch luaType {
	// case "String":
	// 	return genGetGoStringFromValue(varName, luaVarName)
	// default:
	// 	panic(fmt.Sprintf("genGetGoValueFromValue: Unhandled type %s", luaType))
	// }
}

func (tg *TableGenerator) GenConverterFromLuaTable() string {
	return MustExecuteTemplate("GenConverterFromLuaTable", tableConversionTemplate,
		template.FuncMap{
			"genGetGoValueFromLuaField": genGetGoValueFromLuaField,
		},
		struct {
			FuncName       string
			TypeName       string
			Fields         []*FieldDescriptor
			CheckIntegrity string
		}{
			FuncName:       tg.ConvertTableFromLuaName,
			TypeName:       tg.TableTypeName,
			Fields:         tg.Fields.List(),
			CheckIntegrity: tg.TableIntegrityCheckerName,
		})
}

func (f *FieldDescriptor) GenTypeDefField() string {
	return fmt.Sprintf("%s %s", f.GoName, f.Types.GoType())
}

func (tg *TableGenerator) GenTableTypeDefinition() string {
	return MustExecuteTemplate("GenTableTypeDefinition",
		`type {{.TableTypeName}} struct {
	{{- range .Fields.List }}
		{{.GenTypeDefField}}
	{{- end }}
}`,
		template.FuncMap{}, tg)
}

func (tg *TableGenerator) GenPublicTableTypeDefinition() string {
	if len(tg.PublicInterfaceName) > 0 {
		return fmt.Sprintf("type %s %s", tg.PublicTableTypeName, tg.TableTypeName)
	}
	return ""
}

func (tg *TableGenerator) GenPublicTableGetters() string {
	if len(tg.PublicInterfaceName) > 0 {
		t, err := template.New("GenPublicTableGetters").Parse(publicTableGetterTemplate)
		if err != nil {
			panic(err)
		}
		var getters []string
		for _, field := range tg.Fields {
			var buff bytes.Buffer
			err = t.Execute(&buff, struct {
				TableName string
				FuncName  string
				TypeName  string
				FieldName string
			}{
				TableName: tg.PublicTableTypeName,
				FuncName:  field.GetterName,
				TypeName:  field.Types.GoType(),
				FieldName: field.GoName,
			})
			if err != nil {
				panic(err)
			}
			getters = append(getters, buff.String())
		}
		return strings.Join(getters, "\n\n")
	}
	return ""
}

func (tg *TableGenerator) GenPublicTableConverter() string {
	if len(tg.PublicInterfaceName) > 0 {
		t, err := template.New("GenPublicTableConverter").Parse(publicConvertTableTemplate)
		if err != nil {
			panic(err)
		}
		var buff bytes.Buffer
		err = t.Execute(&buff, struct {
			TableName   string
			FuncName    string
			ConvertName string
		}{
			TableName:   tg.PublicTableTypeName,
			FuncName:    tg.PublicConvertTableFromLuaName,
			ConvertName: tg.ConvertTableFromLuaName,
		})
		if err != nil {
			panic(err)
		}
		return buff.String()
	}
	return ""
}

func (tg *TableGenerator) GenPublicNewTable() string {
	if len(tg.PublicInterfaceName) > 0 {
		t, err := template.New("GenPublicNewTable").Parse(publicNewTableTemplate)
		if err != nil {
			panic(err)
		}
		paramDecls := []string{"L *lua.LState"}
		paramUse := []string{"L"}
		if tg.Constructor != nil {
			for _, param := range tg.Constructor.Params {
				paramDecls = append(paramDecls, fmt.Sprintf("%s *%s", param.Name,
					param.Type.LuaType()))
				paramUse = append(paramUse, param.Name)
			}
		}
		var buff bytes.Buffer
		err = t.Execute(&buff, struct {
			FuncName   string
			NewTable   string
			ParamsDecl string
			ParamsUse  string
		}{
			FuncName:   tg.PublicNewTableName,
			NewTable:   tg.NewTableName,
			ParamsDecl: strings.Join(paramDecls, ", "),
			ParamsUse:  strings.Join(paramUse, ", "),
		})
		if err != nil {
			panic(err)
		}
		return buff.String()
	}
	return ""
}

func (tg *TableGenerator) GenPublicLuaNewTable() string {
	if len(tg.PublicInterfaceName) > 0 {
		t, err := template.New("GenPublicLuaNewTable").Parse(publicLuaNewTableTemplate)
		if err != nil {
			panic(err)
		}
		var buff bytes.Buffer
		err = t.Execute(&buff, struct {
			FuncName string
			NewTable string
		}{
			FuncName: tg.PublicLuaNewTableName,
			NewTable: tg.LuaNewTableName,
		})
		if err != nil {
			panic(err)
		}
		return buff.String()
	}
	return ""
}

func (tg *TableGenerator) GenFile() string {
	var methods []string
	for _, m := range tg.Methods {
		methods = append(methods, tg.GenMethod(m))
	}

	functionMapping := tg.GenLuaMappingTable()
	integrityChecker := tg.GenIntegrityChecker()
	newTable := tg.GenNewTable()
	luaNewTable := tg.GenLuaNewTable()
	convertLuaTable := tg.GenConverterFromLuaTable()
	typeDefinition := tg.GenTableTypeDefinition()
	publicTypeDefinition := tg.GenPublicTableTypeDefinition()
	publicGetters := tg.GenPublicTableGetters()
	publicConvertLuaTable := tg.GenPublicTableConverter()
	publicNewTable := tg.GenPublicNewTable()
	publicLuaNewTable := tg.GenPublicLuaNewTable()

	t, err := template.New("file").Parse(libraryBodyTemplate)
	if err != nil {
		panic(err.Error())
	}
	var buff bytes.Buffer
	err = t.Execute(&buff,
		struct {
			PackageName           string
			GoMethods             []string
			FunctionMapping       string
			IntegrityChecker      string
			NewTable              string
			LuaNewTable           string
			PublicNewTable        string
			PublicLuaNewTable     string
			TypeDefinition        string
			PublicTypeDefinition  string
			ConvertLuaTable       string
			PublicGetters         string
			PublicConvertLuaTable string
		}{
			PackageName:           tg.PackageName,
			GoMethods:             methods,
			FunctionMapping:       functionMapping,
			IntegrityChecker:      integrityChecker,
			NewTable:              newTable,
			LuaNewTable:           luaNewTable,
			PublicNewTable:        publicNewTable,
			PublicLuaNewTable:     publicLuaNewTable,
			TypeDefinition:        typeDefinition,
			PublicTypeDefinition:  publicTypeDefinition,
			ConvertLuaTable:       convertLuaTable,
			PublicGetters:         publicGetters,
			PublicConvertLuaTable: publicConvertLuaTable,
		})
	if err != nil {
		panic(err.Error())
	}
	return buff.String()
}
