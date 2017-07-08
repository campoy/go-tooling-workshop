# Testing Go programs

Writing good unit tests is very important not only to show that your
code works while you write it, but also to ensure that future updates
to your code do not break existing functionality.

There's a lot of bibliography on what is a good test, or when tests
should be written, so I will not cover any of those here. Instead we will
simply see how tests are written in Go.

## Running tests: `go test`

In order to run all the tests in an existing package we can use the `test`
subcommand in the `go` tool. For instance, to run all the tests in the
`strings` package of the standard library we can run:

```bash
$ go test strings
ok      strings 0.290s
```

You can add the `-v` flag in order to list all of the tests that are run,
rather  than only those that failed.

Let's now see how to write new tests.

## Writing tests: `"testing"` and the `_test` suffix

Unit tests are functions that match the following conditions:

- they are defined in a file with name ending in `_test.go`
- the function name starts with `Test`
- the function is of type `func(t *testing.T)`

### The `_test` suffix

In a [previous chapter](../../2-building-artifacts/1-go-build.md) we
discussed file name suffixes, and how a file named `foo_windows.go` will
only be compiled for Windows systems.

The same idea applies to the `_test` suffix: any file containing that
suffix will be ignored unless we're running `go test`.

### The `testing` package and `testing.T`

The value received by a test, of type `*testing.T` provides a series of
methods that allows us to flag when a test is failing:

- Use `Error` and `Errorf` to indicate the test failed and continue executing it
- Use `Fatal` and `Fatalf` to indicate the test failed and should not continue

### What package contains the tests?

You have two options depending on what kind of tests you want
to write for a given package `sum` in the directory [sum](sum).

Given a file `sum.go`:
[embedmd]:# (sum/sum.go /package sum/ $)
```go
package sum

// All returns the sum of the given values.
func All(vs ...int) int {
	return recursive(vs)
}

func recursive(vs []int) int {
	if len(vs) == 0 {
		return 0
	}
	return vs[0] + recursive(vs[1:])
}
```

You can define a unit test for `recursive` by creating a new
function in test file that is part of the same package `sum`.

[embedmd]:# (sum/sum_internal_test.go /package sum/ $)
```go
package sum

import "testing"

func TestRecursive(t *testing.T) {
	// Implement the body of this test, calling recursive.
}
```

This is possible because the test belongs to the same package,
so even though `recursive` is not exported it is still
visible to other members of the same package.

Another option, if you want to test only exported elements,
is to define a `sum_test` package:

[embedmd]:# (sum/sum_test.go /package sum_test/ $)
```go
package sum_test

import "testing"

func TestAll(t *testing.T) {
	// Implement the body of this test, calling sum.All.
}
```

As you can see, on this case we need to import the `sum`
package, as you would do if you were a user of the package.

### Exercise: write some tests for the `sum` package

Edit [sum_internal_test.go](sum/sum_internal_test.go) and
[sum_test.go](sum/sum_test.go) in order to make sure the
implementation of `All` and `recursive` is correct.

Run `go test` once you've written them.

## Table driven tests and subtests

Very often you will find yourself writing repetitive tests, let's imagine that you
were writing a test for `strings.Index`, you might end up writing something like:

[embedmd]:# (strings_test.go /TestIndex_Repeat/ /^}/)
```go
TestIndex_Repeat(t *testing.T) {
	if i := strings.Index("hello, world", "hello"); i != 0 {
		t.Errorf(`"hello, world" should contain "hello" at position 0, not %d`, i)
	}
	if i := strings.Index("hello, world", "bye"); i != -1 {
		t.Errorf(`"hello, world" should not contain "hello"`)
	}
	if i := strings.Index("hello, world", "world"); i != 7 {
		t.Errorf(`"hello, world" should contain "hello" at position 7, not %d`, i)
	}
}
```

There's a lot of repeated code, and repeated code means errors can be introduced while
copy pasting very easily.

