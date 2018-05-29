package main

import (
	. "github.com/captainsano/golang-chess/core"
)

var rays, between = Rays()

func main() {
	fen := "rnbqkbnr/ppp2ppp/3p4/4p3/2B1P3/5N2/PPPP1PPP/RNBQK2R"

	board := MakeBoard(fen)

	print(board.Unicode(false, false))
}
