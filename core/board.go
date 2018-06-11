package core

import (
	"regexp"
	"strconv"
	"strings"

	"fmt"

	"github.com/captainsano/golang-chess/util"
)

const (
	StartingFEN            = "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1"
	FENCastlingRegexString = "^(?:-|[KQABCDEFGH]{0,2}[kqabcdefgh]{0,2})"
)

type BaseBoard struct {
	pawns   Bitboard
	knights Bitboard
	bishops Bitboard
	rooks   Bitboard
	queens  Bitboard
	kings   Bitboard

	promoted Bitboard

	occupiedColor []Bitboard
	occupied      Bitboard
}

const (
	StartingBoardFEN = "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR"
)

func NewBaseBoard(fen string) BaseBoard {
	b := BaseBoard{}
	b.occupiedColor = []Bitboard{BBVoid, BBVoid}

	if fen == "" {
		b.Clear()
	} else if fen == StartingBoardFEN || fen == StartingFEN {
		b.Reset()
	} else {
		b.SetFEN(fen)
	}

	return b
}

func NewBaseBoardFromBaseBoard(bb *BaseBoard) BaseBoard {
	b := BaseBoard{}
	b.pawns = bb.pawns
	b.knights = bb.knights
	b.bishops = bb.bishops
	b.rooks = bb.rooks
	b.queens = bb.queens
	b.kings = bb.kings
	b.promoted = bb.promoted
	b.occupiedColor = []Bitboard{BBVoid, BBVoid}
	copy(b.occupiedColor, bb.occupiedColor)
	b.occupied = bb.occupied

	return b
}

func (b *BaseBoard) Reset() {
	b.pawns = BBRank2 | BBRank7
	b.knights = BBB1 | BBG1 | BBB8 | BBG8
	b.bishops = BBC1 | BBF1 | BBC8 | BBF8
	b.rooks = BBCorners
	b.queens = BBD1 | BBD8
	b.kings = BBE1 | BBE8

	b.promoted = BBVoid

	b.occupiedColor[White] = BBRank1 | BBRank2
	b.occupiedColor[Black] = BBRank7 | BBRank8
	b.occupied = BBRank1 | BBRank2 | BBRank7 | BBRank8
}

func (b *BaseBoard) Clear() {
	b.pawns = BBVoid
	b.knights = BBVoid
	b.bishops = BBVoid
	b.rooks = BBVoid
	b.queens = BBVoid
	b.kings = BBVoid

	b.promoted = BBVoid

	b.occupiedColor[White] = BBVoid
	b.occupiedColor[Black] = BBVoid
	b.occupied = BBVoid
}

func (b *BaseBoard) PieceMask(t PieceType, c Color) Bitboard {
	bb := BBVoid
	switch t {
	case Pawn:
		bb = b.pawns
	case Knight:
		bb = b.knights
	case Bishop:
		bb = b.bishops
	case Rook:
		bb = b.rooks
	case Queen:
		bb = b.queens
	case King:
		bb = b.kings
	}

	return bb & b.occupiedColor[c]
}

func (b *BaseBoard) Pieces(t PieceType, c Color) Bitboard {
	return b.PieceMask(t, c)
}

func (b *BaseBoard) PieceAt(s Square) *Piece {
	t := b.PieceTypeAt(s)
	if t == NoPiece {
		return nil
	}

	mask := NewBitboardFromSquare(s)
	if b.occupiedColor[White].IsMaskingBB(mask) {
		return &Piece{t, White}
	}
	return &Piece{t, Black}
}

func (b *BaseBoard) PieceTypeAt(s Square) PieceType {
	mask := NewBitboardFromSquare(s)

	if !b.occupied.IsMaskingBB(mask) {
		return NoPiece
	} else if b.pawns.IsMaskingBB(mask) {
		return Pawn
	} else if b.knights.IsMaskingBB(mask) {
		return Knight
	} else if b.bishops.IsMaskingBB(mask) {
		return Bishop
	} else if b.rooks.IsMaskingBB(mask) {
		return Rook
	} else if b.queens.IsMaskingBB(mask) {
		return Queen
	} else if b.kings.IsMaskingBB(mask) {
		return King
	}

	return NoPiece
}

func (b *BaseBoard) King(c Color) Square {
	mask := b.occupiedColor[c] & b.kings & ^b.promoted

	if mask != 0 {
		return Square(mask.Msb())
	}

	return SquareNone
}

func (b *BaseBoard) Attacks(s Square) Bitboard {
	mask := NewBitboardFromSquare(s)

	if mask.IsMaskingBB(b.pawns) {
		if mask.IsMaskingBB(b.occupiedColor[White]) {
			return PawnAttacks(s, White)
		}

		return PawnAttacks(s, Black)
	}

	if mask.IsMaskingBB(b.knights) {
		return KnightAttacks(s)
	}

	if mask.IsMaskingBB(b.kings) {
		return KingAttacks(s)
	}

	attacks := BBVoid

	if mask.IsMaskingBB(b.bishops) || mask.IsMaskingBB(b.queens) {
		attacks |= DiagAttacks(s)[DiagMasks(s)&b.occupied]
	}

	if mask.IsMaskingBB(b.rooks) || mask.IsMaskingBB(b.queens) {
		attacks |= (RankAttacks(s)[RankMasks(s)&b.occupied]) | (FileAttacks(s)[FileMasks(s)&b.occupied])
	}

	return attacks
}

func (b *BaseBoard) attackersMask(c Color, s Square, occupied Bitboard) Bitboard {
	rankPieces := RankMasks(s) & occupied
	filePieces := FileMasks(s) & occupied
	diagPieces := DiagMasks(s) & occupied

	queensAndRooks := b.queens | b.rooks
	queensAndBishops := b.queens | b.bishops

	attackers := (KingAttacks(s) & b.kings) |
		(KnightAttacks(s) & b.knights) |
		(RankAttacks(s)[rankPieces] & queensAndRooks) |
		(FileAttacks(s)[filePieces] & queensAndRooks) |
		(DiagAttacks(s)[diagPieces] & queensAndBishops) |
		(PawnAttacks(s, c.Swap()) & b.pawns)

	return attackers & b.occupiedColor[c]
}

func (b *BaseBoard) AttackersMask(c Color, s Square) Bitboard {
	return b.attackersMask(c, s, b.occupied)
}

func (b *BaseBoard) IsAttackedBy(c Color, s Square) bool {
	return b.AttackersMask(c, s) != BBVoid
}

func (b *BaseBoard) Attackers(c Color, s Square) Bitboard {
	return b.AttackersMask(c, s)
}

func (b *BaseBoard) PinMask(c Color, s Square) Bitboard {
	kingSq := b.King(c)
	if kingSq == SquareNone {
		return BBAll
	}

	squareMask := NewBitboardFromSquare(s)

	ks := [][]map[Bitboard]Bitboard{fileAttacks, rankAttacks, diagAttacks}
	vs := []Bitboard{b.rooks | b.queens, b.rooks | b.queens, b.bishops | b.queens}

	for i, _ := range ks {
		attacks, sliders := ks[i], vs[i]

		rays := attacks[kingSq][0]
		if rays.IsMaskingBB(squareMask) {
			snipers := rays & sliders & b.occupiedColor[c.Swap()]
			for sniper := range snipers.ScanReversed() {
				if bbBetween[sniper][kingSq]&(b.occupied|squareMask) == squareMask {
					return bbRays[kingSq][sniper]
				}
			}

			break
		}
	}

	return BBAll
}

func (b *BaseBoard) IsPinned(c Color, s Square) bool {
	return b.PinMask(c, s) != BBAll
}

func (b *BaseBoard) RemovePieceAt(s Square) Piece {
	pt := b.PieceTypeAt(s)
	mask := NewBitboardFromSquare(s)

	var color Color
	if b.occupiedColor[White].IsMaskingBB(mask) {
		color = White
	} else if b.occupiedColor[Black].IsMaskingBB(mask) {
		color = Black
	}

	switch pt {
	case Pawn:
		b.pawns ^= mask
	case Knight:
		b.knights ^= mask
	case Bishop:
		b.bishops ^= mask
	case Rook:
		b.rooks ^= mask
	case Queen:
		b.queens ^= mask
	case King:
		b.kings ^= mask
	default:
		return NewPiece(pt, color)
	}

	b.occupied ^= mask
	b.occupiedColor[White] &= ^mask
	b.occupiedColor[Black] &= ^mask

	b.promoted &= ^mask

	return NewPiece(pt, color)
}

func (b *Board) PieceAt(s Square) *Piece {
	return b.baseBoard.PieceAt(s)
}

func (b *BaseBoard) setPieceAt(s Square, pt PieceType, c Color, promoted bool) {
	b.RemovePieceAt(s)

	mask := NewBitboardFromSquare(s)

	switch pt {
	case Pawn:
		b.pawns |= mask
	case Knight:
		b.knights |= mask
	case Bishop:
		b.bishops |= mask
	case Rook:
		b.rooks |= mask
	case Queen:
		b.queens |= mask
	case King:
		b.kings |= mask
	}

	b.occupied ^= mask
	b.occupiedColor[c] ^= mask

	if promoted {
		b.promoted ^= mask
	}
}

func (b *BaseBoard) SetPieceAt(s Square, p *Piece, promoted bool) {
	if p != nil {
		b.setPieceAt(s, p.Type, p.Color, promoted)
	} else {
		b.RemovePieceAt(s)
	}
}

func (b *BaseBoard) FEN(promoted bool) string {
	builder := []string{}
	empty := 0

	for _, square := range squares180 {
		piece := b.PieceAt(square)

		if piece == nil {
			empty++
		} else {
			if empty > 0 {
				builder = append(builder, strconv.Itoa(empty))
				empty = 0
			}

			builder = append(builder, piece.Symbol())

			if promoted && NewBitboardFromSquare(square).IsMaskingBB(b.promoted) {
				builder = append(builder, "~")
			}
		}

		if NewBitboardFromSquare(square).IsMaskingBB(BBFileH) {
			if empty > 0 {
				builder = append(builder, strconv.Itoa(empty))
				empty = 0
			}

			if square != H1 {
				builder = append(builder, "/")
			}
		}
	}

	return strings.Join(builder, "")
}

