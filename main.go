package main

import (
	"fmt"

	. "github.com/captainsano/golang-chess/core"
)

func main() {
	fen := "rnb1kbnr/ppp1q1pp/8/3P1p2/2PP4/3B1p2/PP3PPP/RNBQK2R w KQkq - 0 7"
	// fen := StartingFEN

	board := NewBoard(fen, false)

	fmt.Println("Current Board: ")
	fmt.Println(board.Unicode(false, false))

	fmt.Println("Legal Moves: ")
	for m := range board.GenerateLegalMoves(BBAll, BBAll) {
		fmt.Println(board.San(&m))
	}

}
