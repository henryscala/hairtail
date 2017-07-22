package main

import (
	"bytes"
	"errors"
	"fmt"
	"log"
)

//meta characters
const (
	EscapeChar     string = "\\"
	LeftBraceChar  string = "{"
	RightBraceChar string = "}"
	FillerChar     string = "#"
)

const (
	BlankChars string = " \r\n\t"
	LineFeed   string = "\n"
)

//keywords
const (
	RawTextChar string = "r"
)

var (
	//MetaChars contains all meta characters in slice
	MetaChars = []string{EscapeChar, LeftBraceChar, RightBraceChar, FillerChar}
	//MetaCharMap contains all meta characters in map to be lookup
	MetaCharMap = make(map[string]bool)
)

func init() {
	for _, c := range MetaChars {
		MetaCharMap[c] = true
	}
}

// Chunk denotes a block of text, the block may be length 1 to n in bytes
type Chunk interface {
	GetPosition() int
	GetValue() string
	SetPosition(pos int)
	IsTerminal() bool //return true if it has no children
}

// MetaCharChunk denotes the trunk of length 1, and the content is meta char.
// Escapted meta char should locate in PlainTextChunk instead
type MetaCharChunk struct {
	Position int
	Value    string
}

// GetPosition implements the Chunk interface
func (p *MetaCharChunk) GetPosition() int {
	return p.Position
}

// SetPosition implements the Chunk interface
func (p *MetaCharChunk) SetPosition(pos int) {
	p.Position = pos
}

// GetValue implements the Chunk interface
func (p *MetaCharChunk) GetValue() string {
	return p.Value
}

// GetValue implements the Chunk interface
func (p *MetaCharChunk) IsTerminal() bool {
	return true
}

// String implements the Stringer interface
func (p MetaCharChunk) String() string {
	return fmt.Sprintf("MetaCharChunk{Position:%v,Value:'%v'}", p.Position, p.Value)
}

//ParseChunks is the top level function to Parse input string to Chunks
//there maybe be several passes to finish parsing
func ParseChunks(input string) ([]Chunk, error) {
	var (
		chunks []Chunk
		err    error
	)
	chunks, err = RawTextChunkHandle(input)
	if err != nil {
		return chunks, err
	}
	chunks, err = MetaChunkHandle(chunks)
	if err != nil {
		return chunks, err
	}
	log.Println("-------\nafter MetaChunkHandle:", chunks)
	chunks, err = KeywordChunkHandle(chunks)
	if err != nil {
		return chunks, err
	}
	log.Println("-------\nafter KeywordChunkHandle:", chunks)

	/*
			chunks, err = EmbracedChunkHandle(chunks)
			if err != nil {
				return chunks, err
			}

		chunks, err = KeywordChunkHandle(chunks)
		if err != nil {
			return chunks, err
		}

			chunks, err = InlineChunkHandle(chunks)
			if err != nil {
				return chunks, err
			}

			chunks, err = SectionChunkHandle(chunks)
			if err != nil {
				return chunks, err
			}
	*/

	return chunks, nil
}

//MetaChunkHandle turns the chunk that is PlainTextChunk in inputChunks to MetaCharChunks if any
//It makes use of metaCharChunkHandle
func MetaChunkHandle(inputChunks []Chunk) ([]Chunk, error) {
	var (
		newChunks []Chunk
	)
	for _, chunk := range inputChunks {
		if plainTextChunk, ok := chunk.(*PlainTextChunk); ok {
			subChunks, err := metaCharChunkHandle(plainTextChunk.GetValue())
			if err != nil {
				return newChunks, err
			}
			for _, subChunk := range subChunks {
				subChunk.SetPosition(subChunk.GetPosition() + chunk.GetPosition()) //child's position plus the parent's
			}
			newChunks = append(newChunks, subChunks...)
		} else {
			newChunks = append(newChunks, chunk)
		}
	}
	return newChunks, nil
}

// metaCharChunkHandle converts string to a list of chunks that is used for latter phase handling
// the result list contains PlainTextChunk and MetaCharChunk
func metaCharChunkHandle(s string) ([]Chunk, error) {
	var buf bytes.Buffer
	var escaping = false
	var startPos = 0
	var chunks []Chunk

	for i, arune := range s {
		runeStr := string(arune)
		_, isMetaChar := MetaCharMap[runeStr]
		if escaping {
			if isMetaChar {
				//this meta char is a normal char
				buf.WriteRune(arune)
			} else {
				//last meta char is escape char

				//only store non-empty chunk before the meta char
				if buf.Len() > 0 {
					plainTextChunk := PlainTextChunk{Position: startPos, Value: buf.String()}
					chunks = append(chunks, &plainTextChunk)
				}
				//store the meta char chunk
				metaCharChunk := MetaCharChunk{Position: i - 1, Value: EscapeChar}
				chunks = append(chunks, &metaCharChunk)
				buf.Reset()
				buf.WriteRune(arune) //store the current char to buf
				startPos = i
			}
			escaping = false
			continue
		}

		//from now on, it is not in escaping state
		if isMetaChar {
			if runeStr == EscapeChar {
				escaping = true
				continue // and don't write the EscapeChar to buf
			}

			//only store non-empty chunk before the meta char
			if buf.Len() > 0 {
				plainTextChunk := PlainTextChunk{Position: startPos, Value: buf.String()}
				chunks = append(chunks, &plainTextChunk)
			}
			//store the meta char chunk
			metaCharChunk := MetaCharChunk{Position: i, Value: runeStr}
			chunks = append(chunks, &metaCharChunk)
			buf.Reset()
			startPos = i + 1
		} else {
			buf.WriteRune(arune)
		}
	}

	// there are unhandled chunk
	if buf.Len() > 0 {
		if escaping {
			return chunks, errors.New("come to the end, but it is still in escaping state")
		}
		plainTextChunk := PlainTextChunk{Position: startPos, Value: buf.String()}
		chunks = append(chunks, &plainTextChunk)
	}
	return chunks, nil
}
