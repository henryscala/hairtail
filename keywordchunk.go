package main

import (
	"errors"
	"fmt"
	"log"
	"strings"
)

var (
	errExpectToken     = errors.New("expect token")
	errIndexOutOfBound = errors.New("index out of bound")
	errExpectLBrace    = errors.New("expect Left Brace")
	errExpectRBrace    = errors.New("expect Right Brace")
)

// KeywordChunk denotes meta char '\' followed by a token.
// during handling, the meta char '\' is omitted.
type KeywordChunk struct {
	Position int
	Keyword  string
	Children []Chunk
}

func (p KeywordChunk) IsTerminal() bool {
	return true
}

// String implements the Stringer interface
func (p KeywordChunk) String() string {
	return fmt.Sprintf("keywordChunk{Position:%v,Keyword:'%v',Children:%v}", p.Position, p.Keyword, p.Children)
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

// KeywordChunkHandle parse the inputChunks(which contains MetaCharChunk, PlaintextChunk, RawTextChunk )
func KeywordChunkHandle(inputChunks []Chunk) ([]Chunk, error) {
	var (
		outputChunks []Chunk
		index        int
	)

	for index < len(inputChunks) {
		inputChunk := inputChunks[index]

		if metaCharChunk, isMetaCharChunk := inputChunk.(*MetaCharChunk); isMetaCharChunk {
			if metaCharChunk.GetValue() != EscapeChar {
				outputChunks = append(outputChunks, inputChunk)
				index++
				continue
			}

			token, newIndex, err := consumeToken(inputChunks, index+1)
			if err != nil {
				log.Println(err)
				return outputChunks, err
			}

			log.Println("== handle token:", token[0].GetValue())

			switch token[0].GetValue() {
			case EmphasisFormat, StrongFormat:

				chunks, newIndex, err := consumeEmbracedBlock(inputChunks, newIndex)
				if err != nil {
					log.Println(err)
					return outputChunks, err
				}
				chunks, err = KeywordChunkHandle(chunks) //recursive
				if err != nil {
					log.Println(err)
					return outputChunks, err
				}
				keywordChunk := &KeywordChunk{Position: token[0].GetPosition(),
					Keyword:  token[0].GetValue(),
					Children: []Chunk{chunks[1]},
				}
				outputChunks = append(outputChunks, keywordChunk)
				index = newIndex

			case HyperLink:
				log.Println("handle first embraced block", newIndex)
				chunksUrl, newIndex, err := consumeEmbracedBlock(inputChunks, newIndex)
				if err != nil {
					log.Fatalln(err)
					return outputChunks, err
				}
				chunksUrl, err = KeywordChunkHandle(chunksUrl) //recursive
				if err != nil {
					log.Fatalln(err)
					return outputChunks, err
				}
				log.Println("handle second embraced block", newIndex)
				chunksContent, newIndex, err := consumeEmbracedBlock(inputChunks, newIndex)
				if err != nil {
					log.Fatalln(err)
					return outputChunks, err
				}
				chunksContent, err = KeywordChunkHandle(chunksContent) //recursive
				if err != nil {
					log.Fatalln(err)
					return outputChunks, err
				}

				keywordChunk := &KeywordChunk{Position: token[0].GetPosition(),
					Keyword:  token[0].GetValue(),
					Children: []Chunk{chunksUrl[1], chunksContent[1]},
				}
				outputChunks = append(outputChunks, keywordChunk)
				index = newIndex

			default:
				panic("not implemented")
			}
			continue
		} else {

			outputChunks = append(outputChunks, inputChunk)
		}
		index++
	}

	return outputChunks, nil
}

func consumeEmbracedBlock(inputChunks []Chunk, index int) (chunks []Chunk, newIndex int, err error) {
	if index >= len(inputChunks) {
		log.Fatalln(errIndexOutOfBound)
		return nil, index, errIndexOutOfBound
	}
	leftBraceChunk, ok := inputChunks[index].(*MetaCharChunk)
	if !ok || leftBraceChunk.GetValue() != LeftBraceChar {
		log.Fatalln(errExpectLBrace)
		return nil, index, errExpectLBrace
	}

	var countLeftBrace int = 1
	chunks = append(chunks, leftBraceChunk)
	i := 0
	for i = index + 1; i < len(inputChunks); i++ {
		chunks = append(chunks, inputChunks[i])
		chunk, ok := inputChunks[i].(*MetaCharChunk)
		if ok {
			if chunk.GetValue() == RightBraceChar {
				countLeftBrace--

				if countLeftBrace == 0 {
					break
				}
			} else if chunk.GetValue() == LeftBraceChar {
				countLeftBrace++
			}
		}
	}
	rightBraceChunk, ok := chunks[len(chunks)-1].(*MetaCharChunk)
	if ok && rightBraceChunk.GetValue() == RightBraceChar {
		return chunks, i + 1, nil
	}
	log.Fatalln(errExpectRBrace)
	return chunks, index, errExpectRBrace

}

func cosumeEmbracedToken(inputChunks []Chunk, index int) (chunks []Chunk, newIndex int, err error) {
	chunks1, newIndex, err := consumeEmbracedBlock(inputChunks, index)
	if len(chunks) != 3 {
		return chunks, newIndex, errExpectToken
	}

	chunks2, _, err := consumeToken(inputChunks, index+1)
	if err != nil {
		return nil, index, err
	}

	return []Chunk{chunks1[0], chunks2[1], chunks1[2]}, index + 3, nil
}

func consumeToken(inputChunks []Chunk, index int) (chunks []Chunk, newIndex int, err error) {
	if index >= len(inputChunks) {
		return nil, index, errIndexOutOfBound
	}
	plainTextChunk, ok := inputChunks[index].(*PlainTextChunk)
	if !ok {
		return nil, index, errExpectToken
	}
	text := strings.Trim(plainTextChunk.GetValue(), BlankChars)
	token := gTokenPattern.FindString(text)
	if token == text && len(token) > 0 {
		plainTextChunk.Value = token
		return []Chunk{plainTextChunk}, index + 1, nil
	}
	return nil, index, errExpectToken
}
