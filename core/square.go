package core

import (
	"github.com/captainsano/golang-chess/util"
)

type File uint8

const (
	FileA File = 0
	FileB File = 1
	FileC File = 2
	FileD File = 3
	FileE File = 4
	FileF File = 5
	FileG File = 6
	FileH File = 7

	FileNone File = 8
)

var FileIter = [...]File{FileA, FileB, FileC, FileD, FileE, FileF, FileG, FileH}

func (f File) Name() string {
	switch f {
	case FileA:
		return "a"
	case FileB:
		return "b"
	case FileC:
		return "c"
	case FileD:
		return "d"
	case FileE:
		return "e"
	case FileF:
		return "f"
	case FileG:
		return "g"
	case FileH:
		return "h"
	}

	panic("Invalid file code")
}

type Rank uint8

const (
	Rank1 Rank = 0
	Rank2 Rank = 1
	Rank3 Rank = 2
	Rank4 Rank = 3
	Rank5 Rank = 4
	Rank6 Rank = 5
	Rank7 Rank = 6
	Rank8 Rank = 7

	RankNone Rank = 8
)

var RankReverseIter = [...]Rank{Rank8, Rank7, Rank6, Rank5, Rank4, Rank3, Rank2, Rank1}

func (r Rank) Name() string {
	switch r {
	case Rank1:
		return "1"
	case Rank2:
		return "2"
	case Rank3:
		return "3"
	case Rank4:
		return "4"
	case Rank5:
		return "5"
	case Rank6:
		return "6"
	case Rank7:
		return "7"
	case Rank8:
		return "8"
	}

	panic("Invalid rank code")
}

type Square uint8

const (
	SquareNone Square = 64

	A1 Square = 0
	B1 Square = 1
	C1 Square = 2
	D1 Square = 3
	E1 Square = 4
	F1 Square = 5
	G1 Square = 6
	H1 Square = 7
	A2 Square = 8
	B2 Square = 9
	C2 Square = 10
	D2 Square = 11
	E2 Square = 12
	F2 Square = 13
	G2 Square = 14
	H2 Square = 15
	A3 Square = 16
	B3 Square = 17
	C3 Square = 18
	D3 Square = 19
	E3 Square = 20
	F3 Square = 21
	G3 Square = 22
	H3 Square = 23
	A4 Square = 24
	B4 Square = 25
	C4 Square = 26
	D4 Square = 27
	E4 Square = 28
	F4 Square = 29
	G4 Square = 30
	H4 Square = 31
	A5 Square = 32
	B5 Square = 33
	C5 Square = 34
	D5 Square = 35
	E5 Square = 36
	F5 Square = 37
	G5 Square = 38
	H5 Square = 39
	A6 Square = 40
	B6 Square = 41
	C6 Square = 42
	D6 Square = 43
	E6 Square = 44
	F6 Square = 45
	G6 Square = 46
	H6 Square = 47
	A7 Square = 48
	B7 Square = 49
	C7 Square = 50
	D7 Square = 51
	E7 Square = 52
	F7 Square = 53
	G7 Square = 54
	H7 Square = 55
	A8 Square = 56
	B8 Square = 57
	C8 Square = 58
	D8 Square = 59
	E8 Square = 60
	F8 Square = 61
	G8 Square = 62
	H8 Square = 63
)

var squares = func() []Square {
	sqs := []Square{}
	for i := Square(0); i < 64; i++ {
		sqs = append(sqs, i)
	}
	return sqs
}()

var squares180 = func() []Square {
	sqs180 := []Square{}
	for _, i := range squares {
		sqs180 = append(sqs180, i.Mirror())
	}
	return sqs180
}()

func NewSquare(file File, rank Rank) Square {
	return Square(uint(rank)*8 + uint(file))
}

func FileFromName(name string) File {
	switch name {
	case "a":
		return FileA
	case "b":
		return FileB
	case "c":
		return FileC
	case "d":
		return FileD
	case "e":
		return FileE
	case "f":
		return FileF
	case "g":
		return FileG
	case "h":
		return FileH
	}

	return FileNone
}

func RankFromName(name string) Rank {
	switch name {
	case "1":
		return Rank1
	case "2":
		return Rank2
	case "3":
		return Rank3
	case "4":
		return Rank4
	case "5":
		return Rank5
	case "6":
		return Rank6
	case "7":
		return Rank7
	case "8":
		return Rank8
	}

	return RankNone
}

// TODO: Optimize with ASCII value computation
func NewSquareFromName(name string) Square {
	file := FileFromName(name[0:1])
	if file == FileNone {
		return SquareNone
	}

	rank := RankFromName(name[1:2])
	if rank == RankNone {
		return SquareNone
	}

	return NewSquare(file, rank)
}

func (s Square) File() File {
	return File(uint8(s) & 7)
}

func (s Square) Rank() Rank {
	return Rank(uint8(s) >> 3)
}

func (s Square) Name() string {
	return s.File().Name() + s.Rank().Name()
}

func (s *Square) AddDelta(d int) {
	*s = Square(int(*s) + d)
}

func (s Square) IsInRange() bool {
	x := int(s)
	return x >= 0 && x < 64
}

func (a Square) Distance(b Square) int {
	x := util.AbsInt(int(a.File()) - int(b.File()))
	y := util.AbsInt(int(a.Rank()) - int(b.Rank()))
	return util.MaxInt(x, y)
}

func SquareDistance(a, b Square) int {
	return a.Distance(b)
}

func (s Square) Index() uint8 {
	return uint8(s)
}

func (s Square) Mirror() Square {
	return Square(s ^ 0x38)
}
