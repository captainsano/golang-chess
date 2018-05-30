package core

import (
	"strings"
)

type PieceType uint8

const (
	NoPiece PieceType = 0
	Pawn    PieceType = 1
	Knight  PieceType = 2
	Bishop  PieceType = 3
	Rook    PieceType = 4
	Queen   PieceType = 5
	King    PieceType = 6
)

func NewPieceType(i uint8) PieceType {
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

func (p PieceType) Symbol() string {
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

func (p PieceType) Name() string {
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

type Piece struct {
	Type  PieceType
	Color Color
}

func NewPiece(t PieceType, c Color) Piece {
	return Piece{Type: t, Color: c}
}

func NewPieceFromSymbol(s string) Piece {
	switch s {
	case "R":
		return Piece{Rook, White}
	case "r":
		return Piece{Rook, Black}
	case "N":
		return Piece{Knight, White}
	case "n":
		return Piece{Knight, Black}
	case "B":
		return Piece{Bishop, White}
	case "b":
		return Piece{Bishop, Black}
	case "Q":
		return Piece{Queen, White}
	case "q":
		return Piece{Queen, Black}
	case "K":
		return Piece{King, White}
	case "k":
		return Piece{King, Black}
	case "P":
		return Piece{Pawn, White}
	case "p":
		return Piece{Pawn, Black}
	}

	panic("Invalid piece symbol")
}

func (p *Piece) Symbol() string {
	if p.Color == White {
		return strings.ToUpper(p.Type.Symbol())
	}

	return p.Type.Symbol()
}

func (p *Piece) UnicodeSymbol(invertColor bool) string {
	s := p.Symbol()

	if invertColor && strings.ToUpper(s) == s {
		return UnicodePieceSymbol(strings.ToLower(s))
	}

	if invertColor && strings.ToLower(s) != s {
		return UnicodePieceSymbol(strings.ToUpper(s))
	}

	return UnicodePieceSymbol(s)
}

// @TODO: Implement hash function
