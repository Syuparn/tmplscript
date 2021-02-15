package main

import (
	"bufio"
	"fmt"
	"os"
	"text/template"

	"github.com/Masterminds/sprig"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintf(os.Stderr, "template must be passed to argument")
		return
	}
	tmplStr := os.Args[1]

	// NOTE: FuncMap is for html/template, TxtFuncMap is for text/template
	tmpl, err := template.New("tmpl").Funcs(sprig.TxtFuncMap()).Parse(tmplStr)
	if err != nil {
		fmt.Fprintf(os.Stderr, "syntax error:\n%v", err)
		return
	}

	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		line := scanner.Text()

		err = tmpl.Execute(os.Stdout, line)
		if err != nil {
			fmt.Fprintf(os.Stderr, "runtime error:\n%v", err)
			return
		}
	}
}
