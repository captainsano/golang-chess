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
		a := NewNormalMove(A1, A2)
		b := NewNormalMove(A1, A2)
		c := NewPromotionMove(H7, H8, Bishop)
		d1 := NewNormalMove(H7, H8)
		d2 := NewNormalMove(H7, H8)

		if a != b {
			t.Errorf("Move not equal %v %v", a, b)
		}

		if b != a {
			t.Errorf("Move not equal %v %v", b, a)
		}

		if d1 != d2 {
			t.Errorf("Move not equal %v %v", d1, d2)
		}

		if a == c {
			t.Errorf("Move equal %v %v", a, c)
		}

		if c == d1 {
			t.Errorf("Move equal %v %v", c, d1)
		}

		if b == d1 {
			t.Errorf("Move equal %v %v", b, d1)
		}

		if (d1 != d2) == true {
			t.Errorf("Move equal %v %v", d1, d2)
		}
	})

	t.Run("UCI parsing", func(t *testing.T) {
		table := []string{"b5c7", "e7e8q", "P@e4", "B@f4"}

		for _, u := range table {
			m := NewMoveFromUci(u)
			if m.Uci() != u {
				t.Errorf("Error in UCI move: %v", u)
			}
		}
	})

	t.Run("Invalid UCI", func(t *testing.T) {
		// @TODO
	})
}
