package main

import (
	"fmt"
)

// AnchorChunk denotes anchors(inner article link) in article.
type AnchorChunk struct {
	Position int
	Id       string
	Value    string
}

func (p AnchorChunk) IsTerminal() bool {
	return true
}

// String implements the Stringer interface
func (p AnchorChunk) String() string {
	return fmt.Sprintf("AnchorChunk{Position: %d, Id: %v, Value: %v }",
		p.GetPosition(), p.Id, p.GetValue())
}

// GetPosition implements the Chunk interface
func (p *AnchorChunk) GetPosition() int {
	return p.Position
}

// SetPosition implements the Chunk interface
func (p *AnchorChunk) SetPosition(pos int) {
	p.Position = pos
}

// GetValue implements the Chunk interface
func (p *AnchorChunk) GetValue() string {
	return p.Value
}
