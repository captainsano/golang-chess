package core

import (
	"strconv"
	"strings"
)

const (
	StartingFEN = "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1"
)

type Board struct {
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

func MakeBoard(fen string) Board {
	b := Board{}
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

func (b *Board) Reset() {
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

func (b *Board) Clear() {
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

func (b *Board) PieceMask(t PieceType, c Color) Bitboard {
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

func (b *Board) Pieces(t PieceType, c Color) Bitboard {
	return b.PieceMask(t, c)
}

func (b *Board) PieceAt(s Square) *Piece {
	t := b.PieceTypeAt(s)
	if t == NoPiece {
		return nil
	}

	mask := BBSquare(s)
	if b.occupiedColor[White].IsMaskingBB(mask) {
		return &Piece{t, White}
	}
	return &Piece{t, Black}
}

func (b *Board) PieceTypeAt(s Square) PieceType {
	mask := BBSquare(s)

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

func (b *Board) King(c Color) *Square {
	mask := b.occupiedColor[c] & b.kings & ^b.promoted

	if mask != 0 {
		s := Square(mask.Msb())
		return &s
	}

	return nil
}

func (b *Board) Attacks(s Square) Bitboard {
	mask := BBSquare(s)

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

func (b *Board) AttackersMask(c Color, s Square) Bitboard {
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

func (b *Board) IsAttackedBy(c Color, s Square) bool {
	return b.AttackersMask(c, s) != 0
}

func (b *Board) Attackers(c Color, s Square) Bitboard {
	return b.AttackersMask(c, s)
}

func (b *Board) PinMask(c Color, s Square) Bitboard {
	kingSq := b.King(c)
	if kingSq == nil {
		return BBAll
	}

	squareMask := MakeBitboardFromSquare(s)

	ks := [][]map[Bitboard]Bitboard{fileAttacks, rankAttacks, diagAttacks}
	vs := []Bitboard{b.rooks | b.queens, b.rooks | b.queens, b.bishops | b.queens}

	for i, _ := range ks {
		attacks, sliders := ks[i], vs[i]

		rays := attacks[*kingSq][0]
		if rays.IsMaskingBB(squareMask) {
			snipers := rays & sliders & b.occupiedColor[c.Swap()]
			for sniper := range snipers.ScanReversed() {
				if bbBetween[sniper][*kingSq]&(b.occupied|squareMask) == squareMask {
					return bbRays[*kingSq][sniper]
				}
			}

			break
		}
	}

	return BBAll
}

func (b *Board) IsPinned(c Color, s Square) bool {
	return b.PinMask(c, s) != BBAll
}

func (b *Board) RemovePieceAt(s Square) Piece {
	pt := b.PieceTypeAt(s)
	mask := MakeBitboardFromSquare(s)

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
		return MakePiece(pt, color)
	}

	b.occupied ^= mask
	b.occupiedColor[White] &= ^mask
	b.occupiedColor[Black] &= ^mask

	b.promoted &= ^mask

	return MakePiece(pt, color)
}

func (b *Board) _setPieceAt(s Square, pt PieceType, c Color, promoted bool) {
	b.RemovePieceAt(s)

	mask := MakeBitboardFromSquare(s)

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

func (b *Board) SetPieceAt(s Square, p *Piece, promoted bool) {
	if p != nil {
		b._setPieceAt(s, p.Type, p.Color, promoted)
	} else {
		b.RemovePieceAt(s)
	}
}

func (b *Board) FEN(promoted bool) string {
	builder := []string{}
	empty := 0

	for _, square := range squares180 {
		piece := b.PieceAt(square)

		if piece == nil {
			empty += 1
		} else {
			if empty > 0 {
				builder = append(builder, string(empty))
				empty = 0
			}

			builder = append(builder, piece.Symbol())

			if promoted && MakeBitboardFromSquare(square).IsMaskingBB(b.promoted) {
				builder = append(builder, "~")
			}
		}

		if MakeBitboardFromSquare(square).IsMaskingBB(BBFileH) {
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

func (b *Board) SetFEN(fen string) {
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
				fieldSum += 1
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
			piece := MakePieceFromSymbol(c)
			b.SetPieceAt(squares180[squareIndex], &piece, false)
			squareIndex += 1
		} else if c == "~" {
			b.promoted |= MakeBitboardFromSquare(squares[squares180[squareIndex-1]])
		}
	}
}

func (b *Board) PieceMap() map[Square]*Piece {
	result := make(map[Square]*Piece)
	for s := range b.occupied.ScanReversed() {
		p := b.PieceAt(Square(s))
		cp := MakePiece(p.Type, p.Color)
		result[Square(s)] = &cp
	}
	return result
}

func (b *Board) SetPieceMap(pm map[Square]*Piece) {
	b.Clear()
	for s, p := range pm {
		cp := MakePiece(p.Type, p.Color)
		b.SetPieceAt(s, &cp, false)
	}
}

// func SetChess960Pos() {

// }

func (b *Board) Ascii() string {
	builder := []string{}

	for _, square := range squares180 {
		piece := b.PieceAt(square)

		if piece != nil {
			builder = append(builder, piece.Symbol())
		} else {
			builder = append(builder, ".")
		}

		if MakeBitboardFromSquare(square).IsMaskingBB(BBFileH) {
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
func (b *Board) Unicode(invertColor, borders bool) string {
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
			square := MakeSquare(File(file), Rank(rank))

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