func (b *BaseBoard) SetFEN(fen string) {
	fen = strings.TrimSpace(fen)
	if strings.Contains(fen, " ") {
		panic("expected position part of fen, got multiple parts")
	}

	// Ensure the FEN is valid
	rows := strings.Split(fen, "/")
	if len(rows) != 8 {
		panic("expected 8 rows in position part of fen")
	}

	fenNumbers := map[string]bool{"1": true, "2": true, "3": true, "4": true, "5": true, "6": true, "7": true, "8": true}
	fenPieces := map[string]bool{"p": true, "n": true, "b": true, "r": true, "q": true, "k": true}

	// Validate each row.
	for _, row := range rows {
		fieldSum := 0
		previousWasDigit := false
		previousWasPiece := false

		for _, c := range strings.Split(row, "") {
			if _, ok := fenNumbers[c]; ok {
				if previousWasDigit {
					panic("two subsequent digits in position part of fen")
				}

				i, _ := strconv.Atoi(c)
				fieldSum += i
				previousWasDigit = true
				previousWasPiece = false
			} else if c == "~" {
				if !previousWasPiece {
					panic("~ not after piece in position part of fen")
				}
				previousWasDigit = false
				previousWasPiece = false
			} else if _, ok := fenPieces[strings.ToLower(c)]; ok {
				fieldSum++
				previousWasDigit = false
				previousWasPiece = true
			} else {
				panic("invalid character in position part of fen")
			}
		}

		if fieldSum != 8 {
			panic("expected 8 columns per row in position part of fen")
		}
	}

	// Clear the board.
	b.Clear()

	// Put pieces on the board.
	squareIndex := 0
	for _, c := range strings.Split(fen, "") {
		if _, ok := fenNumbers[c]; ok {
			i, _ := strconv.Atoi(c)
			squareIndex += i
		} else if _, ok := fenPieces[strings.ToLower(c)]; ok {
			piece := NewPieceFromSymbol(c)
			b.SetPieceAt(squares180[squareIndex], &piece, false)
			squareIndex++
		} else if c == "~" {
			b.promoted |= NewBitboardFromSquare(squares[squares180[squareIndex-1]])
		}
	}
}

func (b *BaseBoard) PieceMap() map[Square]*Piece {
	result := make(map[Square]*Piece)
	for s := range b.occupied.ScanReversed() {
		p := b.PieceAt(Square(s))
		cp := NewPiece(p.Type, p.Color)
		result[Square(s)] = &cp
	}
	return result
}

func (b *BaseBoard) SetPieceMap(pm map[Square]*Piece) {
	b.Clear()
	for s, p := range pm {
		cp := NewPiece(p.Type, p.Color)
		b.SetPieceAt(s, &cp, false)
	}
}

func (b *BaseBoard) SetChess960Pos(sharnagl int) {
	if sharnagl < 0 || sharnagl > 960 {
		panic("invalid position")
	}

	n, bw := sharnagl/4, sharnagl%4
	n, bb := n/4, n%4
	n, q := n/6, n%6

	n1, n2 := 0, 0
	for n1 = 0; n1 < 4; n1++ {
		n2 := n + (3-n1)*(4-n1)
		if n1 < n2 && 1 <= n2 && n2 <= 4 {
			break
		}
	}

	// Bishops
	bwFile := File(bw*2 + 1)
	bbFile := File(bb * 2)
	b.bishops = (NewBitboardFromFile(bwFile) | NewBitboardFromFile(bbFile)) & BBBackRanks

	// Queens.
	qFile := q
	if util.MinInt(int(bwFile), int(bbFile)) <= qFile {
		qFile += 1
	}
	if util.MaxInt(int(bwFile), int(bbFile)) <= qFile {
		qFile += 1
	}

	b.queens = NewBitboardFromFile(File(qFile)) & BBBackRanks

	used := map[int]bool{int(bwFile): true, int(bbFile): true, qFile: true}

	// Knights.
	b.knights = BBVoid
	for i := 0; i < 8; i++ {
		if !used[i] {
			if n1 == 0 || n2 == 0 {
				b.knights |= NewBitboardFromFile(File(i)) & BBBackRanks
				used[i] = true
			}
			n1--
			n2--
		}
	}

	// RKR.
	for i := 0; i < 8; i++ {
		if !used[i] {
			b.rooks = NewBitboardFromFile(File(i)) & BBBackRanks
			used[i] = true
			break
		}
	}
	for i := 0; i < 8; i++ {
		if !used[i] {
			b.kings = NewBitboardFromFile(File(i)) & BBBackRanks
			used[i] = true
			break
		}
	}
	for i := 0; i < 8; i++ {
		if !used[i] {
			b.rooks = NewBitboardFromFile(File(i)) & BBBackRanks
			used[i] = true
			break
		}
	}

	// Finalize
	b.pawns = BBRank2 | BBRank7
	b.occupiedColor[White] = BBRank1 | BBRank2
	b.occupiedColor[Black] = BBRank7 | BBRank8
	b.occupied = BBRank1 | BBRank2 | BBRank7 | BBRank8
	b.promoted = BBVoid
}

func (b *BaseBoard) Chess960Pos() int {
	if !b.occupiedColor[White].IsMaskingBB(BBRank1 | BBRank2) {
		return -1
	}
	if !b.occupiedColor[Black].IsMaskingBB(BBRank7 | BBRank8) {
		return -1
	}
	if !b.pawns.IsMaskingBB(BBRank2 | BBRank7) {
		return -1
	}
	if b.promoted != BBVoid {
		return -1
	}

	if b.bishops.PopCount() != 4 {
		return -1
	}
	if b.rooks.PopCount() != 4 {
		return -1
	}
	if b.knights.PopCount() != 4 {
		return -1
	}
	if b.queens.PopCount() != 2 {
		return -1
	}
	if b.kings.PopCount() != 2 {
		return -1
	}

	if (BBRank1&b.knights)<<56 != BBRank8&b.knights {
		return -1
	}
	if (BBRank1&b.bishops)<<56 != BBRank8&b.bishops {
		return -1
	}
	if (BBRank1&b.rooks)<<56 != BBRank8&b.rooks {
		return -1
	}
	if (BBRank1&b.queens)<<56 != BBRank8&b.queens {
		return -1
	}
	if (BBRank1&b.kings)<<56 != BBRank8&b.kings {
		return -1
	}

	x := b.bishops & (2 + 8 + 32 + 128)
	if x == BBVoid {
		return -1
	}
	bs1 := x.Lsb() - 1 // 2
	ccPos := bs1
	x = b.bishops & (1 + 4 + 16 + 64)
	if x != BBVoid {
		return -1
	}
	bs2 := x.Lsb() * 2
	ccPos += bs2

	// Algorithm from ChessX, src/database/bitboard.cpp, r2254.
	q := 0
	qf := false
	n0 := 0
	n1 := 0
	n0f := false
	n1f := false
	rf := 0
	n0s := []int{0, 4, 7, 9}

	for square := A1; square < H1+1; square++ {
		bb := NewBitboardFromSquare(square)
		if bb.IsMaskingBB(b.queens) {
			qf = true
		} else if bb.IsMaskingBB(b.rooks) || bb.IsMaskingBB(b.kings) {
			if bb.IsMaskingBB(b.kings) {
				if rf != 1 {
					return -1
				}
			} else {
				rf++
			}

			if !qf {
				q++
			}

			if !n0f {
				n0++
			} else if !n1f {
				n1++
			}
		} else if bb.IsMaskingBB(b.knights) {
			if !qf {
				q++
			}

			if !n0f {
				n0f = true
			} else if !n1f {
				n1f = true
			}
		}
	}

	if n0 < 4 && n1f && qf {
		ccPos += q * 16
		krn := n0s[n0] + n1
		ccPos += krn * 96
		return ccPos
	}

	return -1
}

func (b *BaseBoard) Ascii() string {
	builder := []string{}

	for _, square := range squares180 {
		piece := b.PieceAt(square)

		if piece != nil {
			builder = append(builder, piece.Symbol())
		} else {
			builder = append(builder, ".")
		}

		if NewBitboardFromSquare(square).IsMaskingBB(BBFileH) {
			if square != H1 {
				builder = append(builder, "\n")
			}
		} else {
			builder = append(builder, " ")
		}
	}

	return strings.Join(builder, "")
}

// TODO: Rendering with borders borders is screwed up
func (b *BaseBoard) Unicode(invertColor, borders bool) string {
	builder := []string{}

	for rank := 7; rank >= 0; rank-- {
		if borders {
			builder = append(builder, "  ")
			builder = append(builder, strings.Repeat("-", 17))
			builder = append(builder, "\n")

			builder = append(builder, Rank(rank).Name())
			builder = append(builder, " ")
		}

		for file := 0; file < 8; file++ {
			square := NewSquare(File(file), Rank(rank))

			if borders {
				builder = append(builder, "|")
			} else if file > 0 {
				builder = append(builder, " ")
			}

			piece := b.PieceAt(square)
			if piece != nil {
				builder = append(builder, piece.UnicodeSymbol(invertColor))
			} else {
				builder = append(builder, ".")
			}
		}

		if borders {
			builder = append(builder, "|")
		}

		if borders || rank > 0 {
			builder = append(builder, "\n")
		}
	}

	if borders {
		builder = append(builder, "  ")
		builder = append(builder, strings.Repeat("-", 17))
		builder = append(builder, "\n")
		builder = append(builder, "   a b c d e f g h")
	}

	return strings.Join(builder, "")
}

type BoardState struct {
	BaseBoard

	turn           Color
	castlingRights Bitboard
	epSquare       Square
	halfMoveClock  uint
	fullMoveNumber uint
}

func NewBoardStateFromBoard(b *Board) BoardState {
	bs := BoardState{}

	bs.pawns = b.baseBoard.pawns
	bs.knights = b.baseBoard.knights
	bs.bishops = b.baseBoard.bishops
	bs.rooks = b.baseBoard.rooks
	bs.queens = b.baseBoard.queens
	bs.kings = b.baseBoard.kings

	bs.occupiedColor = []Bitboard{BBVoid, BBVoid}
	bs.occupiedColor[White] = b.baseBoard.occupiedColor[White]
	bs.occupiedColor[Black] = b.baseBoard.occupiedColor[Black]
	bs.occupied = b.baseBoard.occupied

	bs.promoted = b.baseBoard.promoted

	bs.turn = b.turn
	bs.castlingRights = b.castlingRights
	bs.epSquare = b.epSquare
	bs.halfMoveClock = b.halfMoveClock
	bs.fullMoveNumber = b.fullMoveNumber

	return bs
}

