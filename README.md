# tmplscript
[![GoTest](https://github.com/Syuparn/tmplscript/actions/workflows/test.yml/badge.svg)](https://github.com/Syuparn/tmplscript/actions/workflows/test.yml)
[![MIT License](https://img.shields.io/badge/license-MIT-blue.svg?style=flat)](LICENSE)

executable go-template command (like awk and jq!)

# install

Get binary [here](https://github.com/Syuparn/tmplscript/releases).

You can also install from `go get`.

```bash
$ go get -u github.com/syuparn/tmplscript
```

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

In REPL mode, you can use histories and function autocompletes. See [peterh/liner](https://github.com/peterh/liner) for details.

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
