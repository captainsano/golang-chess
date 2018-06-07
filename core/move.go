package core

import (
	"fmt"
	"strings"
)

type MoveError struct {
	error
	description string
}

func (e *MoveError) Error() string {
	return e.description
}

type Move struct {
	FromSquare Square
	ToSquare   Square
	Promotion  PieceType
	Drop       PieceType
}

func NewMove(fromSquare, toSquare Square, promotion, drop PieceType) (*Move, error) {
	return &Move{fromSquare, toSquare, promotion, drop}, nil
}

func NewPromotionMove(fromSquare, toSquare Square, promotion PieceType) (*Move, error) {
	return &Move{fromSquare, toSquare, promotion, NoPiece}, nil
}

func NewNormalMove(fromSquare, toSquare Square) (*Move, error) {
	return &Move{fromSquare, toSquare, NoPiece, NoPiece}, nil
}

func NewNullMove() (*Move, error) {
	return &Move{SquareNone, SquareNone, NoPiece, NoPiece}, nil
}

func NewDropMove(square Square, drop PieceType) (*Move, error) {
	return &Move{square, square, NoPiece, drop}, nil
}

func NewMoveFromUci(uci string) (*Move, error) {
	if uci == "0000" {
		return NewNullMove()
	}

	if len(uci) == 4 && uci[1] == '@' {
		drop := NewPieceFromSymbol(string(uci[0])).Type
		square := NewSquareFromName(uci[2:])
		return NewDropMove(square, drop)
	}

	if len(uci) == 4 {
		return NewNormalMove(NewSquareFromName(uci[0:2]), NewSquareFromName(uci[2:4]))
	}

	if len(uci) == 5 {
		promotion := NewPieceFromSymbol(string(uci[4])).Type
		return NewPromotionMove(NewSquareFromName(uci[0:2]), NewSquareFromName(uci[2:4]), promotion)
	}

	return nil, &MoveError{description: "Invalid uci string:" + uci}
}

func (m *Move) Uci() string {
	if m.Drop != NoPiece {
		return strings.ToUpper(m.Drop.Symbol()) + "@" + m.ToSquare.Name()
	}

	if m.Promotion != NoPiece {
		return m.FromSquare.Name() + m.ToSquare.Name() + m.Promotion.Symbol()
	}

	if m.IsNotNull() {
		return m.FromSquare.Name() + m.ToSquare.Name()
	}

	return "0000"
}

func (m *Move) IsNotNull() bool {
	return m.FromSquare != SquareNone || m.ToSquare != SquareNone || m.Promotion != NoPiece || m.Drop != NoPiece
}

func (m *Move) String() string {
	return m.Uci()
}

func (m *Move) Hash() string {
	return fmt.Sprintf("%d%d%s%s", m.FromSquare, m.ToSquare, m.Promotion.Symbol(), m.Drop.Symbol())
}
