package main

import (
	"fmt"

	. "github.com/captainsano/golang-chess/core"
)

func oldMain() {
	// fen := "rnb1kbnr/ppp1q1pp/8/3P1p2/2PP4/3B1p2/PP3PPP/RNBQK2R w KQkq - 0 7"
	// fen := "rnbqkbnr/pppppppp/8/8/4P3/8/PPPP1PPP/RNBQKBNR b KQkq - 1 1"
	board := NewDefaultBoard()

	printStatus(&board)
	printAfterMove(&board, "e2e4")
	printAfterMove(&board, "c7c5")
	printAfterMove(&board, "g1f3")
	printAfterMove(&board, "d7d6")
	printAfterMove(&board, "d2d4")
	printAfterMove(&board, "c5d4")
	printAfterMove(&board, "f3d4")
	printAfterMove(&board, "g8f6")
	printAfterMove(&board, "b1c3")
	printAfterMove(&board, "a7a6")
	printAfterMove(&board, "f1e2")
	printAfterMove(&board, "e7e5")
	printAfterMove(&board, "d4f3")
	printAfterMove(&board, "f8e7")
	printAfterMove(&board, "e1g1")
	printAfterMove(&board, "e8g8")
}

func printAfterMove(board *Board, uci string) {
	m, _ := NewMoveFromUci(uci)
	board.Push(m)
	fmt.Println(board.FEN(false, "legal", NoPiece))
	printStatus(board)
}

func printStatus(b *Board) {
	fmt.Println()
	fmt.Println("Current Board: ")
	fmt.Println(b.Unicode(false, false))

	fmt.Print("Legal Moves: ")
	for m := range b.GenerateLegalMoves(BBAll, BBAll) {
		fmt.Print(b.San(&m), " ")
	}
	fmt.Println()
}

func main() {
	fen := "8/8/8/R2nkn2/8/8/2K5/8 b - - 0 1"
	b := NewBoardFromFEN(fen, false)

	m1, _ := NewMoveFromUci("f5e3")
	b.Push(m1)

	fmt.Println("Current position:")
	fmt.Println(b)

	for m := range b.GenerateLegalMoves(BBAll, BBAll) {
		fmt.Println("--> Legal Move: ", m.Uci())
	}
}
