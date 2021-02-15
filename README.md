# tmplscript
executable go-template command (like awk and jq!)

# usage

```bash
$ go get -u github.com/syuparn/tmplscript
$ echo "hello" | tmplscript '{{print . ", " "world!"}}'
hello, world!
```

# functions

Functions in [Sprig](http://masterminds.github.io/sprig/) can be used.

(Of course [template standard functions](https://golang.org/pkg/text/template/#hdr-Functions) too!)
