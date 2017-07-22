package main

import (
	"errors"
	"fmt"
	"log"
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

	isFollowedChunksValid := func( index, num int ) error {
		for i:=index; i< index+num; i++ {
			_, isEmbracedChunk := inputChunks[i].(*EmbracedChunk)
			if !isEmbracedChunk {
				return errors.New("inline format followed by chunks other than EmbracedChunk")
			}
		}
		return nil
	}

	accumulateChunks := func(index, num int) []Chunk {
		return append([]Chunk{}, inputChunks[index:index+num]... )
	}

	//first pass, convert keywordChunk to inlineChunk
	for i := 0; i < len(inputChunks); i++ {
		chunk := inputChunks[i]
		keywordChunk, isKeywordChunk := chunk.(*KeywordChunk)

		if !isKeywordChunk {
			outputChunks = append(outputChunks, chunk)
			continue
		}
		inlineFormatDescription, isInlineFormat := gInlineFormatKeywordMap[keywordChunk.Keyword]
		if  !isInlineFormat{
			outputChunks = append(outputChunks, chunk)
			continue
		}

		if i + inlineFormatDescription.NumEmbracedBlock >= len(inputChunks) {
			return outputChunks, fmt.Errorf("inline format does not follow enough number(%d) of chunks ", inlineFormatDescription.NumEmbracedBlock)
		}

		err = isFollowedChunksValid(i+1, inlineFormatDescription.NumEmbracedBlock)
		if err != nil  {
			return outputChunks, err
		}

		inlineChunk := &InlineChunk{Position: keywordChunk.Position,
			Keyword:  keywordChunk.Keyword,
			Children: accumulateChunks(i+1,inlineFormatDescription.NumEmbracedBlock),
		}

		inlineChunk.Children, err = InlineChunkHandle(inlineChunk.Children) //recursive call
		if err != nil {
			return outputChunks, errors.New("inline chunk handle failed in children")
		}

		outputChunks = append(outputChunks, inlineChunk)
		i+=inlineFormatDescription.NumEmbracedBlock //jump the embraced chunk
	}
	log.Println("inlineChunk before render:", outputChunks)
	//second pass, render inlineChunk. If the inlineChunk neighbors with PlainTextChunk, then merge them.
	outputChunks, err = InlineChunkListRender(outputChunks)
	return outputChunks, err
}
