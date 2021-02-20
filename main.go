package main

import (
	"bufio"
	"bytes"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"text/template"
)

var (
	fileName   = flag.String("f", "", "use template file")
	runsREPL   = flag.Bool("i", false, "run interactive REPL instead")
	leftDelim  = flag.String("ldelim", "{{", "specify left deliminater")
	rightDelim = flag.String("rdelim", "}}", "specify right deliminater")
)

func main() {
	flag.Parse()

	tmpl := newTemplate(*leftDelim, *rightDelim)

	if *runsREPL {
		runREPLMode(tmpl)
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

		runPipeMode(tmpl, string(b))
		return
	}

	if flag.Arg(0) == "" {
		fmt.Fprintf(os.Stderr, "template must be passed to argument\n")
		return
	}

	runPipeMode(tmpl, flag.Arg(0))
}

func runPipeMode(tmplGen *template.Template, tmplStr string) {
	tmpl, err := tmplGen.Parse(tmplStr)
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

func runREPLMode(tmplGen *template.Template) {
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

func newTemplate(leftDelim, rightDelim string) *template.Template {
	return template.New("tmpl").Delims(leftDelim, rightDelim).Funcs(funcMap())
}
