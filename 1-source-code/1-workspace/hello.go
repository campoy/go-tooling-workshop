package main

import (
	"fmt"
	"foo"
	"math"

	_ "github.com/golang/example/stringutil"
)

func main() {
	fmt.Println("Pi is", math.Pi)
	fmt.Println("Foo is", foo.Foo)
}
