package main

import (
	"fmt"

	. "github.com/captainsano/golang-chess/core"
)

var rays, between = Rays()

func main() {
	fen := StartingFEN

	board := NewBoard(fen, false)

	// TODO: Debug why the white pieces are black
	fmt.Println("Current Board: ")
	fmt.Println(board.Ascii())

	for m := range board.GenerateAllPseudoLegalMoves() {
		fmt.Println("--> ", m.Uci())
	}
}
