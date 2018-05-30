package main

import (
	. "github.com/captainsano/golang-chess/core"
)

var rays, between = Rays()

func main() {
	fen := StartingBoardFEN

	board := NewBaseBoard(fen)

	board.RemovePieceAt(E2)
	p := NewPiece(Pawn, White)
	board.SetPieceAt(E4, &p, false)

	print(board.Unicode(false, false))
}
