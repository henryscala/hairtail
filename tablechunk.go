package main

import (
	"fmt"
)

// TableChunk denotes table
type TableChunk struct {
	Position  int
	Id        string
	Caption   string //optional
	Numbering string //optional Numbering before Caption
	Cells     [][]Chunk
}

// String implements the Stringer interface
func (p TableChunk) String() string {
	return fmt.Sprintf("TableChunk{Position: %d, Id: %v, Caption: %v, Cells: %v}",
		p.GetPosition(), p.Id, p.Caption, p.Cells)
}

// GetPosition implements the Chunk interface
func (p *TableChunk) GetPosition() int {
	return p.Position
}

// SetPosition implements the Chunk interface
func (p *TableChunk) SetPosition(pos int) {
	p.Position = pos
}

// GetValue implements the Chunk interface
func (p *TableChunk) GetValue() string {
	return p.Id
}

func (p *TableChunk) GetId() string {
	return p.Id
}

func (p *TableChunk) GetCaption() string {
	return p.Caption
}

func (p *TableChunk) SetCaption(c string) {
	p.Caption = c
}
func (p *TableChunk) SetNumbering(c string) {
	p.Numbering = c
}
func (p TableChunk) GetNumbering() string {
	return p.Numbering
}
