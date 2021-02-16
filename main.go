package main

import (
	"bufio"
	"bytes"
	"flag"
	"fmt"
	"os"
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
	// NOTE: FuncMap is for html/template, TxtFuncMap is for text/template
	tmpl, err := template.New("tmpl").Funcs(sprig.TxtFuncMap()).Parse(tmplStr)
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
	tmplGen := template.New("tmpl").Funcs(sprig.TxtFuncMap())
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
