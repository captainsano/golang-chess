package core

type Move struct {
	FromSquare Square
	ToSquare   Square
	Promotion  PieceType
	// Drop?
}

func NewMove(fromSquare, toSquare Square, promotion PieceType) Move {
	return Move{
		FromSquare: fromSquare,
		ToSquare:   toSquare,
		Promotion:  promotion,
	}
}

// func NewMoveFromUci(uci string) Move {

// }

func NewNullMove() Move {
	return Move{
		FromSquare: SquareNone,
		ToSquare:   SquareNone,
	}
}

func (m *Move) IsNull() bool {
	return m.FromSquare == SquareNone || m.ToSquare == SquareNone
}

func (m *Move) Uci() string {
	if m.Promotion != NoPiece {
		return m.FromSquare.Name() + m.ToSquare.Name() + m.Promotion.Symbol()
	}

	return m.FromSquare.Name() + m.ToSquare.Name()
}
