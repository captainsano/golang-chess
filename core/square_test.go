package core

import (
	"testing"
)

func TestSquareLabel(t *testing.T) {
	for squareIndex := 0; squareIndex < 64; squareIndex++ {
		s := Square(squareIndex)

		file := s.File()
		rank := s.Rank()

		if MakeSquare(file, rank) != s {
			t.Errorf("(file, rank) to square conversion error: (%v, %v) - %v - %v", file, rank, MakeSquare(file, rank), s)
		}

		fileName := file.Name()
		rankName := rank.Name()

		if fileName+rankName != s.Name() {
			t.Errorf("filename+rankName not equal to square name")
		}
	}
}

func TestSquareDistance(t *testing.T) {
	tests := map[int]int{
		MakeSquare(0, 0).Distance(MakeSquare(0, 0)): 0,
		MakeSquare(0, 0).Distance(MakeSquare(7, 7)): 7,
		MakeSquare(7, 7).Distance(MakeSquare(0, 0)): 7,
		MakeSquare(0, 0).Distance(MakeSquare(0, 7)): 7,
		MakeSquare(0, 7).Distance(MakeSquare(0, 0)): 7,
		MakeSquare(3, 3).Distance(MakeSquare(3, 4)): 1,
	}

	for a, e := range tests {
		if a != e {
			t.Errorf("Distance error: got %v, expected %v", a, e)
		}
	}
}
