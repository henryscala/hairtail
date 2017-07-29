package main

import (
	"fmt"
)

// ReferToChunk denotes anchors(inner article link) in article.
type ReferToChunk struct {
	Position int
	Id       string
	Value    string //if required, the value is get from the referred to chunk
}

// String implements the Stringer interface
func (p ReferToChunk) String() string {
	return fmt.Sprintf("ReferToChunk{Position: %d, Id: %v, Value: %v }",
		p.GetPosition(), p.Id, p.GetValue())
}

// GetPosition implements the Chunk interface
func (p *ReferToChunk) GetPosition() int {
	return p.Position
}

// SetPosition implements the Chunk interface
func (p *ReferToChunk) SetPosition(pos int) {
	p.Position = pos
}

// GetValue implements the Chunk interface
func (p *ReferToChunk) GetValue() string {
	return p.Value
}
