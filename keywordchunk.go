package main

import (
	"errors"
	"fmt"
	"log"
	"runtime/debug"
	"strings"
)

var (
	errExpectToken       = errors.New("expect token")
	errIndexOutOfBound   = errors.New("index out of bound")
	errExpectLBrace      = errors.New("expect Left Brace")
	errExpectRBrace      = errors.New("expect Right Brace")
	errExpectPlainText   = errors.New("expect Plain Text")
	errExpectRawText     = errors.New("expect Raw Text")
	errExpectListItem    = errors.New("expect List Item ")
	errExpectChunkWithId = errors.New("expect chunk with specific Id ")
)

// KeywordChunk denotes meta char '\' followed by a token.
// during handling, the meta char '\' is omitted.
type KeywordChunk struct {
	Position int
	Keyword  string
	Value    string
	Children []Chunk
}

// String implements the Stringer interface
func (p KeywordChunk) String() string {
	return fmt.Sprintf("keywordChunk{Position:%v,Keyword:%v,Value:%v,Children:%v}", p.Position, p.Keyword, p.Value, p.Children)
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
	return p.Value
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
				return nil, err
			}

			switch token[0].GetValue() {
			case EmphasisFormat, StrongFormat:
				outputChunks, index, err = inlineBlockOneParamHandle(token[0], inputChunks, outputChunks, newIndex)

			case HyperLink:
				outputChunks, index, err = hyperLinkBlockHandle(token[0], inputChunks, outputChunks, newIndex)
			case ImageKeyword:
				outputChunks, index, err = imageBlockHandle(token[0], inputChunks, outputChunks, newIndex)
			case InlineCode:
				outputChunks, index, err = inlineCodeBlockHandle(token[0], inputChunks, outputChunks, newIndex)
			case TableCellDelimiterKeyword, ListItemMark, SectionIndexKeyword, ImageIndexKeyword, TableIndexKeyword, OrderListIndexKeyword, BulletListIndexKeyword, CodeIndexKeyword:
				outputChunks, index, err = simpleKeywordHandle(token[0], inputChunks, outputChunks, newIndex)
			case AnchorBlock:
				outputChunks, index, err = anchorBlockHandle(token[0], inputChunks, outputChunks, newIndex)
			case ReferToBlock:
				outputChunks, index, err = referToBlockHandle(token[0], inputChunks, outputChunks, newIndex)
			case TitleKeyword, SubTitleKeyword, AuthorKeyword, CreateDateKeyword, ModifyDateKeyword, KeywordsKeyword:
				outputChunks, index, err = metaKeywordHandle(token[0], inputChunks, outputChunks, newIndex)
			case BlockCode:
				outputChunks, index, err = blockCodeBlockHandle(token[0], inputChunks, outputChunks, newIndex)
			case SectionHeader, SectionHeader1, SectionHeader2, SectionHeader3,
				SectionHeader4, SectionHeader5, SectionHeader6:
				outputChunks, index, err = sectionBlockHandle(token[0], inputChunks, outputChunks, newIndex)
			case OrderList, BulletList:
				outputChunks, index, err = listBlockHandle(token[0], inputChunks, outputChunks, newIndex)
			case TableKeyword:
				outputChunks, index, err = tableBlockHandle(token[0], inputChunks, outputChunks, newIndex)
			case CaptionKeyword:
				outputChunks, index, err = captionBlockHandle(token[0], inputChunks, outputChunks, newIndex)
			default:
				log.Fatal("not implemented")
				panic("not implemented")
			}
			if err != nil {
				log.Fatalln(err)
				return nil, err
			}
			continue
		} else {
			outputChunks = append(outputChunks, inputChunk)
		}
		index++
	}

	return outputChunks, nil
}

func referToBlockHandle(token Chunk, inputChunks, outputChunks []Chunk, index int) (newOutputChunks []Chunk, newIndex int, err error) {
	tokenChunks, newIndex, err := consumeEmbracedToken(inputChunks, index)

	if err != nil {
		log.Fatalln(err)
		return outputChunks, index, err
	}

	chunk := &ReferToChunk{
		Position: token.GetPosition(),
		Id:       tokenChunks[1].GetValue(),
		Value:    "",
	}
	keywordChunk := &KeywordChunk{Position: token.GetPosition(),
		Keyword:  token.GetValue(),
		Children: []Chunk{chunk},
	}

	outputChunks = append(outputChunks, keywordChunk)

	return outputChunks, newIndex, nil
}

