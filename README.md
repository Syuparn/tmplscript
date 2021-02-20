# tmplscript
[![GoTest](https://github.com/Syuparn/tmplscript/actions/workflows/test.yml/badge.svg)](https://github.com/Syuparn/tmplscript/actions/workflows/test.yml)

executable go-template command (like awk and jq!)

# usage

```bash
$ go get -u github.com/syuparn/tmplscript
$ echo "hello" | tmplscript '{{print . ", " "world!"}}'
hello, world!
# read from file instead
$ seq 15 | tmplscript -f example/fizzbuzz.tmpl
1
2
fizz
...
```

REPL mode

```bash
$ tmplscript -i
tmpl:1> {{add 2 3}}
5
tmpl:2> {{$i := "Hello"}}

tmpl:3> {{print $i ", world!"}}
Hello, world!
tmpl:4> ^C
```

# functions

Functions in [Sprig](http://masterminds.github.io/sprig/) are available.

(Of course [template standard functions](https://golang.org/pkg/text/template/#hdr-Functions) too!)

Also, you can use `searchFunc` to search defined functions and
`docFunc` to check function's arguments.

```bash
$ tmplscript -i
tmpl:1> {{searchFunc "a"}}
[abbrev adler32sum atoi ago add1 add append abbrevboth]
tmpl:2> {{docFunc "atoi"}}
atoi string -> (int)
```
