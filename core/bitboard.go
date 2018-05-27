package core

import (
	"strings"
)

type Bitboard uint64

const (
	BBVoid Bitboard = 0x0
	BBAll  Bitboard = 0xffffffffffffffff
)

const (
	BBA1 Bitboard = 1 << A1
	BBB1 Bitboard = 1 << B1
	BBC1 Bitboard = 1 << C1
	BBD1 Bitboard = 1 << D1
	BBE1 Bitboard = 1 << E1
	BBF1 Bitboard = 1 << F1
	BBG1 Bitboard = 1 << G1
	BBH1 Bitboard = 1 << H1

	BBA2 Bitboard = 1 << A2
	BBB2 Bitboard = 1 << B2
	BBC2 Bitboard = 1 << C2
	BBD2 Bitboard = 1 << D2
	BBE2 Bitboard = 1 << E2
	BBF2 Bitboard = 1 << F2
	BBG2 Bitboard = 1 << G2
	BBH2 Bitboard = 1 << H2

	BBA3 Bitboard = 1 << A3
	BBB3 Bitboard = 1 << B3
	BBC3 Bitboard = 1 << C3
	BBD3 Bitboard = 1 << D3
	BBE3 Bitboard = 1 << E3
	BBF3 Bitboard = 1 << F3
	BBG3 Bitboard = 1 << G3
	BBH3 Bitboard = 1 << H3

	BBA4 Bitboard = 1 << A4
	BBB4 Bitboard = 1 << B4
	BBC4 Bitboard = 1 << C4
	BBD4 Bitboard = 1 << D4
	BBE4 Bitboard = 1 << E4
	BBF4 Bitboard = 1 << F4
	BBG4 Bitboard = 1 << G4
	BBH4 Bitboard = 1 << H4

	BBA5 Bitboard = 1 << A5
	BBB5 Bitboard = 1 << B5
	BBC5 Bitboard = 1 << C5
	BBD5 Bitboard = 1 << D5
	BBE5 Bitboard = 1 << E5
	BBF5 Bitboard = 1 << F5
	BBG5 Bitboard = 1 << G5
	BBH5 Bitboard = 1 << H5

	BBA6 Bitboard = 1 << A6
	BBB6 Bitboard = 1 << B6
	BBC6 Bitboard = 1 << C6
	BBD6 Bitboard = 1 << D6
	BBE6 Bitboard = 1 << E6
	BBF6 Bitboard = 1 << F6
	BBG6 Bitboard = 1 << G6
	BBH6 Bitboard = 1 << H6

	BBA7 Bitboard = 1 << A7
	BBB7 Bitboard = 1 << B7
	BBC7 Bitboard = 1 << C7
	BBD7 Bitboard = 1 << D7
	BBE7 Bitboard = 1 << E7
	BBF7 Bitboard = 1 << F7
	BBG7 Bitboard = 1 << G7
	BBH7 Bitboard = 1 << H7

	BBA8 Bitboard = 1 << A8
	BBB8 Bitboard = 1 << B8
	BBC8 Bitboard = 1 << C8
	BBD8 Bitboard = 1 << D8
	BBE8 Bitboard = 1 << E8
	BBF8 Bitboard = 1 << F8
	BBG8 Bitboard = 1 << G8
	BBH8 Bitboard = 1 << H8

	BBCorners = BBA1 | BBH1 | BBA8 | BBH8

	BBLightsquares Bitboard = 0x55aa55aa55aa55aa
	BBDarkSquares  Bitboard = 0xaa55aa55aa55aa55

	BBFileA Bitboard = 0x0101010101010101 << 0
	BBFileB Bitboard = 0x0101010101010101 << 1
	BBFileC Bitboard = 0x0101010101010101 << 2
	BBFileD Bitboard = 0x0101010101010101 << 3
	BBFileE Bitboard = 0x0101010101010101 << 4
	BBFileF Bitboard = 0x0101010101010101 << 5
	BBFileG Bitboard = 0x0101010101010101 << 6
	BBFileH Bitboard = 0x0101010101010101 << 7

	BBRank1 Bitboard = 0xff << (8 * 0)
	BBRank2 Bitboard = 0xff << (8 * 1)
	BBRank3 Bitboard = 0xff << (8 * 2)
	BBRank4 Bitboard = 0xff << (8 * 3)
	BBRank5 Bitboard = 0xff << (8 * 4)
	BBRank6 Bitboard = 0xff << (8 * 5)
	BBRank7 Bitboard = 0xff << (8 * 6)
	BBRank8 Bitboard = 0xff << (8 * 7)

	BBBackRanks Bitboard = BBRank1 | BBRank8
)

