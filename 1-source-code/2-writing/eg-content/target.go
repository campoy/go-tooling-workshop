package main

import "fmt"

func main() { fmt.Printf("Error: %v", fmt.Errorf("%s", "Whoops!")) }
