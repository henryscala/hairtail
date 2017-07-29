package main

import (
	"errors"
	"fmt"
	"math"
	"strings"
)

// PlainTextChunk denotes a block of plain text without any meta char
type PlainTextChunk struct {
	Position int
	Value    string
}

// GetPosition implements the Chunk interface
func (p *PlainTextChunk) GetPosition() int {
	return p.Position
}

// SetPosition implements the Chunk interface
func (p *PlainTextChunk) SetPosition(pos int) {
	p.Position = pos
}

// GetValue implements the Chunk interface
func (p *PlainTextChunk) GetValue() string {
	return p.Value
}

// Divide the chunk to two paragraphs separated by two \n or \r\n
func (p *PlainTextChunk) ToParagraphList() (paragraphList []string) {
	list := gParagraphDivider.Split(p.GetValue(), math.MaxInt64)
	for _, paragraph := range list {
		//ignore blank line
		if len(strings.Trim(paragraph, BlankChars)) == 0 {
			continue
		}
		paragraphList = append(paragraphList, paragraph)
	}
	return
}

// Divide the chunk to two parts by \n, \r\n, \r
func (p *PlainTextChunk) FirstLineRestLines() (firstLine Chunk, restLines Chunk, err error) {
	parts := strings.SplitN(p.GetValue(), LineFeed, 2)
	if len(parts) < 1 {
		err = errors.New("PlainTextChunk has less than one line")
		return
	}
	firstLine = &PlainTextChunk{Position: p.GetPosition(), Value: strings.Trim(parts[0], BlankChars)}
	if len(parts) > 1 {
		restLines = &PlainTextChunk{Position: p.GetPosition() + len(parts[0]) + 1, Value: parts[1]}
	}
	return
}

// String implements the Stringer interface
func (p PlainTextChunk) String() string {
	return fmt.Sprintf("PlainTextChunk{Position:%v,Value:'%v'}", p.Position, p.Value)
}
