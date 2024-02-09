package maintest

import (
	"fmt"
	"os"
)

func main() {
	fmt.Println("Welcome")

	defer fmt.Println("Bye bye")

	os.Exit(3)
}
