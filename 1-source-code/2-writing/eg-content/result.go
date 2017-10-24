package main

import (
	"errors"
	"fmt"
)

func main() { fmt.Printf("Error: %v", errors.New("Whoops!")) }
