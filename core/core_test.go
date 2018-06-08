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

func TestBoard(t *testing.T) {
	t.Run("Default position", func(t *testing.T) {
		b := NewDefaultBoard()

		if *(b.PieceAt(B1)) != NewPieceFromSymbol("N") {
			t.Errorf("pieceAt failed")
		}

		if b.FEN(false, "legal", NoPiece) != StartingFEN {
			t.Errorf("FEN generation failed actual: %v expected %v", b.FEN(false, "legal", NoPiece), StartingFEN)
		}

		if b.Turn() != White {
			t.Errorf("Turn matching failed")
		}
	})

	t.Run("Empty", func(t *testing.T) {
		b := NewBoard(false)

		if b.FEN(false, "legal", NoPiece) != "8/8/8/8/8/8/8/8 w - - 0 1" {
			t.Errorf("creating empty board failed")
		}

		// TODO: Use go-cmp
		b1 := NewBoard(false)
		if !b.Equal(&b1) {
			t.Errorf("empty board equality failed")
		}
	})

	// TODO: EPD
	// t.Run("Test from epd", func(t *testing.T) {
	// 	baseEpd := "rnbqkb1r/ppp1pppp/5n2/3P4/8/8/PPPP1PPP/RNBQKBNR w KQkq -"
	// 	b, ops := NewBoardFromEpd(baseEpd)

	// 	if ops["ce"] != 55 {
	// 		t.Errorf("EPD operation not matching")
	// 	}

	// 	if b.FEN(false, "legal", NoPiece) != baseEpd+" 0 1" {
	// 		t.Errorf("FEN not matching EPD")
	// 	}
	// })

	t.Run("Move making", func(t *testing.T) {
		b := NewDefaultBoard()
		move, _ := NewNormalMove(E2, E4)
		b.Push(move)

		if *move != *(b.Peek()) {
			t.Errorf("moves not matching")
		}
	})

	t.Run("FEN", func(t *testing.T) {
		b := NewDefaultBoard()

		if b.FEN(false, "legal", NoPiece) != StartingFEN {
			t.Error("FEN not matching")
		}

		fen := "6k1/pb3pp1/1p2p2p/1Bn1P3/8/5N2/PP1q1PPP/6K1 w - - 0 24"
		b.SetFEN(fen)
		if b.FEN(false, "legal", NoPiece) != fen {
			t.Error("FEN not matching")
		}

		m, _ := NewMoveFromUci("f3d2")
		b.Push(m)
		if b.FEN(false, "legal", NoPiece) != "6k1/pb3pp1/1p2p2p/1Bn1P3/8/8/PP1N1PPP/6K1 b - - 0 24" {
			t.Error("FEN not matching")
		}
	})

	t.Run("XFEN", func(t *testing.T) {
		t.Run("normal", func(t *testing.T) {
			xfen := "rn2k1r1/ppp1pp1p/3p2p1/5bn1/P7/2N2B2/1PPPPP2/2BNK1RR w Gkq - 4 11"
			b := NewBoardFromFEN(xfen, true)

			if b.CastlingRights() != BBG1|BBA8|BBG8 {
				t.Errorf("castling rights not matching")
			}
			if b.CleanCastlingRights() != BBG1|BBA8|BBG8 {
				t.Errorf("clean castling rights not matching")
			}
			if b.ShredderFEN("legal", NoPiece) != "rn2k1r1/ppp1pp1p/3p2p1/5bn1/P7/2N2B2/1PPPPP2/2BNK1RR w Gga - 4 11" {
				t.Errorf("shredder fen not matching, actual: %v", b.ShredderFEN("legal", NoPiece))
			}
			if b.FEN(false, "legal", NoPiece) != xfen {
				t.Errorf("fen not matching")
			}
			if !b.HasCastlingRights(White) {
				t.Errorf("has castling rights not matching")
			}
			if !b.HasCastlingRights(Black) {
				t.Errorf("has castling rights not matching")
			}
			if !b.HasKingsideCastlingRights(Black) {
				t.Errorf("has castling rights not matching")
			}
			if !b.HasKingsideCastlingRights(White) {
				t.Errorf("has castling rights not matching")
			}
			if !b.HasQueensideCastlingRights(Black) {
				t.Errorf("has castling rights not matching")
			}
			if b.HasQueensideCastlingRights(White) {
				t.Errorf("has castling rights not matching")
			}
		})

		t.Run("Chess960 #284", func(t *testing.T) {
			b := NewBoardFromFEN("rkbqrbnn/pppppppp/8/8/8/8/PPPPPPPP/RKBQRBNN w - - 0 1", true)
			b.castlingRights = b.baseBoard.rooks

			if !b.CleanCastlingRights().IsMaskingBB(BBA1) {
				t.Errorf("Chess960 castling rights not matching")
			}
			if b.FEN(false, "legal", NoPiece) != "rkbqrbnn/pppppppp/8/8/8/8/PPPPPPPP/RKBQRBNN w KQkq - 0 1" {
				t.Errorf("Chess960 FEN not matching")
			}
			if b.ShredderFEN("legal", NoPiece) != "rkbqrbnn/pppppppp/8/8/8/8/PPPPPPPP/RKBQRBNN w EAea - 0 1" {
				t.Errorf("Chess960 Shredder FEN not matching")
			}
		})

		t.Run("Valid enpassant square on illegal board", func(t *testing.T) {
			fen := "8/8/8/pP6/8/8/8/8 w - a6 0 1"
			b := NewBoardFromFEN(fen, false)
			if b.FEN(false, "legal", NoPiece) != fen {
				t.Errorf("enpassant with invalid FEN not matching")
			}
		})

		t.Run("Illegal enpassant square on illegal board", func(t *testing.T) {
			fen := "1r6/8/8/pP6/8/8/8/1K6 w - a6 0 1"
			b := NewBoardFromFEN(fen, false)
			if b.FEN(false, "legal", NoPiece) != "1r6/8/8/pP6/8/8/8/1K6 w - - 0 1" {
				t.Errorf("illegal enpassant with invalid FEN not matching, actual: %v", b.FEN(false, "legal", NoPiece))
			}
		})
	})

	t.Run("FEN enpassant", func(t *testing.T) {
		b := NewDefaultBoard()
		b.PushSan("e4")

		if b.FEN(false, "fen", NoPiece) != "rnbqkbnr/pppppppp/8/8/4P3/8/PPPP1PPP/RNBQKBNR b KQkq e3 0 1" {
			t.Errorf("FEN not matching")
		}

		if b.FEN(false, "xfen", NoPiece) != "rnbqkbnr/pppppppp/8/8/4P3/8/PPPP1PPP/RNBQKBNR b KQkq - 0 1" {
			t.Errorf("FEN not matching")
		}
	})

	t.Run("Get Set", func(t *testing.T) {
		b := NewDefaultBoard()

		if *(b.PieceAt(B1)) != NewPieceFromSymbol("N") {
			t.Errorf("piece at square failed")
		}

		b.RemovePieceAt(E2)
		if b.PieceAt(E2) != nil {
			t.Errorf("piece at square failed")
		}

		p1 := NewPieceFromSymbol("r")
		b.SetPieceAt(E4, &p1, false)
		if b.PieceAt(E4).Type != Rook {
			t.Errorf("piece at square failed")
		}

		b.SetPieceAt(F1, nil, false)
		if b.PieceAt(F1) != nil {
			t.Errorf("piece at square failed")
		}

		p2 := NewPieceFromSymbol("Q")
		b.SetPieceAt(H7, &p2, true)
		if !b.baseBoard.promoted.IsMaskingBB(NewBitboardFromSquare(H7)) {
			t.Errorf("piece at square promoted failed")
		}
	})

	t.Run("Test pawn captures", func(t *testing.T) {
		b := NewDefaultBoard()

		// Kings gambit
		var m *Move
		m, _ = NewMoveFromUci("e2e4")
		b.Push(m)
		m, _ = NewMoveFromUci("e7e5")
		b.Push(m)
		m, _ = NewMoveFromUci("f2f4")
		b.Push(m)

		// Accepted
		var found bool
		m, _ = NewMoveFromUci("e5f4")

		found = false
		for move := range b.GeneratePseudoLegalMoves() {
			if *m == move {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("Pawn capture not listed in pseudo legal")
		}

		found = false
		for move := range b.GenerateLegalMoves(BBAll, BBAll) {
			if *m == move {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("Pawn capture not listed in legal")
		}

		b.Push(m)
		if *m != *(b.Pop()) {
			t.Errorf("Popped move not equal")
		}
	})

	t.Run("Pawn move generation", func(t *testing.T) {
		b := NewBoardFromFEN("8/2R1P3/8/2pp4/2k1r3/P7/8/1K6 w - - 1 55", false)

		pseudoLegalMovesCount := 0
		for range b.GeneratePseudoLegalMoves() {
			pseudoLegalMovesCount++
		}

		if pseudoLegalMovesCount != 16 {
			t.Errorf("Pawn moves generation failed")
		}
	})

	t.Run("Single step pawn move", func(t *testing.T) {
		b := NewDefaultBoard()

		a3, _ := NewMoveFromUci("a2a3")

		var found bool
		found = false
		for move := range b.GeneratePseudoLegalMoves() {
			if *a3 == move {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("Pawn capture not listed in pseudo legal")
		}

		found = false
		for move := range b.GenerateLegalMoves(BBAll, BBAll) {
			if *a3 == move {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("Pawn capture not listed in legal")
		}

		b.Push(a3)
		b.Pop()

		if b.FEN(false, "legal", NoPiece) != StartingFEN {
			t.Errorf("Single step pawn move FEN failed")
		}
	})
}
