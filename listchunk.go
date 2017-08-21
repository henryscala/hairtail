package main

import (
	"fmt"
)

type ListItem struct {
	Value []Chunk //may contain inline chunks
}

// ListChunk denotes a order/bullet list
type ListChunk struct {
	Position  int
	ListType  string
	Id        string
	Caption   string //optional
	Numbering string //optional Numbering before Caption
	Items     []*ListItem
}

// String implements the Stringer interface
func (p ListChunk) String() string {
	return fmt.Sprintf("ListChunk{Position: %d, ListType %v, Id: %v, Caption: %v, Items: %v}",
		p.GetPosition(), p.ListType, p.Id, p.Caption, p.Items)
}

// GetPosition implements the Chunk interface
func (p *ListChunk) GetPosition() int {
	return p.Position
}

// SetPosition implements the Chunk interface
func (p *ListChunk) SetPosition(pos int) {
	p.Position = pos
}

// GetValue implements the Chunk interface
func (p *ListChunk) GetValue() string {
	return p.Id
}
func (p *ListChunk) GetId() string {
	return p.Id
}

func (p *ListChunk) GetCaption() string {
	return p.Caption
}

func (p *ListChunk) SetCaption(c string) {
	p.Caption = c
}

func (p *ListChunk) SetNumbering(c string) {
	p.Numbering = c
}
func (p ListChunk) GetNumbering() string {
	return p.Numbering
}
