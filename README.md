# Intro

hvue is a [GopherJS](https://github.com/gopherjs/gopherjs) and wasm wrapper
for the [Vue](https://vuejs.org/) Javascript framework.

WASM code is under development.

# Install

## Install Go 1.11

See https://golang.org/dl/.

## Install hvue's wasm branch

```bash
cd path/to/github.com # in your $GOPATH
mkdir huckridgesw
cd huckridgesw
git clone git@github.com:HuckRidgeSW/hvue.git
cd hvue
git checkout --track origin/wasm
```

# Examples & Demos

## Overview

***Note: As of this writing, only examples 01-07 have been converted & tested to run under both GopherJS and go/wasm.***

Generally speaking, the [examples](https://github.com/HuckRidgeSW/hvue/tree/master/examples)
follow the examples in the Vue [guide](https://vuejs.org/v2/guide/).

[01-introduction](https://github.com/HuckRidgeSW/hvue/tree/master/examples/01-introduction)
has examples from the Vue [Introduction](https://vuejs.org/v2/guide/index.html) page.

[02-lifecycle](https://github.com/HuckRidgeSW/hvue/tree/master/examples/02-lifecycle)
demos Vue [lifecycle hooks](https://vuejs.org/v2/guide/instance.html#Instance-Lifecycle-Hooks)
but does not correspond to any specific example on that page.

[03-computed-basic](https://github.com/HuckRidgeSW/hvue/tree/master/examples/03-computed-basic)
and [04-computed-with-setter](https://github.com/HuckRidgeSW/hvue/tree/master/examples/04-computed-with-setter)
have examples from [Computed Properties and Watchers](https://vuejs.org/v2/guide/computed.html).

And so on.  Links are in the code.

## Running the examples

### GopherJS

```bash
cd path/to/github.com/huckridgesw/hvue
echo "var hvue_wasm = false;" > examples/maybe_wasm.js
gopherjs serve github.com/huckridgesw/hvue # listens on 8080
```

and then
- http://localhost:8080/examples/01-introduction/
- http://localhost:8080/examples/02-lifecycle/
- http://localhost:8080/examples/ for more

### WASM

```bash
cd path/to/github.com/huckridgesw/hvue
echo "var hvue_wasm = true;" > examples/maybe_wasm.js
go run examples/server/main.go # Listens on 8081
cd examples/??-???? # some examples directory
GOARCH=wasm GOOS=js go build -o ${PWD##*/}.wasm main.go # compile wasm
```

and then
- http://localhost:8081/examples/01-introduction/
- http://localhost:8081/examples/02-lifecycle/
- http://localhost:8081/examples/ for more

Remember to recompile after any changes.  There's no facility yet to
auto-build yet (a-la `gopherjs build -w` or `gopherjs serve`).

### Switching from GopherJS to WASM and back

- Do the appropriate `"echo "var hvue_wasm = ?;" > examples/maybe_wasm.js`.
  (See above.)
- Be sure to do "shift-cmd-R" (Chrome, macOS; other browsers / OSes will vary)
  to reload without using the cache, to get the new `maybe_wasm.js` and/or new
  wasm.  (Actually I'm not sure you need that to get new wasm, since it's
  loaded via an explicit `fetch()` call, but it's probably not a bad idea.)
  Alternatively, in Chrome you can open the developer console, go to the
  network tab, and check "disable cache".  (AIUI only works while said console
  window is open.)

# GoDoc

http://godoc.org/github.com/HuckRidgeSW/hvue

