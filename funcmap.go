package main

import (
	"fmt"
	"reflect"
	"sort"
	"strings"
	"text/template"

	sprig "github.com/Masterminds/sprig/v3"
)

func funcMap() template.FuncMap {
	// NOTE: FuncMap is for html/template, TxtFuncMap is for text/template
	funcMap := sprig.TxtFuncMap()

	// add (meta-)functions to describe functions
	funcMap["searchFunc"] = searchFunc(funcMap)
	funcMap["docFunc"] = docFunc(funcMap)
	return funcMap
}

func searchFunc(funcMap template.FuncMap) func(string) []string {
	return func(prefix string) []string {
		keys := []string{}
		for k := range funcMap {
			if strings.HasPrefix(k, prefix) {
				keys = append(keys, k)
			}
		}

		// sort alphabetically
		sort.Strings(keys)

		return keys
	}
}

func docFunc(funcMap template.FuncMap) func(string) string {
	return func(name string) string {
		f, ok := funcMap[name]
		if !ok {
			return fmt.Sprintf("function %s is not defined (or embedded)", name)
		}

		rt := reflect.TypeOf(f)
		if rt.Kind() != reflect.Func {
			return fmt.Sprintf("%s is not a function", name)
		}

		paramTypes := []string{}
		for i := 0; i < rt.NumIn(); i++ {
			paramTypes = append(paramTypes, rt.In(i).String())
		}

		if rt.IsVariadic() {
			paramTypes[len(paramTypes)-1] = toVariadic(paramTypes[len(paramTypes)-1])
		}

		returnTypes := []string{}
		for i := 0; i < rt.NumOut(); i++ {
			returnTypes = append(returnTypes, rt.Out(i).String())
		}

		paramList := strings.Join(paramTypes, " ")
		returnList := strings.Join(returnTypes, " ")

		if paramList == "" {
			return fmt.Sprintf("%s -> (%s)", name, returnList)
		}
		return fmt.Sprintf("%s %s -> (%s)", name, paramList, returnList)
	}
}

func toVariadic(typeStr string) string {
	// replace prefix `[]` to `...` in variadic parameter
	return "..." + strings.TrimPrefix(typeStr, "[]")
}
