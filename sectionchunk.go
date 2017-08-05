package main

import (
	"fmt"
)

// SectionChunk denotes Sections in article. It is nested structure.
type SectionChunk struct {
	Position  int
	Level     int //1 2 .. 6
	Id        string
	Caption   string
	Numbering string //optional Numbering before Caption
	Children  []Chunk
}

// String implements the Stringer interface
func (p SectionChunk) String() string {
	return fmt.Sprintf("SectionChunk{Position: %d, Level: %d, Id: %v, Caption: %v, Numbering:%v,Children: %v}",
		p.GetPosition(), p.Level, p.Id, p.Caption, p.Numbering, p.Children)
}

// GetPosition implements the Chunk interface
func (p *SectionChunk) GetPosition() int {
	return p.Position
}

// SetPosition implements the Chunk interface
func (p *SectionChunk) SetPosition(pos int) {
	p.Position = pos
}

// GetValue implements the Chunk interface
func (p *SectionChunk) GetValue() string {

	return p.Caption

}
func (p *SectionChunk) SetNumbering(c string) {
	p.Numbering = c
}
