package main

import (
	"errors"
	"fmt"
	"log"
	"runtime/debug"
	"strings"
)

var (
	errExpectToken     = errors.New("expect token")
	errIndexOutOfBound = errors.New("index out of bound")
	errExpectLBrace    = errors.New("expect Left Brace")
	errExpectRBrace    = errors.New("expect Right Brace")
	errExpectPlainText = errors.New("expect Plain Text")
	errExpectRawText   = errors.New("expect Raw Text")
	errExpectListItem  = errors.New("expect List Item ")
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
			case InlineCode:
				chunks, newIndex1, err := consumeEmbracedBlock(inputChunks, newIndex)
				if err == nil {

					keywordChunk := &KeywordChunk{Position: token[0].GetPosition(),
						Keyword:  token[0].GetValue(),
						Children: []Chunk{chunks[1]},
					}
					outputChunks = append(outputChunks, keywordChunk)
					newIndex = newIndex1
					index = newIndex
				} else {
					//InlineCode content may be either EmbracedBlock or RawTextBlock
					if newIndex >= len(inputChunks) {
						log.Fatalln(errIndexOutOfBound)
						return outputChunks, errIndexOutOfBound
					}
					rawTextChunk, ok := inputChunks[newIndex].(*RawTextChunk)
					if !ok {
						log.Fatalln(errExpectRawText)
						return outputChunks, errExpectRawText
					}

					keywordChunk := &KeywordChunk{Position: token[0].GetPosition(),
						Keyword:  token[0].GetValue(),
						Children: []Chunk{rawTextChunk},
					}

					outputChunks = append(outputChunks, keywordChunk)
					newIndex++
					index = newIndex
				}
			case ListItemMark:
				keywordChunk := &KeywordChunk{Position: token[0].GetPosition(),
					Keyword: token[0].GetValue(),
				}
				outputChunks = append(outputChunks, keywordChunk)

				index = newIndex

			case BlockCode:

				tokenChunks, newIndex, err := cosumeEmbracedToken(inputChunks, newIndex)

				if err != nil {
					log.Fatalln(err)
					return outputChunks, err
				}
				chunksContent, newIndex1, err := consumeEmbracedBlock(inputChunks, newIndex)
				if err == nil {
					blockCodeChunk := &BlockCodeChunk{
						Position: token[0].GetPosition(),
						Id:       tokenChunks[1].GetValue(),
						Value:    chunksContent[1].GetValue(),
					}
					keywordChunk := &KeywordChunk{Position: token[0].GetPosition(),
						Keyword:  token[0].GetValue(),
						Children: []Chunk{blockCodeChunk},
					}

					outputChunks = append(outputChunks, keywordChunk)
					index = newIndex1
				} else {
					//BlockCode content may be either EmbracedBlock or RawTextBlock
					if newIndex >= len(inputChunks) {
						log.Fatalln(errIndexOutOfBound)
						return outputChunks, errIndexOutOfBound
					}
					rawTextChunk, ok := inputChunks[newIndex].(*RawTextChunk)
					if !ok {
						log.Fatalln(errExpectRawText)
						return outputChunks, errExpectRawText
					}

					blockCodeChunk := &BlockCodeChunk{
						Position: token[0].GetPosition(),
						Id:       tokenChunks[1].GetValue(),
						Value:    rawTextChunk.GetValue(),
					}
					keywordChunk := &KeywordChunk{Position: token[0].GetPosition(),
						Keyword:  token[0].GetValue(),
						Children: []Chunk{blockCodeChunk},
					}

					outputChunks = append(outputChunks, keywordChunk)
					newIndex++
					index = newIndex
				}
			case SectionHeader, SectionHeader1, SectionHeader2, SectionHeader3,
				SectionHeader4, SectionHeader5, SectionHeader6:

				header := token[0].GetValue()
				tokenChunks, newIndex, err := cosumeEmbracedToken(inputChunks, newIndex)

				if err != nil {
					log.Fatalln(err)
					return outputChunks, err
				}
				if newIndex > len(inputChunks) {
					log.Fatalln(errIndexOutOfBound)
					return outputChunks, errIndexOutOfBound
				}

				plainTextChunk, ok := inputChunks[newIndex].(*PlainTextChunk)
				newIndex++
				if !ok {
					log.Fatalln(errExpectPlainText)
					return outputChunks, errExpectPlainText
				}
				firstLineChunk, restLineChunk, err := plainTextChunk.FirstLineRestLines()
				if err != nil {
					log.Fatalln(errExpectPlainText)
					return outputChunks, errExpectPlainText
				}
				//update the plainTextChunk in-place
				if restLineChunk == nil {
					plainTextChunk.Value = ""
				} else {
					plainTextChunk.Value = restLineChunk.GetValue()
					plainTextChunk.Position = restLineChunk.GetPosition()
				}
				level := gSectionLevel[header]

				sectionChunk := &SectionChunk{Position: token[0].GetPosition(),
					Level: level, Caption: firstLineChunk.GetValue(),
					Id: tokenChunks[1].GetValue(),
				}
				keywordChunk := &KeywordChunk{Position: token[0].GetPosition(),
					Keyword:  token[0].GetValue(),
					Children: []Chunk{sectionChunk},
				}
				outputChunks = append(outputChunks, keywordChunk)
				outputChunks = append(outputChunks, plainTextChunk)
				index = newIndex
			case OrderList, BulletList:
				tokenChunks, newIndex, err := cosumeEmbracedToken(inputChunks, newIndex)
				if err != nil {
					log.Fatalln(err)
					return outputChunks, err
				}
				chunksContent, newIndex, err := consumeEmbracedBlock(inputChunks, newIndex)
				if err != nil {
					log.Fatalln(err)
					return outputChunks, err
				}

				chunksContent, err = KeywordChunkHandle(chunksContent)
				if err != nil {
					log.Fatalln(err)
					return outputChunks, err
				}

				listChunk := &ListChunk{Position: token[0].GetPosition(),
					Id:       tokenChunks[1].GetValue(),
					ListType: token[0].GetValue(),
				}

				items, _, err := consumeListItems(chunksContent[1:len(chunksContent)-1], 0)
				if err != nil {
					log.Fatalln(err)
					return outputChunks, err
				}
				listChunk.Items = items
				keywordChunk := &KeywordChunk{Position: token[0].GetPosition(),
					Keyword:  token[0].GetValue(),
					Children: []Chunk{listChunk},
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

func consumeListItem(inputChunks []Chunk, index int) (item *ListItem, newIndex int, err error) {
	ignoreBlank := func() {

		for index < len(inputChunks) {
			plainTextChunk, ok := inputChunks[index].(*PlainTextChunk)
			if ok && len(strings.Trim(plainTextChunk.Value, BlankChars)) == 0 {
				index++
			} else {
				return
			}
		}

	}

	ignoreBlank()

	if index >= len(inputChunks) {
		return nil, index, errIndexOutOfBound
	}

	listItemChunk, ok := inputChunks[index].(*KeywordChunk)
	if !ok || listItemChunk.Keyword != ListItemMark {
		return nil, index, errExpectListItem
	}

	const (
		ListEnd int = iota
		ListItemFound
		OrderListFound
		BulletListFound
	)

	findNextListItem := func(idx int) (int, int) {
		var i int

		for i = idx; i < len(inputChunks); i++ {
			chunk, ok := inputChunks[i].(*KeywordChunk)
			if ok {
				if chunk.Keyword == ListItemMark {
					return i, ListItemFound
				}

				if chunk.Keyword == OrderList {
					return i, OrderListFound
				}

				if chunk.Keyword == BulletList {
					return i, BulletListFound
				}

			}
		}
		return i, ListEnd
	}

	if index+1 >= len(inputChunks) {
		return nil, index, errIndexOutOfBound
	}

	item = &ListItem{}

	newIndex, option := findNextListItem(index + 1)

	switch option {
	case ListEnd, ListItemFound:
		item.Value = append(item.Value, inputChunks[index+1:newIndex]...) //no difference for now
	case OrderListFound, BulletListFound:
		item.Value = append(item.Value, inputChunks[index+1:newIndex]...) //no difference for now
	}
	return item, newIndex, nil
}

func consumeListItems(inputChunks []Chunk, index int) (items []*ListItem, newIndex int, err error) {
	for index < len(inputChunks) {

		item, newIndex, err := consumeListItem(inputChunks, index)
		if err != nil {
			debug.PrintStack()
			log.Fatalln(err)
			return items, newIndex, err
		}
		items = append(items, item)
		index = newIndex
	}
	return
}

func consumeEmbracedBlock(inputChunks []Chunk, index int) (chunks []Chunk, newIndex int, err error) {
	if index >= len(inputChunks) {
		log.Println(errIndexOutOfBound)
		return nil, index, errIndexOutOfBound
	}
	leftBraceChunk, ok := inputChunks[index].(*MetaCharChunk)
	if !ok || leftBraceChunk.GetValue() != LeftBraceChar {
		log.Println(errExpectLBrace)
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
	log.Println(errExpectRBrace)
	return chunks, index, errExpectRBrace

}

func cosumeEmbracedToken(inputChunks []Chunk, index int) (chunks []Chunk, newIndex int, err error) {
	chunks1, newIndex, err := consumeEmbracedBlock(inputChunks, index)
	if len(chunks1) != 3 {
		fmt.Println("===== inputChunks", inputChunks)
		fmt.Println("===== chunks1", chunks1)
		debug.PrintStack()
		log.Fatalln(errExpectToken)
		return chunks, newIndex, errExpectToken
	}

	chunks2, _, err := consumeToken(inputChunks, index+1)
	if err != nil {
		log.Fatalln(err)
		return nil, index, err
	}

	return []Chunk{chunks1[0], chunks2[0], chunks1[2]}, index + 3, nil
}

func consumeToken(inputChunks []Chunk, index int) (chunks []Chunk, newIndex int, err error) {
	if index >= len(inputChunks) {
		log.Fatalln(errIndexOutOfBound)
		return nil, index, errIndexOutOfBound
	}
	plainTextChunk, ok := inputChunks[index].(*PlainTextChunk)
	if !ok {
		log.Fatalln(errExpectToken)
		return nil, index, errExpectToken
	}

	text := strings.TrimLeft(plainTextChunk.GetValue(), BlankChars)
	delta := len(plainTextChunk.GetValue()) - len(text)
	token := gTokenPattern.FindString(text)
	if strings.HasPrefix(text, token) && len(token) > 0 {
		newPlainText := &PlainTextChunk{}
		newPlainText.Position = plainTextChunk.GetPosition()
		newPlainText.Value = token

		plainTextChunk.SetPosition(plainTextChunk.GetPosition() + delta + len(token))
		plainTextChunk.Value = plainTextChunk.Value[delta+len(token):]
		if len(strings.Trim(plainTextChunk.Value, BlankChars)) == 0 {
			newIndex = index + 1
		} else {
			newIndex = index //the index does not change because we splitted the plainTextChunk
		}
		return []Chunk{newPlainText}, newIndex, nil
	}
	debug.PrintStack()
	log.Fatalln(errExpectToken)
	return nil, index, errExpectToken
}
