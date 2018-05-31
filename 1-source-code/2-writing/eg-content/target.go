package main

import "fmt"

func main() {
	if err := run(); err != nil {
		fmt.Printf("could not run: %v", err)
	}
}

func run() error {
	return fmt.Errorf("%s", "didn't want to run")
}