In Go, it is recommended to use what we call a table drive test. Rather than repeating
the code, we keep all the commonalities across all the repetitions and extract the
differences into a slice of test cases. The previous example would be something like:

[embedmd]:# (strings_test.go /TestIndex_Table/ /^}/)
```go
TestIndex_Table(t *testing.T) {
	tt := []struct {
		text string
		sub  string
		idx  int
	}{
		{"hello, world", "hello", 0},
		{"hello, world", "bye", -1},
		{"hello, world", "world", 7},
	}
	for _, tc := range tt {
		if idx := strings.Index(tc.text, tc.sub); idx != tc.idx {
			if tc.idx >= 0 {
				t.Errorf("%s should contain %s at position %d, not %d", tc.text, tc.sub, tc.idx, idx)
			} else {
				t.Errorf("%s should not contain %s", tc.text, tc.sub)
			}
		}
	}
}
```

Even better, since Go 1.7, we can use subtests, which allow you to have a better control over what
test cases inside of a test you run. You need to give a name to your test case, which is a good
practice anyway, and call `t.Run`.

[embedmd]:# (strings_test.go /TestIndex_Subtest/ /^}/)
```go
TestIndex_Subtest(t *testing.T) {
	tt := []struct {
		name string
		text string
		sub  string
		idx  int
	}{
		{"first character", "hello, world", "hello", 0},
		{"not found", "hello, world", "bye", -1},
		{"last character", "hello, world", "world", 7},
	}
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			if idx := strings.Index(tc.text, tc.sub); idx != tc.idx {
				if tc.idx >= 0 {
					t.Fatalf("%s should contain %s at position %d, not %d", tc.text, tc.sub, tc.idx, idx)
				}
				t.Fatalf("%s should not contain %s", tc.text, tc.sub)
			}
		})
	}
}
```

_Note_: you can use `Fatalf` with subtests and only the current subtest will fail, unlike in the
previous loops where later tests would be ignored.

If you run `go test -v` you'll see the subtest that were ran.

```bash
$ go test -v
--- PASS: TestIndex_Subtest (0.00s)
    --- PASS: TestIndex_Subtest/first_character (0.00s)
    --- PASS: TestIndex_Subtest/not_found (0.00s)
    --- PASS: TestIndex_Subtest/last_character (0.00s)
```

### Exercise: subtests

Make your previous test even better by adding more use cases via
subtests.

## Tests as examples

When you write tests in Go using the `_test` packages, you
are somehow showing how to use the package. This is so
useful that Go provides a way to convert those tests into
examples.

Examples are functions whose name starts with `Example` and
do not receive any parameters, nor return any values.

The name of an example contains the information of what the example is demoing. In order to write an example for:

- a function `Foo`, write `ExampleFoo`,
- a method `Foo` on type `Bar`, write `ExampleFooBar`,
- a type `Foo`, write `ExampleFoo`

During tests we can verify the results by writing code, in examples
we have a more limited solution: by checking the expected output
of executing the example.

To do this, we simply add a comment of the form:

```go
func ExmampleHelloWorld() {
    fmt.Println("hello, world")
    // Output:
    // hello, world
}
```

### Exercise: writing examples

Write an example for `Sum` and see it running on your local `godoc` web
server.
Use an `Output` comment to make sure that if you run `go test` and the
output doesn't match the example will fail.

### Exercise (optional): testing an `http.Handler`

In order to test an `http.Handler` or an `http.HandlerFunc` we need to provide
an `http.Request` and an `http.ResponseWriter`. The first one is easy to create,
but the second one is a bit trickier.

For a slightly more advanced exercise try to write a test for the `handler`
function defined in our [webserver](../webserver) program.

## Congratulations

You're now able to write unit tests and examples, and you're able to run them
on any package you wish. That's pretty awesome!

Next we're going to learn how to figure out what parts of your code you should
be testing with [code coverage](2-code-coverage.md).
