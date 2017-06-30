[![Go Report Card](https://goreportcard.com/badge/github.com/campoy/go-tooling-workshop)](https://goreportcard.com/report/github.com/campoy/go-tooling-workshop)

# Go Tooling in Action

Hi, and welcome to "Go Tooling in Action". This is a multi hour workshop,
so get ready for some intense learning!

After going through the whole content, you will know about what tools can
help you better write Go code, how to build artifacts from that code, and
how to understand the performance of your code once it's running.

You should be relatively familiar with Go, even though we won't be writing
much code. Maybe it's time to check out the [Go Tour](https://tour.golang.org).

For a shorter and sweeter version of this workshop, you can watch the video
that I made for a conference, and inspired me to create this workshop.

<div style="text-align:center">
    <a href="https://www.youtube.com/watch?v=uBjoTxosSys">
        <img src="https://img.youtube.com/vi/uBjoTxosSys/0.jpg" alt="Go Tooling in Action">
        <p>Go Tooling in Action</p>
    </a>
</div>

## Software requirements

To go through this you will need the following:

1. You have installed the [Go Programming Language](https://golang.org).
1. We will be using [Visual Studio Code](https://code.visualstudio.com/) it's free and open source.
1. At some point we'll also use [delve](https://github.com/derekparker/delve/tree/master/Documentation/installation), so install it now if you think you won't have good WiFi later.
1. Finally, we'll need [GraphViz](http://www.graphviz.org/Download..php).

The rest of the software we'll use is quick to install through `go get`, so
do not worry yet.

## Contents

The workshop is for now composed of three independent sections:

- [1: Source code management](1-source-code/README.md)
- [2: Building artifacts from code](2-building-artifacts/README.md)
- [3: Dynamic program analysis](3-dynamic-analysis/README.md)

In the future one more section might be added regarding monitoring of running
systems, but for now that topic is out of the scope of this workshop.

## Issues

This workshop is very new, so some things might be missing or wrong.

If you find anything that seems broken, please file an issue. Or even better,
send a pull request! You will need to sign a CLA, you'll get the info once
you send the PR.

## Resources

These are places where you can find more information for Go:

- [golang.org](https://golang.org)
- [godoc.org](https://godoc.org), where you can find the documentation for any package.
- [The Go Programming Language Blog](https://blog.golang.org)

My favorite aspect of Go is its community, and you are now part of it too. Welcome!

As a newcomer to the Go community you might have questions or get blocked at some point.
This is completely normal, and we're here to help you.
Some of the places where gophers tend to hang out are:

- [The Go Forum](https://forum.golangbridge.org/)
- #go-nuts IRC channel at [Freenode](https://freenode.net/)
- Gophers’ community on [Slack](https://gophers.slack.com/messages/general/) (signup [here](https://invite.slack.golangbridge.org/) for an account).
- [@golang](https://twitter.com/golang) and [#golang](https://twitter.com/search?q=%23golang) on Twitter.
- [Go+ community](https://plus.google.com/u/1/communities/114112804251407510571) on Google Plus.
- [Go user meetups](https://go-meetups.appspot.com/)
- golang-nuts [mailing list](https://groups.google.com/forum/?fromgroups#!forum/golang-nuts)
- Go community [Wiki](https://github.com/golang/go/wiki)

### Disclaimer

This is not an official Google product (experimental or otherwise), it is just
code that happens to be owned by Google.
