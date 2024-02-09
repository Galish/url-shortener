package main

import (
	"fmt"
	"os"
)

func main() {
	fmt.Println("Welcome")

	defer fmt.Println("Bye bye")

	os.Exit(3) // want "illegal `os.Exit` call"
}
