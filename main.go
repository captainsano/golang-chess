package main

import (
	"fmt"

	"github.com/captainsano/golang-chess/core"
)

func main() {
	rays, between := core.Rays()

	from := core.A1
	to := core.A5

	fmt.Println("Rays: ")
	fmt.Println(rays[from][to].Ascii())

	fmt.Println()

	fmt.Println("Between: ")
	fmt.Println(between[from][to].Ascii())
}
