# Performance profiling

So far we've seen how to verify the correctness of our programs, aka as our programs do what they're supposed to.
To do so we've used `golint`, `go vet`, and `go test`.

Correctness is necessary, but not sufficient. We also need our programs to
be fast. And how do we measure that?

Let's work again with the web server program we used before
for debugging. There is a slightly modified version of the code in the
[webserver](webserver) directory.

[embedmd]:# (webserver/main.go /package main/ $)
```go
package main

import (
	"fmt"
	"log"
	"net/http"
	"regexp"
)

func main() {
	http.HandleFunc("/", handler)
	log.Printf("listening on localhost:8080")
	log.Fatal(http.ListenAndServe("localhost:8080", nil))
}

func handler(w http.ResponseWriter, r *http.Request) {
	if name, ok := isGopher(r.URL.Path[1:]); ok {
		fmt.Fprintf(w, "hello, %s", name)
		return
	}
	fmt.Fprintln(w, "hello, stranger")
}

func isGopher(email string) (string, bool) {
	re := regexp.MustCompile("^([[:alpha:]]+)@golang.org$")
	match := re.FindStringSubmatch(email)
	if len(match) == 2 {
		return match[1], true
	}
	return "", false
}
```

## Basic performance analysis

A pretty straight forward way to measure the performance of any job
is to call it and track the time it takes to finish. It is not a perfect
measure, but it is a useful one to start with.

In the case of a program that runs for a limited time (unlike servers, which
run indefintely) I recommend using the `time` command.

```bash
$ go get github.com/golang/example/hello
$ time $GOPATH/bin/hello
Hello, Go examples!
$GOPATH/bin/hello  0.00s user 0.00s system 15% cpu 0.016 total
```

Unfortunately, this doesn't work for web servers, where we
would need to track how long it takes to answer a single request.

There's many tools that provide this functionality, such as
[Apache Bench](https://httpd.apache.org/docs/2.4/programs/ab.html). In this
workshop we will use [go-wrk](https://github.com/tsliwowicz/go-wrk), which can be installed with `go get`:

```bash
$ go get github.com/tsliwowicz/go-wrk
```

Start running the `webserver` binary, then run `go-wrk` to get an analysis of
the server's performance.

```bash
$ cd 3-dynamic-analysis/3-profiling/webserver
$ go run main.go
```

In a different terminal:

```bash
$ go-wrk -d 5 http://localhost:8080
Running 5s test @ http://localhost:8080
  10 goroutine(s) running concurrently
164677 requests in 4.914899534s, 18.22MB read
Requests/sec:           33505.67
Transfer/sec:           3.71MB
Avg Req Time:           298.456µs
Fastest Request:        78.387µs
Slowest Request:        10.08497ms
Number of Errors:       0
```

Great! We now know that we're able to handle around 33505 requests
per second and that the average time per request is 298µs.

But it seems the slowest request took over 10ms!

## Benchmarks

With the previous exercise we were able to put a number on how fast, or
slow, our server performs. But this doesn't necessarily helps us figuring
out what piece of code we should change. We need a deeper analysis.

How can we measure the performance of specific pieces of code? It's almost
like unit testing but for performance analysis. Well, those are benchmarks.

In Go, benchmarks are very similar to unit tests. Simply run a function
named `BenchmarkXXX`, where `XXX` refers to what you're benchmarking,
with type `func (*testing.B)`.

How fast is the function `strings.Fields`? We can write
benchmark for that:

[embedmd]:# (strings_test.go /func Benchmark/ /^}/)
```go
func BenchmarkFields(b *testing.B) {
	for i := 0; i < b.N; i++ {
		strings.Fields("hello, dear friend")
	}
}
```

You can then run this benchmark with `go test`, providing the `-bench`
flag and a regular expression so only the tests whose name match will
be executed. To execute all of them, you can use `-bench=.`.

Note that we repeat the operation we want to benchmark (in this case
`strings.Field`) `b.N` times. This value will be changed by `go test`
until the resulting times are statistically significant.

```bash
$ go test -bench=.
goos: darwin
goarch: amd64
pkg: github.com/campoy/go-tooling-workshop/3-dynamic-analysis/3-profiling
BenchmarkFields-8       20000000                95.7 ns/op
PASS
ok      github.com/campoy/go-tooling-workshop/3-dynamic-analysis/3-profiling    2.030s
```

The number of memory allocations is also a good indicator of how well a function
might perform at scale, and you can obtain it by adding `-benchmem`.

```bash
$ go test -bench=. -benchmem
goos: darwin
goarch: amd64
pkg: github.com/campoy/go-tooling-workshop/3-dynamic-analysis/3-profiling
BenchmarkFields-8               20000000               103 ns/op              48 B/op          1 allocs/op
PASS
```

### A note on compiler optimizations

Sometimes your benchmarks might be way faster that you would expect. When
that happens do not be too happy, it might be caused by a compiler optimization
that simply removed the function call you thought you were measuring.


In order to avoid that, it is a good idea to store the result of the function call
in a package variable. This will force the compiler to keep the call.

[embedmd]:# (strings_test.go /var res/ $)
```go
var res []string

func BenchmarkFields_Escape(b *testing.B) {
	for i := 0; i < b.N; i++ {
		res = strings.Fields("hello, dear friend")
	}
}
```

### Exercise: benchmarks

Write a benchmark for the `isGopher` function in [webserver/main.go](webserver/main.go).

Do not change the implementation of `isGopher` just yet.
How fast is it? How many allocations per operation does it require?

### Exercise: benchmarks and alternative implementations

Write a benchmark for the [sum](../2-testing/sum) package we tested before.
Try writing an alternative implementation of `Sum` that uses iteration rather
than recursion. Which one is faster? By how much?

Is there a big difference in performance between `for _, v := range vs` and
`for i := 0; i < len(vs); i++`?  Why do you think that is?

## Congratulations

You just learned how to find out how well a server performs, and how to measure
specific parts of your code via benchmarks.

This is good, but it still doesn't help us understand which piece of code is the
one that we should optimize first! To do so, we'll learn how to use pprof during
the [next chapter](2-pprof.md).
