# Intro

hvue is a [GopherJS](https://github.com/gopherjs/gopherjs) and wasm wrapper
for the [Vue](https://vuejs.org/) Javascript framework.

[This (the "master" branch)](https://github.com/HuckRidgeSW/hvue/tree/master)
is the go/wasm + GopherJS version.  It uses
[gopherwasm](https://github.com/gopherjs/gopherwasm) to provide a compatability
layer between go/wasm and GopherJS.  go/wasm is patterned on GopherJS, but
doesn't have all of its capabilities and language-specific "magic".  In
particular, go/wasm doesn't have GopherJS's "dual struct/object" magic, which
allows you to embed a *js.Object in a struct, define struct fields with
`js:"jsName"` tags, and have the compiler automatically change references to
those fields into references to fields in the inner *js.Object.  So to access a
JavaScript object in go/wasm, you have to use a "naked" js.Value and either use
`thing.Get("jsField")` (and related functions) everywhere (ew) or write access
functions (less ew).  You can also write GopherJS in the same style, and
gopherwasm creates a compatability layer so the go/wasm style compiles under
GopherJS.

The GopherJS-only version is tagged as
[v1](https://github.com/HuckRidgeSW/hvue/tree/v1), and also as
[gopherjs](https://github.com/HuckRidgeSW/hvue/tree/gopherjs).

The [wasm](https://github.com/HuckRidgeSW/hvue/tree/wasm) branch still exists,
because I shared it pretty widely, and I want those links to keep working for a
while.

So if you want to use the go/wasm code, and/or also use this exact code in
GopherJS, use this branch.  If you want the GopherJS-only code, use the
[v1](https://github.com/HuckRidgeSW/hvue/tree/v1) or
[gopherjs](https://github.com/HuckRidgeSW/hvue/tree/gopherjs) branch.  **Click
on over to the README in that branch for installation instructions.**  They may
need modification, since they still date to when hvue's "master" branch was
GopherJS-only.

# Install

## Install Go 1.11

See https://golang.org/dl/.

## Install hvue

(Side note: If you skipped it, please make sure you've read the Intro above
about the difference between this (the go/wasm + GopherJS code), and previous
GopherJS-only versions.)

```bash
cd path/to/github.com # in your $GOPATH
mkdir huckridgesw
cd huckridgesw
git clone git@github.com:HuckRidgeSW/hvue.git
```

# Examples & Demos

## Overview

Generally speaking, the
[examples](https://github.com/HuckRidgeSW/hvue/tree/master/examples) follow the
examples in the Vue [guide](https://vuejs.org/v2/guide/).  Some don't, because
the Guide has changed since I wrote the examples.  But most of them do.

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

Remember to recompile after any changes.  There's no facility yet to auto-build
(a-la `gopherjs build -w` or `gopherjs serve`).

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

