package main

import (
	"fmt"
)

// InlineChunk denotes chunk for inline fomrmat

type InlineChunk struct {
	Position int
	Keyword  string
	Children []Chunk
}

func (p InlineChunk) IsTerminal() bool {
	if len(p.Children) == 0 {
		panic("inline format shall always has children")
	}
	if len(p.Children) == 1 {
		return true
	}
	return false
}

// String implements the Stringer interface
func (p InlineChunk) String() string {
	return fmt.Sprintf("InlineChunk{Position:%v,Keyword:%v, Children:%v}", p.Position, p.Keyword, p.Children)
}

// GetPosition implements the Chunk interface
func (p *InlineChunk) GetPosition() int {
	return p.Position
}

// SetPosition implements the Chunk interface
func (p *InlineChunk) SetPosition(pos int) {
	p.Position = pos
}

// GetValue implements the Chunk interface
func (p *InlineChunk) GetValue() string {
	panic("should no use,  except debug")
	return p.Keyword
}
