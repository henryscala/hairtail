package main

import (
	"bytes"
	"errors"
	"fmt"
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

// RawTextChunk denotes raw text without any escape inside. Even meta chars in side it are treated as raw text.
type RawTextChunk struct {
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

// GetPosition implements the Chunk interface
func (p *RawTextChunk) GetPosition() int {
	return p.Position
}

// SetPosition implements the Chunk interface
func (p *RawTextChunk) SetPosition(pos int) {
	p.Position = pos
}

// GetValue implements the Chunk interface
func (p *RawTextChunk) GetValue() string {
	return p.Value
}

// GetValue implements the Chunk interface
func (p *RawTextChunk) IsTerminal() bool {
	return true
}

// String implements the Stringer interface
func (p RawTextChunk) String() string {
	return fmt.Sprintf("RawTextChunk{Position:%v,Value:'%v'}", p.Position, p.Value)
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

	return chunks, nil
}

// RawTextChunkHandle handles raw text chunks, chunks other than RawTextChunk will be stored as PlainTextChunk
// the result list contains RawTextChunk and PlainTextChunk
func RawTextChunkHandle(input string) ([]Chunk, error) {
	//the states of the handling
	//it is not possible to compare whether a variable equals a closure(a closure variable can only compare to nil),
	//so use these states to know the present state
	const (
		waitEscapeChar          int = iota //Escape Char is \
		waitRawTextChar                    //RawTextChar is r
		waitFillCharOrLeftBrace            //Fille char is ~
		waitRawTextEnd                     //RawTextEnd is calculated
	)

	type stateHandle func(int, rune) (stateHandle, error)
	var (
		waitEscapeCharStateHandle          stateHandle
		waitRawTextCharStateHandle         stateHandle
		waitFillCharOrLeftBraceStateHandle stateHandle
		waitRawTextEndStateHandle          stateHandle
		buf                                bytes.Buffer
		chunks                             []Chunk
		numIgnoreChar                      int
		startPos                           int
		rawTextFillStr                     string
		state                              int
	)

	waitEscapeCharStateHandle = func(pos int, arune rune) (stateHandle, error) {
		astr := string(arune)
		if astr == EscapeChar {
			state = waitRawTextChar
			return waitRawTextCharStateHandle, nil //state changes
		}
		buf.WriteRune(arune)
		return waitEscapeCharStateHandle, nil //state not changes
	}

	waitRawTextCharStateHandle = func(pos int, arune rune) (stateHandle, error) {
		astr := string(arune)
		if astr == RawTextChar {
			rawTextFillStr = "" //reset
			state = waitFillCharOrLeftBrace
			return waitFillCharOrLeftBraceStateHandle, nil //state changes
		}
		buf.WriteString(EscapeChar) //write the EscapeChar in last position
		buf.WriteRune(arune)        //write the current Char in this position
		state = waitEscapeChar
		return waitEscapeCharStateHandle, nil //state changes to original
	}

	waitFillCharOrLeftBraceStateHandle = func(pos int, arune rune) (stateHandle, error) {
		astr := string(arune)
		if astr == FillerChar {
			rawTextFillStr += FillerChar
			return waitFillCharOrLeftBraceStateHandle, nil // state not change
		}
		if astr == LeftBraceChar {
			if buf.Len() > 0 {
				chunk := &PlainTextChunk{Position: startPos, Value: buf.String()}
				chunks = append(chunks, chunk) //store the previous plain text chunk
			}
			buf.Reset()
			startPos = pos + 1
			state = waitRawTextEnd
			return waitRawTextEndStateHandle, nil //state changes
		}
		return nil, errors.New("expect either raw text filler char(~) or LeftBraceChar")
	}

	waitRawTextEndStateHandle = func(pos int, arune rune) (stateHandle, error) {
		rawTextRightDelimiter := RightBraceChar + rawTextFillStr
		totalLen := len(input)
		delimiterLen := len(rawTextRightDelimiter)
		if pos+delimiterLen > totalLen {
			return nil, errors.New("come to the end, but does not find the end raw text chunk")
		}
		if input[pos:pos+delimiterLen] == rawTextRightDelimiter {
			if buf.Len() > 0 {
				chunk := &RawTextChunk{Position: startPos, Value: buf.String()}
				chunks = append(chunks, chunk) //store the Raw text chunk
			}
			buf.Reset()
			startPos = pos + 1
			numIgnoreChar = delimiterLen - 1
			state = waitEscapeChar
			return waitEscapeCharStateHandle, nil //state changes to original
		}
		buf.WriteRune(arune)
		return waitRawTextEndStateHandle, nil //state not change
	}

	handler := waitEscapeCharStateHandle
	state = waitEscapeChar

	for i, arune := range input {
		if numIgnoreChar > 0 {
			numIgnoreChar--
			continue
		}
		newHandler, err := handler(i, arune)
		if err != nil {
			return chunks, err
		}
		handler = newHandler
	}

	//check if the current state is legal
	if state != waitEscapeChar {
		return nil, fmt.Errorf("illegal state after handling the whole input text state=%d", state)
	}

	//check if there are remaining containts in the buf at last
	if buf.Len() > 0 {
		chunk := &PlainTextChunk{Position: startPos, Value: buf.String()}
		chunks = append(chunks, chunk) //store the previous plain text chunk
	}

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
