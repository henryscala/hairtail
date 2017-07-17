package main

import (
	"errors"
	"fmt"
	"strings"
)

// KeywordChunk denotes meta char '\' followed by a token.
// during handling, the meta char '\' is omitted.
type KeywordChunk struct {
	Position int
	Keyword  string
}

func (p KeywordChunk) IsTerminal() bool {
	return true
}

// String implements the Stringer interface
func (p KeywordChunk) String() string {
	return fmt.Sprintf("keywordChunk{Position:%v,Keyword:'%v'}", p.Position, p.Keyword)
}

// GetPosition implements the Chunk interface
func (p *KeywordChunk) GetPosition() int {
	return p.Position
}

// SetPosition implements the Chunk interface
func (p *KeywordChunk) SetPosition(pos int) {
	p.Position = pos
}

// GetValue implements the Chunk interface
func (p *KeywordChunk) GetValue() string {
	return p.Keyword
}

// KeywordChunkHandle parse the inputChunks(which contains nested chunks)
// and combine the metachar \ followed by token as KeywordChunk
func KeywordChunkHandle(inputChunks []Chunk) ([]Chunk, error) {
	var (
		outputChunks []Chunk
		err          error
	)

	for i := 0; i < len(inputChunks); i++ {
		inputChunk := inputChunks[i]
		//it is non-terminal
		if embracedChunk, isEmbracedChunk := inputChunk.(*EmbracedChunk); isEmbracedChunk {
			embracedChunk.Children, err = KeywordChunkHandle(embracedChunk.Children) //recursive call
			if err != nil {
				return outputChunks, err
			}
			outputChunks = append(outputChunks, embracedChunk)
			continue
		}

		if !inputChunk.IsTerminal() {
			panic(fmt.Sprintf("%T:%v -- %v", inputChunk, inputChunk, "from here, the chunk should be terminal"))
		}

		if metaCharChunk, isMetaCharChunk := inputChunk.(*MetaCharChunk); isMetaCharChunk {
			if metaCharChunk.GetValue() != EscapeChar {
				outputChunks = append(outputChunks, inputChunk)
				continue
			}

			if i+1 >= len(inputChunks) {
				return outputChunks, errors.New("Escape char followed by nothing")
			}
			nextChunk := inputChunks[i+1]
			plainTextChunk, isPlainTextChunk := nextChunk.(*PlainTextChunk)
			if !isPlainTextChunk {
				return outputChunks, errors.New("Escape char should followed by plain text so that we can get a token")
			}
			token := gTokenPattern.FindString(plainTextChunk.GetValue())
			if len(token) == 0 || !strings.HasPrefix(plainTextChunk.GetValue(), token) {
				return outputChunks, errors.New("no token is found following Escape char")
			}

			keywordChunk := &KeywordChunk{Position: plainTextChunk.GetPosition(), Keyword: token}
			outputChunks = append(outputChunks, keywordChunk)
			plainTextChunk.Value = plainTextChunk.Value[len(token):]
			plainTextChunk.SetPosition(plainTextChunk.GetPosition() + len(token))
			//if there is remaining contents
			if len(plainTextChunk.Value) > 0 {
				outputChunks = append(outputChunks, plainTextChunk)
			}
			i++ //jump the nextChunk, because it is plainText and is handled
			continue
		}

		outputChunks = append(outputChunks, inputChunk)
	}

	return outputChunks, nil
}
