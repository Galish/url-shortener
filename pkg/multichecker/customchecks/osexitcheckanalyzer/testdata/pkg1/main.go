package main

import (
	"fmt"
	"os"
)

func main() {
	fmt.Println("Welcome!")

	defer os.Exit(0) // want "illegal `os.Exit` call"

	defer fmt.Println("Bye bye")

	os.Exit(3) // want "illegal `os.Exit` call"

	exit()

	func() {
		os.Exit(0) // want "illegal `os.Exit` call"
	}()
}

func exit() {
	os.Exit(0)
}
