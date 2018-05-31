package core

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
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

func (b *BaseBoard) AttackersMask(c Color, s Square) Bitboard {
	rank_pieces := RankMasks(s) & b.occupied
	file_pieces := FileMasks(s) & b.occupied
	diag_pieces := DiagMasks(s) & b.occupied

	queens_and_rooks := b.queens | b.rooks
	queens_and_bishops := b.queens | b.bishops

	attackers := (KingAttacks(s) & b.kings) |
		(KnightAttacks(s) & b.knights) |
		(RankAttacks(s)[rank_pieces] & queens_and_rooks) |
		(FileAttacks(s)[file_pieces] & queens_and_rooks) |
		(DiagAttacks(s)[diag_pieces] & queens_and_bishops) |
		(PawnAttacks(s, c.Swap()) & b.pawns)

	return attackers & b.occupiedColor[c]
}

func (b *BaseBoard) IsAttackedBy(c Color, s Square) bool {
	return b.AttackersMask(c, s) != 0
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

func (b *BaseBoard) _setPieceAt(s Square, pt PieceType, c Color, promoted bool) {
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
		b._setPieceAt(s, p.Type, p.Color, promoted)
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
				builder = append(builder, string(empty))
				empty = 0
			}

			builder = append(builder, piece.Symbol())

			if promoted && NewBitboardFromSquare(square).IsMaskingBB(b.promoted) {
				builder = append(builder, "~")
			}
		}

		if NewBitboardFromSquare(square).IsMaskingBB(BBFileH) {
			if empty > 0 {
				builder = append(builder, string(empty))
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
	stack     []string

	turn           Color
	castlingRights Bitboard
	epSquare       Square
	halfMoveClock  uint
	fullMoveNumber uint
}

func NewBoard(fen string, chess960 bool) Board {
	board := Board{}

	board.aliases = []string{"Standard", "Chess", "Classical", "Normal"}
	board.uciVariant = "chess"
	board.startingFen = StartingFEN
	board.connectedKings = false
	board.oneKing = true
	board.capturesCompulsory = false

	board.chess960 = chess960

	board.moveStack = []Move{}
	board.stack = []string{}

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

func (b *Board) Reset() {
	b.turn = White
	b.castlingRights = BBCorners
	b.epSquare = SquareNone
	b.halfMoveClock = 0
	b.fullMoveNumber = 1

	b.baseBoard.Reset()
	// b.clearStack()
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
	b.stack = []string{}
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

func (b *Board) GeneratePseudoLegalMoves(fromMask, toMask Bitboard) chan Move {
	ch := make(chan Move)

	go func() {
		ourPieces := b.baseBoard.occupiedColor[b.turn]

		// Generate piece moves
		nonPawns := ourPieces & ^b.baseBoard.pawns & fromMask
		for fromSquare := range nonPawns.ScanReversed() {
			fmt.Println("--> Evaluating from Square: ", fromSquare)
			moves := b.baseBoard.Attacks(Square(fromSquare)) & ^ourPieces & toMask
			for toSquare := range moves.ScanReversed() {
				ch <- NewMove(Square(fromSquare), Square(toSquare), NoPiece)
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
			close(ch)
			return
		}

		// Generate captures
		capturers := pawns
		for fromSquare := range capturers.ScanReversed() {
			targets := pawnAttacks[b.turn][fromSquare] & b.baseBoard.occupiedColor[b.turn.Swap()] & toMask
			for toSquare := range targets.ScanReversed() {
				if Square(toSquare).Rank() == 0 || Square(toSquare).Rank() == 7 {
					ch <- NewMove(Square(fromSquare), Square(toSquare), Queen)
					ch <- NewMove(Square(fromSquare), Square(toSquare), Rook)
					ch <- NewMove(Square(fromSquare), Square(toSquare), Bishop)
					ch <- NewMove(Square(fromSquare), Square(toSquare), Knight)
				} else {
					ch <- NewMove(Square(fromSquare), Square(toSquare), NoPiece)
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
				ch <- NewMove(fromSquare, Square(toSquare), Queen)
				ch <- NewMove(fromSquare, Square(toSquare), Rook)
				ch <- NewMove(fromSquare, Square(toSquare), Bishop)
				ch <- NewMove(fromSquare, Square(toSquare), Knight)
			} else {
				ch <- NewMove(fromSquare, Square(toSquare), NoPiece)
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
			ch <- NewMove(fromSquare, Square(toSquare), NoPiece)
		}

		// Generate enpassant captures
		if b.epSquare != SquareNone {
			for move := range b.generatePseudoLegalEp(fromMask, toMask) {
				ch <- move
			}
		}

		close(ch)
	}()

	return ch
}

func (b *Board) GenerateAllPseudoLegalMoves() chan Move {
	return b.GeneratePseudoLegalMoves(BBAll, BBAll)
}

func (b *Board) generatePseudoLegalEp(fromMask, toMask Bitboard) chan Move {
	ch := make(chan Move)

	go func() {
		if (b.epSquare != SquareNone) || !NewBitboardFromSquare(b.epSquare).IsMaskingBB(toMask) {
			close(ch)
			return
		}

		if NewBitboardFromSquare(b.epSquare).IsMaskingBB(b.baseBoard.occupied) {
			close(ch)
			return
		}

		capturers := b.baseBoard.pawns & b.baseBoard.occupiedColor[b.turn] & fromMask & PawnAttacks(b.epSquare, b.turn.Swap())
		if b.turn == White {
			capturers &= BBRank4
		} else {
			capturers &= BBRank3
		}

		for capturer := range capturers.ScanReversed() {
			ch <- NewMove(Square(capturer), b.epSquare, NoPiece)
		}

		close(ch)
	}()

	return ch
}

func (b *Board) generatePseudoLegalCaptures(fromMask, toMask Bitboard) chan Move {
	ch := make(chan Move)

	go func() {
		for m := range b.GeneratePseudoLegalMoves(fromMask, (toMask & b.baseBoard.occupiedColor[b.turn.Swap()])) {
			ch <- m
		}

		for m := range b.generatePseudoLegalEp(fromMask, toMask) {
			ch <- m
		}

		close(ch)
	}()

	return ch
}

func (b *Board) generateCastlingMoves(fromMask, toMask Bitboard) chan Move {
	ch := make(chan Move)

	// TODO: To be implemented

	return ch
}

func (b *Board) IsCheck() bool {
	kingSquare := b.baseBoard.King(b.turn)
	return kingSquare != SquareNone && b.baseBoard.IsAttackedBy(b.turn.Swap(), kingSquare)
}

func (b *Board) IsIntoCheck(m *Move) bool {
	// kingSquare := b.baseBoard.King(b.turn)
	// if kingSquare == SquareNone {
	// 	return false
	// }

	// checkers := b.baseBoard.AttackersMask(b.turn.Swap(), kingSquare)
	// if checkers != BBVoid {
	// 	if
	// }

	return false
}

func (b *Board) WasIntoCheck() bool {
	return false
}

func (b *Board) IsPseudoLegal(m *Move) bool {
	return false
}

func (b *Board) IsLegal(m *Move) bool {
	return false
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

func (b *Board) IsGameOver() bool {
	return false
}

func (b *Board) Result(claimDraw bool) string {
	return "*"
}

func (b *Board) IsCheckmate() bool {
	return false
}

func (b *Board) IsStalemate() bool {
	return false
}

func (b *Board) IsInsufficientMaterial() bool {
	return false
}

func (b *Board) IsSeventyFiveMoves() bool {
	return false
}

func (b *Board) IsFiveFoldRepetition() bool {
	return false
}

func (b *Board) CanClaimDraw() bool {
	return false
}

func (b *Board) CanClaimFiftyMoves() bool {
	return false
}

func (b *Board) CanClaimThreefoldRepetition() bool {
	return false
}

func (b *Board) pushCapture(m *Move, captureSquare Square, pt PieceType, wasPromoted bool) {
}

func (b *Board) Push(m *Move) {

}

func (b *Board) Pop() *Move {
	return nil
}

func (b *Board) Peek() {

}

func (b *Board) FEN(shredder bool, enPassant string, promoted PieceType) string {
	// return strings.Join([]string{
	// 	b.epd(shredder, enPassant, promoted),
	// 	string(b.halfMoveClock),
	// 	string(b.fullMoveNumber),
	// }, " ")
	return ""
}

func (b *Board) ShredderFen(enPassant string, promoted PieceType) string {
	// return strings.Join([]string{
	// 	b.epd(true, enPassant, promoted),
	// 	string(b.halfMoveClock),
	// 	string(b.fullMoveNumber),
	// }, " ")
	return ""
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
		color := White
		if flag == strings.ToUpper(flag) {
			color = White
		} else {
			color = Black
		}

		flag = strings.ToLower(flag)
		backRank := BBVoid
		if color == White {
			backRank = BBRank1
		} else {
			backRank = BBRank8
		}
		rooks := b.baseBoard.occupiedColor[color] & b.baseBoard.rooks & backRank
		kingSquare := b.baseBoard.King(color)

		if flag == "q" {
			if kingSquare != SquareNone && rooks.Lsb() < int(kingSquare) {
				b.castlingRights |= rooks & ^rooks
			} else {
				b.castlingRights |= BBFileA ^ backRank
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

func (b *Board) SetPieceMap(pm map[Square]*Piece) {
	b.baseBoard.SetPieceMap(pm)
	b.clearStack()
}

func (b *Board) Ascii() string {
	// TODO: Add other FEN params
	return b.baseBoard.Ascii()
}

func (b *Board) Unicode(invertColor, borders bool) string {
	// TODO: Add other FEN params
	return b.baseBoard.Unicode(invertColor, borders)
}
