package goref

import (
	"fmt"
	"go/token"
)

// A Position is similar to token.Position in that it gives an
// absolute position within a file, but it may also denote the Pos +
// End concept of token.Pos.
//
// The End is optional. If NoPos is used as the End, Position only
// contains file:line:column.
//
// The Pos is not optional and must resolve to file:line:column.
type Position struct {
	File       string
	PosL, PosC int
	EndL, EndC int
}

const (
	// NoPos represents a mising position. It has the same
	// semantics as token.NoPos. It may be used in place of an
	// "end" Pos if the end of an identifier isn't known.
	NoPos = token.NoPos
)

func (p Position) String() string {
	if p == (Position{}) {
		return "-"
	}
	if p.EndL >= 0 {
		return fmt.Sprintf("%s:%d:%d-%d:%d", p.File, p.PosL, p.PosC, p.EndL, p.EndC)
	}
	return fmt.Sprintf("%s:%d:%d", p.File, p.PosL, p.PosC)
}

// NewPosition creates a Position from a token.FileSet and a pair of
// Pos in that FileSet. It will panic if both Pos are not from the
// same Filename.
func NewPosition(fset *token.FileSet, pos, end token.Pos) Position {
	ppos := fset.Position(pos)
	if end == token.NoPos {
		return Position{
			File: ppos.Filename,
			PosL: ppos.Line,
			PosC: ppos.Column,
			EndL: -1,
			EndC: -1,
		}
	}
	pend := fset.Position(end)
	if ppos.Filename != pend.Filename {
		panic("Invalid pair of {pos,end} for NewPosition: pos and end come from different files!")
	}
	return Position{
		File: ppos.Filename,
		PosL: ppos.Line,
		PosC: ppos.Column,
		EndL: pend.Line,
		EndC: pend.Column,
	}
}