func anchorBlockHandle(token Chunk, inputChunks, outputChunks []Chunk, index int) (newOutputChunks []Chunk, newIndex int, err error) {
	tokenChunks, newIndex, err := consumeEmbracedToken(inputChunks, index)

	if err != nil {
		log.Fatalln(err)
		return outputChunks, index, err
	}
	chunksContent, newIndex, err := consumeEmbracedBlock(inputChunks, newIndex)
	if err != nil {
		log.Fatal(err)
		return outputChunks, index, err
	}

	chunksContent, err = KeywordChunkHandle(chunksContent) //recursive
	if err != nil {
		log.Fatalln(err)
		return outputChunks, index, err
	}

	if len(chunksContent) < 2 {
		log.Fatalln(errExpectRBrace)
		return outputChunks, index, errExpectRBrace
	}

	var value string
	if len(chunksContent) == 2 {
		value = ""
	} else {
		value = chunksContent[1].GetValue()
	}

	chunk := &AnchorChunk{
		Position: token.GetPosition(),
		Id:       tokenChunks[1].GetValue(),
		Value:    value,
	}
	keywordChunk := &KeywordChunk{Position: token.GetPosition(),
		Keyword:  token.GetValue(),
		Children: []Chunk{chunk},
	}

	outputChunks = append(outputChunks, keywordChunk)

	return outputChunks, newIndex, nil
}

//inlineBlockOneParamHandle handles inline format keyword followed by one pair of brace
func inlineBlockOneParamHandle(token Chunk, inputChunks, outputChunks []Chunk, index int) (newOutputChunks []Chunk, newIndex int, err error) {
	chunks, newIndex, err := consumeEmbracedBlock(inputChunks, index)
	if err != nil {
		log.Println(err)
		return outputChunks, index, err
	}
	chunks, err = KeywordChunkHandle(chunks) //recursive
	if err != nil {
		log.Println(err)
		return outputChunks, index, err
	}
	keywordChunk := &KeywordChunk{Position: token.GetPosition(),
		Keyword:  token.GetValue(),
		Children: []Chunk{chunks[1]},
	}
	outputChunks = append(outputChunks, keywordChunk)

	return outputChunks, newIndex, nil
}

func hyperLinkBlockHandle(token Chunk, inputChunks, outputChunks []Chunk, index int) (newOutputChunks []Chunk, newIndex int, err error) {
	chunksUrl, newIndex, err := consumeEmbracedBlock(inputChunks, index)
	if err != nil {
		log.Fatalln(err)
		return outputChunks, index, err
	}
	chunksUrl, err = KeywordChunkHandle(chunksUrl) //recursive
	if err != nil {
		log.Fatalln(err)
		return outputChunks, index, err
	}

	chunksContent, newIndex, err := consumeEmbracedBlock(inputChunks, newIndex)
	if err != nil {
		log.Fatalln(err)
		return outputChunks, index, err
	}
	chunksContent, err = KeywordChunkHandle(chunksContent) //recursive
	if err != nil {
		log.Fatalln(err)
		return outputChunks, index, err
	}

	keywordChunk := &KeywordChunk{Position: token.GetPosition(),
		Keyword:  token.GetValue(),
		Children: []Chunk{chunksUrl[1], chunksContent[1]},
	}
	outputChunks = append(outputChunks, keywordChunk)
	return outputChunks, newIndex, nil
}

func captionBlockHandle(token Chunk, inputChunks, outputChunks []Chunk, index int) (newOutputChunks []Chunk, newIndex int, err error) {
	tokenChunks, newIndex, err := consumeEmbracedToken(inputChunks, index)
	if err != nil {
		log.Fatalln(err)
		return outputChunks, index, err
	}

	chunksContent, newIndex, err := consumeEmbracedBlock(inputChunks, newIndex)
	if err != nil {
		log.Fatalln(err)
		return outputChunks, index, err
	}
	chunksContent, err = KeywordChunkHandle(chunksContent) //recursive
	if err != nil {
		log.Fatalln(err)
		return outputChunks, index, err
	}

	keywordChunk := &KeywordChunk{Position: token.GetPosition(),
		Keyword:  token.GetValue(),
		Children: []Chunk{tokenChunks[1], chunksContent[1]}, //id + caption
	}
	outputChunks = append(outputChunks, keywordChunk)
	return outputChunks, newIndex, nil
}

