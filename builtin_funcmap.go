package main

import (
	// HACK import internal built-in functions to peek their signatures
	"reflect"

	_ "text/template"
	_ "unsafe"
)

//go:linkname reflectedBuiltinFuncs text/template.builtinFuncs
func reflectedBuiltinFuncs() map[string]reflect.Value
