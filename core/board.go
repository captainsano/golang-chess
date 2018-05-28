package core

type FEN string

const (
	StartingFEN FEN = "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1"
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
	StartingBoardFEN FEN = "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR"
)

func MakeBoard(fen FEN) Board {
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
	return BBAll
}

func (b *Board) Pin(c Color, s Square) Bitboard {
	return BBVoid
}

func (b *Board) SetFEN(fen FEN) {

}