type Board struct {
	baseBoard BaseBoard

	aliases     []string
	uciVariant  string
	startingFen string

	// TODO: tbw/tbz

	connectedKings     bool
	oneKing            bool
	capturesCompulsory bool

	chess960 bool

	moveStack []Move
	stack     []BoardState

	turn           Color
	castlingRights Bitboard
	epSquare       Square
	halfMoveClock  uint
	fullMoveNumber uint
}

func NewBoard(chess960 bool) Board {
	return NewBoardFromFEN("8/8/8/8/8/8/8/8 w - - 0 1", chess960)
}

func NewBoardFromFEN(fen string, chess960 bool) Board {
	board := Board{}

	board.aliases = []string{"Standard", "Chess", "Classical", "Normal"}
	board.uciVariant = "chess"
	board.startingFen = StartingFEN
	board.connectedKings = false
	board.oneKing = true
	board.capturesCompulsory = false

	board.chess960 = chess960

	board.moveStack = []Move{}
	board.stack = []BoardState{}

	board.baseBoard = NewBaseBoard(strings.Fields(fen)[0])
	if fen == "" {
		board.Clear()
	} else if fen == StartingFEN {
		board.Reset()
	} else {
		board.SetFEN(fen)
	}

	return board
}

func NewDefaultBoard() Board {
	return NewBoardFromFEN(StartingFEN, false)
}

func NewBoardFromBoard(b *Board) Board {
	board := Board{}

	copy(board.aliases, b.aliases)
	board.uciVariant = b.uciVariant
	board.startingFen = b.startingFen
	board.connectedKings = b.connectedKings
	board.oneKing = b.oneKing
	board.capturesCompulsory = b.capturesCompulsory

	board.chess960 = b.chess960

	copy(board.moveStack, b.moveStack)
	copy(board.stack, b.stack)

	board.baseBoard = NewBaseBoardFromBaseBoard(&b.baseBoard)

	board.turn = b.turn
	board.castlingRights = b.castlingRights
	board.epSquare = b.epSquare
	board.halfMoveClock = b.halfMoveClock
	b.fullMoveNumber = b.fullMoveNumber

	return board
}

func (a *Board) Equal(b *Board) bool {
	// TODO: Improve this with value comparisons
	return a.transpositionKey() == b.transpositionKey()
}

func (b *Board) Turn() Color {
	return b.turn
}

func (b *Board) CastlingRights() Bitboard {
	return b.castlingRights
}

func (b *Board) Reset() {
	b.turn = White
	b.castlingRights = BBCorners
	b.epSquare = SquareNone
	b.halfMoveClock = 0
	b.fullMoveNumber = 1

	b.baseBoard.Reset()
	b.clearStack()
}

func (b *Board) Clear() {
	b.turn = White
	b.castlingRights = BBVoid
	b.epSquare = SquareNone
	b.halfMoveClock = 0
	b.fullMoveNumber = 1

	b.baseBoard.Clear()
	b.clearStack()
}

func (b *Board) clearStack() {
	b.moveStack = []Move{}
	b.stack = []BoardState{}
}

func (b *Board) RemovePieceAt(s Square) Piece {
	piece := b.baseBoard.RemovePieceAt(s)
	b.clearStack()
	return piece
}

func (b *Board) SetPieceAt(s Square, p *Piece, promoted bool) {
	b.baseBoard.SetPieceAt(s, p, promoted)
	b.clearStack()
}

func (b *Board) generatePseudoLegalMoves(fromMask, toMask Bitboard) chan Move {
	ch := make(chan Move)

	go func(b Board) {
		defer close(ch)

		ourPieces := b.baseBoard.occupiedColor[b.turn]

		// Generate piece moves
		nonPawns := ourPieces & ^b.baseBoard.pawns & fromMask
		for fromSquare := range nonPawns.ScanReversed() {
			moves := b.baseBoard.Attacks(Square(fromSquare)) & ^ourPieces & toMask
			for toSquare := range moves.ScanReversed() {
				m, _ := NewNormalMove(Square(fromSquare), Square(toSquare))
				ch <- *m
			}
		}

		// Generate castling moves
		if fromMask.IsMaskingBB(b.baseBoard.kings) {
			for move := range b.generateCastlingMoves(fromMask, toMask) {
				ch <- move
			}
		}

		// The remaining moves are pawn moves
		pawns := b.baseBoard.pawns & b.baseBoard.occupiedColor[b.turn] & fromMask
		if pawns == BBVoid {
			return
		}

		// Generate captures
		capturers := pawns
		for fromSquare := range capturers.ScanReversed() {
			targets := pawnAttacks[b.turn][fromSquare] & b.baseBoard.occupiedColor[b.turn.Swap()] & toMask
			for toSquare := range targets.ScanReversed() {
				if Square(toSquare).Rank() == 0 || Square(toSquare).Rank() == 7 {
					var m *Move
					m, _ = NewPromotionMove(Square(fromSquare), Square(toSquare), Queen)
					ch <- *m
					m, _ = NewPromotionMove(Square(fromSquare), Square(toSquare), Rook)
					ch <- *m
					m, _ = NewPromotionMove(Square(fromSquare), Square(toSquare), Bishop)
					ch <- *m
					m, _ = NewPromotionMove(Square(fromSquare), Square(toSquare), Knight)
					ch <- *m
				} else {
					m, _ := NewNormalMove(Square(fromSquare), Square(toSquare))
					ch <- *m
				}
			}
		}

		// Prepare pawn advance generation
		singleMoves, doubleMoves := BBVoid, BBVoid
		if b.turn == White {
			singleMoves = (pawns << 8) & ^b.baseBoard.occupied
			doubleMoves = (singleMoves << 8) & ^b.baseBoard.occupied & (BBRank3 | BBRank4)
		} else {
			singleMoves = (pawns >> 8) & ^b.baseBoard.occupied
			doubleMoves = (singleMoves >> 8) & ^b.baseBoard.occupied & (BBRank6 | BBRank5)
		}
		singleMoves &= toMask
		doubleMoves &= toMask

		// Generate single pawn moves
		for toSquare := range singleMoves.ScanReversed() {
			fromSquare := Square(toSquare)
			if b.turn == Black {
				fromSquare += 8
			} else {
				fromSquare -= 8
			}

			if Square(toSquare).Rank() == 0 || Square(toSquare).Rank() == 7 {
				var m *Move
				m, _ = NewPromotionMove(fromSquare, Square(toSquare), Queen)
				ch <- *m
				m, _ = NewPromotionMove(fromSquare, Square(toSquare), Rook)
				ch <- *m
				m, _ = NewPromotionMove(fromSquare, Square(toSquare), Bishop)
				ch <- *m
				m, _ = NewPromotionMove(fromSquare, Square(toSquare), Knight)
				ch <- *m
			} else {
				m, _ := NewNormalMove(fromSquare, Square(toSquare))
				ch <- *m
			}
		}

		// Generate double pawn moves
		for toSquare := range doubleMoves.ScanReversed() {
			fromSquare := Square(toSquare)
			if b.turn == Black {
				fromSquare += 16
			} else {
				fromSquare -= 16
			}

			m, _ := NewNormalMove(fromSquare, Square(toSquare))
			ch <- *m
		}

		// Generate enpassant captures
		if b.epSquare != SquareNone {
			for move := range b.generatePseudoLegalEp(fromMask, toMask) {
				ch <- move
			}
		}
	}(NewBoardFromBoard(b))

	return ch
}

func (b *Board) GeneratePseudoLegalMoves() chan Move {
	return b.generatePseudoLegalMoves(BBAll, BBAll)
}

func (b *Board) generatePseudoLegalEp(fromMask, toMask Bitboard) chan Move {
	ch := make(chan Move)

	go func(b Board) {
		defer close(ch)

		if (b.epSquare == SquareNone) || !NewBitboardFromSquare(b.epSquare).IsMaskingBB(toMask) {
			return
		}

		if NewBitboardFromSquare(b.epSquare).IsMaskingBB(b.baseBoard.occupied) {
			return
		}

		capturers := b.baseBoard.pawns & b.baseBoard.occupiedColor[b.turn] & fromMask & PawnAttacks(b.epSquare, b.turn.Swap())
		if b.turn == White {
			capturers &= BBRank5
		} else {
			capturers &= BBRank4
		}

		for capturer := range capturers.ScanReversed() {
			m, _ := NewNormalMove(Square(capturer), b.epSquare)
			ch <- *m
		}
	}(NewBoardFromBoard(b))

	return ch
}

func (b *Board) generatePseudoLegalCaptures(fromMask, toMask Bitboard) chan Move {
	ch := make(chan Move)

	go func(b Board) {
		defer close(ch)

		for m := range b.generatePseudoLegalMoves(fromMask, (toMask & b.baseBoard.occupiedColor[b.turn.Swap()])) {
			ch <- m
		}

		for m := range b.generatePseudoLegalEp(fromMask, toMask) {
			ch <- m
		}
	}(NewBoardFromBoard(b))

	return ch
}

func (b *Board) IsCheck() bool {
	kingSquare := b.baseBoard.King(b.turn)
	return kingSquare != SquareNone && b.baseBoard.IsAttackedBy(b.turn.Swap(), kingSquare)
}

func (b *Board) IsIntoCheck(m *Move) bool {
	kingSquare := b.baseBoard.King(b.turn)
	if kingSquare == SquareNone {
		return false
	}

	checkers := b.baseBoard.AttackersMask(b.turn.Swap(), kingSquare)
	if checkers != BBVoid {
		isIn := false
		for move := range b.generateEvasions(kingSquare, checkers, NewBitboardFromSquare(m.FromSquare), NewBitboardFromSquare(m.ToSquare)) {
			if move == *m {
				isIn = true
				break
			}
		}

		if !isIn {
			return true
		}
	}

	return !b.isSafe(kingSquare, b.sliderBlockers(kingSquare), m)
}

