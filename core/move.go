package core

import (
	"fmt"
	"strings"
)

type Move struct {
	FromSquare Square
	ToSquare   Square
	Promotion  PieceType
	Drop       PieceType
}

func NewMove(fromSquare, toSquare Square, promotion, drop PieceType) Move {
	return Move{
		FromSquare: fromSquare,
		ToSquare:   toSquare,
		Promotion:  promotion,
		Drop:       drop,
	}
}

func NewPromotionMove(fromSquare, toSquare Square, promotion PieceType) Move {
	return Move{
		FromSquare: fromSquare,
		ToSquare:   toSquare,
		Promotion:  promotion,
		Drop:       NoPiece,
	}
}

func NewNormalMove(fromSquare, toSquare Square) Move {
	return Move{
		FromSquare: fromSquare,
		ToSquare:   toSquare,
		Promotion:  NoPiece,
		Drop:       NoPiece,
	}
}

func NewNullMove() Move {
	return Move{
		FromSquare: SquareNone,
		ToSquare:   SquareNone,
		Promotion:  NoPiece,
		Drop:       NoPiece,
	}
}

func NewDropMove(fromSquare, toSquare Square, drop PieceType) Move {
	return Move{
		FromSquare: SquareNone,
		ToSquare:   SquareNone,
		Promotion:  NoPiece,
		Drop:       drop,
	}
}

func NewMoveFromUci(uci string) Move {
	if uci == "0000" {
		return NewNullMove()
	}

	if len(uci) == 4 && uci[1] == '@' {
		drop := NewPieceFromSymbol(string(uci[0])).Type
		square := NewSquareFromName(uci[2:])
		return NewDropMove(square, square, drop)
	}

	if len(uci) == 4 {
		return NewNormalMove(NewSquareFromName(uci[0:2]), NewSquareFromName(uci[2:4]))
	}

	if len(uci) == 5 {
		promotion := NewPieceFromSymbol(string(uci[4])).Type
		return NewPromotionMove(NewSquareFromName(uci[0:2]), NewSquareFromName(uci[2:4]), promotion)
	}

	panic("Expected UCI string to be of length 4 or 5")
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
