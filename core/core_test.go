package core

import (
	"testing"
)

// Utility function to check if the given move is in legal moves
func isInLegalMoves(m *Move, b *Board) bool {
	for move := range b.GenerateLegalMoves(BBAll, BBAll) {
		if *m == move {
			return true
		}
	}

	return false
}

// Utility function to check if the given move is in pseudo legal moves
func isInPsuedoLegalMoves(m *Move, b *Board) bool {
	for move := range b.GeneratePseudoLegalMoves() {
		if *m == move {
			return true
		}
	}

	return false
}

// Util function to count the number of moves emitted in the channel
func lenMoveChan(mc chan Move) int {
	count := 0
	for range mc {
		count++
	}
	return count
}

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
		m, _ = NewMoveFromUci("e5f4")
		if !isInPsuedoLegalMoves(m, &b) {
			t.Errorf("Pawn capture not listed in pseudo legal")
		}

		if !isInLegalMoves(m, &b) {
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

		if !isInPsuedoLegalMoves(a3, &b) {
			t.Errorf("Pawn capture not listed in pseudo legal")
		}
		if !isInLegalMoves(a3, &b) {
			t.Errorf("Pawn capture not listed in legal")
		}

		b.Push(a3)
		b.Pop()

		if b.FEN(false, "legal", NoPiece) != StartingFEN {
			t.Errorf("Single step pawn move FEN failed")
		}
	})

	t.Run("Castling", func(t *testing.T) {
		b := NewBoardFromFEN("r3k2r/8/8/8/8/8/8/R3K2R w KQkq - 1 1", false)

		var m, x *Move

		// Let white castle short
		m, _ = b.parseSan("O-O")
		x, _ = NewMoveFromUci("e1g1")
		if *m != *x || b.San(m) != "O-O" || !isInLegalMoves(m, &b) {
			t.Errorf("white castling kingside failed")
		}
		b.Push(m)

		// Let black castle long
		m, _ = b.parseSan("O-O-O")
		if b.San(m) != "O-O-O" || !isInLegalMoves(m, &b) {
			t.Errorf("black castling queenside failed")
		}
		b.Push(m)
		if b.FEN(false, "legal", NoPiece) != "2kr3r/8/8/8/8/8/8/R4RK1 w - - 3 2" {
			t.Errorf("black castling queenside failed")
		}

		// Undo both castling moves
		b.Pop()
		b.Pop()
		if b.FEN(false, "legal", NoPiece) != "r3k2r/8/8/8/8/8/8/R3K2R w KQkq - 1 1" {
			t.Errorf("castling undo failed")
		}

		// Let white castle long
		m, _ = b.parseSan("O-O-O")
		if b.San(m) != "O-O-O" || !isInLegalMoves(m, &b) {
			t.Errorf("white castling queenside failed")
		}
		b.Push(m)

		// Let black castle short
		m, _ = b.parseSan("O-O")
		if b.San(m) != "O-O" || !isInLegalMoves(m, &b) {
			t.Errorf("black castling kingside failed")
		}
		b.Push(m)
		if b.FEN(false, "legal", NoPiece) != "r4rk1/8/8/8/8/8/8/2KR3R w - - 3 2" {
			t.Errorf("black castling queenside failed")
		}

		// Undo both castling moves
		b.Pop()
		b.Pop()
		if b.FEN(false, "legal", NoPiece) != "r3k2r/8/8/8/8/8/8/R3K2R w KQkq - 1 1" {
			t.Errorf("castling undo failed")
		}
	})

	t.Run("960 castling", func(t *testing.T) {
		fen := "3r1k1r/4pp2/8/8/8/8/8/4RKR1 w Gd - 1 1"
		b := NewBoardFromFEN(fen, true)

		var m *Move

		// Let white do the king side swap
		m, _ = b.parseSan("O-O")
		if b.San(m) != "O-O" || m.FromSquare != F1 || m.ToSquare != G1 || !isInLegalMoves(m, &b) {
			t.Errorf("960 castling failed")
		}
		b.Push(m)
		if b.ShredderFEN("legal", NoPiece) != "3r1k1r/4pp2/8/8/8/8/8/4RRK1 b d - 2 1" {
			t.Errorf("960 castling failed shredder FEN match")
		}

		// Black cannot castle kingside
		m, _ = NewMoveFromUci("e8h8")
		if isInLegalMoves(m, &b) {
			t.Errorf("960 castling failed")
		}

		// Let black castle on queenside
		m, _ = b.parseSan("O-O-O")
		if b.San(m) != "O-O-O" || m.FromSquare != F8 || m.ToSquare != D8 || !isInLegalMoves(m, &b) {
			t.Errorf("960 black queenside castling failed")
		}
		b.Push(m)
		if b.ShredderFEN("legal", NoPiece) != "2kr3r/4pp2/8/8/8/8/8/4RRK1 w - - 3 2" {
			t.Errorf("960 fen not matching")
		}

		// Restore initial position
		b.Pop()
		b.Pop()
		if b.ShredderFEN("legal", NoPiece) != fen {
			t.Errorf("960 fen not matching")
		}

		fen = "Qr4k1/4pppp/8/8/8/8/8/R5KR w Hb - 0 1"
		b = NewBoardFromFEN(fen, true)

		// White can just hop the rook over
		m, _ = b.parseSan("O-O")
		if b.San(m) != "O-O" || m.FromSquare != G1 || m.ToSquare != H1 || !isInLegalMoves(m, &b) {
			t.Errorf("960 castling failed")
		}
		b.Push(m)
		if b.ShredderFEN("legal", NoPiece) != "Qr4k1/4pppp/8/8/8/8/8/R4RK1 b b - 1 1" {
			t.Errorf("960 fen not matching")
		}

		// Black can not castle queenside or kingside
		if len(b.generateCastlingMoves(BBAll, BBAll)) != 0 {
			t.Error("960 black should not have any castling moves")
		}

		// Restore initial position
		b.Pop()
		if b.ShredderFEN("legal", NoPiece) != fen {
			t.Errorf("960 fen not matching")
		}
	})

	t.Run("selective castling", func(t *testing.T) {
		b := NewBoardFromFEN("r3k2r/pppppppp/8/8/8/8/PPPPPPPP/R3K2R w KQkq - 0 1", false)

		// King not selected
		if lenMoveChan(b.generateCastlingMoves(BBAll & ^b.baseBoard.kings, BBAll)) != 0 {
			t.Error("King not selected failed")
		}

		// Rook on h1 not selected
		if lenMoveChan(b.generateCastlingMoves(BBAll, BBAll & ^BBH1)) != 1 {
			t.Error("Rook on h1 selected failed")
		}
	})

	t.Run("Castling right not destroyed (bug)", func(t *testing.T) {
		// A rook move from H8 to H1 was only taking whites possible castling rights away.
		b := NewBoardFromFEN("2r1k2r/2qbbpp1/p2pp3/1p3PP1/Pn2P3/1PN1B3/1P3QB1/1K1R3R b k - 0 22", false)
		b.PushSan("Rxh1")
		if b.epd(false, "legal", NoPiece) != "2r1k3/2qbbpp1/p2pp3/1p3PP1/Pn2P3/1PN1B3/1P3QB1/1K1R3r w - -" {
			t.Errorf("fen not matching")
		}
	})

	// TODO: Status
	// t.Run("invalid castling rights", func(t *testing.T) {
	// 	// KQkq is not valid in this standard chess position.
	// 	b := NewBoardFromFEN("1r2k3/8/8/8/8/8/8/R3KR2 w KQkq - 0 1", false)
	// 	if b.Status() != StatusBadCastlingRights ||
	// 		b.FEN(false, "legal", NoPiece) != "1r2k3/8/8/8/8/8/8/R3KR2 w Q - 0 1" ||
	// 		!b.HasQueensideCastlingRights(White) ||
	// 		b.HasKingsideCastlingRights(White) ||
	// 		b.HasQueensideCastlingRights(Black) ||
	// 		b.HasKingsideCastlingRights(Black) {
	// 		t.Error("castling rights failed")
	// 	}

	// 	b = NewBoardFromFEN("4k2r/8/8/8/8/8/8/R1K5 w KQkq - 0 1", true)
	// 	if b.Status() != StatusBadCastlingRights || b.FEN(false, "legal", NoPiece) != "4k2r/8/8/8/8/8/8/R1K5 w Qk - 0 1" {
	// 		t.Error("castling rights failed")
	// 	}

	// 	b = NewBoardFromFEN("1r2k3/8/1p6/8/8/5P2/8/1R2KR2 w KQkq - 0 1", true)
	// 	if b.Status() != StatusBadCastlingRights || b.FEN(false, "legal", NoPiece) != "1r2k3/8/1p6/8/8/5P2/8/1R2KR2 w KQq - 0 1" {
	// 		t.Error("castling rights failed")
	// 	}
	// })

	t.Run("960 different king and rook file", func(t *testing.T) {
		// Theoretically this position (with castling rights) can not be reached
		// with a series of legal moves from one of the 960 starting positions.
		// Decision: We don't care. Neither does Stockfish or lichess.org.
		fen := "1r1k1r2/5p2/8/8/8/8/3N4/R5KR b KQkq - 0 1"
		b := NewBoardFromFEN(fen, true)
		if b.FEN(false, "legal", NoPiece) != fen {
			t.Errorf("fen not matching")
		}
	})

	t.Run("960 prevented castle", func(t *testing.T) {
		b := NewBoardFromFEN("4k3/8/8/1b6/8/8/8/5RKR w KQ - 0 1", true)
		m, _ := NewMoveFromUci("g1f1")
		if b.IsLegal(m) {
			t.Error("expected move to be legal")
		}
	})

	t.Run("Insufficient material", func(t *testing.T) {
		// starting position
		b := NewDefaultBoard()
		if b.IsInsufficientMaterial() {
			t.Error("insufficient material failed")
		}

		// King vs. King + 2 bishops of the same color.
		b = NewBoardFromFEN("k1K1B1B1/8/8/8/8/8/8/8 w - - 7 32", false)
		if !b.IsInsufficientMaterial() {
			t.Error("insufficient material failed")
		}

		// Add bishop of opposite color for the weaker side.
		p := NewPieceFromSymbol("b")
		b.SetPieceAt(B8, &p, false)
		if b.IsInsufficientMaterial() {
			t.Error("insufficient material failed")
		}
	})

	t.Run("Promotion with check", func(t *testing.T) {
		b := NewBoardFromFEN("8/6P1/2p5/1Pqk4/6P1/2P1RKP1/4P1P1/8 w - - 0 1", false)
		m, _ := NewMoveFromUci("g7g8q")
		b.Push(m)
		if !b.IsCheck() || b.FEN(false, "legal", NoPiece) != "6Q1/8/2p5/1Pqk4/6P1/2P1RKP1/4P1P1/8 b - - 0 1" {
			t.Error("promotion with check failed")
		}

		b = NewBoardFromFEN("8/8/8/3R1P2/8/2k2K2/3p4/r7 b - - 0 82", false)
		b.PushSan("d1=Q+")
		if b.FEN(false, "legal", NoPiece) != "8/8/8/3R1P2/8/2k2K2/8/r2q4 w - - 0 83" {
			t.Errorf("promotion with check failed, \n got: %v \n expected: %v", b.FEN(false, "legal", NoPiece), "8/8/8/3R1P2/8/2k2K2/8/r2q4 w - - 0 83")
		}
	})

}
