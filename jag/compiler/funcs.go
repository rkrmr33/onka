package compiler

import (
	"io/ioutil"
	"text/template"

	"github.com/Masterminds/sprig"
)

func funcMap() template.FuncMap {
	m := sprig.TxtFuncMap()

	extra := map[string]interface{}{}
	extra["import"] = importFunc

	for k, v := range extra {
		m[k] = v
	}

	return m
}

func importFunc(path string) (string, error) {
	builtin, exists := builtins[path]
	if exists {
		return builtin, nil
	}

	data, err := ioutil.ReadFile(path)
	if err != nil {
		return "", err
	}

	return string(data), nil
}
