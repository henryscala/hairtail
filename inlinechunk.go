package main

import (
	"errors"
	"fmt"
)

// InlineChunk denotes chunk for inline fomrmat

type InlineChunk struct {
	Position int
	Keyword  string
	Children []Chunk
}

func (p InlineChunk) IsTerminal() bool {
	if len(p.Children) == 0 {
		panic("inline format shall always has children")
	}
	if len(p.Children) == 1 {
		return true
	}
	return false
}

// String implements the Stringer interface
func (p InlineChunk) String() string {
	return fmt.Sprintf("InlineChunk{Position:%v,Keyword:%v, Children:%v}", p.Position, p.Keyword, p.Children)
}

// GetPosition implements the Chunk interface
func (p *InlineChunk) GetPosition() int {
	return p.Position
}

// SetPosition implements the Chunk interface
func (p *InlineChunk) SetPosition(pos int) {
	p.Position = pos
}

// GetValue implements the Chunk interface
func (p *InlineChunk) GetValue() string {
	panic("should no use,  except debug")
	return p.Keyword
}

// KeywordChunkHandle parse the inputChunks that is keywordChunk.
// if it is KeywordChunk and is for inline-format, convert it
func InlineChunkHandle(inputChunks []Chunk) ([]Chunk, error) {
	var outputChunks []Chunk
	var err error
	//first pass, convert keywordChunk to inlineChunk
	for i := 0; i < len(inputChunks); i++ {
		chunk := inputChunks[i]
		keywordChunk, isKeywordChunk := chunk.(*KeywordChunk)

		if !isKeywordChunk {
			outputChunks = append(outputChunks, chunk)
			continue
		}

		if _, ok := gInlineFormatKeywordMap[keywordChunk.Keyword]; !ok {
			outputChunks = append(outputChunks, chunk)
			continue
		}

		if i+1 >= len(inputChunks) {
			return outputChunks, errors.New("inline format does not follow a chunk")
		}
		embracedChunk, isEmbracedChunk := inputChunks[i+1].(*EmbracedChunk)
		if !isEmbracedChunk {
			return outputChunks, errors.New("inline format does not follow a EmbracedChunk")
		}

		inlineChunk := &InlineChunk{Position: keywordChunk.Position,
			Keyword:  keywordChunk.Keyword,
			Children: embracedChunk.Children,
		}
		inlineChunk.Children, err = InlineChunkHandle(inlineChunk.Children) //recursive call
		if err != nil {
			return outputChunks, errors.New("inline chunk handle failed in children")
		}
		outputChunks = append(outputChunks, inlineChunk)
		i++ //jump the embraced chunk
	}
	//second pass, render inlineChunk. If the inlineChunk neighbors with PlainTextChunk, then merge them.
	outputChunks, err = InlineChunkListRender(outputChunks)
	return outputChunks, err
}
