package main

import (
	"fmt"
)

// BlockCodeChunk denotes a block of code
type BlockCodeChunk struct {
	Position int
	Id       string
	Caption  string //optional
	Value    string
}

func (p BlockCodeChunk) IsTerminal() bool {
	return true
}

// String implements the Stringer interface
func (p BlockCodeChunk) String() string {
	return fmt.Sprintf("BlockCodeChunk{Position: %d, Id: %v, Caption: %v, Value: %v}",
		p.GetPosition(), p.Id, p.Caption, p.Value)
}

// GetPosition implements the Chunk interface
func (p *BlockCodeChunk) GetPosition() int {
	return p.Position
}

// SetPosition implements the Chunk interface
func (p *BlockCodeChunk) SetPosition(pos int) {
	p.Position = pos
}

// GetValue implements the Chunk interface
func (p *BlockCodeChunk) GetValue() string {
	return p.Value
}
