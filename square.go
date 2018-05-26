package main

import "github.com/captainsano/golang-chess/util"

var squares = func() []int {
	s := []int{}
	for i := 0; i < 64; i++ {
		s[i] = i
	}
	return s
}()

func Squares(i int) int {
	return squares[uint(i)]
}

const (
	A1 = 0
	B1 = 1
	C1 = 2
	D1 = 3
	E1 = 4
	F1 = 5
	G1 = 6
	H1 = 7
	A2 = 8
	B2 = 9
	C2 = 10
	D2 = 11
	E2 = 12
	F2 = 13
	G2 = 14
	H2 = 15
	A3 = 16
	B3 = 17
	C3 = 18
	D3 = 19
	E3 = 20
	F3 = 21
	G3 = 22
	H3 = 23
	A4 = 24
	B4 = 25
	C4 = 26
	D4 = 27
	E4 = 28
	F4 = 29
	G4 = 30
	H4 = 31
	A5 = 32
	B5 = 33
	C5 = 34
	D5 = 35
	E5 = 36
	F5 = 37
	G5 = 38
	H5 = 39
	A6 = 40
	B6 = 41
	C6 = 42
	D6 = 43
	E6 = 44
	F6 = 45
	G6 = 46
	H6 = 47
	A7 = 48
	B7 = 49
	C7 = 50
	D7 = 51
	E7 = 52
	F7 = 53
	G7 = 54
	H7 = 55
	A8 = 56
	B8 = 57
	C8 = 58
	D8 = 59
	E8 = 60
	F8 = 61
	G8 = 62
	H8 = 63
)

func Square(fileIndex, rankIndex int) int {
	return Squares(rankIndex*8 + fileIndex)
}

func SquareFile(square int) int {
	return square & 7
}

func SquareRank(square int) int {
	return square >> 3
}

func SquareName(square int) string {
	return FileNames(SquareFile(square)) + RankNames(SquareRank(square))
}

func SquareDistance(a, b int) int {
	x := util.AbsInt(SquareFile(a) - SquareFile(b))
	y := util.AbsInt(SquareRank(a) - SquareRank(b))
	return util.MaxInt(x, y)
}

func SquareMirror(s int) int {
	return s ^ 0x38
}