func (b *Board) WasIntoCheck() bool {
	kingSquare := b.baseBoard.King(b.turn.Swap())
	return kingSquare != SquareNone && b.baseBoard.IsAttackedBy(b.turn, kingSquare)
}

func (b *Board) IsPseudoLegal(m *Move) bool {
	// Null moves are not pseudo legal
	if !m.IsNotNull() {
		return false
	}

	// Drop moves are not pseudo legal
	if m.Drop != NoPiece {
		return false
	}

	// Source square must not be vacant
	pieceType := b.baseBoard.PieceTypeAt(m.FromSquare)
	if pieceType == NoPiece {
		return false
	}

	// Get square masks
	fromMask := NewBitboardFromSquare(m.FromSquare)
	toMask := NewBitboardFromSquare(m.ToSquare)

	// check turn
	if !b.baseBoard.occupiedColor[b.turn].IsMaskingBB(fromMask) {
		return false
	}

	// only pawns can promote and only on the back rank
	if m.Promotion != NoPiece {
		if pieceType != Pawn {
			return false
		}

		if b.turn == White && m.ToSquare.Rank() != 7 {
			return false
		} else if b.turn == Black && m.ToSquare.Rank() != 0 {
			return false
		}
	}

	// Handle castling
	if pieceType == King {
		for move := range b.generateCastlingMoves(BBAll, BBAll) {
			if move == *m {
				return true
			}
		}
	}

	// Destination square cannot be occupied
	if b.baseBoard.occupiedColor[b.turn].IsMaskingBB(toMask) {
		return false
	}

	// Handle pawn moves
	if pieceType == Pawn {
		for move := range b.generatePseudoLegalMoves(fromMask, toMask) {
			if move == *m {
				return true
			}
		}

		return false
	}

	return b.baseBoard.Attacks(m.FromSquare)&toMask != BBVoid
}

func (b *Board) IsLegal(m *Move) bool {
	return !b.IsVariantEnd() && b.IsPseudoLegal(m) && !b.IsIntoCheck(m)
}

func (b *Board) IsVariantEnd() bool {
	return false
}

func (b *Board) IsVariantLoss() bool {
	return false
}

func (b *Board) IsVariantWin() bool {
	return false
}

func (b *Board) IsVariantDraw() bool {
	return false
}

func (b *Board) IsGameOver(claimDraw bool) bool {
	// 75 move rule
	if b.IsSeventyFiveMoves() {
		return true
	}

	// Insufficient material
	if b.IsInsufficientMaterial() {
		return true
	}

	// stalemate or checkmate
	hasLegalMoves := 0
	for range b.GenerateLegalMoves(BBAll, BBAll) {
		hasLegalMoves++
		break
	}

	if hasLegalMoves == 0 {
		return true
	}

	// Fivefold repetition
	if b.IsFiveFoldRepetition() {
		return true
	}

	// claim draw
	if claimDraw && b.CanClaimDraw() {
		return true
	}

	return false
}

func (b *Board) Result(claimDraw bool) string {
	// Chess variant support
	if b.IsVariantLoss() {
		if b.turn == White {
			return "0-1"
		}
		return "1-0"
	} else if b.IsVariantWin() {
		if b.turn == White {
			return "1-0"
		}
		return "0-1"
	} else if b.IsVariantDraw() {
		return "1/2-1/2"
	}

	// Checkmate
	if b.IsCheckmate() {
		if b.turn == White {
			return "0-1"
		}
		return "1-0"
	}

	// Draw claimed
	if claimDraw && b.CanClaimDraw() {
		return "1/2-1/2"
	}

	// 75 move rule or fivefold repetition
	if b.IsSeventyFiveMoves() || b.IsFiveFoldRepetition() {
		return "1/2-1/2"
	}

	// Insufficient material
	if b.IsInsufficientMaterial() {
		return "1/2-1/2"
	}

	// Stalemate
	hasLegalMoves := 0
	for range b.GenerateLegalMoves(BBAll, BBAll) {
		hasLegalMoves++
		break
	}

	if hasLegalMoves == 0 {
		return "1/2-1/2"
	}

	// Undetermined
	return "*"
}

func (b *Board) IsCheckmate() bool {
	if !b.IsCheck() {
		return false
	}

	legalMovesCount := 0
	for range b.GenerateLegalMoves(BBAll, BBAll) {
		legalMovesCount++
		break
	}

	return legalMovesCount == 0
}

func (b *Board) IsStalemate() bool {
	if b.IsCheck() {
		return false
	}

	if b.IsVariantEnd() {
		return false
	}

	hasLegalMoves := 0
	for range b.GenerateLegalMoves(BBAll, BBAll) {
		hasLegalMoves++
		break
	}

	return hasLegalMoves == 0
}

func (b *Board) IsInsufficientMaterial() bool {
	// Enough material to mate.
	if b.baseBoard.pawns != BBVoid || b.baseBoard.rooks != BBVoid || b.baseBoard.queens != BBVoid {
		return false
	}

	// A single knight or a single bishop.
	if b.baseBoard.occupied.PopCount() <= 3 {
		return true
	}

	// More than a single knight.
	if b.baseBoard.knights != BBVoid {
		return false
	}

	// All bishops on the same color.
	if b.baseBoard.bishops&BBDarkSquares == BBVoid {
		return true
	} else if b.baseBoard.bishops&BBLightsquares == BBVoid {
		return true
	}

	return false
}

func (b *Board) IsSeventyFiveMoves() bool {
	if b.halfMoveClock >= 150 {
		legalMoves := 0
		for range b.GenerateLegalMoves(BBAll, BBAll) {
			legalMoves++
			break
		}

		if legalMoves != 0 {
			return true
		}
	}

	return false
}

func (b *Board) IsFiveFoldRepetition() bool {
	transpositionKey := b.transpositionKey()

	if len(b.moveStack) < 16 {
		return false
	}

	switchYard := []*Move{}

	for i := 0; i < 4; i++ {
		for j := 0; j < 4; j++ {
			switchYard = append(switchYard, b.Pop())
		}

		if b.transpositionKey() != transpositionKey {
			for i := range switchYard {
				b.Push(switchYard[len(switchYard)-i-1])
			}

			return false
		}
	}

	for i := range switchYard {
		b.Push(switchYard[len(switchYard)-i-1])
	}

	return true
}

func (b *Board) CanClaimDraw() bool {
	return b.CanClaimFiftyMoves() || b.CanClaimThreefoldRepetition()
}

func (b *Board) CanClaimFiftyMoves() bool {
	if b.halfMoveClock >= 100 {
		legalMoves := 0
		for range b.GenerateLegalMoves(BBAll, BBAll) {
			legalMoves++
			break
		}

		if legalMoves != 0 {
			return true
		}
	}

	return false
}

func (b *Board) CanClaimThreefoldRepetition() bool {
	transpositionKey := b.transpositionKey()
	transpositions := map[string]int{transpositionKey: 1}

	// Count positions.
	switchyard := []*Move{}
	for len(b.moveStack) != 0 {
		m := b.Pop()
		switchyard = append(switchyard, m)

		if b.isIrreversible(m) {
			break
		}

		transpositions[b.transpositionKey()]++
	}

	for i := range switchyard {
		b.Push(switchyard[len(switchyard)-i-1])
	}

	// Threefold repetition occured.
	if transpositions[transpositionKey] >= 3 {
		return true
	}

	// The next legal move is a threefold repetition.
	for move := range b.GenerateLegalMoves(BBAll, BBAll) {
		b.Push(&move)

		if transpositions[b.transpositionKey()] >= 2 {
			b.Pop()
			return true
		}

		b.Pop()
	}

	return false
}

func (b *Board) pushCapture(m *Move, captureSquare Square, pt PieceType, wasPromoted bool) {
	// Noop
}

func (b *Board) Push(move *Move) {
	b.stack = append(b.stack, NewBoardStateFromBoard(b)) // Capture the board state
	b.moveStack = append(b.moveStack, *move)             // TODO: Make a defensive copy

	m := b.toChess960(move)

	// Reset en passant square
	epSquare := b.epSquare
	b.epSquare = SquareNone

	// Increment move counters
	b.halfMoveClock++
	if b.turn == Black {
		b.fullMoveNumber++
	}

	// On a null move, simply swap the turn
	if !m.IsNotNull() {
		b.turn = b.turn.Swap()
		return
	}

	// Drops
	if m.Drop != NoPiece {
		p := NewPiece(m.Drop, b.turn)
		b.baseBoard.SetPieceAt(m.ToSquare, &p, false)
		b.turn = b.turn.Swap()
		return
	}

	// Zero the half move clock
	if b.isZeroing(m) {
		b.halfMoveClock = 0
	}

	fromMask := NewBitboardFromSquare(m.FromSquare)
	toMask := NewBitboardFromSquare(m.ToSquare)

	promoted := b.baseBoard.promoted.IsMaskingBB(fromMask)
	piece := b.baseBoard.RemovePieceAt(m.FromSquare)
	captureSquare := m.ToSquare
	capturedPieceType := b.baseBoard.PieceTypeAt(m.ToSquare)

	// Update castling rights
	b.castlingRights = b.CleanCastlingRights() & ^toMask & ^fromMask
	if piece.Type == King && !promoted {
		if b.turn == White {
			b.castlingRights &= ^BBRank1
		} else {
			b.castlingRights &= ^BBRank8
		}
	} else if capturedPieceType == King && !b.baseBoard.promoted.IsMaskingBB(toMask) {
		if b.turn == White && m.ToSquare.Rank() == 7 {
			b.castlingRights &= ^BBRank8
		} else if b.turn == Black && m.ToSquare.Rank() == 0 {
			b.castlingRights &= ^BBRank1
		}
	}

	// Handle special pawn moves
	if piece.Type == Pawn {
		diff := int(m.ToSquare - m.FromSquare)

		if diff == 16 && m.FromSquare.Rank() == 1 {
			b.epSquare = m.FromSquare + 8
		} else if diff == -16 && m.FromSquare == 6 {
			b.epSquare = m.FromSquare - 8
		} else if m.ToSquare == epSquare && (util.AbsInt(diff) == 7 || util.AbsInt(diff) == 9) && capturedPieceType == NoPiece {
			// Remove pawns captured en passant
			if b.turn == White {
				captureSquare := Square(epSquare - 8)
				capturedPieceType = b.baseBoard.RemovePieceAt(captureSquare).Type
			} else {
				captureSquare := Square(epSquare + 8)
				capturedPieceType = b.baseBoard.RemovePieceAt(captureSquare).Type
			}
		}
	}

	// Promotion
	if m.Promotion != NoPiece {
		promoted = true
		piece.Type = m.Promotion
	}

	// Castling
	castling := BBVoid
	if piece.Type == King {
		castling = b.baseBoard.occupiedColor[b.turn] & toMask
	}

	if castling != BBVoid {
		aSide := m.ToSquare.File() < m.FromSquare.File()

		b.baseBoard.RemovePieceAt(m.FromSquare)
		b.baseBoard.RemovePieceAt(m.ToSquare)

		if aSide {
			if b.turn == White {
				b.baseBoard.setPieceAt(C1, King, b.turn, false)
				b.baseBoard.setPieceAt(D1, Rook, b.turn, false)
			} else {
				b.baseBoard.setPieceAt(C8, King, b.turn, false)
				b.baseBoard.setPieceAt(D8, Rook, b.turn, false)
			}
		} else {
			if b.turn == White {
				b.baseBoard.setPieceAt(G1, King, b.turn, false)
				b.baseBoard.setPieceAt(F1, Rook, b.turn, false)
			} else {
				b.baseBoard.setPieceAt(G8, King, b.turn, false)
				b.baseBoard.setPieceAt(F8, Rook, b.turn, false)
			}
		}
	}

	// Put the piece on the target square
	if castling == BBVoid && piece.Type != NoPiece {
		wasPromoted := b.baseBoard.promoted.IsMaskingBB(toMask)
		b.baseBoard.setPieceAt(m.ToSquare, piece.Type, b.turn, wasPromoted)

		if capturedPieceType != NoPiece {
			b.pushCapture(m, captureSquare, capturedPieceType, wasPromoted)
		}
	}

	// Swap turn
	b.turn = b.turn.Swap()
}

