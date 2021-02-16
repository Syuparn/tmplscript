package main

import (
	"bufio"
	"bytes"
	"flag"
	"fmt"
	"os"
	"reflect"
	"strings"
	"text/template"

	"github.com/Masterminds/sprig"
)

var (
	runsREPL = flag.Bool("i", false, "run interactive REPL instead")
)

func main() {
	flag.Parse()

	if *runsREPL {
		runREPLMode()
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
		// print only last line, which corresponds to the latest input
		// NOTE: break line is unneccessary because it already exists
		fmt.Print(lastLine(out))

		tmplStr = tmplStr + line
		lineNum++
	}
}

func lastLine(out *bytes.Buffer) string {
	last := ""

	for {
		line, err := out.ReadString('\n')
		if err != nil {
			return last
		}
		last = line
	}
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
	// TODO: impl
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

		returnTypes := []string{}
		for i := 0; i < rt.NumOut(); i++ {
			returnTypes = append(returnTypes, rt.Out(i).String())
		}

		paramList := strings.Join(paramTypes, " ")
		returnList := strings.Join(returnTypes, " ")
		return fmt.Sprintf("%s %s -> (%s)", name, paramList, returnList)
	}
}
