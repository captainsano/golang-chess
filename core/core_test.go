package core

import (
	"testing"
)

func TestSquare(t *testing.T) {
	for _, sq := range squares {
		file := sq.File()
		rank := sq.Rank()

		if NewSquare(file, rank) != sq {
			t.Errorf("Square not equal %v %v", NewSquare(file, rank), sq)
		}

		if NewSquare(file, rank).Name() != file.Name()+rank.Name() {
			t.Errorf("Square name not equal %v %v %v", NewSquare(file, rank), file.Name(), rank.Name())
		}
	}
}

func TestShifts(t *testing.T) {
	assertShift := func(t *testing.T, bbSq Bitboard) {
		c := bbSq.PopCount()

		if !(c <= 1) {
			t.Errorf("ShiftDown Failed %v", bbSq)
		}

		if c != (bbSq & BBAll).PopCount() {
			t.Errorf("ShiftDown Failed %v", bbSq)
		}
	}

	t.Run("Shift Down", func(t *testing.T) {
		for _, sq := range squares {
			bbSq := NewBitboardFromSquare(sq)
			bbSq.ShiftDown()
			assertShift(t, bbSq)
		}
	})

	t.Run("Shift 2 Down", func(t *testing.T) {
		for _, sq := range squares {
			bbSq := NewBitboardFromSquare(sq)
			bbSq.Shift2Down()
			assertShift(t, bbSq)
		}
	})

	t.Run("Shift Up", func(t *testing.T) {
		for _, sq := range squares {
			bbSq := NewBitboardFromSquare(sq)
			bbSq.ShiftUp()
			assertShift(t, bbSq)
		}
	})

	t.Run("Shift 2 Up", func(t *testing.T) {
		for _, sq := range squares {
			bbSq := NewBitboardFromSquare(sq)
			bbSq.Shift2Up()
			assertShift(t, bbSq)
		}
	})

	t.Run("Shift Right", func(t *testing.T) {
		for _, sq := range squares {
			bbSq := NewBitboardFromSquare(sq)
			bbSq.ShiftRight()
			assertShift(t, bbSq)
		}
	})

	t.Run("Shift 2 Right", func(t *testing.T) {
		for _, sq := range squares {
			bbSq := NewBitboardFromSquare(sq)
			bbSq.Shift2Right()
			assertShift(t, bbSq)
		}
	})

	t.Run("Shift Left", func(t *testing.T) {
		for _, sq := range squares {
			bbSq := NewBitboardFromSquare(sq)
			bbSq.ShiftLeft()
			assertShift(t, bbSq)
		}
	})

	t.Run("Shift 2 Left", func(t *testing.T) {
		for _, sq := range squares {
			bbSq := NewBitboardFromSquare(sq)
			bbSq.Shift2Left()
			assertShift(t, bbSq)
		}
	})

	t.Run("Shift Up Left", func(t *testing.T) {
		for _, sq := range squares {
			bbSq := NewBitboardFromSquare(sq)
			bbSq.ShiftUpLeft()
			assertShift(t, bbSq)
		}
	})

	t.Run("Shift Up Right", func(t *testing.T) {
		for _, sq := range squares {
			bbSq := NewBitboardFromSquare(sq)
			bbSq.ShiftUpRight()
			assertShift(t, bbSq)
		}
	})

	t.Run("Shift Down Left", func(t *testing.T) {
		for _, sq := range squares {
			bbSq := NewBitboardFromSquare(sq)
			bbSq.ShiftDownLeft()
			assertShift(t, bbSq)
		}
	})

	t.Run("Shift Down Right", func(t *testing.T) {
		for _, sq := range squares {
			bbSq := NewBitboardFromSquare(sq)
			bbSq.ShiftDownRight()
			assertShift(t, bbSq)
		}
	})
}

func TestMove(t *testing.T) {

	t.Run("Equality", func(t *testing.T) {
		a, _ := NewNormalMove(A1, A2)
		b, _ := NewNormalMove(A1, A2)
		c, _ := NewPromotionMove(H7, H8, Bishop)
		d1, _ := NewNormalMove(H7, H8)
		d2, _ := NewNormalMove(H7, H8)

		if *a != *b {
			t.Errorf("Move not equal %v %v", *a, *b)
		}

		if *b != *a {
			t.Errorf("Move not equal %v %v", *b, *a)
		}

		if *d1 != *d2 {
			t.Errorf("Move not equal %v %v", *d1, *d2)
		}

		if *a == *c {
			t.Errorf("Move equal %v %v", *a, *c)
		}

		if *c == *d1 {
			t.Errorf("Move equal %v %v", *c, *d1)
		}

		if *b == *d1 {
			t.Errorf("Move equal %v %v", *b, *d1)
		}

		if (*d1 != *d2) == true {
			t.Errorf("Move equal %v %v", *d1, *d2)
		}
	})

	t.Run("UCI parsing", func(t *testing.T) {
		table := []string{"b5c7", "e7e8q", "P@e4", "B@f4"}

		for _, u := range table {
			m, _ := NewMoveFromUci(u)
			if m == nil || m.Uci() != u {
				t.Errorf("Error in UCI move: %v", u)
			}
		}
	})

	t.Run("Invalid UCI", func(t *testing.T) {
		table := []string{"", "N", "z1g3", "Q@g9"}

		for _, u := range table {
			m, err := NewMoveFromUci("")
			if err == nil || m != nil {
				t.Errorf("Expected invalid move %v", u)
			}
		}
	})
}

func TestPiece(t *testing.T) {
	t.Run("Equality", func(t *testing.T) {
		a := NewPiece(Bishop, White)
		b := NewPiece(King, Black)
		c := NewPiece(King, White)
		d1 := NewPiece(Bishop, White)
		d2 := NewPiece(Bishop, White)

		table := []struct {
			x, y *Piece
			eq   bool
		}{
			{&a, &d1, true},
			{&d1, &a, true},
			{&d1, &d2, true},
			{&a, &b, false},
			{&b, &c, false},
			{&b, &d1, false},
			{&a, &c, false},
			{&d1, &d2, true},
		}

		for _, c := range table {
			if (*(c.x) == *(c.y)) != c.eq {
				t.Errorf("Error in piece equality test")
			}

			if (c.x.Symbol() == c.y.Symbol()) != c.eq {
				t.Errorf("Error in piece symbol equality test")
			}
		}
	})

	t.Run("Symbol", func(t *testing.T) {
		wn := NewPieceFromSymbol("N")
		if wn.Type != Knight || wn.Color != White || wn.Symbol() != "N" {
			t.Errorf("Piece from symbol failed")
		}

		bq := NewPieceFromSymbol("q")
		if bq.Type != Queen || bq.Color != Black || bq.Symbol() != "q" {
			t.Errorf("Piece from symbol failed")
		}
	})
}