func BBRank(r Rank) Bitboard {
	return 0xff << (8 * r)
}

func BBFile(f File) Bitboard {
	return 0x0101010101010101 << f
}

func BBSquare(s Square) Bitboard {
	return 1 << s.Index()
}

func MakeBitboard(b uint64) Bitboard {
	return Bitboard(b)
}

func MakeBitboardFromSquareIndex(i uint8) Bitboard {
	return Bitboard(1 << i)
}

func MakeBitboardFromSquare(s Square) Bitboard {
	return MakeBitboardFromSquareIndex(s.Index())
}

func MakeBitboardFromFileIndex(i uint8) Bitboard {
	return Bitboard(0x0101010101010101 << i)
}

func MakeBitboardFromRankIndex(i uint8) Bitboard {
	return Bitboard(0xff << (8 * i))
}

func (a Bitboard) IsMaskingBB(b Bitboard) bool {
	return (uint64(a) & uint64(b)) > 0
}

func (b *Bitboard) Add(s Square) {
	mask := MakeBitboardFromSquareIndex(s.Index())
	*b |= mask
}

func (b *Bitboard) Remove(s Square) {
	mask := ^MakeBitboardFromSquareIndex(s.Index())
	*b &= mask
}

func (b *Bitboard) Clear() {
	*b = BBVoid
}

func (b *Bitboard) ShiftDown() {
	*b = *b >> 8
}

func (b *Bitboard) Shift2Down() {
	*b = *b >> 16
}

func (b *Bitboard) ShiftUp() {
	*b = *b << 8
}

func (b *Bitboard) Shift2Up() {
	*b = *b << 16
}

func (b *Bitboard) ShiftRight() {
	*b = (*b << 1) & ^BBFileA & BBAll
}

func (b *Bitboard) Shift2Right() {
	*b = (*b << 2) & ^BBFileA & ^BBFileB & BBAll
}

func (b *Bitboard) ShiftLeft() {
	*b = (*b >> 1) & ^BBFileH
}

func (b *Bitboard) Shift2Left() {
	*b = (*b >> 2) & ^BBFileG & ^BBFileH
}

func (b *Bitboard) ShiftUpLeft() {
	*b = (*b << 7) & ^BBFileH & BBAll
}

func (b *Bitboard) ShiftUpRight() {
	*b = (*b << 9) & ^BBFileA & BBAll
}

func (b *Bitboard) ShiftDownLeft() {
	*b = (*b >> 9) & ^BBFileH
}

func (b *Bitboard) ShiftDownRight() {
	*b = (*b >> 7) & ^BBFileA
}

func (b Bitboard) Ascii() string {
	ascii := []string{}

	for _, r := range RankReverseIter {
		for _, f := range FileIter {
			sq := MakeSquare(f, r)
			mask := MakeBitboardFromSquareIndex(sq.Index())

			if b.IsMaskingBB(mask) {
				ascii = append(ascii, "1")
			} else {
				ascii = append(ascii, ".")
			}

			if mask.IsMaskingBB(BBFileH) {
				if sq != H1 {
					ascii = append(ascii, "\n")
				}
			} else {
				ascii = append(ascii, " ")
			}
		}
	}

	return strings.Join(ascii, "")
}

func SlidingAttacks(square Square, occupied Bitboard, deltas []int) Bitboard {
	attacks := BBVoid

	for _, d := range deltas {
		sq := square

		for {
			sq.AddDelta(d)
			sq2 := sq
			sq2.AddDelta(-d)
			if !sq.IsInRange() || sq.Distance(sq2) > 2 {
				break
			}

			mask := MakeBitboardFromSquare(sq)

			attacks |= mask

			if occupied.IsMaskingBB(mask) {
				break
			}
		}
	}

	return attacks
}

// ------- Bitboard utility functions --------
func KnightAttacks(s Square) Bitboard {
	return SlidingAttacks(s, BBAll, []int{
		17, 15, 10, 6, -17, -15, -10, -6,
	})
}

func KingAttacks(s Square) Bitboard {
	return SlidingAttacks(s, BBAll, []int{
		9, 8, 7, 1, -9, -8, -7, -1,
	})
}

