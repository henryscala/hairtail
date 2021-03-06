package main

import (
	"fmt"
)

// ImageChunk denotes table
type ImageChunk struct {
	Position  int
	Id        string
	Caption   string //optional
	Numbering string //optional Numbering before Caption
	Src       string
}

// String implements the Stringer interface
func (p ImageChunk) String() string {
	return fmt.Sprintf("ImageChunk{Position: %d, Id: %v, Caption: %v, Src: %v}",
		p.GetPosition(), p.Id, p.Caption, p.Src)
}

// GetPosition implements the Chunk interface
func (p *ImageChunk) GetPosition() int {
	return p.Position
}

// SetPosition implements the Chunk interface
func (p *ImageChunk) SetPosition(pos int) {
	p.Position = pos
}

// GetValue implements the Chunk interface
func (p *ImageChunk) GetValue() string {
	return p.Id
}

func (p *ImageChunk) GetId() string {
	return p.Id
}

func (p *ImageChunk) GetCaption() string {
	return p.Caption
}

func (p *ImageChunk) SetCaption(c string) {
	p.Caption = c
}

func (p *ImageChunk) SetNumbering(c string) {
	p.Numbering = c
}
func (p ImageChunk) GetNumbering() string {
	return p.Numbering
}
