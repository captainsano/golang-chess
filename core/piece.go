package core

type Piece uint8

const (
	Pawn   Piece = 1
	Knight Piece = 2
	Bishop Piece = 3
	Rook   Piece = 4
	Queen  Piece = 5
	King   Piece = 6
)

func MakePiece(i uint8) Piece {
	switch i {
	case 1:
		return Pawn
	case 2:
		return Knight
	case 3:
		return Bishop
	case 4:
		return Rook
	case 5:
		return Queen
	case 6:
		return King
	}

	panic("Invalid piece code")
}

func (p Piece) Symbol() string {
	switch p {
	case Pawn:
		return "p"
	case Knight:
		return "n"
	case Bishop:
		return "b"
	case Rook:
		return "r"
	case Queen:
		return "q"
	case King:
		return "k"
	}

	panic("Invalid piece")
}

func (p Piece) Name() string {
	switch p {
	case Pawn:
		return "pawn"
	case Knight:
		return "knight"
	case Bishop:
		return "bishop"
	case Rook:
		return "rook"
	case Queen:
		return "queen"
	case King:
		return "king"
	}

	panic("Invalid piece")
}

func UnicodePieceSymbol(fenPiece string) string {
	switch fenPiece {
	case "R":
		return "♖"
	case "r":
		return "♜"
	case "N":
		return "♘"
	case "n":
		return "♞"
	case "B":
		return "♗"
	case "b":
		return "♝"
	case "Q":
		return "♕"
	case "q":
		return "♛"
	case "K":
		return "♔"
	case "k":
		return "♚"
	case "P":
		return "♙"
	case "p":
		return "♟"
	}

	return ""
}
