# Detecting mistakes early on

How many times have I written `fmt.Println` when I meant to say `fmt.Printf`?
Way too many to be funny!

These mistakes are not detected by the compiler, since `fmt.Println("hello, %s", name)`
is valid Go. But wouldn't it be great if a tool could tell me about it before deploying
to production and polluting my logs?

Let's see some tools that will reduce the number of mistakes you push to prod!

## Vetting your code with `go vet`

`go vet` detects many classes of mistakes doing a static analysis of your code. It is
powerful enough to detect useful cases, but does not guarantee that your code is bug free.

For instance, if we have the following Go code in a file named `main.go`:

```go
package main

import "fmt"

func main() {
	name := "francesc"
	fmt.Println("hello, %s", name)
}
```

We can vet for mistakes running the following command:

```bash
$ go vet main.go
main.go:7: possible formatting directive in Println call
exit status 1
```

Do you see where the mistake is?

You can learn more about the powers of `go vet` by running `go doc cmd/vet`.

Can you see the mistakes on this code?

```go
package main

import (
	"fmt"
	"time"
)

type Person struct {
	First_Name string `json:name`
	Last_Name  string `json:"-"`
}

func (this Person) Full_Name() string {
	return fmt.Sprint("%s %s", this.First_Name, this.Last_Name)
}

func main() {
	p := Person{"Francesc", "Campoy"}
	for i := 0; ; i++ {
		time.Sleep(time.Second)
		fmt.Printf("hello, %s", p.Full_Name)
	}
	fmt.Println("done")
}
```

Copy the code onto VSCode and see how green lines appear on every mistake.
Fix all mistakes until `go vet` gives not output.

## Can we make our code even better?

At this point, the code above doesn't not given any `go vet` warning. Great!
Does that mean the code is perfect? Well, `go vet` only detects code that
seems to have the wrong behavior, but it doesn't enforce any style guides or
other conventions that all gophers have agreed on.

To enforce those we have `golint`!

## Enforcing conventions with `golint`

Let's run `golint` on the code you ended up having after fixing all the `go vet`
warnings.

If you didn't fix those and want to continue, feel free to copy the code inside
of the `details` block below.

<details>

```go
package main

import (
	"fmt"
	"time"
)

type Person struct {
	First_Name string `json:"name"`
	Last_Name  string `json:"-"`
}

func (this Person) Full_Name() string {
	return fmt.Sprintf("%s %s", this.First_Name, this.Last_Name)
}

func main() {
	p := Person{"Francesc", "Campoy"}
	for i := 0; ; i++ {
		time.Sleep(time.Second)
		fmt.Printf("hello, %s", p.Full_Name())
	}
}
```

</details>

Let's run `golint` on it:

```bash
$ golint main.go
main.go:8:6: exported type Person should have comment or be unexported
main.go:9:2: don't use underscores in Go names; struct field First_Name should be FirstName
main.go:10:2: don't use underscores in Go names; struct field Last_Name should be LastName
main.go:13:1: exported method Person.Full_Name should have comment or be unexported
main.go:13:1: receiver name should be a reflection of its identity; don't use generic names such as "this" or "self"
main.go:13:20: don't use underscores in Go names; method Full_Name should be FullName
```

Woah! So many things to fix!

Why didn't VSCode show this before? Well, I agree with that. Let's add `golint`
to VSCode.

To do so open the *Settings* page (`Command` + `,` on Mac) and find the property
named `"go.lintOnSave"`. Edit it to be `"true"`.

```json
"go.lintOnSave": "true",
```

Now when you edit and save your code both `go vet` and `golint` will be used.

As an exercise fix all of the `golint` warnings too.

## More early error detection

As for editing, many tools have been created that can be used to analyze your code
early on and detect mistakes. We will cover just one more, one that will detect
what the problem is with the following code.

[embedmd]:# (errcheck.go /package main/ $)
```go
package main

import (
	"fmt"
	"net/http"
)

func main() {
	http.HandleFunc("/", handle)
	http.ListenAndServe(":80", nil)
}

func handle(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "hello")
}
```

Can you see the problem? Try running this web server.
You might see no output at all. Why?

Turns out `http.ListenAndServer` returns an `error` and we're not checking it!
This kind of mistakes can actually cause issues in your production code, so it
is incredibly important that no errors are *silently* ignored.

We're lucky, `errcheck` checks exactly that!

# TODO

explain how to install and use errcheck, also how to integrate it with VSCode

## Congratulations

You are now able to produce code with less mistakes, because you're able to detect
them as soon as you write them!

Great, so we know how to modify code, what's next? Reading it! There's a lot of
tools that help gopher navigate unknown codebases, and we'll see them in the
[next section](../3-reading/1-godoc.md).