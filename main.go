package main

import (
	"bufio"
	"bytes"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"reflect"
	"strings"
	"text/template"

	"github.com/Masterminds/sprig"
)

var (
	fileName = flag.String("f", "", "use template file")
	runsREPL = flag.Bool("i", false, "run interactive REPL instead")
)

func main() {
	flag.Parse()

	if *runsREPL {
		runREPLMode()
		return
	}

	if *fileName != "" {
		fp, err := os.Open(*fileName)
		if err != nil {
			fmt.Fprintf(os.Stderr, "failed to open %s: %v\n", *fileName, err)
			return
		}

		b, err := ioutil.ReadAll(fp)
		if err != nil {
			fmt.Fprintf(os.Stderr, "failed to read %s: %v\n", *fileName, err)
		}

		runPipeMode(string(b))
		return
	}

	if flag.Arg(0) == "" {
		fmt.Fprintf(os.Stderr, "template must be passed to argument\n")
		return
	}

	runPipeMode(flag.Arg(0))
}

func runPipeMode(tmplStr string) {
	tmpl, err := template.New("tmpl").Funcs(funcMap()).Parse(tmplStr)
	if err != nil {
		fmt.Fprintf(os.Stderr, "template error:\n%v\n", err)
		return
	}

	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		line := scanner.Text()

		err = tmpl.Execute(os.Stdout, line)
		if err != nil {
			fmt.Fprintf(os.Stderr, "runtime error:\n%v\n", err)
			return
		}
	}
}

func runREPLMode() {
	tmplGen := template.New("tmpl").Funcs(funcMap())
	scanner := bufio.NewScanner(os.Stdin)
	tmplStr := ""
	lineNum := 1
	// NOTE: save previous output of template
	// to print only added output, which corresponds to the latest input
	previousOutStr := ""

	for {
		// show prompt
		fmt.Printf("tmpl:%d> ", lineNum)

		ok := scanner.Scan()
		if !ok {
			return
		}

		line := scanner.Text() + "\n"

		// NOTE: whole history is necessary to refer previous variable statement
		tmpl, err := tmplGen.Parse(tmplStr + line)
		if err != nil {
			fmt.Fprintf(os.Stderr, "template error:\n%v\n", err)
			continue
		}

		out := new(bytes.Buffer)
		err = tmpl.Execute(out, nil)
		if err != nil {
			fmt.Fprintf(os.Stderr, "runtime error:\n%v\n", err)
			continue
		}
		// print only added output, which corresponds to the latest input
		// NOTE: break line is unneccessary because it already exists
		outStr := out.String()
		fmt.Print(diffStr(outStr, previousOutStr))

		tmplStr = tmplStr + line
		previousOutStr = outStr
		lineNum++
	}
}

func diffStr(newStr, previousStr string) string {
	return newStr[len(previousStr):]
}

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

		// add `...` to variadic parameter
		if rt.IsVariadic() {
			paramTypes[len(paramTypes)-1] = "..." + paramTypes[len(paramTypes)-1]
		}

		returnTypes := []string{}
		for i := 0; i < rt.NumOut(); i++ {
			returnTypes = append(returnTypes, rt.Out(i).String())
		}

		paramList := strings.Join(paramTypes, " ")
		returnList := strings.Join(returnTypes, " ")
		return fmt.Sprintf("%s %s -> (%s)", name, paramList, returnList)
	}
}