func (b *Board) Pop() *Move {
	var move Move
	move, b.moveStack = b.moveStack[len(b.moveStack)-1], b.moveStack[:len(b.moveStack)-1]
	var state BoardState
	state, b.stack = b.stack[len(b.stack)-1], b.stack[:len(b.stack)-1]

	b.baseBoard.pawns = state.pawns
	b.baseBoard.knights = state.knights
	b.baseBoard.bishops = state.bishops
	b.baseBoard.rooks = state.rooks
	b.baseBoard.queens = state.queens
	b.baseBoard.kings = state.kings

	b.baseBoard.occupiedColor[White] = state.occupiedColor[White]
	b.baseBoard.occupiedColor[Black] = state.occupiedColor[Black]
	b.baseBoard.occupied = state.occupied

	b.baseBoard.promoted = state.promoted

	b.turn = state.turn
	b.castlingRights = state.castlingRights
	b.epSquare = state.epSquare
	b.halfMoveClock = state.halfMoveClock
	b.fullMoveNumber = state.fullMoveNumber

	return &move
}

func (b *Board) Peek() *Move {
	return &b.moveStack[len(b.moveStack)-1]
}

func (b *Board) CastlingShredderFen() string {
	castlingRights := b.CleanCastlingRights()
	if castlingRights == BBVoid {
		return "-"
	}

	builder := []string{}

	for s := range (castlingRights & BBRank1).ScanReversed() {
		builder = append(builder, strings.ToUpper(Square(s).File().Name()))
	}

	for s := range (castlingRights & BBRank8).ScanReversed() {
		builder = append(builder, Square(s).File().Name())
	}

	return strings.Join(builder, "")
}

func (b *Board) CastlingXFEN() string {
	builder := []string{}

	for _, color := range []Color{White, Black} {
		kingSquare := b.baseBoard.King(color)
		if kingSquare == SquareNone {
			continue
		}

		kingFile := kingSquare.File()
		backRank := BBRank8
		if color == White {
			backRank = BBRank1
		}

		for rookSquare := range (b.CleanCastlingRights() & backRank).ScanReversed() {
			rookFile := Square(rookSquare).File()
			aSide := rookFile < kingFile

			otherRooks := b.baseBoard.occupiedColor[color] & b.baseBoard.rooks & backRank & ^NewBitboardFromSquare(Square(rookSquare))

			ch := "k"
			if aSide {
				ch = "q"
			}

			for other := range otherRooks.ScanReversed() {
				if (Square(other).File() < rookFile) == aSide {
					ch = rookFile.Name()
					break
				}
			}

			if color == White {
				builder = append(builder, strings.ToUpper(ch))
			} else {
				builder = append(builder, ch)
			}
		}
	}

	if len(builder) > 0 {
		return strings.Join(builder, "")
	}

	return "-"
}

func (b *Board) hasPseudoLegalEnPassant() bool {
	if b.epSquare != SquareNone {
		for range b.generatePseudoLegalEp(BBAll, BBAll) {
			return true
		}
	}

	return false
}

func (b *Board) hasLegalEnPassant() bool {
	if b.epSquare != SquareNone {
		for range b.generateLegalEp(BBAll, BBAll) {
			return true
		}
	}

	return false
}

func (b *Board) FEN(shredder bool, enPassant string, promoted PieceType) string {
	return fmt.Sprintf("%s %d %d", b.epd(shredder, enPassant, promoted), b.halfMoveClock, b.fullMoveNumber)
}

func (b *Board) ShredderFEN(enPassant string, promoted PieceType) string {
	return fmt.Sprintf("%s %d %d", b.epd(true, enPassant, promoted), b.halfMoveClock, b.fullMoveNumber)
}

func (b *Board) SetFEN(fen string) {
	parts := strings.Fields(fen)
	if len(parts) != 6 {
		panic("FEN string should consist of 6 parts")
	}

	if !(parts[1] == "w" || parts[1] == "b") {
		panic("Expected 'w' or 'b' for turn part of fen")
	}

	if parts[3] != "-" {
		sq := NewSquareFromName(parts[3])
		if sq == SquareNone {
			panic("Invalid enpassant square name")
		}
	}

	halfMoveClock, err := strconv.Atoi(parts[4])
	if err != nil || halfMoveClock < 0 {
		panic("Halfmove clock invalid or cannot be negative")
	}

	fullMoveNumber, err := strconv.Atoi(parts[5])
	if err != nil || fullMoveNumber < 0 {
		panic("fullmove number invalid or cannot be negative")
	}

	b.baseBoard.SetFEN(parts[0])

	// set turn
	if parts[1] == "w" {
		b.turn = White
	} else {
		b.turn = Black
	}

	b.setCastlingFEN(parts[2])

	if parts[3] == "-" {
		b.epSquare = SquareNone
	} else {
		b.epSquare = NewSquareFromName(parts[3])
	}

	b.halfMoveClock = uint(halfMoveClock)
	b.fullMoveNumber = uint(fullMoveNumber)

	b.clearStack()
}

func (b *Board) setCastlingFEN(castlingFen string) {
	if len(castlingFen) == 0 || castlingFen == "-" {
		b.castlingRights = BBVoid
		return
	}

	if matches, _ := regexp.MatchString(FENCastlingRegexString, castlingFen); !matches {
		panic("Invalid castling fen")
	}

	b.castlingRights = BBVoid

	for _, flag := range strings.Split(castlingFen, "") {
		color := Black
		if flag == strings.ToUpper(flag) {
			color = White
		}

		flag = strings.ToLower(flag)
		backRank := BBRank8
		if color == White {
			backRank = BBRank1
		}

		rooks := b.baseBoard.occupiedColor[color] & b.baseBoard.rooks & backRank
		kingSquare := b.baseBoard.King(color)

		if flag == "q" {
			if kingSquare != SquareNone && rooks.Lsb() < int(kingSquare) {
				b.castlingRights |= rooks & -rooks
			} else {
				b.castlingRights |= BBFileA & backRank
			}
		} else if flag == "k" {
			rook := rooks.Msb()
			if kingSquare != SquareNone && int(kingSquare) < rook {
				b.castlingRights |= NewBitboardFromSquare(Square(rook))
			} else {
				b.castlingRights |= BBFileH & backRank
			}
		} else {
			b.castlingRights |= NewBitboardFromFile(FileFromName(flag)) & backRank
		}
	}

	b.clearStack()
}

func (b *Board) SetBoardFEN(fen string) {
	b.baseBoard.SetFEN(fen)
	b.clearStack()
}

func (b *Board) SetPieceMap(pm map[Square]*Piece) {
	b.baseBoard.SetPieceMap(pm)
	b.clearStack()
}

func (b *Board) SetChess960Pos(sharnagl int) {
	b.baseBoard.SetChess960Pos(sharnagl)
	b.chess960 = true
	b.turn = White
	b.castlingRights = b.baseBoard.rooks
	b.epSquare = SquareNone
	b.halfMoveClock = 0
	b.fullMoveNumber = 1

	b.clearStack()
}

func (b *Board) Chess960Pos(ignoreTurn, ignoreCastling, ignoreCounters bool) int {
	if b.epSquare != SquareNone {
		return -1
	}

	if !ignoreTurn {
		if b.turn != White {
			return -1
		}
	}

	if !ignoreCastling {
		if !b.CleanCastlingRights().IsMaskingBB(b.baseBoard.rooks) {
			return -1
		}
	}

	if !ignoreCounters {
		if b.fullMoveNumber != 1 || b.halfMoveClock != 0 {
			return -1
		}
	}

	return b.baseBoard.Chess960Pos()
}

// func (b *Board) epdOperations(operations []struct {
// 	opcode  string
// 	operand *interface{}
// }) string {
// 	epd := []string{}
// 	firstOp := true
// 	for _, op := range operations {
// 		if !firstOp {
// 			epd = append(epd, "")
// 		}

// 		firstOp = false

// 		if op.operand == nil {
// 			epd = append(epd, ";")
// 			continue
// 		}

// 		opcode, operand := op.opcode, *op.operand

