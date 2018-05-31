package main

import (
	"errors"
	"fmt"
)

func main() {
	if err := run(); err != nil {
		fmt.Printf("could not run: %v", err)
	}
}

func run() error {
	return errors.New("didn't want to run")
}
