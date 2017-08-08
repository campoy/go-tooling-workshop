# Managing your workspace

Now that we understand how a Go workspace is structured, let's see what tools we can use
to analyze it, navigate it, and keep it up to date.

## Using `go list`

The `go` command has many options, and we'll see many of them during this workshop, but `list`
is definitely not the most well-known one.

`go list` allows you to obtain information about your workspace and the packages stored in it.

_NOTE_: in order to list all the packages in the standard library you can simply run `go list std`.

For instance, given that you've ran `go get github.com/golang/example/hello` in the chapter before,
you should be able to list all the packages under `github.com/golang/example` by running:

```bash
$ go list github.com/golang/example/...    # remember that ... means "and everything below"
github.com/golang/example/appengine-hello
github.com/golang/example/gotypes
github.com/golang/example/gotypes/defsuses
github.com/golang/example/gotypes/doc
github.com/golang/example/gotypes/hello
github.com/golang/example/gotypes/hugeparam
github.com/golang/example/gotypes/implements
github.com/golang/example/gotypes/lookup
github.com/golang/example/gotypes/nilfunc
github.com/golang/example/gotypes/pkginfo
github.com/golang/example/gotypes/skeleton
github.com/golang/example/gotypes/typeandvalue
github.com/golang/example/hello
github.com/golang/example/outyet
github.com/golang/example/stringutil
github.com/golang/example/template
```

This can be pretty useful, but there's much more you can do!

To list all the information related to `"github.com/golang/example/hello"` you could use the `-f` flag, which by default has the value `"{{.ImportPath}}"`.

```bash
$ go list -f '{{.}}' github.com/golang/example/hello
{/Users/campoy/src/github.com/golang/example/hello github.com/golang/example/hello  main  /Users/campoy/bin/hello  false false true build ID mismatch /Users/campoy  false [hello.go] [] [] [] [] [] [] [] [] [] [] [] [] [] []
[] [] [] [fmt github.com/golang/example/stringutil] [errors fmt github.com/golang/example/stringutil internal/cpu internal/poll internal/race io math os reflect runtime runtime/internal/atomic runtime/internal/sys strconv sy
nc sync/atomic syscall time unicode/utf8 unsafe] false <nil> [] [] [] [] []}
```

