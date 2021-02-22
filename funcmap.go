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
	funcMap["searchFunc"] = searchFunc(funcMap, reflectedBuiltinFuncs())
	funcMap["docFunc"] = docFunc(funcMap, reflectedBuiltinFuncs())
	return funcMap
}

func searchFunc(
	funcMap template.FuncMap,
	builtInFuncMap map[string]reflect.Value,
) func(string) []string {
	return func(prefix string) []string {
		keys := []string{}
		for k := range funcMap {
			if strings.HasPrefix(k, prefix) {
				keys = append(keys, k)
			}
		}

		// search built-in funcs ("and", "call", "print", etc.)
		for k := range builtInFuncMap {
			if strings.HasPrefix(k, prefix) {
				keys = append(keys, k)
			}
		}

		// sort alphabetically
		sort.Strings(keys)

		return keys
	}
}

func docFunc(
	funcMap template.FuncMap,
	builtInFuncMap map[string]reflect.Value,
) func(string) string {
	return func(name string) string {
		rt, ok := findElement(name, funcMap, builtInFuncMap)
		if !ok {
			return fmt.Sprintf("function %s is not defined (or embedded)", name)
		}

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

func findElement(
	name string,
	funcMap template.FuncMap,
	builtInFuncMap map[string]reflect.Value,
) (reflect.Type, bool) {
	elem, ok := funcMap[name]
	if ok {
		return reflect.TypeOf(elem), true
	}

	builtInElem, ok := builtInFuncMap[name]
	if ok {
		return builtInElem.Type(), true
	}

	return nil, false
}

func toVariadic(typeStr string) string {
	// replace prefix `[]` to `...` in variadic parameter
	return "..." + strings.TrimPrefix(typeStr, "[]")
}
