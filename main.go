package main

import (
	"fmt"

	. "github.com/captainsano/golang-chess/core"
)

var rays, between = Rays()

func printAttacks(from, to Square) {
	x := rays[from][to]

	fmt.Println("--> ", from.Name(), " to ", to.Name())
	fmt.Println(x.Ascii())
	fmt.Println()
}

func printBetween(from, to Square) {
	x := between[from][to]

	fmt.Println("--> ", from.Name(), " to ", to.Name())
	fmt.Println(x.Ascii())
	fmt.Println()
}

func main() {
	fmt.Println("Attacks: ")
	printAttacks(B1, G1)

	fmt.Println("Between: ")
	printBetween(B1, G1)
}
