package main

import (
	"fmt"
	"strings"
)

type TypeDescriptor interface {
	GetDefaultValue() string
	Equals(TypeDescriptor) bool
	Type() string
	GoType() string
	LuaType() string
	LuaTypeCheck(varname string) string
	IsTable() bool
}

type TypeDescriptors []TypeDescriptor

func (tds TypeDescriptors) Contains(t TypeDescriptor) bool {
	for _, td := range tds {
		if t.Equals(td) {
			return true
		}
	}
	return false
}

func (tds TypeDescriptors) String() string {
	var types []string
	for _, td := range tds {
		types = append(types, td.Type())
	}
	return strings.Join(types, "|")
}

func (tds TypeDescriptors) HasTable() bool {
	for _, td := range tds {
		if td.IsTable() {
			return true
		}
	}
	return false
}

func (tds TypeDescriptors) GetType() TypeDescriptor {
	for _, td := range tds {
		if td.IsTable() {
			return td
		}
	}
	return nil
}

func (tds TypeDescriptors) GetDefaultValue() string {
	if len(tds) == 1 {
		return tds[0].GetDefaultValue()
	} else if tds.HasTable() {
		for _, td := range tds {
			if td.IsTable() {
				return td.GetDefaultValue()
			}
		}
	}
	panic("GetDefaultValue called on empty type")
}

func (tds TypeDescriptors) GoType() string {
	if len(tds) == 1 {
		return tds[0].GoType()
	} else if tds.HasTable() {
		for _, td := range tds {
			if td.IsTable() {
				return td.GoType()
			}
		}
	}
	panic("GoType called on empty type")
}

type StringDescriptor struct{}

type TableDescriptor struct {
	InnerType TypeDescriptors
}

type CustomTypeDescriptor struct {
	TypeName           string
	GoTypeName         string
	DefaultValue       string
	TableIntegrityFunc string
	ConverterFunc      string
}

func NewCustomTypeDescriptor(tg *TableGenerator) *CustomTypeDescriptor {
	res := &CustomTypeDescriptor{
		TypeName:           tg.TableName,
		GoTypeName:         "*" + tg.TableTypeName,
		DefaultValue:       fmt.Sprintf("%s(L)", tg.NewTableName),
		TableIntegrityFunc: tg.TableIntegrityCheckerName,
		ConverterFunc:      tg.ConvertTableFromLuaName,
	}
	return res
}

func NewTypeDescriptor(ty string) TypeDescriptor {
	if typeIsTable(ty) {
		return &TableDescriptor{
			InnerType: NewTypeDescriptors(getInnerType(ty)),
		}
	} else {
		switch ty {
		case "String":
			return &StringDescriptor{}
		}
	}
	for _, deps := range Dependencies {
		if deps.TableName == ty {
			return NewCustomTypeDescriptor(deps)
		}
	}
	panic(fmt.Errorf("Could not parse type %s", ty))
}

func NewTypeDescriptors(tys string) TypeDescriptors {
	var tyds []TypeDescriptor
	for _, ty := range strings.Split(tys, "|") {
		tyds = append(tyds, NewTypeDescriptor(ty))
	}
	return tyds
}

func (tds TypeDescriptors) IsUnique() bool {
	return len(tds) == 1
}

func (tds TypeDescriptors) GetTable() *TableDescriptor {
	return (tds[0].(*TableDescriptor))
}

func (tds TypeDescriptors) GenTypeCheckCond(varName string) string {
	var checks []string

	for _, ty := range tds {
		checks = append(checks, ty.LuaTypeCheck(varName))
	}
	return fmt.Sprintf("!(%s)", strings.Join(checks, " || "))
}

func (s *StringDescriptor) GetDefaultValue() string {
	return `lua.LString("")`
}

func (s *TableDescriptor) GetDefaultValue() string {
	return `L.NewTable()`
}

func (s *CustomTypeDescriptor) GetDefaultValue() string {
	return s.DefaultValue
}

func (s *StringDescriptor) Equals(ty TypeDescriptor) bool {
	switch ty.(type) {
	case *StringDescriptor:
		return true
	default:
		return false
	}
}

func (s *TableDescriptor) Equals(ty TypeDescriptor) bool {
	switch ty.(type) {
	case *TableDescriptor:
		return true
	default:
		return false
	}
}

func (s *CustomTypeDescriptor) Equals(ty TypeDescriptor) bool {
	switch t := ty.(type) {
	case *CustomTypeDescriptor:
		return s.TypeName == t.TypeName
	default:
		return false
	}
}

func (s *StringDescriptor) Type() string    { return "String" }
func (s *StringDescriptor) LuaType() string { return "lua.LString" }
func (s *StringDescriptor) LuaTypeCheck(varname string) string {
	return fmt.Sprintf("(%s.Type() == lua.LTString)", varname)
}
func (s *StringDescriptor) GoType() string { return "string" }
func (s *StringDescriptor) String() string { return s.Type() }

func (s *TableDescriptor) Type() string    { return fmt.Sprintf("Table(%s)", s.InnerType) }
func (s *TableDescriptor) LuaType() string { return "lua.LTable" }
func (s *TableDescriptor) LuaTypeCheck(varname string) string {
	return fmt.Sprintf("(%s.Type() == lua.LTTable)", varname)
}
func (s *TableDescriptor) GoType() string                { return fmt.Sprintf("[]%s", s.InnerType.GoType()) }
func (s *TableDescriptor) String() string                { return s.Type() }
func (s *TableDescriptor) GetInnerType() TypeDescriptors { return s.InnerType }

func (s *CustomTypeDescriptor) Type() string    { return s.TypeName }
func (s *CustomTypeDescriptor) LuaType() string { return "lua.LTable" }
func (s *CustomTypeDescriptor) LuaTypeCheck(varname string) string {
	return fmt.Sprintf("(%s(L, %s.(*lua.LTable)) == nil)", s.TableIntegrityFunc, varname)
}
func (s *CustomTypeDescriptor) GoType() string { return s.GoTypeName }
func (s *CustomTypeDescriptor) String() string { return s.Type() }

func (s *StringDescriptor) IsTable() bool     { return false }
func (s *TableDescriptor) IsTable() bool      { return true }
func (s *CustomTypeDescriptor) IsTable() bool { return false }