func imageBlockHandle(token Chunk, inputChunks, outputChunks []Chunk, index int) (newOutputChunks []Chunk, newIndex int, err error) {
	tokenChunks, newIndex, err := consumeEmbracedToken(inputChunks, index)

	if err != nil {
		log.Fatalln(err)
		return outputChunks, index, err
	}
	chunksContent, newIndex, err := consumeEmbracedBlock(inputChunks, newIndex)

	if err != nil {
		log.Fatalln(err)
		return outputChunks, index, err
	}

	imageChunk := &ImageChunk{Id: tokenChunks[1].GetValue(),
		Src:      chunksContent[1].GetValue(),
		Position: token.GetPosition(),
	}

	keywordChunk := &KeywordChunk{Position: token.GetPosition(),
		Keyword:  token.GetValue(),
		Children: []Chunk{imageChunk},
	}
	outputChunks = append(outputChunks, keywordChunk)
	return outputChunks, newIndex, nil
}

func inlineCodeBlockHandle(token Chunk, inputChunks, outputChunks []Chunk, index int) (newOutputChunks []Chunk, newIndex int, err error) {
	chunks, newIndex1, err := consumeEmbracedBlock(inputChunks, index)
	if err == nil {

		keywordChunk := &KeywordChunk{Position: token.GetPosition(),
			Keyword:  token.GetValue(),
			Children: []Chunk{chunks[1]},
		}
		outputChunks = append(outputChunks, keywordChunk)
		newIndex = newIndex1

	} else {
		//InlineCode content may be either EmbracedBlock or RawTextBlock
		if index >= len(inputChunks) {
			log.Fatalln(errIndexOutOfBound)
			return outputChunks, index, errIndexOutOfBound
		}
		rawTextChunk, ok := inputChunks[index].(*RawTextChunk)
		if !ok {
			log.Fatalln(errExpectRawText)
			return outputChunks, index, errExpectRawText
		}

		keywordChunk := &KeywordChunk{Position: token.GetPosition(),
			Keyword:  token.GetValue(),
			Children: []Chunk{rawTextChunk},
		}

		outputChunks = append(outputChunks, keywordChunk)
		newIndex = index + 1

	}
	return outputChunks, newIndex, nil
}

func blockCodeBlockHandle(token Chunk, inputChunks, outputChunks []Chunk, index int) (newOutputChunks []Chunk, newIndex int, err error) {
	tokenChunks, newIndex, err := consumeEmbracedToken(inputChunks, index)

	if err != nil {
		log.Fatalln(err)
		return outputChunks, index, err
	}
	chunksContent, newIndex1, err := consumeEmbracedBlock(inputChunks, newIndex)
	if err == nil {
		blockCodeChunk := &BlockCodeChunk{
			Position: token.GetPosition(),
			Id:       tokenChunks[1].GetValue(),
			Value:    chunksContent[1].GetValue(),
		}
		keywordChunk := &KeywordChunk{Position: token.GetPosition(),
			Keyword:  token.GetValue(),
			Children: []Chunk{blockCodeChunk},
		}

		outputChunks = append(outputChunks, keywordChunk)
		newIndex = newIndex1
	} else {
		//BlockCode content may be either EmbracedBlock or RawTextBlock
		if newIndex >= len(inputChunks) {
			log.Fatalln(errIndexOutOfBound)
			return outputChunks, index, errIndexOutOfBound
		}
		rawTextChunk, ok := inputChunks[newIndex].(*RawTextChunk)
		if !ok {
			log.Fatalln(errExpectRawText)
			return outputChunks, index, errExpectRawText
		}

		blockCodeChunk := &BlockCodeChunk{
			Position: token.GetPosition(),
			Id:       tokenChunks[1].GetValue(),
			Value:    rawTextChunk.GetValue(),
		}
		keywordChunk := &KeywordChunk{Position: token.GetPosition(),
			Keyword:  token.GetValue(),
			Children: []Chunk{blockCodeChunk},
		}

		outputChunks = append(outputChunks, keywordChunk)
		newIndex++

	}
	return outputChunks, newIndex, nil
}

