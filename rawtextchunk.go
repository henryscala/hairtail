package main

import (
	"bytes"
	"errors"
	"fmt"
)

// RawTextChunk denotes raw text without any escape inside. Even meta chars in side it are treated as raw text.
type RawTextChunk struct {
	Position int
	Value    string
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

// RawTextChunkHandle handles raw text chunks, chunks other than RawTextChunk will be stored as PlainTextChunk
// the result list contains RawTextChunk and PlainTextChunk
func RawTextChunkHandle(input string) ([]Chunk, error) {
	//the states of the handling
	//it is not possible to compare whether a variable equals a closure(a closure variable can only compare to nil),
	//so use these states to know the present state
	const (
		waitEscapeChar          int = iota //Escape Char is \
		waitRawTextChar                    //RawTextChar is r
		waitFillCharOrLeftBrace            //Fille char is #
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

	//check if there are remaining content in the buf at last
	if buf.Len() > 0 {
		chunk := &PlainTextChunk{Position: startPos, Value: buf.String()}
		chunks = append(chunks, chunk) //store the previous plain text chunk
	}

	return chunks, nil
}
