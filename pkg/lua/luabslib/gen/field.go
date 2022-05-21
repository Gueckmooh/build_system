package main

import "fmt"

type Field struct {
	Name string
	Type Type
}

func (f *Field) String() string {
	name := f.Name
	if name == "" {
		name = "<anonymous>"
	}
	return fmt.Sprintf("%s %s", f.Name, f.Type.GoString())
}