// 		// Value is empty
// 		if operand == nil {
// 			epd = append(epd, ";")
// 			continue
// 		}

// 		// Value is a move
// 		operandMap, ok := operand.(map[string]string)
// 		if ok {
// 			from, fok := operandMap["from_square"]
// 			to, tok := operandMap["to_square"]
// 			promotion, pok := operandMap["promotion"]
// 			if fok && tok && pok {
// 				epd = append(epd, " ")
// 				epd = append(epd, b.San(operandMap))
// 				epd = append(epd, ";")
// 				continue
// 			}
// 		}

// 		operandStr, ok := operand.(string)
// 		if ok {
// 			// Value is int
// 			epdInt, err := strconv.ParseInt(operandStr, 10, 0)
// 			if err == nil {
// 				epd = append(epd, " ")
// 				epd = append(epd, operandStr)
// 				epd = append(epd, ";")
// 				continue
// 			}

// 			// Value is float
// 			epdFloat, err := strconv.ParseFloat(operandStr, 0)
// 			if err == nil {
// 				epd = append(epd, " ")
// 				epd = append(epd, operandStr)
// 				epd = append(epd, ";")
// 				continue
// 			}
// 		}

// 		// value is a slice of moves
// 		operandMoves, ok := operand.([]Move)
// 		if ok {
// 			position := b
// 			if opcode == "pv" {
// 				_b := NewBoard(b.ShredderFEN("legal", NoPiece), false)
// 				position = &_b
// 			}

// 			firstMove := operandMoves[0]
// 			if firstMove.FromSquare != SquareNone && firstMove.ToSquare != SquareNone && firstMove.Promotion != NoPiece {
// 				epd = append(epd, " ")
// 				epd = append(epd, position.San(firstMove))
// 				if opcode == "pv" {
// 					position.Push(firstMove)
// 				}

// 				for _, move := range operandMoves[1:] {
// 					epd = append(epd, " ")
// 					epd = append(epd, position.San(move))
// 					if opcode == "pv" {
// 						position.Push(move)
// 					}
// 				}
// 			}

// 			epd = append(epd, ";")
// 			continue
// 		}

// 		operandStr, ok = operand.(string)
// 		if ok {
// 			epd = append(epd, " \"")
// 			epd = append(epd, strings.Replace(strings.Replace(strings.Replace(strings.Replace(operandStr, "\r", "", -1), "\n", " ", -1), "\\", "\\\\", -1), ";", "\\s", -1))
// 			epd = append(epd, "\";")
// 		}
// 	}

// 	return strings.Join(epd, "")
// }

func (b *Board) epd(shredder bool, enPassant string, promoted PieceType, epdOperations ...interface{}) string {
	epd := []string{}

	epd = append(epd, b.baseBoard.FEN(promoted != NoPiece))
	if b.turn == White {
		epd = append(epd, "w")
	} else {
		epd = append(epd, "b")
	}
	if shredder {
		epd = append(epd, b.CastlingShredderFen())
	} else {
		epd = append(epd, b.CastlingXFEN())
	}

	if enPassant == "fen" {
		if b.epSquare != SquareNone {
			epd = append(epd, b.epSquare.Name())
		} else {
			epd = append(epd, "-")
		}
	} else if enPassant == "xfen" {
		if b.hasPseudoLegalEnPassant() {
			epd = append(epd, b.epSquare.Name())
		} else {
			epd = append(epd, "-")
		}
	} else {
		if b.hasLegalEnPassant() {
			epd = append(epd, b.epSquare.Name())
		} else {
			epd = append(epd, "-")
		}
	}

	if len(epdOperations) > 0 {
		// @TODO: EPD Operations
	}

	return strings.Join(epd, " ")
}

// TODO: Parse EPD ops

// TODO: Set EPD

func (b *Board) San(move *Move) string {
	return b.algebraic(move, false)
}

func (b *Board) VariationSan(variation []Move) (string, error) {
	board := NewBoardFromBoard(b)
	san := []string{}

	for _, m := range variation {
		if !board.IsLegal(&m) {
			return "", &MoveError{description: "Illegal move " + m.Uci()}
		}

		if board.turn == White {
			san = append(san, fmt.Sprintf("%d. %s", board.fullMoveNumber, board.San(&m)))
		} else if len(san) == 0 {
			san = append(san, fmt.Sprintf("%d...%s", board.fullMoveNumber, board.San(&m)))
		} else {
			san = append(san, board.San(&m))
		}

		board.Push(&m)
	}

	return strings.Join(san, " "), nil
}

func (b *Board) Lan(move *Move) string {
	return b.algebraic(move, true)
}

func (b *Board) algebraic(move *Move, long bool) string {
	if !move.IsNotNull() {
		return "--"
	}

	san := ""

	// Look ahead for check/checkmate
	b.Push(move)
	isCheck := b.IsCheck()
	isCheckmate := (isCheck && b.IsCheckmate()) || b.IsVariantLoss() || b.IsVariantWin()
	b.Pop()

	// Drops
	if move.Drop != NoPiece {
		san = ""
		if move.Drop != Pawn {
			san = strings.ToUpper(move.Drop.Symbol())
		}
		san += "@" + move.ToSquare.Name()
	}

	// Castling.
	if b.IsCastling(move) {
		if move.ToSquare < move.FromSquare {
			san = "O-O-O"
		} else {
			san = "O-O"
		}
	}

	if move.Drop != NoPiece || b.IsCastling(move) {
		if isCheckmate {
			return san + "#"
		}
		if isCheck {
			return san + "+"
		}
		return san
	}

	piece := b.baseBoard.PieceTypeAt(move.FromSquare)
	capture := b.IsCapture(move)

	if piece == Pawn {
		san = ""
	} else {
		san = strings.ToUpper(piece.Symbol())
	}

	if long {
		san += move.FromSquare.Name()
	} else if piece != Pawn {
		// Get ambiguous move candidates.
		// Relevant candidates: not exactly the current move,
		// but to the same square.
		others := BBVoid
		fromMask := b.baseBoard.PieceMask(piece, b.turn)
		fromMask &= ^NewBitboardFromSquare(move.FromSquare)
		toMask := NewBitboardFromSquare(move.ToSquare)
		for candidate := range b.GenerateLegalMoves(fromMask, toMask) {
			others |= NewBitboardFromSquare(candidate.FromSquare)
		}

		// Disambiguate.
		if others != BBVoid {
			row, column := false, false

			if others.IsMaskingBB(NewBitboardFromRank(move.FromSquare.Rank())) {
				column = true
			}

			if others.IsMaskingBB(NewBitboardFromFile(move.FromSquare.File())) {
				row = true
			} else {
				column = true
			}

			if column {
				san += move.FromSquare.File().Name()
			}
			if row {
				san += move.FromSquare.Rank().Name()
			}
		}
	} else if capture {
		san += move.FromSquare.File().Name()
	}

	// Captures
	if capture {
		san += "x"
	} else if long {
		san += "-"
	}

	// Destination square
	san += move.ToSquare.Name()

	// Promotion
	if move.Promotion != NoPiece {
		san += "=" + strings.ToUpper(move.Promotion.Symbol())
	}

	// Add check or checkmate suffix.
	if isCheckmate {
		san += "#"
	} else if isCheck {
		san += "+"
	}

	return san
}

var sanRegexp *regexp.Regexp
var fenCastlingRegexp *regexp.Regexp

func init() {
	sanRegexp = regexp.MustCompile("^([NBKRQ])?([a-h])?([1-8])?[\\-x]?([a-h][1-8])(=?[nbrqkNBRQK])?(\\+|#)?\\z")
	fenCastlingRegexp = regexp.MustCompile("^(?:-|[KQABCDEFGH]{0,2}[kqabcdefgh]{0,2})\\z")
}

type SanParseError struct {
	error
	description string
}

func (e SanParseError) Error() string { return e.description }

func (b *Board) parseSan(san string) (*Move, error) {
	// Castling
	if _, ok := (map[string]bool{"O-O": true, "O-O+": true, "O-O#": true})[san]; ok {
		for m := range b.generateCastlingMoves(BBAll, BBAll) {
			if b.IsKingsideCastling(&m) {
				return &m, nil
			}
		}

		return nil, SanParseError{description: "Invalid kingside castling expression"}
	} else if _, ok := (map[string]bool{"O-O-O": true, "O-O-O+": true, "O-O-O#": true})[san]; ok {
		for m := range b.generateCastlingMoves(BBAll, BBAll) {
			if b.IsQueensideCastling(&m) {
				return &m, nil
			}
		}

		return nil, SanParseError{description: "Invalid queenside castling expression"}
	}

	// Match normal moves
	match := sanRegexp.MatchString(san)
	if !match {
		// Null moves
		if san == "--" || san == "Z0" {
			return NewNullMove()
		}

		return nil, SanParseError{description: "Invalid san " + san}
	}

	// Get target square
	matches := sanRegexp.FindStringSubmatch(san)
	toSquare := matches[4]
	toMask := NewBitboardFromSquare(NewSquareFromName(toSquare))

	// Get the promotion type
	promotion := NoPiece
	if len(matches) > 4 && len(matches[5]) > 0 {
		p := strings.ToLower(matches[5][len(matches[5])-1:])
		promotion = NewPieceFromSymbol(p).Type
	}

	// Filter by piece type
	pieceType := NoPiece
	fromMask := b.baseBoard.pawns
	if len(matches) > 0 && len(matches[1]) > 0 {
		pieceType = NewPieceFromSymbol(strings.ToLower(matches[1])).Type
		fromMask = b.baseBoard.PieceMask(pieceType, b.turn)
	}

	// Filter by source file
	if len(matches) > 1 && len(matches[2]) > 0 {
		fromMask &= NewBitboardFromFile(FileFromName(matches[2]))
	}

	// Filter by source rank
	if len(matches) > 2 && len(matches[3]) > 0 {
		fromMask &= NewBitboardFromRank(RankFromName(matches[3]))
	}

	// Match legal moves
	m, _ := NewNullMove()
	matchedMove := *m
	for move := range b.GenerateLegalMoves(fromMask, toMask) {

		if move.Promotion != promotion {
			continue
		}

		if matchedMove.FromSquare != SquareNone || matchedMove.ToSquare != SquareNone {
			return nil, SanParseError{description: "Ambiguous SAN " + san + " " + b.FEN(false, "legal", NoPiece)}
		}

		matchedMove = move
	}

	if !matchedMove.IsNotNull() {
		return nil, SanParseError{description: "Illegal SAN " + san + " " + b.FEN(false, "legal", NoPiece)}
	}

	return &matchedMove, nil
}

