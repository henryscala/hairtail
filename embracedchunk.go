package main

import (
	"fmt"
	"errors"
)

// EmbracedChunk denotes content inside braces {}
// EmbraceChunk may be terminal or non-terminal
// terminal chunk contains only one PlainTextChunk/RawTextChunk as child, while non-terminal chunk contains more Children chunks.
// and the contained chunks may be EmbracedChunks too. So it is a tree structure.
type EmbracedChunk struct {
	Position int
	Children []Chunk
}

func (p EmbracedChunk) IsTerminal() bool {
	if len(p.Children) == 0 {
		panic("A embraced Chunk cannot have 0 child")
	}
	if len(p.Children) != 1 {
		return false
	}
	chunk := p.Children[0]
	 _, isPlainTextChunk := chunk.(*PlainTextChunk)
	_, isRawTextChunk := chunk.(*PlainTextChunk)
	if isPlainTextChunk || isRawTextChunk {
		return true
	}
	return false
}

func (p EmbracedChunk) getTerminalValue() string {
	if ! p.IsTerminal() {
		panic("it is not terminal")
	}

	chunk := p.Children[0]
	return chunk.GetValue()
}

// String implements the Stringer interface
func (p EmbracedChunk) String() string {
	return fmt.Sprintf("EmbracedChunk{Position:%v,Children:'%v'}", p.Position, p.Children)
}

// GetPosition implements the Chunk interface
func (p *EmbracedChunk) GetPosition() int {
	return p.Position
}

// SetPosition implements the Chunk interface
func (p *EmbracedChunk) SetPosition(pos int)  {
	p.Position = pos
}

// GetValue implements the Chunk interface
func (p *EmbracedChunk) GetValue() string {
	if p.IsTerminal() {
		return p.getTerminalValue()
	}
	//TODO how to handle non-terminal value
	panic("not implemented, should not reach here")
	return ""
}

//EmbracedChunkHandle parse the inputChunks
// and combine Chunks between LeftBrace(MetaCharChunk) and RightBrace(MetaCharChunk) as EmbracedChunk
func EmbracedChunkHandle(inputChunks []Chunk) ([]Chunk, error) {
	var (
		chunkStack []Chunk
	)

	//find  MetaCharChunk that contains LeftBraceChar from top of the stack to the bottom
	findLeftBraceChunk := func() (int, error) {
		for i:= len(chunkStack)-1; i>=0; i-- {
			metaCharChunk, isMetaCharChunk := chunkStack[i].(*MetaCharChunk)
			if isMetaCharChunk && metaCharChunk.GetValue() == LeftBraceChar {
				return i, nil
			}
		}
		return -1, errors.New("Unbalanced LeftBraceChunk and RightBraceChunk")
	}

	for _,chunk := range inputChunks {
		metaCharChunk, isMetaCharChunk := chunk.(*MetaCharChunk)
		if isMetaCharChunk && metaCharChunk.GetValue() == RightBraceChar {
			leftBraceIndex, err := findLeftBraceChunk()
			if err!= nil {
				return chunkStack, err
			}
			var embracedChunks []Chunk //new slice, not piggyback on original slice
			embracedChunks  = append(embracedChunks, chunkStack[leftBraceIndex+1:]...)
			//ignore empty chunk (no contents between LeftBrace and RightBrace
			if len(embracedChunks) == 0 {
				continue
			}

			//combine chunks from leftBrace to rightBrace
			embracedChunk := &EmbracedChunk{Position:embracedChunks[0].GetPosition(), Children:embracedChunks }
			//pop handled chunks
			chunkStack = chunkStack[:leftBraceIndex]

			//push to the stack
			chunkStack = append(chunkStack, embracedChunk)
		} else {
			chunkStack = append(chunkStack, chunk)
		}
	}
	return chunkStack, nil
}
