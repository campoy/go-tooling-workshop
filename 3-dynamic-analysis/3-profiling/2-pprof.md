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