func (b *Board) PushSan(san string) (*Move, error) {
	move, err := b.parseSan(san)
	if err != nil {
		return nil, err
	}

	b.Push(move)
	return move, nil
}

func (b *Board) parseUci(uci string) (*Move, error) {
	move, err := NewMoveFromUci(uci)
	if err != nil {
		return nil, err
	}

	move = b.toChess960(move)
	move = b.fromChess960(b.chess960, move.FromSquare, move.ToSquare, move.Promotion, move.Drop)

	if !b.IsLegal(move) {
		return nil, &MoveError{description: "Illegal uci " + uci}
	}

	return move, nil
}

func (b *Board) PushUci(uci string) (*Move, error) {
	move, err := b.parseUci(uci)
	if err != nil {
		return nil, err
	}

	b.Push(move)
	return move, nil
}

func (b *Board) Ascii() string {
	// TODO: Add other FEN params
	return b.baseBoard.Ascii()
}

func (b *Board) Unicode(invertColor, borders bool) string {
	// TODO: Add other FEN params
	return b.baseBoard.Unicode(invertColor, borders)
}

func (b *Board) sliderBlockers(kingSquare Square) Bitboard {
	rooks_and_queens := b.baseBoard.rooks | b.baseBoard.queens
	bishops_and_queens := b.baseBoard.bishops | b.baseBoard.queens

	snipers := ((rankAttacks[kingSquare][0] & rooks_and_queens) |
		(fileAttacks[kingSquare][0] & rooks_and_queens) |
		(diagAttacks[kingSquare][0] & bishops_and_queens))

	blockers := BBVoid

	for sniper := range (snipers & b.baseBoard.occupiedColor[b.turn.Swap()]).ScanReversed() {
		b := bbBetween[kingSquare][sniper] & b.baseBoard.occupied

		// Add to blockers if exactly one piece in-between.
		if (b & NewBitboardFromSquare(Square(b.Msb()))) == b {
			blockers |= b
		}
	}

	return blockers & b.baseBoard.occupiedColor[b.turn]
}

func (b *Board) CleanCastlingRights() Bitboard {
	if len(b.stack) > 0 {
		return b.castlingRights
	}

	castling := b.castlingRights & b.baseBoard.rooks
	whiteCastling := castling & BBRank1 & b.baseBoard.occupiedColor[White]
	blackCastling := castling & BBRank8 & b.baseBoard.occupiedColor[Black]

	if !b.chess960 {
		whiteCastling &= BBA1 | BBH1
		blackCastling &= BBA8 | BBH8

		if (b.baseBoard.occupiedColor[White] & b.baseBoard.kings & ^b.baseBoard.promoted & BBE1) == BBVoid {
			whiteCastling = BBVoid
		}
		if (b.baseBoard.occupiedColor[Black] & b.baseBoard.kings & ^b.baseBoard.promoted & BBE8) == BBVoid {
			blackCastling = BBVoid
		}

		return whiteCastling | blackCastling
	}

	// Thinks must be on the back rank
	whiteKingMask := b.baseBoard.occupiedColor[White] & b.baseBoard.kings & BBRank1 & ^b.baseBoard.promoted
	blackKingMask := b.baseBoard.occupiedColor[Black] & b.baseBoard.kings & BBRank8 & ^b.baseBoard.promoted
	if whiteKingMask == BBVoid {
		whiteCastling = BBVoid
	}
	if blackKingMask == BBVoid {
		blackCastling = BBVoid
	}

	// There are only two ways of castling, a-side and h-side and the king must be between the rooks
	whiteASide := whiteCastling & -whiteCastling
	whiteHSide := BBVoid
	if whiteCastling != BBVoid {
		whiteHSide = NewBitboardFromSquare(Square(whiteCastling.Msb()))
	}

	if whiteASide != BBVoid && whiteASide.Msb() > whiteKingMask.Msb() {
		whiteASide = BBVoid
	}
	if whiteHSide != BBVoid && whiteHSide.Msb() < whiteKingMask.Msb() {
		whiteHSide = BBVoid
	}

	blackASide := blackCastling & -blackCastling
	blackHSide := BBVoid
	if blackCastling != BBVoid {
		blackHSide = NewBitboardFromSquare(Square(blackCastling.Msb()))
	}

	if blackASide != BBVoid && blackASide.Msb() > blackKingMask.Msb() {
		blackASide = BBVoid
	}
	if blackHSide != BBVoid && blackHSide.Msb() < blackKingMask.Msb() {
		blackHSide = BBVoid
	}

	return blackASide | blackHSide | whiteASide | whiteHSide
}

func (b *Board) HasCastlingRights(c Color) bool {
	backrank := BBRank1
	if c == Black {
		backrank = BBRank8
	}
	return b.CleanCastlingRights().IsMaskingBB(backrank)
}

func (b *Board) HasKingsideCastlingRights(c Color) bool {
	backrank := BBRank1
	if c == Black {
		backrank = BBRank8
	}
	kingMask := b.baseBoard.kings & b.baseBoard.occupiedColor[c] & backrank & ^b.baseBoard.promoted
	if kingMask == BBVoid {
		return false
	}

	castlingRights := b.CleanCastlingRights() & backrank
	for castlingRights != BBVoid {
		rook := castlingRights & -castlingRights

		if rook > kingMask {
			return true
		}

		castlingRights = castlingRights & (castlingRights - 1)
	}

	return false
}

func (b *Board) HasQueensideCastlingRights(c Color) bool {
	backrank := BBRank1
	if c == Black {
		backrank = BBRank8
	}
	kingMask := b.baseBoard.kings & b.baseBoard.occupiedColor[c] & backrank & ^b.baseBoard.promoted
	if kingMask == BBVoid {
		return false
	}

	castlingRights := b.CleanCastlingRights() & backrank
	for castlingRights != BBVoid {
		rook := castlingRights & -castlingRights

		if rook < kingMask {
			return true
		}

		castlingRights = castlingRights & (castlingRights - 1)
	}

	return false
}

func (b *Board) HasChess960CastlingRights() bool {
	chess960 := b.chess960
	b.chess960 = true
	castlingRights := b.CleanCastlingRights()
	b.chess960 = chess960

	// Standard chess castling rights can only be on the standard
	// starting rook squares.
	if castlingRights.IsMaskingBB(^BBCorners) {
		return true
	}

	// If there are any castling rights in standard chess, the king must be
	// on e1 or e8.
	if castlingRights.IsMaskingBB(BBRank1) && (b.baseBoard.occupiedColor[White]&b.baseBoard.kings&BBE1) == BBVoid {
		return true
	}

	if castlingRights.IsMaskingBB(BBRank8) && (b.baseBoard.occupiedColor[Black]&b.baseBoard.kings&BBE8) == BBVoid {
		return true
	}

	return false
}

// TODO: Status
func (b *Board) Status() uint64 {
	return 0
}

func (b *Board) isSafe(king Square, blockers Bitboard, move *Move) bool {
	if move.FromSquare == king {
		if b.IsCastling(move) {
			return true
		}

		return !b.baseBoard.IsAttackedBy(b.turn.Swap(), move.ToSquare)
	} else if b.IsEnPassant(move) {
		return b.baseBoard.PinMask(b.turn, move.FromSquare).IsMaskingBB(NewBitboardFromSquare(move.ToSquare)) && !b.epSkewered(king, move.FromSquare)
	}

	return !blockers.IsMaskingBB(NewBitboardFromSquare(move.FromSquare)) || bbRays[move.FromSquare][move.ToSquare].IsMaskingBB(NewBitboardFromSquare(king))
}

func (b *Board) epSkewered(kingSquare, capturer Square) bool {
	lastDouble := b.epSquare
	if b.turn == White {
		lastDouble -= 8
	} else {
		lastDouble += 8
	}

	occupancy := (b.baseBoard.occupied & ^NewBitboardFromSquare(lastDouble) & ^NewBitboardFromSquare(capturer) | NewBitboardFromSquare(b.epSquare))

	horizontalAttackers := b.baseBoard.occupiedColor[b.turn.Swap()] & (b.baseBoard.rooks | b.baseBoard.queens)
	if rankAttacks[kingSquare][(rankMasks[kingSquare] & occupancy)].IsMaskingBB(horizontalAttackers) {
		return true
	}

	diagonalAttackers := b.baseBoard.occupiedColor[b.turn.Swap()] & (b.baseBoard.bishops & b.baseBoard.queens)
	if diagAttacks[kingSquare][(diagMasks[kingSquare] & occupancy)].IsMaskingBB(diagonalAttackers) {
		return true
	}

	return false
}

func (b *Board) generateEvasions(kingSquare Square, checkers, fromMask, toMask Bitboard) chan Move {
	ch := make(chan Move)

	go func(b Board) {
		defer close(ch)

		sliders := checkers & (b.baseBoard.bishops | b.baseBoard.rooks | b.baseBoard.queens)

		attacked := BBVoid
		for checker := range sliders.ScanReversed() {
			attacked |= bbRays[kingSquare][checker] & ^NewBitboardFromSquare(Square(checker))
		}

		if NewBitboardFromSquare(kingSquare).IsMaskingBB(fromMask) {
			for toSquare := range (kingAttacks[kingSquare] & ^b.baseBoard.occupiedColor[b.turn] & ^attacked & toMask).ScanReversed() {
				m, _ := NewNormalMove(kingSquare, Square(toSquare))
				ch <- *m
			}
		}

		checker := Square(checkers.Msb())
		if NewBitboardFromSquare(checker) == checkers {
			// capture or block a single checker
			target := bbBetween[kingSquare][checker] | checkers

			for move := range b.generatePseudoLegalMoves(^b.baseBoard.kings&fromMask, target&toMask) {
				ch <- move
			}

			// Capture the checking pawn en passant (avoid duplicate)
			if b.epSquare != SquareNone && !NewBitboardFromSquare(b.epSquare).IsMaskingBB(target) {
				lastDouble := b.epSquare
				if b.turn == White {
					lastDouble -= 8
				} else {
					lastDouble += 8
				}

				if lastDouble == checker {
					for move := range b.generatePseudoLegalEp(fromMask, toMask) {
						ch <- move
					}
				}
			}
		}
	}(NewBoardFromBoard(b))

	return ch
}

