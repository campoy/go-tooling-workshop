# Understanding performance with pprof

In order to better understand how a program performs we can use a similar
approach to the one we used with code coverage, back in the
[testing chapter](../2-testing/2-code-coverage.md).

The idea here is we will periodicallly check what a program is running,
exactly what line of code it's executing at that given time. By doing this
often enough and for a long enough period of time, we can figure out what lines
of code we spend the longest time executing!

This is, again, not a perfect tool as Filippo Valsorda explains during
[this talk](https://speakerdeck.com/filosottile/you-latency-and-profiling-at-gophercon-india-2017) but it is a very good indication for many cases. So let's learn how to use
it with Go.

## Generating pprof profiles from benchmarks

The easiest way to use pprof is by adding the `-cpuprofile` flag to any benchmark.

Let's add that flag to the same benchmark for `strings.Fields` we saw before.

```bash
$ go test -bench=. -cpuprofile=cpu.pprof
goos: darwin
goarch: amd64
pkg: github.com/campoy/go-tooling-workshop/3-dynamic-analysis/3-profiling
BenchmarkFields-8               20000000               118 ns/op
BenchmarkFields_Escape-8        20000000               112 ns/op
PASS
ok      github.com/campoy/go-tooling-workshop/3-dynamic-analysis/3-profiling    5.051s
```

In addition to the benchmark output, two new files were created.
- cpu.pprof:  It is a binary file, so no need to try to read it manually. Instead, we're going to use the `pprof` tool to read it.
- x.test: where x corresponds to the name of the current directory. It's the test binary that runs the benchmarks.

## Reading pprof profiles with `go tool pprof`

Once we have the test binary and the pprof profile, we can use the `pprof`
tool under `go tool` to analyze the profile:

```bash
$ go tool pprof 3-profiling.test cpu.pprof
File: 3-profiling.test
Type: cpu
Time: Jun 29, 2017 at 3:30pm (EDT)
Duration: 5.03s, Total samples = 4.31s (85.70%)
Entering interactive mode (type "help" for commands, "o" for options)
(pprof)
```

This is an interactive console, the most important commands in my opinion are:

- top: shows the lines of code that were sampled the most often.

```bash
(pprof) top
Showing nodes accounting for 4.14s, 96.06% of 4.31s total
Dropped 25 nodes (cum <= 0.02s)
Showing top 10 nodes out of 44
      flat  flat%   sum%        cum   cum%
     1.70s 39.44% 39.44%      1.70s 39.44%  runtime.kevent /Users/campoy/src/golang.org/x/go/src/runtime/sys_darwin_amd64.s
     1.01s 23.43% 62.88%         4s 92.81%  strings.Fields /Users/campoy/src/golang.org/x/go/src/strings/strings.go
     0.57s 13.23% 76.10%      0.57s 13.23%  runtime.heapBitsSetType /Users/campoy/src/golang.org/x/go/src/runtime/mbitmap.go
     0.32s  7.42% 83.53%      2.84s 65.89%  runtime.mallocgc /Users/campoy/src/golang.org/x/go/src/runtime/malloc.go
     0.18s  4.18% 87.70%      0.18s  4.18%  runtime.nextFreeFast /Users/campoy/src/golang.org/x/go/src/runtime/malloc.go
     0.10s  2.32% 90.02%      0.10s  2.32%  runtime.mach_semaphore_signal /Users/campoy/src/golang.org/x/go/src/runtime/sys_darwin_amd64.s
     0.09s  2.09% 92.11%      2.99s 69.37%  runtime.makeslice /Users/campoy/src/golang.org/x/go/src/runtime/slice.go
     0.09s  2.09% 94.20%      0.09s  2.09%  runtime.memclrNoHeapPointers /Users/campoy/src/golang.org/x/go/src/runtime/memclr_amd64.s
     0.04s  0.93% 95.13%      0.04s  0.93%  runtime.mach_semaphore_timedwait /Users/campoy/src/golang.org/x/go/src/runtime/sys_darwin_amd64.s
     0.04s  0.93% 96.06%      0.04s  0.93%  runtime.mallocgc /Users/campoy/src/golang.org/x/go/src/runtime/runtime1.go
```

- list strings.Fields: shows how much time we spent per line of code in `strings.Fields`.

```bash
(pprof) list strings.Fields
Total: 4.31s
ROUTINE ======================== strings.Fields in /Users/campoy/src/golang.org/x/go/src/strings/strings.go
     1.01s         4s (flat, cum) 92.81% of Total
         .          .    310:var asciiSpace = [256]uint8{'\t': 1, '\n': 1, '\v': 1, '\f': 1, '\r': 1, ' ': 1}
         .          .    311:
         .          .    312:// Fields splits the string s around each instance of one or more consecutive white space
         .          .    313:// characters, as defined by unicode.IsSpace, returning an array of substrings of s or an
         .          .    314:// empty list if s contains only white space.
      20ms       20ms    315:func Fields(s string) []string {
         .          .    316:   // First count the fields.
         .          .    317:   // This is an exact count if s is ASCII, otherwise it is an approximation.
         .          .    318:   n := 0
         .          .    319:   wasSpace := 1
         .          .    320:   // setBits is used to track which bits are set in the bytes of s.
         .          .    321:   setBits := uint8(0)
      50ms       50ms    322:   for i := 0; i < len(s); i++ {
      30ms       30ms    323:           r := s[i]
      60ms       60ms    324:           setBits |= r
     180ms      180ms    325:           isSpace := int(asciiSpace[r])
     160ms      160ms    326:           n += wasSpace & ^isSpace
         .          .    327:           wasSpace = isSpace
...
```

- web: show a graph with all the functions called and the time we spent on each.

```bash
(pprof) web
Failed to execute dot. Is Graphviz installed? Error: exec: "dot": executable file not found in $PATH
```

Oh, it seems we're missing something! Make sure you installed [Graphviz](http://www.graphviz.org/) first.

```bash
$ go tool pprof 3-profiling.test cpu.pprof
File: 3-profiling.test
Type: cpu
Time: Jun 29, 2017 at 3:30pm (EDT)
Duration: 5.03s, Total samples = 4.31s (85.70%)
Entering interactive mode (type "help" for commands, "o" for options)
(pprof) web
```

This will generate an SVG file and open it. Make sure your computer is set up
so your `.svg` files are opened with something that works.

![pprof web](web.png)

We can see that most of the time is spent creating a slice `runtime.makeslice`.

### Exercise: visualizing the result of pprof on `isGopher`

Using the same benchmark you wrote on the previous chapter, run it again with
`-cpuprofile` in order to generate a CPU profile and visualize it with pprof.

Do not modify the code of `isGopher` just yet, simply analyze it and try to understand
what we could improve.

<details>

Admitting that you wrote a benchmark somehow looking like this:

```go
package main

import "testing"

var name string
var ok bool

func BenchmarkIsGopher(b *testing.B) {
	for i := 0; i < b.N; i++ {
		name, ok = isGopher("campoy@golang.org")
	}
}
```

You should be able to run it with:

```bash
$ go test -bench=. -cpuprofile pprof.cpu
goos: darwin
goarch: amd64
pkg: github.com/campoy/go-tooling-workshop/3-dynamic-analysis/3-profiling/webserver
BenchmarkIsGopher-8       200000              9044 ns/op
PASS
ok      github.com/campoy/go-tooling-workshop/3-dynamic-analysis/3-profiling/webserver  2.107s
```

Finally, using `go tool pprof` you can visualize the analysis as a graph.

```bash
$ go tool pprof webserver.test pprof.cpu
File: webserver.test
Type: cpu
Time: Jun 29, 2017 at 4:09pm (EDT)
Duration: 2.07s, Total samples = 2.23s (107.72%)
Entering interactive mode (type "help" for commands, "o" for options)
(pprof) web
```

You should see a graph that starts with something similar to this graph.

![webserver pprof](webserver-pprof.png)

</details>

## Generating profiles from running HTTP servers

Benchmarks are great, but you need to write them. Wouldn't it be cool to obtain CPU
profiles directly from a running web server? That's exactly what the
[net/http/pprof](https://golang.org/pkg/net/http/pprof) provides.

Simply add an import statement of that package. You will need to import it with an
underscore (`_`) as the package name, since we're importing for the side effects of
the initialization rather than to refer to any of the identifiers defined in the package.

```go
import _ "net/http/pprof"
```

Add it to the webserver we were working with previously, so we can analyze the whole endpoint
rather than just `isGopher`. Who knows, maybe there's something we can optimize there too!

Once you've added that import statement, recompile the program and run it again.
Now you can visit http://localhost:8080/debug/pprof/ in order to have a snapshot
of the current activity of the server. Explore this dashboard.

This is how most processes work at Google! They expose a web server with an
endpoint that exposes metrics we can create dashboards on. Pretty simple, but
powerful.

Click on the `full goroutine stack dump` link to see what every goroutine is
executing at this point. This endpoint is the one we're going to use to let
`go tool pprof` analyze the performance of a running HTTP server.

```bash
$ go tool pprof -seconds 5 http://localhost:8080/debug/pprof/profile
Fetching profile over HTTP from http://localhost:8080/debug/pprof/profile?seconds=5
Please wait... (5s)
Saved profile in /Users/campoy/pprof/pprof.samples.cpu.001.pb.gz
Type: cpu
Time: Jun 29, 2017 at 9:36pm (EDT)
Duration: 5s, Total samples = 0
Entering interactive mode (type "help" for commands, "o" for options)
(pprof)
```

This command just analyzed 5 seconds of the activity by the HTTP server behind
that endpoint, and now it's ready to expose the data.
Let's try running `top` to have an idea of the content of the profile.

```
(pprof) top
Showing nodes accounting for 0, 0% of 0 total
      flat  flat%   sum%        cum   cum%
```

0% is not a huge percent ... what hapenned? Well, during the time we were analyzing
the server there was no traffic, so nothing got recorded.

Let's try again, this time sending some traffic with `go-wrk` first, then starting
`go tool pprof`.

In a terminal run:

```bash
$ go-wrk -d 30 http://localhost:8080
```

This will take around a minute to run, so in the meanwhile we can run
`go tool pprof`.

```bash
$ go tool pprof -seconds 5 http://localhost:8080/debug/pprof/profile
Fetching profile over HTTP from http://localhost:8080/debug/pprof/profile?seconds=5
Please wait... (5s)
Saved profile in /Users/campoy/pprof/pprof.samples.cpu.001.pb.gz
Type: cpu
Time: Jun 29, 2017 at 9:36pm (EDT)
Duration: 5s, Total samples = 0
Entering interactive mode (type "help" for commands, "o" for options)
(pprof) top
13.92s of 15.30s total (90.98%)
Dropped 185 nodes (cum <= 0.08s)
Showing top 10 nodes out of 100 (cum >= 0.08s)
      flat  flat%   sum%        cum   cum%
    10.79s 70.52% 70.52%     10.80s 70.59%  syscall.Syscall
     0.71s  4.64% 75.16%      0.71s  4.64%  runtime.usleep
     0.64s  4.18% 79.35%      0.64s  4.18%  runtime.mach_semaphore_wait
     0.63s  4.12% 83.46%      0.63s  4.12%  runtime.kevent
     0.46s  3.01% 86.47%      0.46s  3.01%  runtime.mach_semaphore_signal
     0.28s  1.83% 88.30%      0.29s  1.90%  runtime.freedefer
     0.16s  1.05% 89.35%      0.16s  1.05%  nanotime
     0.09s  0.59% 89.93%      0.42s  2.75%  runtime.mallocgc
     0.08s  0.52% 90.46%      0.08s  0.52%  runtime.duffcopy
     0.08s  0.52% 90.98%      0.08s  0.52%  runtime.heapBitsSetType
```

Cool, this time the profile got some data. Not easy to understand though, so let's
`web`.

![webserver pprof from live traffic](webserver-live-pprof.png)

Ok ... do you understand anything? I know I don't!
Maybe there's a better way to visualize performance profiles from pprof.

Good news, there is, and that's what we're going to learn next!

## Congratulations

You're able to obtain CPU profiles from benchmarks and from running servers.
You even know how to generate cool graphics showing how much time we spend on
each function call ... unfortunately those graphics are not easy to parse.

Let's learn a different way to display that data on the
[next chapter](3-flame-graphs.md).
