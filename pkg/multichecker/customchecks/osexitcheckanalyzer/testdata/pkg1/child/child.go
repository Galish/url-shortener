package child

import (
	"fmt"
	"os"
)

func main() {
	fmt.Println("Welcome!")

	defer os.Exit(0)

	defer fmt.Println("Bye bye")

	os.Exit(0)

	exit()

	func() {
		os.Exit(0)
	}()
}

func exit() {
	os.Exit(0)
}
