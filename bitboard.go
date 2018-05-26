package main

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

func BBSquares(i int) Bitboard {
	return Bitboard(1 << uint(i))
}

func BBFiles(i int) Bitboard {
	return Bitboard(0x0101010101010101 << uint(i))
}

func BBRanks(i int) Bitboard {
	return Bitboard(0xff << (8 * uint(i)))
}

// TODO: Bitboard functions