func metaKeywordHandle(token Chunk, inputChunks, outputChunks []Chunk, index int) (newOutputChunks []Chunk, newIndex int, err error) {

	plainTextChunk, ok := inputChunks[index].(*PlainTextChunk)

	if !ok {
		log.Fatalln(errExpectPlainText)
		return outputChunks, index, errExpectPlainText
	}
	firstLineChunk, restLineChunk, err := plainTextChunk.FirstLineRestLines()
	if err != nil {
		log.Fatalln(errExpectPlainText)
		return outputChunks, index, errExpectPlainText
	}
	//update the plainTextChunk in-place
	if restLineChunk == nil {
		plainTextChunk.Value = ""
	} else {
		plainTextChunk.Value = restLineChunk.GetValue()
		plainTextChunk.Position = restLineChunk.GetPosition()
	}

	keywordChunk := &KeywordChunk{Position: token.GetPosition(),
		Keyword: token.GetValue(),
		Value:   firstLineChunk.GetValue(),
	}
	outputChunks = append(outputChunks, keywordChunk)
	outputChunks = append(outputChunks, plainTextChunk)
	newIndex = index + 1
	return outputChunks, newIndex, nil
}

func sectionBlockHandle(token Chunk, inputChunks, outputChunks []Chunk, index int) (newOutputChunks []Chunk, newIndex int, err error) {
	header := token.GetValue()
	tokenChunks, newIndex, err := consumeEmbracedToken(inputChunks, index)

	if err != nil {
		log.Fatalln(err)
		return outputChunks, index, err
	}
	if newIndex > len(inputChunks) {
		log.Fatalln(errIndexOutOfBound)
		return outputChunks, index, errIndexOutOfBound
	}

	plainTextChunk, ok := inputChunks[newIndex].(*PlainTextChunk)
	newIndex++
	if !ok {
		log.Fatalln(errExpectPlainText)
		return outputChunks, index, errExpectPlainText
	}
	firstLineChunk, restLineChunk, err := plainTextChunk.FirstLineRestLines()
	if err != nil {
		log.Fatalln(errExpectPlainText)
		return outputChunks, index, errExpectPlainText
	}
	//update the plainTextChunk in-place
	if restLineChunk == nil {
		plainTextChunk.Value = ""
	} else {
		plainTextChunk.Value = restLineChunk.GetValue()
		plainTextChunk.Position = restLineChunk.GetPosition()
	}
	level := gSectionLevel[header]

	sectionChunk := &SectionChunk{Position: token.GetPosition(),
		Level: level, Caption: firstLineChunk.GetValue(),
		Id: tokenChunks[1].GetValue(),
	}
	keywordChunk := &KeywordChunk{Position: token.GetPosition(),
		Keyword:  token.GetValue(),
		Children: []Chunk{sectionChunk},
	}
	outputChunks = append(outputChunks, keywordChunk)
	outputChunks = append(outputChunks, plainTextChunk)
	return outputChunks, newIndex, nil
}

func listBlockHandle(token Chunk, inputChunks, outputChunks []Chunk, index int) (newOutputChunks []Chunk, newIndex int, err error) {
	tokenChunks, newIndex, err := consumeEmbracedToken(inputChunks, index)
	if err != nil {
		log.Fatalln(err)
		return outputChunks, index, err
	}
	chunksContent, newIndex, err := consumeEmbracedBlock(inputChunks, newIndex)
	if err != nil {
		log.Fatalln(err)
		return outputChunks, index, err
	}

	chunksContent, err = KeywordChunkHandle(chunksContent)
	if err != nil {
		log.Fatalln(err)
		return outputChunks, index, err
	}

	listChunk := &ListChunk{Position: token.GetPosition(),
		Id:       tokenChunks[1].GetValue(),
		ListType: token.GetValue(),
	}

	items, _, err := consumeListItems(chunksContent[1:len(chunksContent)-1], 0)
	if err != nil {
		log.Fatalln(err)
		return outputChunks, index, err
	}
	listChunk.Items = items
	keywordChunk := &KeywordChunk{Position: token.GetPosition(),
		Keyword:  token.GetValue(),
		Children: []Chunk{listChunk},
	}
	outputChunks = append(outputChunks, keywordChunk)
	return outputChunks, newIndex, nil
}