func PawnAttacks(s Square, c Color) Bitboard {
	if c == White {
		return SlidingAttacks(s, BBAll, []int{7, 9})
	}

	return SlidingAttacks(s, BBAll, []int{-7, -9})
}

func Edges(s Square) Bitboard {
	return ((BBRank1 | BBRank8) & ^BBRank(s.Rank())) | ((BBFileA | BBFileH) & ^BBFile(s.File()))
}

func CarryRippler(mask Bitboard) chan Bitboard {
	ch := make(chan Bitboard)

	go func(m Bitboard) {
		subset := BBVoid
		for {
			ch <- subset
			subset = (subset - mask) & mask
			if subset == 0 {
				close(ch)
				break
			}
		}
	}(mask)

	return ch
}

func AttackTable(deltas []int) ([]Bitboard, []map[Bitboard]Bitboard) {
	maskTable := []Bitboard{}
	attackTable := []map[Bitboard]Bitboard{}

	for sq := Square(0); sq < 64; sq++ {
		attacks := make(map[Bitboard]Bitboard)

		mask := SlidingAttacks(sq, BBVoid, deltas) & Edges(sq)
		for subset := range CarryRippler(mask) {
			attacks[subset] = SlidingAttacks(sq, subset, deltas)
		}

		attackTable = append(attackTable, attacks)
		maskTable = append(maskTable, mask)
	}

	return maskTable, attackTable
}

var diagMasks []Bitboard

func DiagMask(s Square) Bitboard {
	if diagMasks != nil {
		return diagMasks[s]
	}

	diagMasks, _ = AttackTable([]int{-9, -7, 7, 9})
	return diagMasks[s]
}

var diagAttacks []map[Bitboard]Bitboard

func DiagAttacks(s Square) map[Bitboard]Bitboard {
	if diagAttacks != nil {
		return diagAttacks[s]
	}

	_, diagAttacks = AttackTable([]int{-9, -7, 7, 9})
	return diagAttacks[s]
}

var fileMasks []Bitboard

func FileMask(s Square) Bitboard {
	if fileMasks != nil {
		return fileMasks[s]
	}

	fileMasks, _ = AttackTable([]int{-8, 8})
	return fileMasks[s]
}

var fileAttacks []map[Bitboard]Bitboard

func FileAttack(s Square) map[Bitboard]Bitboard {
	if fileAttacks != nil {
		return fileAttacks[s]
	}

	_, fileAttacks = AttackTable([]int{-8, 8})
	return fileAttacks[s]
}

var rankMasks []Bitboard

func RankMask(s Square) Bitboard {
	if rankMasks != nil {
		return rankMasks[s]
	}

	rankMasks, _ = AttackTable([]int{-1, 1})
	return rankMasks[s]
}

var rankAttacks []map[Bitboard]Bitboard

func RankAttack(s Square) map[Bitboard]Bitboard {
	if rankAttacks != nil {
		return rankAttacks[s]
	}

	_, rankAttacks = AttackTable([]int{-1, 1})
	return rankAttacks[s]
}

func Rays() ([][]Bitboard, [][]Bitboard) {
	rays := [][]Bitboard{}
	between := [][]Bitboard{}

	for a := Square(0); a < 64; a++ {
		bbA := BBSquare(a)

		rays_row := []Bitboard{}
		between_row := []Bitboard{}

		for b := Square(0); b < 64; b++ {
			bbB := BBSquare(b)

			if DiagAttacks(a)[0].IsMaskingBB(bbB) {
				rays_row = append(rays_row, ((DiagAttacks(a)[0] & DiagAttacks(b)[0]) | bbA | bbB))
				between_row = append(between_row, (DiagAttacks(a)[DiagMask(a)&bbB] & DiagAttacks(b)[DiagMask(b)&bbA]))
			} else if RankAttack(a)[0].IsMaskingBB(bbB) {
				rays_row = append(rays_row, RankAttack(a)[0]|bbA)
				between_row = append(between_row, RankAttack(a)[RankMask(a)&bbB]&RankAttack(b)[RankMask(b)&bbA])
			} else if FileAttack(a)[0].IsMaskingBB(bbB) {
				rays_row = append(rays_row, FileAttack(a)[0]|bbA)
				between_row = append(between_row, FileAttack(a)[FileMask(a)&bbB]&FileAttack(b)[FileMask(b)&bbA])
			} else {
				rays_row = append(rays_row, 0)
				between_row = append(between_row, 0)
			}
		}

		rays = append(rays, rays_row)
		between = append(between, between_row)
	}

	return rays, between
}
