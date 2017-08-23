package main

import (
	"fmt"
)

// BlockTexChunk denotes a block of tex formulas
type BlockTexChunk struct {
	Position  int
	Id        string
	Caption   string //optional
	Numbering string //optional Numbering before Caption
	Value     string
}

// String implements the Stringer interface
func (p BlockTexChunk) String() string {
	return fmt.Sprintf("BlockTexChunk{Position: %d, Id: %v, Caption: %v, Value: %v}",
		p.GetPosition(), p.Id, p.Caption, p.Value)
}

// GetPosition implements the Chunk interface
func (p *BlockTexChunk) GetPosition() int {
	return p.Position
}

// SetPosition implements the Chunk interface
func (p *BlockTexChunk) SetPosition(pos int) {
	p.Position = pos
}

// GetValue implements the Chunk interface
func (p *BlockTexChunk) GetValue() string {
	return p.Value
}

func (p *BlockTexChunk) GetId() string {
	return p.Id
}

func (p *BlockTexChunk) GetCaption() string {
	return p.Caption
}

func (p *BlockTexChunk) SetCaption(c string) {
	p.Caption = c
}
func (p *BlockTexChunk) SetNumbering(c string) {
	p.Numbering = c
}
func (p BlockTexChunk) GetNumbering() string {
	return p.Numbering
}