By using the template language documented in the
[`"text/template"`](https://golang.org/pkg/text/template) package, you can extract
specific pieces of information.

For instance, you could list all the import statements in `"github.com/golang/example/hello"`:

```bash
$ go list -f '{{.Imports}}' github.com/golang/example/hello
[fmt github.com/golang/example/stringutil]
```

Once you learn a bit more of the Go templating language you'll be able to make your output look better:

```bash
$ go list -f '{{join .Imports "\n"}}' github.com/golang/example/hello
fmt
github.com/golang/example/stringutil
```

Alternatively, you can also get all of the information linked to package in JSON,
by adding the `-json` flag.

```bash
$ go list -json github.com/golang/example/hello
{
	"Dir": "/Users/campoy/src/github.com/golang/example/hello",
	"ImportPath": "github.com/golang/example/hello",
	"Name": "main",
	"Target": "/Users/campoy/bin/hello",
	"Stale": true,
	"StaleReason": "build ID mismatch",
	"Root": "/Users/campoy",
	"GoFiles": [
		"hello.go"
	],
	"Imports": [
		"fmt",
		"github.com/golang/example/stringutil"
	],
	"Deps": [
		"errors",
		"fmt",
		"github.com/golang/example/stringutil",
		"internal/cpu",
		"internal/poll",
		"internal/race",
		"io",
		"math",
		"os",
		"reflect",
		"runtime",
		"runtime/internal/atomic",
		"runtime/internal/sys",
		"strconv",
		"sync",
		"sync/atomic",
		"syscall",
		"time",
		"unicode/utf8",
		"unsafe"
	]
}
```

Note that if you use the `-json` flag the `-f` flag will be ignored.

### Exercise with `go list`

Can you count how many `.go` files are under `github.com/golang/example` using `go list`?

To see one possible solution, expand the `Details` below.

<details>

List all of the files with `go list`:

```bash
$ go list -f '{{join .GoFiles "\n"}}' github.com/golang/example/...
app.go
weave.go
main.go
main.go
hello.go
main.go
main.go
lookup.go
main.go
main.go
main.go
main.go
hello.go
main.go
reverse.go
main.go
```

If you're in Linux or Mac you can add `wc -l` to get a number:

```bash
$ go list -f '{{join .GoFiles "\n"}}' github.com/golang/example/... | wc -l
16
```
</details>

## Multiple paths in `GOPATH`

In a similar way to how the `PATH` environment variable can contain multiple directories, which are used to find binaries
in your system, `GOPATH` can contain multiple directories which will be used consecutively to find the packages you import.

The `go` tool will try to find a package in the first component of `GOPATH`, and fallback to the second one only if the
package wasn't found yet. Whenever you run `go get` the source code will be stored under the first `GOPATH` component.

This setup is not used often, but it is important to know it exists. We will not cover it in any detail during the workshop.

## Special directories: vendor/internal

There are two directory names that the `go` tool interprets in special ways: `vendor` and `internal`.
Let's learn why they exist and when to use them.

### The `vendor` directory

The `vendor` directory is used to keep third party dependencies to your project.

Whenever an import path `"example.com/foo"` is found in Go code, the `go` tool will navigate the directory
tree towards the root, looking for a directory containing a `vendor` directory.

When a `vendor` directory is found, the `go` tool will look into `vendor/example.com/foo`, if that directory
exists it will be used as the package imported with the path `"example.com/foo"`.
If it doesn't exist the `go` tool will continue upwards looking for other `vendor` directories, and eventually
in `GOPATH`.

This means that given a folder structure such as:

```
GOPATH
\- src
  \- github.com
    \- campoy
      |- go-tooling-workshop
      |  |- foo
      |  | \- main.go
      |  \- vendor
      |    \- github.com
      |      \- campoy
      |        \- bar
      |          \- bar.go
      \- bar
        \- bar.go
```

An import path of `github.com/campoy/bar` inside of `main.go` will be resolved to the package inside of `vendor`,
but the same import statement in a file outside of `go-tooling-workshop` would resolved to `$GOPATH/src/github.com/campoy/bar`.

#### Exercise: redefine Pi

Haven't you ever wished that `Pi` was more rational? Something easier to remember? Maybe ... `3`?

Create a new directory under your `GOPATH` containing the `main.go` file below.

[embedmd]:# (pi/main.go /package main/ $)
```go
package main

import (
	"fmt"
	"math"
)

func main() {
	fmt.Println("Pi is", math.Pi)
}
```

Where would you add a `vendor` directory that redefined `math.Pi` to be `3`?

<details>

You could create a `math` directory in many points:

- `$GOPATH/src/vendor/math`
- `$GOPATH/src/github.com/vendor/math`
- `$GOPATH/src/github.com/campoy/vendor/math`
- `$GOPATH/src/github.com/campoy/go-tooling-workshop/vendor/math`
- `$GOPATH/src/github.com/campoy/go-tooling-workshop/1-source-code/1-workspace/vendor/math`
- `$GOPATH/src/github.com/campoy/go-tooling-workshop/1-source-code/1-workspace/pi/vendor/math`

Where do you think it's best to do this?

</details>

### The `internal` directory

Go visibility rules for types are quite straight forward:

- Identifiers that start with a capital letter are visible to everyone.
- All other identifiers are visible only from the package that defines them.

With this two rules it is impossible to provide an identifier that is visible from only a subset
of other packages, that is why the `go` tool supports `internal` directories.

The packages inside of an `internal` directory are only accessible to those packages that are
siblings of that `internal` package or are contained by one of those siblings.

For instance, given this directory structure:

```
GOPATH
\- src
  \- github.com
    \- campoy
      |- go-tooling-workshop
      |  |- foo
      |  | \- main.go
      |  \- internal
      |    \- bar
      |       \- bar.go
      \- other
        \- main.go
```

We could import `github.com/campoy/go-tooling-workshop/internal/bar` from the `main.go` in package `foo`,
but not from the one in package `other`.

#### Exercise: find internal packages in the standard library

Go's standard library uses `internal` directories in order to organize and share code in a more organized
way inside of the standard library without adding API surface.

Can you find any of those packages? Where are they used from?

<details>

There's many of these directories, one of them is `image/internal/imageutil`.

```
package imageutil
    import "image/internal/imageutil"

    Package imageutil contains code shared by image-related packages.

FUNCTIONS

func DrawYCbCr(dst *image.RGBA, r image.Rectangle, src *image.YCbCr, sp image.Point) (ok bool)
    DrawYCbCr draws the YCbCr source image on the RGBA destination image
    with r.Min in dst aligned with sp in src. It reports whether the draw
    was successful. If it returns false, no dst pixels were changed.

    This function assumes that r is entirely within dst's bounds and the
    translation of r from dst coordinate space to src coordinate space is
    entirely within src's bounds.
```

Using `go list` and `grep` we can find where that package is being used from:

```bash
$ go list -f '{{.ImportPath}}: {{.Imports}}' image/... | grep "internal"
image/draw: [image image/color image/internal/imageutil]
image/internal/imageutil: [image]
image/jpeg: [bufio errors image image/color image/internal/imageutil io]
```

</details>

## Congratulations

You're now an expert in Go workspaces! You know how to organize your workspace under `GOPATH`,
how to use `go list` to navigate dependencies, and how `vendor` and `internal` can be used
to organize your projects.

Let's see now what tools we can use to manage the third party dependencies in `vendor` in the [next lesson](3-deps.md).
