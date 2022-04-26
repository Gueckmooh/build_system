package main

import (
	"bytes"
	"text/template"
)

func ExecuteTemplate(name string, temp string, funcs template.FuncMap, data any) (string, error) {
	t, err := template.New(name).Funcs(funcs).Parse(temp)
	if err != nil {
		return "", err
	}
	var buff bytes.Buffer
	err = t.Execute(&buff, data)
	if err != nil {
		return "", err
	}
	return buff.String(), nil
}

func MustExecuteTemplate(name string, temp string, funcs template.FuncMap, data any) string {
	t, err := template.New(name).Funcs(funcs).Parse(temp)
	if err != nil {
		panic(err)
	}
	var buff bytes.Buffer
	err = t.Execute(&buff, data)
	if err != nil {
		panic(err)
	}
	return buff.String()
}
