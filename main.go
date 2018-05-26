package main

type Color int

const (
	Black Color = 0
	White Color = 1
)

var colorNames = []string{"black", "white"}

func ColorNames(i int) string {
	return colorNames[i]
}

type Piece int

const (
	Pawn   Piece = 1
	Knight Piece = 2
	Bishop Piece = 3
	Rook   Piece = 4
	Queen  Piece = 5
	King   Piece = 6
)

var pieceSymbol = []string{"", "p", "n", "b", "r", "q", "k"}

func PieceSymbol(i int) string {
	return pieceSymbol[i]
}

var pieceNames = []string{"", "pawn", "knight", "bishop", "rook", "queen", "king"}

func PieceNames(i int) string {
	return pieceNames[i]
}

var unicodePieceSymbols = map[string]string{
	"R": "♖", "r": "♜",
	"N": "♘", "n": "♞",
	"B": "♗", "b": "♝",
	"Q": "♕", "q": "♛",
	"K": "♔", "k": "♚",
	"P": "♙", "p": "♟",
}

func UnicodePieceSymbols(s string) string {
	return unicodePieceSymbols[s]
}

var fileNames = []string{"a", "b", "c", "d", "e", "f", "g", "h"}

func FileNames(i int) string {
	return fileNames[i]
}

var rankNames = []string{"1", "2", "3", "4", "5", "6", "7", "8"}

func RankNames(i int) string {
	return rankNames[i]
}

/* FEN Parsing */
const (
	StartingFen      = "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1"
	StartingBoardFen = "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR"
)

func main() {
}