func tableBlockHandle(token Chunk, inputChunks, outputChunks []Chunk, index int) (newOutputChunks []Chunk, newIndex int, err error) {
	tokenChunks, newIndex, err := consumeEmbracedToken(inputChunks, index)
	if err != nil {
		log.Fatalln(err)
		return outputChunks, index, err
	}
	chunksContent, newIndex, err := consumeEmbracedBlock(inputChunks, newIndex)
	if err != nil {
		log.Fatalln(err)
		return outputChunks, index, err
	}

	chunksContent, err = KeywordChunkHandle(chunksContent[1 : len(chunksContent)-1])
	if err != nil {
		log.Fatalln(err)
		return outputChunks, index, err
	}
	tableChunk := &TableChunk{
		Position: token.GetPosition(),
		Id:       tokenChunks[1].GetValue(),
	}

	row := []Chunk{}
	for i := 0; i < len(chunksContent); i++ {
		chunk := chunksContent[i]
		if keywordChunk, ok := chunk.(*KeywordChunk); ok && keywordChunk.Keyword == TableCellDelimiterKeyword {
			continue //ignore it
		}
		if plainTextChunk, ok := chunk.(*PlainTextChunk); ok {
			firstLine, restLine, err := plainTextChunk.FirstLineRestLines()
			if err != nil {
				return newOutputChunks, index, err
			}
			//empty line is ignored
			if len(strings.Trim(firstLine.GetValue(), BlankChars)) > 0 {
				row = append(row, firstLine)
			}

			if restLine == nil {
				continue

			}
			if len(row) > 0 {
				tableChunk.Cells = append(tableChunk.Cells, row)
			}

			row = []Chunk{}
			if len(strings.Trim(restLine.GetValue(), BlankChars)) > 0 {
				row = append(row, restLine)
			}
		}
	}
	if len(row) > 0 {
		tableChunk.Cells = append(tableChunk.Cells, row)
	}

	keywordChunk := &KeywordChunk{Position: token.GetPosition(),
		Keyword:  token.GetValue(),
		Children: []Chunk{tableChunk},
	}
	outputChunks = append(outputChunks, keywordChunk)

	return outputChunks, newIndex, nil
}

//only keyword itself, no following blocks
func simpleKeywordHandle(token Chunk, inputChunks, outputChunks []Chunk, index int) (newOutputChunks []Chunk, newIndex int, err error) {
	keywordChunk := &KeywordChunk{Position: token.GetPosition(),
		Keyword: token.GetValue(),
	}
	outputChunks = append(outputChunks, keywordChunk)
	newOutputChunks = outputChunks
	newIndex = index
	err = nil
	return
}

func ignoreBlank(inputChunks []Chunk, index int) (newIndex int) {

	for index < len(inputChunks) {
		plainTextChunk, ok := inputChunks[index].(*PlainTextChunk)
		if ok && len(strings.Trim(plainTextChunk.Value, BlankChars)) == 0 {
			index++
		} else {
			break
		}
	}
	newIndex = index
	return

}

func consumeListItem(inputChunks []Chunk, index int) (item *ListItem, newIndex int, err error) {

	index = ignoreBlank(inputChunks, index)

	if index >= len(inputChunks) {
		return nil, index, errIndexOutOfBound
	}

	listItemChunk, ok := inputChunks[index].(*KeywordChunk)
	if !ok || listItemChunk.Keyword != ListItemMark {
		log.Println("== listItemChunk : ", listItemChunk)
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
		item.Value = append(item.Value, inputChunks[index+1:newIndex]...)
	case OrderListFound, BulletListFound:
		newIndex++ //also contain the nested list
		item.Value = append(item.Value, inputChunks[index+1:newIndex]...)
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
	index = ignoreBlank(inputChunks, index)
	if index >= len(inputChunks) {
		log.Println(errIndexOutOfBound)
		return nil, index, errIndexOutOfBound
	}
	leftBraceChunk, ok := inputChunks[index].(*MetaCharChunk)
	if !ok || leftBraceChunk.GetValue() != LeftBraceChar {

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

func consumeEmbracedToken(inputChunks []Chunk, index int) (chunks []Chunk, newIndex int, err error) {
	chunks1, newIndex, err := consumeEmbracedBlock(inputChunks, index)
	if len(chunks1) < 3 { //there might be empty plain-text here, so the len may be > 3

		debug.PrintStack()
		log.Fatalln(errExpectToken)
		return chunks, newIndex, errExpectToken
	}

	chunks2, newIndex2, err := consumeToken(inputChunks, index+1)
	if err != nil {
		log.Fatalln(err)
		return nil, index, err
	}
	if newIndex2 >= newIndex {
		return nil, index, errExpectPlainText
	}

	return []Chunk{chunks1[0], chunks2[0], chunks1[len(chunks1)-1]}, newIndex, nil
}

func consumeToken(inputChunks []Chunk, index int) (chunks []Chunk, newIndex int, err error) {

	index = ignoreBlank(inputChunks, index)
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