func (b *Board) GenerateLegalMoves(fromMask, toMask Bitboard) chan Move {
	ch := make(chan Move)

	go func(b Board) {
		defer close(ch)

		if b.IsVariantEnd() {
			return
		}

		kingMask := b.baseBoard.kings & b.baseBoard.occupiedColor[b.turn]
		if kingMask != BBVoid {
			king := Square(kingMask.Msb())
			blockers := b.sliderBlockers(king)
			checkers := b.baseBoard.AttackersMask(b.turn.Swap(), king)

			if checkers != BBVoid {
				for move := range b.generateEvasions(king, checkers, fromMask, toMask) {
					if b.isSafe(king, blockers, &move) {
						ch <- move
					}
				}
			} else {
				for move := range b.generatePseudoLegalMoves(fromMask, toMask) {
					if b.isSafe(king, blockers, &move) {
						ch <- move
					}
				}
			}
		} else {
			for move := range b.generatePseudoLegalMoves(fromMask, toMask) {
				ch <- move
			}
		}
	}(NewBoardFromBoard(b))

	return ch
}

func (b *Board) generateLegalEp(fromMask, toMask Bitboard) chan Move {
	ch := make(chan Move)

	go func(b Board) {
		defer close(ch)

		if b.IsVariantEnd() {
			return
		}

		for move := range b.generatePseudoLegalEp(fromMask, toMask) {
			if !b.IsIntoCheck(&move) {
				ch <- move
			}
		}
	}(NewBoardFromBoard(b))

	return ch
}

func (b *Board) generateLegalCaptures(fromMask, toMask Bitboard) chan Move {
	ch := make(chan Move)

	go func() {
		defer close(ch)
		for move := range b.GenerateLegalMoves(fromMask, toMask&b.baseBoard.occupiedColor[b.turn.Swap()]) {
			ch <- move
		}
		for move := range b.generateLegalEp(fromMask, toMask) {
			ch <- move
		}
	}()

	return ch
}

func (b *Board) attackedForKing(path, occupied Bitboard) bool {
	for sq := range path.ScanReversed() {
		if b.baseBoard.attackersMask(b.turn.Swap(), Square(sq), occupied) != BBVoid {
			return true
		}
	}

	return false
}

func (b *Board) castlingUncoversRankAttack(rookMask Bitboard, kingTo Square) Bitboard {
	if kingTo == SquareNone {
		return BBVoid
	}

	rankPieces := rankMasks[kingTo] & (b.baseBoard.occupied ^ rookMask)
	sliders := (b.baseBoard.queens | b.baseBoard.rooks) & b.baseBoard.occupiedColor[b.turn.Swap()]
	return rankAttacks[kingTo][rankPieces] & sliders
}

func (b *Board) generateCastlingMoves(fromMask, toMask Bitboard) chan Move {
	ch := make(chan Move)

	go func(b Board) {
		defer close(ch)

		if b.IsVariantEnd() {
			return
		}

		backRank := BBRank8
		if b.turn == White {
			backRank = BBRank1
		}

		kingMask := b.baseBoard.occupiedColor[b.turn] & b.baseBoard.kings & ^b.baseBoard.promoted & backRank & fromMask
		kingMask = kingMask & -kingMask
		if kingMask == BBVoid || b.attackedForKing(kingMask, b.baseBoard.occupied) {
			return
		}

		bbC := BBFileC & backRank
		bbD := BBFileD & backRank
		bbF := BBFileF & backRank
		bbG := BBFileG & backRank

		for candidate := range (b.CleanCastlingRights() & backRank & toMask).ScanReversed() {
			rookMask := NewBitboardFromSquare(Square(candidate))
			aSide := rookMask < kingMask

			emptyForRook := BBVoid
			emptyForKing := BBVoid

			kingTo := SquareNone

			if aSide {
				kingTo := Square(bbC.Msb())
				if !rookMask.IsMaskingBB(bbD) {
					emptyForRook = bbBetween[candidate][bbD.Msb()] | bbD
				}
				if !kingMask.IsMaskingBB(bbC) {
					emptyForKing = bbBetween[kingMask.Msb()][kingTo] | bbC
				}
			} else {
				kingTo := Square(bbG.Msb())
				if !rookMask.IsMaskingBB(bbF) {
					emptyForRook = bbBetween[candidate][bbF.Msb()] | bbF
				}
				if !kingMask.IsMaskingBB(bbG) {
					emptyForKing = bbBetween[kingMask.Msb()][kingTo] | bbG
				}
			}

			if !((((b.baseBoard.occupied ^ kingMask ^ rookMask) & (emptyForKing | emptyForRook)) != BBVoid) ||
				b.attackedForKing(emptyForKing, (b.baseBoard.occupied^kingMask)) ||
				b.castlingUncoversRankAttack(rookMask, kingTo) != BBVoid) {
				m := b.fromChess960(b.chess960, Square(kingMask.Msb()), Square(candidate), NoPiece, NoPiece)
				ch <- *m
			}
		}
	}(NewBoardFromBoard(b))

	return ch
}

func (b *Board) IsEnPassant(move *Move) bool {
	return b.epSquare == move.ToSquare &&
		b.baseBoard.pawns.IsMaskingBB(NewBitboardFromSquare(move.FromSquare)) &&
		(util.AbsInt(int(move.ToSquare)-int(move.FromSquare)) == 7 || util.AbsInt(int(move.ToSquare)-int(move.FromSquare)) == 9) &&
		!b.baseBoard.occupied.IsMaskingBB(NewBitboardFromSquare(move.ToSquare))
}

func (b *Board) IsCapture(move *Move) bool {
	return NewBitboardFromSquare(move.ToSquare).IsMaskingBB(b.baseBoard.occupiedColor[b.turn.Swap()]) || b.IsEnPassant(move)
}

func (b *Board) isZeroing(move *Move) bool {
	return NewBitboardFromSquare(move.FromSquare).IsMaskingBB(b.baseBoard.pawns) ||
		NewBitboardFromSquare(move.ToSquare).IsMaskingBB(b.baseBoard.occupiedColor[b.turn.Swap()])
}

func (b *Board) isIrreversible(move *Move) bool {
	backrank := BBRank1
	if b.turn == Black {
		backrank = BBRank8
	}
	cr := b.CleanCastlingRights() & backrank

	return b.isZeroing(move) ||
		(cr != BBVoid && (NewBitboardFromSquare(move.FromSquare)&b.baseBoard.kings & ^b.baseBoard.promoted) != BBVoid) ||
		cr.IsMaskingBB(NewBitboardFromSquare(move.FromSquare)) ||
		cr.IsMaskingBB(NewBitboardFromSquare(move.ToSquare))
}

func (b *Board) IsCastling(m *Move) bool {
	if b.baseBoard.kings.IsMaskingBB(NewBitboardFromSquare(m.FromSquare)) {
		diff := int(m.FromSquare.File()) - int(m.ToSquare.File())
		return util.AbsInt(diff) > 1 || (b.baseBoard.rooks&b.baseBoard.occupiedColor[b.turn]&NewBitboardFromSquare(m.ToSquare) != BBVoid)
	}

	return false
}

func (b *Board) IsKingsideCastling(move *Move) bool {
	return b.IsCastling(move) && move.ToSquare.File() > move.FromSquare.File()
}

func (b *Board) IsQueensideCastling(move *Move) bool {
	return b.IsCastling(move) && move.ToSquare.File() < move.FromSquare.File()
}

func (b *Board) fromChess960(chess960 bool, fromSquare, toSquare Square, promotion, drop PieceType) *Move {
	if !chess960 && drop == NoPiece {
		if fromSquare == E1 && b.baseBoard.kings.IsMaskingBB(BBE1) {
			if toSquare == H1 {
				m, _ := NewNormalMove(E1, G1)
				return m
			} else if toSquare == A1 {
				m, _ := NewNormalMove(E1, C1)
				return m
			}
		} else if fromSquare == E8 && b.baseBoard.kings.IsMaskingBB(BBE8) {
			if toSquare == H8 {
				m, _ := NewNormalMove(E8, G8)
				return m
			} else if toSquare == A8 {
				m, _ := NewNormalMove(E8, C8)
				return m
			}
		}
	}

	m, _ := NewMove(fromSquare, toSquare, promotion, drop)
	return m
}

func (b *Board) toChess960(m *Move) *Move {
	if m.FromSquare == E1 && b.baseBoard.kings.IsMaskingBB(BBE1) {
		if m.ToSquare == G1 && !b.baseBoard.rooks.IsMaskingBB(BBG1) {
			m, _ := NewNormalMove(E1, H1)
			return m
		} else if m.ToSquare == C1 && !b.baseBoard.rooks.IsMaskingBB(BBC1) {
			m, _ := NewNormalMove(E1, A1)
			return m
		}
	} else if m.FromSquare == E8 && b.baseBoard.kings.IsMaskingBB(BBE8) {
		if m.ToSquare == G8 && !b.baseBoard.rooks.IsMaskingBB(BBG8) {
			m, _ := NewNormalMove(E8, H8)
			return m
		} else if m.ToSquare == C8 && !b.baseBoard.rooks.IsMaskingBB(BBC8) {
			m, _ := NewNormalMove(E8, A8)
			return m
		}
	}

	return m
}

func (b *Board) transpositionKey() string {
	return fmt.Sprintf(
		"%d%d%d%d%d%d%d%d%d%d%d",
		b.baseBoard.pawns,
		b.baseBoard.knights,
		b.baseBoard.bishops,
		b.baseBoard.rooks,
		b.baseBoard.queens,
		b.baseBoard.kings,
		b.baseBoard.occupiedColor[White],
		b.baseBoard.occupiedColor[Black],
		b.turn,
		b.CleanCastlingRights(),
		b.epSquare,
	)
}

func (b Board) String() string {
	return b.Ascii()
}
