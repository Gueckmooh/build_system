package main

import (
	"bytes"
	"fmt"
	"regexp"
	"strings"
	"text/template"
)

type FieldDescriptor struct {
	Name       string
	GoName     string
	Type       string
	GoType     string
	GetterName string
}

func NewFieldDescriptor(f *Field) *FieldDescriptor {
	var name string
	if f.Private {
		name = fmt.Sprintf("_%s_", f.Name)
	} else {
		name = f.Name
	}
	return &FieldDescriptor{
		Name:       name,
		GoName:     fmt.Sprintf("__v%s", f.Name),
		Type:       f.Type,
		GoType:     luaTypeToGoType(f.Type),
		GetterName: fmt.Sprintf("get%s", snakeCaseToCamelCase(f.Name)),
	}
}

type MethodDescriptor struct {
	Name    string
	LuaName string
	Kind    string
	Type    string
	Target  string
}

func NewMethodDescriptor(m *Method, tableName string) *MethodDescriptor {
	name := fmt.Sprintf("__lua_%s_%s", tableName, m.Name)
	return &MethodDescriptor{
		Name:    name,
		LuaName: m.Name,
		Kind:    m.Kind,
		Type:    m.Type,
		Target:  m.Target,
	}
}

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
	Fields                        map[string]*FieldDescriptor
	Methods                       map[string]*MethodDescriptor
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
	}
	for _, field := range t.Fields.Field {
		tg.Fields[field.Name] = NewFieldDescriptor(&field)
	}
	for _, method := range t.Methods.Method {
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

func (tg *TableGenerator) FieldExists(name string) bool {
	_, ok := tg.Fields[name]
	return ok
}

func (tg *TableGenerator) GetField(name string) *FieldDescriptor {
	return tg.Fields[name]
}

func genTypeCheckCond(typesString string, varName string) string {
	types := strings.Split(typesString, ",")
	var typesToCheck []string
	for _, t := range types {
		switch t {
		case "String":
			typesToCheck = append(typesToCheck, "lua.LTString")
		default:
			if typeIsTable(t) {
				typesToCheck = append(typesToCheck, "lua.LTTable")
			} else {
				panic(fmt.Sprintf("genTypeCheckCond: unhandled type %s", t))
			}
		}
	}
	var typesChecks []string
	for _, t := range typesToCheck {
		typesChecks = append(typesChecks, fmt.Sprintf("%s.Type() == %s", varName, t))
	}
	return fmt.Sprintf("!(%s)", strings.Join(typesChecks, " || "))
}

var templateFuncMap = template.FuncMap{
	"genTypeCheck":      genTypeCheck,
	"genTypeCheckError": genTypeCheckError,
	"genAppendForTypes": genAppendForTypes,
	"genFieldInit":      genFieldInit,
	"print":             func(s ...string) string { return fmt.Sprintf("%s", strings.Join(s, " -- ")) },
}

func genFieldInit(field *FieldDescriptor, tableName string) string {
	var linit string
	switch field.Type {
	case "String":
		linit = "lua.LString(\"\")"
	default:
		if typeIsTable(field.Type) {
			linit = "L.NewTable()"
		} else {
			panic("Unhandled type")
		}
	}
	return fmt.Sprintf(`L.SetField(%s, "%s", %s)`, tableName, field.Name, linit)
}

func genAppendForTypes(typesString string, varName string, tabName string) string {
	types := strings.Split(typesString, ",")
	var appenders []string
	for _, t := range types {
		var temp string
		switch t {
		case "String":
			temp = `if {{.VarName}}.Type() == lua.LTString {
	{{.TabName}}.Append({{.VarName}})
}`
		default:
			if typeIsTable(t) {
				temp = `if {{.VarName}}.Type() == lua.LTTable {
	L.ForEach({{.VarName}}.(*lua.LTable), func(_, v lua.LValue) {
		{{.TabName}}.Append(v)
	})
}`
			} else {
				panic(fmt.Sprintf("genAppendForTypes: unhandled type %s", t))
			}
		}
		t, err := template.New("genAppendForTypes").Parse(temp)
		if err != nil {
			panic(err.Error())
		}
		var buff bytes.Buffer
		err = t.Execute(&buff, struct {
			TabName string
			VarName string
		}{
			TabName: tabName,
			VarName: varName,
		})
		if err != nil {
			panic(err.Error())
		}
		appenders = append(appenders, buff.String())
	}
	return strings.Join(appenders, " else ")
}

func genTypeCheck(typesString string, varName string) string {
	templateFuncMap := template.FuncMap{
		"genTypeCheckCond": genTypeCheckCond,
		"genTypeCheck":     genTypeCheck,
	}
	var firstCheck string
	{
		t, err := template.New("typeCheck").Funcs(templateFuncMap).Parse(typeCheckTemplate)
		if err != nil {
			panic(err.Error())
		}
		var buff bytes.Buffer
		err = t.Execute(&buff, struct {
			VarType string
			VarName string
		}{
			VarType: typesString,
			VarName: varName,
		})
		if err != nil {
			panic(err.Error())
		}
		firstCheck = buff.String()
	}
	var secondCheck string
	if typesHasTable(typesString) {
		types := getInnerType(typesString)
		t, err := template.New("typeCheckTable").Funcs(templateFuncMap).Parse(tableTypeCheckTemplate)
		if err != nil {
			panic(err.Error())
		}
		var buff bytes.Buffer
		err = t.Execute(&buff, struct {
			VarType string
			VarName string
		}{
			VarType: types,
			VarName: varName,
		})
		if err != nil {
			panic(err.Error())
		}
		secondCheck = buff.String()
	}
	if len(secondCheck) > 0 {
		return fmt.Sprintf("%s\n%s", firstCheck, secondCheck)
	}
	return firstCheck
}

func genTypeCheckError(typesString string, varName string, errName string) string {
	templateFuncMap := template.FuncMap{
		"genTypeCheckCond":  genTypeCheckCond,
		"genTypeCheckError": genTypeCheckError,
	}
	var firstCheck string
	{
		t, err := template.New("typeCheckError").Funcs(templateFuncMap).Parse(typeCheckErrorTemplate)
		if err != nil {
			panic(err.Error())
		}
		var buff bytes.Buffer
		err = t.Execute(&buff, struct {
			VarType string
			VarName string
			ErrName string
		}{
			VarType: typesString,
			VarName: varName,
			ErrName: errName,
		})
		if err != nil {
			panic(err.Error())
		}
		firstCheck = buff.String()
	}
	var secondCheck string
	if typesHasTable(typesString) {
		types := getInnerType(typesString)
		t, err := template.New("tableTypeCheckError").Funcs(templateFuncMap).Parse(tableTypeCheckErrorTemplate)
		if err != nil {
			panic(err.Error())
		}
		var buff bytes.Buffer
		err = t.Execute(&buff, struct {
			VarType string
			VarName string
		}{
			VarType: types,
			VarName: varName,
		})
		if err != nil {
			panic(err.Error())
		}
		secondCheck = buff.String()
	}
	if len(secondCheck) > 0 {
		return fmt.Sprintf("%s\n%s", firstCheck, secondCheck)
	}
	return firstCheck
}

func (tg *TableGenerator) GenSetter(m *MethodDescriptor) string {
	if !tg.FieldExists(m.Target) {
		panic("!tg.FieldExists(m.Target)")
	}
	f := tg.GetField(m.Target)

	ut, err := template.New("method").Funcs(templateFuncMap).Parse(methodSetterTemplate)
	if err != nil {
		panic(err.Error())
	}
	var buff bytes.Buffer
	err = ut.Execute(&buff, struct {
		MethodName string
		FieldName  string
		ValueType  string
	}{
		MethodName: m.Name,
		FieldName:  f.Name,
		ValueType:  m.Type,
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

func typesHasTable(types string) bool {
	for _, t := range strings.Split(types, ",") {
		if typeIsTable(t) {
			return true
		}
	}
	return false
}

func getInnerType(types string) string {
	for _, t := range strings.Split(types, ",") {
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
	if !typeIsTable(f.Type) {
		panic(fmt.Sprintf("Cannot append to non table field: %s", f.Type))
	}
	ut, err := template.New("method").Funcs(templateFuncMap).Parse(methodAppendTemplate)
	if err != nil {
		panic(err.Error())
	}
	var buff bytes.Buffer
	err = ut.Execute(&buff, struct {
		MethodName string
		FieldName  string
		ValueType  string
		TargetType string
	}{
		MethodName: m.Name,
		FieldName:  f.Name,
		ValueType:  m.Type,
		TargetType: f.Type,
	})
	if err != nil {
		panic(err.Error())
	}
	return buff.String()
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
	t, err := template.New("GenIntegrityChecker").Funcs(templateFuncMap).Parse(checkTableIntegrityTemplate)
	if err != nil {
		panic(err.Error())
	}
	var fields []*FieldDescriptor
	for _, v := range tg.Fields {
		fields = append(fields, v)
	}
	var buff bytes.Buffer
	err = t.Execute(&buff, struct {
		FuncName string
		Fields   []*FieldDescriptor
	}{
		FuncName: tg.TableIntegrityCheckerName,
		Fields:   fields,
	})
	if err != nil {
		panic(err.Error())
	}
	return buff.String()
}

func (tg *TableGenerator) GenNewTable() string {
	t, err := template.New("GenNewTable").Funcs(templateFuncMap).Parse(newTableTemplate)
	if err != nil {
		panic(err.Error())
	}
	var fields []*FieldDescriptor
	for _, v := range tg.Fields {
		fields = append(fields, v)
	}
	var buff bytes.Buffer
	err = t.Execute(&buff, struct {
		FuncName        string
		Fields          []*FieldDescriptor
		FunctionMapping string
	}{
		FuncName:        tg.NewTableName,
		Fields:          fields,
		FunctionMapping: tg.LuaMappingName,
	})
	if err != nil {
		panic(err.Error())
	}
	return buff.String()
}

func (tg *TableGenerator) GenLuaNewTable() string {
	t, err := template.New("GenNewTable").Funcs(templateFuncMap).Parse(luaNewTableTemplate)
	if err != nil {
		panic(err.Error())
	}
	var fields []*FieldDescriptor
	for _, v := range tg.Fields {
		fields = append(fields, v)
	}
	var buff bytes.Buffer
	err = t.Execute(&buff, struct {
		FuncName    string
		NewFuncName string
	}{
		FuncName:    tg.LuaNewTableName,
		NewFuncName: tg.NewTableName,
	})
	if err != nil {
		panic(err.Error())
	}
	return buff.String()
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
	t, err := template.New("luaMapping").Parse(totoTemplate)
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

func genGetGoTableFromLuaField(varName, fieldName, tableName, luaType string) string {
	t, err := template.New("genGetGoTableFromLuaField").Funcs(template.FuncMap{
		"genGetGoValueFromValue": genGetGoValueFromValue,
	}).Parse(getLuaTableFromTableFieldTemplate)
	if err != nil {
		panic(err)
	}
	var buff bytes.Buffer
	err = t.Execute(&buff, struct {
		VarName   string
		TableName string
		FieldName string
		LuaType   string
		GoType    string
	}{
		VarName:   varName,
		TableName: tableName,
		FieldName: fieldName,
		LuaType:   luaType,
		GoType:    luaTypeToGoType(luaType),
	})
	if err != nil {
		panic(err)
	}
	return buff.String()
}

func genGetGoStringFromValue(varName, luaVarName string) string {
	t, err := template.New("genGetGoStringFromValue").Parse(getLuaStringFromValueTemplate)
	if err != nil {
		panic(err)
	}
	var buff bytes.Buffer
	err = t.Execute(&buff, struct {
		VarName    string
		LuaVarName string
	}{
		VarName:    varName,
		LuaVarName: luaVarName,
	})
	if err != nil {
		panic(err)
	}
	return buff.String()
}

func genGetGoValueFromLuaField(varName, fieldName, tableName, luaType string) string {
	switch luaType {
	case "String":
		return genGetGoStringFromLuaField(varName, fieldName, tableName)
	default:
		if typeIsTable(luaType) {
			return genGetGoTableFromLuaField(varName, fieldName, tableName, getInnerType(luaType))
		}
		panic(fmt.Sprintf("genGetGoValueFromLuaField: Unhandled type %s", luaType))
	}
}

func genGetGoValueFromValue(varName, luaVarName, luaType string) string {
	switch luaType {
	case "String":
		return genGetGoStringFromValue(varName, luaVarName)
	default:
		panic(fmt.Sprintf("genGetGoValueFromValue: Unhandled type %s", luaType))
	}
}

func (tg *TableGenerator) GenConverterFromLuaTable() string {
	t, err := template.New("GenConverterFromLuaTable").Funcs(template.FuncMap{
		"genGetGoValueFromLuaField": genGetGoValueFromLuaField,
	}).Parse(tableConversionTemplate)
	if err != nil {
		panic(err)
	}
	var fields []*FieldDescriptor
	for _, v := range tg.Fields {
		fields = append(fields, v)
	}
	var buff bytes.Buffer
	err = t.Execute(&buff, struct {
		FuncName       string
		TypeName       string
		Fields         []*FieldDescriptor
		CheckIntegrity string
	}{
		FuncName:       tg.ConvertTableFromLuaName,
		TypeName:       tg.TableTypeName,
		Fields:         fields,
		CheckIntegrity: tg.TableIntegrityCheckerName,
	})
	if err != nil {
		panic(err)
	}
	return buff.String()
}

func genTypeDefField(f *FieldDescriptor) string {
	return fmt.Sprintf("%s %s", f.GoName, f.GoType)
}

func (tg *TableGenerator) GenTableTypeDefinition() string {
	t, err := template.New("GenTableTypeDefinition").Funcs(template.FuncMap{
		"genTypeDefField": genTypeDefField,
	}).Parse(tableTypeDefTemplate)
	if err != nil {
		panic(err)
	}
	var fields []*FieldDescriptor
	for _, v := range tg.Fields {
		fields = append(fields, v)
	}
	var buff bytes.Buffer
	err = t.Execute(&buff, struct {
		TypeName string
		Fields   []*FieldDescriptor
	}{
		TypeName: tg.TableTypeName,
		Fields:   fields,
	})
	if err != nil {
		panic(err)
	}
	return buff.String()
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
				TypeName:  field.GoType,
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
		var buff bytes.Buffer
		err = t.Execute(&buff, struct {
			FuncName string
			NewTable string
		}{
			FuncName: tg.PublicNewTableName,
			NewTable: tg.NewTableName,
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
