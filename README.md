# tmplscript
executable go-template command (like awk and jq!)

# usage

```bash
$ go get -u github.com/syuparn/tmplscript
$ echo "hello" | tmplscript '{{print . ", " "world!"}}'
hello, world!
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

Functions in [Sprig](http://masterminds.github.io/sprig/) can be used.

(Of course [template standard functions](https://golang.org/pkg/text/template/#hdr-Functions) too!)
