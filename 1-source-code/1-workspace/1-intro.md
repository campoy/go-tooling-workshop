# Workspace Introduction

Go has a particular way of structuring your workspace. This might be weird at first,
but I personally adopted this style for all my projects now (no matter the language)
I couldn't be happier.

## Understanding GOPATH

All the components of a Go workspace are included under a single directory
(or more, but for now let's say it's only one) that we call `GOPATH`.

There's an environment variable that allows you to set what workspace root
should be in your environment. But if it's not defined, you will get a default
value that you can check running

```bash
$ go env GOPATH
```

Under this directory you will find three subdirectories:

- `src`: this includes *all* of your Go source code, not only your project but also all its dependencies.
- `bin`: whenever you compile a Go program with `go install` or `go get` (we'll see them later) the resulting binary will be stored here.
- `pkg`: this is simply a cache for compiled packages, you really shouldn't care about this.

## Import paths, directories, and repositories

Now that we understand how Go code is stored, let's see how you reference a given
package from Go code.

This is done with an `import` statement, which imports an `import path` that is
simply a `string` according to the language specification.

[embedmd]:# (../../vendor/github.com/golang/example/hello/hello.go /package main/ /\)/)
```go
package main

import (
	"fmt"

	"github.com/golang/example/stringutil"
)
```

As you can see in the code above ([hello.go](../../vendor/github.com/golang/example/hello/hello.go)),
import paths take two forms.

The Go tool will first try to find a match to the import path in the standard library and your `GOPATH`.
If it doesn't find the package, it considers the path to be a URL pointing to a repository where
the code can be found.

The `go` tool is able to fetch code from most repositories, as long as they support one of the following:

- `git` - Git, download at http://git-scm.com/downloads
- `svn` - Subversion, download at: http://subversion.apache.org/packages.html
- `hg` - Mercurial, download at https://www.mercurial-scm.org/downloads
- `bzr` - Bazaar, download at http://wiki.bazaar.canonical.com/Download

For example, git is used for Github, hg is used for Bitbucket, etc.

Note that `"github.com/golang/example/stringutil"` is not technically pointing to any code on GitHub,
as the actual URL to the code would also include `https` and the branch (probably `master`).
The `go` tool is able to figure these differences out, and it will store the code under `$GOPATH/src/{import_path}`.


### Exercise: `go get`

Check what is current value for your `GOPATH` with `go env GOPATH`. If you'd rather
store your Go code somewhere else, simply declare an environment variable named
`GOPATH` to override this value.

Now run this command:

```bash
$ go get github.com/golang/example/hello
```

This will create the three directories mentioned previously, explore their contents.
- Can you find where the source code was stored?
- Can you find where the `hello` binary was stored? Try running it.

## Advanced options for `go get`

On top of the basic `go get` usage that we saw above, there's some flags that might interest you.

- `go get -d`: download the code, but do not compile anything.
- `go get -u`: even if the code is already stored in `GOPATH`, download the latest version.
- `go get -v`: enable verbose mode.

If you run `go get` with an import path that doesn't contain any Go code, you will see an error:

```bash
$ go get github.com/campoy/go-web-workshop
can't load package: package github.com/campoy/go-web-workshop: no buildable Go source files in /Users/campoy/src/github.com/campoy/go-web-workshop
```

In this case it is because the Go code is in directories under that root. You can request to
install all packages under that import path by adding three trailing dots `...`.

```bash
$ go get github.com/campoy/go-web-workshop/...
```

This works with all Go tools, it might be handy for testing too.

## Congratulations

You now know how to organize a Go workspace, download dependencies, and find the resulting binaries.

Let's learn more about Go workspace and the tooling that associated to it in the [next lesson](2-tooling.md).
