package core

type Move struct {
	FromSquare Square
	ToSquare   Square
	Promotion  *Piece
	// Drop?
}

func NewMove(fromSquare, toSquare Square, promotion *Piece) Move {
	var p *Piece
	if promotion != nil {
		*p = NewPiece(promotion.Type, promotion.Color)
	}

	return Move{
		FromSquare: fromSquare,
		ToSquare:   toSquare,
		Promotion:  p,
	}
}

func NewNullMove() Move {
	return Move{
		FromSquare: 0,
		ToSquare:   0,
	}
}

func (m *Move) Uci() string {
	if m.Promotion != nil {
		return m.FromSquare.Name() + m.ToSquare.Name() + m.Promotion.Symbol()
	}

	return m.FromSquare.Name() + m.ToSquare.Name()
}
