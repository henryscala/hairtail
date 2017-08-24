package main

import (
	"bytes"
	"errors"
	"fmt"
	"log"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
)

type WithId interface {
	GetId() string
}
type WithCaption interface {
	GetCaption() string
	SetCaption(string)
}
type WithIdCaption interface {
	WithId
	WithCaption
}

type withNumbering interface {
	SetNumbering(c string)
	GetNumbering() string
}
type WithIdCaptionNumbering interface {
	WithIdCaption
	withNumbering
}

// Chunk denotes a block of text, the block may be length 1 to n in bytes
type Chunk interface {
	GetPosition() int
	GetValue() string
	SetPosition(pos int)
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
		log.Fatalln(err)
		return chunks, err
	}
	chunks, err = MetaChunkHandle(chunks)
	if err != nil {
		log.Fatalln(err)
		return chunks, err
	}

	chunks, err = KeywordChunkHandle(chunks)
	if err != nil {
		log.Fatalln(err)
		return chunks, err
	}

	chunks, err = IncludeChunkHandle(chunks)
	if err != nil {
		log.Fatalln(err)
		return chunks, err
	}

	chunks, err = CaptionChunkHandle(chunks)
	if err != nil {
		log.Fatalln(err)
		return chunks, err
	}
	chunks, err = SectionChunkHandle(chunks)
	if err != nil {
		log.Fatalln(err)
		return chunks, err
	}
	chunks, err = ChunkWithNumberingHandle(chunks)
	if err != nil {
		log.Fatalln(err)
		return chunks, err
	}

	gInlineRenderMode = true
	//first render inlineChunk, so that there is not extra <p> around inlineChunk
	chunks, err = InlineChunkListRender(chunks)
	if err != nil {
		log.Fatalln(err)
		return chunks, err
	}

	return chunks, nil
}

//SectionChunkHandle set numbering of SectionChunk
func SectionChunkHandle(inputChunks []Chunk) ([]Chunk, error) {
	sectionChunkList := []*SectionChunk{}
	var levels []int
	var levelMap = make(map[int]bool)
	var levelNumberingMap = make(map[int]int)

	storeLevel := func(level int) {
		if levelMap[level] {
			return
		}
		levels = append(levels, level)
		levelMap[level] = true
	}

	getLevelIndex := func(level int) int {
		for i := 0; i < len(levels); i++ {
			if levels[i] == level {
				return i
			}
		}
		log.Fatalln("should not reach here")
		panic("should not reach here")
	}

	resetLevelLowerThan := func(level int) {
		from := getLevelIndex(level)
		for i := from + 1; i < len(levels); i++ {
			levelNumberingMap[levels[i]] = 0
		}
	}

	calcNumbering := func(level int) string {
		var buf bytes.Buffer
		idx := getLevelIndex(level)
		for i := 0; i <= idx; i++ {
			if i != 0 {
				buf.WriteString(".")
			}
			buf.WriteString(strconv.Itoa(levelNumberingMap[levels[i]]))
		}
		return buf.String()
	}

	for _, c := range inputChunks {
		if keywordChunk, ok := c.(*KeywordChunk); ok {
			level, isSection := gSectionLevel[keywordChunk.Keyword]
			if isSection {
				sectionChunkList = append(sectionChunkList, keywordChunk.Children[0].(*SectionChunk))
				storeLevel(level)
			}
		}
		sort.Ints(levels)
	}

	var bufSectionIndex bytes.Buffer

	for _, sectionChunk := range sectionChunkList {
		level := sectionChunk.Level
		resetLevelLowerThan(level)
		levelNumberingMap[level]++
		//generate numbering for this section
		sectionChunk.Numbering = calcNumbering(level)
		//generate index for the whole doc
		err := gSectionIndexTemplate.Execute(&bufSectionIndex, sectionChunk)
		if err != nil {
			return nil, err
		}
	}
	gDoc.SectionIndex = bufSectionIndex.String()

	return inputChunks, nil
}

//Chunk with numbering handle
func ChunkWithNumberingHandle(inputChunks []Chunk) ([]Chunk, error) {
	numberingMap := make(map[string]int)       //to generate numbering
	indexMap := make(map[string]*bytes.Buffer) //to generate index

	for _, chunk := range inputChunks {
		keywordChunk, ok := chunk.(*KeywordChunk)
		if !ok {
			continue
		}
		if gChunkWithCaptionMap[keywordChunk.Keyword] {
			chunkWithIdCaptionNumbering := keywordChunk.Children[0].(WithIdCaptionNumbering)
			//only chunks that the caption has been set, we set numbering for them
			if len(chunkWithIdCaptionNumbering.GetCaption()) > 0 {
				numberingMap[keywordChunk.Keyword]++
				prefix := getKeywordName(keywordChunk.Keyword)
				//set numbering
				chunkWithIdCaptionNumbering.SetNumbering(prefix + " " + strconv.Itoa(numberingMap[keywordChunk.Keyword]) + ": ")

				//generate index for this type
				buf := indexMap[keywordChunk.Keyword]
				if nil == buf {
					buf = new(bytes.Buffer)
					indexMap[keywordChunk.Keyword] = buf
				}
				err := gGlobalIndexTemplate.Execute(buf, chunkWithIdCaptionNumbering)
				if err != nil {
					return nil, err
				}
			}
		}
	}

	//set indices to global doc obj
	for _, keyword := range gChunkWithCaptionList {
		buf := indexMap[keyword]
		if buf != nil {
			switch keyword {
			case OrderList:
				gDoc.OrderListIndex = buf.String()
			case BulletList:
				gDoc.BulletListIndex = buf.String()
			case TableKeyword:
				gDoc.TableIndex = buf.String()
			case BlockCode:
				gDoc.CodeIndex = buf.String()
			case BlockTex:
				gDoc.MathIndex = buf.String()
			case ImageKeyword:
				gDoc.ImageIndex = buf.String()
			}
		}
	}

	return inputChunks, nil
}

//CaptionChunkHandle filter the Caption chunk, and set caption to the chunk it refers to
func CaptionChunkHandle(inputChunks []Chunk) ([]Chunk, error) {
	outputChunks := []Chunk{}
	idToChunk := make(map[string]Chunk)
	captionChunks := []Chunk{}
	for _, chunk := range inputChunks {
		keywordChunk, ok := chunk.(*KeywordChunk)
		if !ok {
			outputChunks = append(outputChunks, chunk)
			continue
		}
		if gChunkWithCaptionMap[keywordChunk.Keyword] {
			chunkWithIdCaption := keywordChunk.Children[0].(WithIdCaption)
			idToChunk[chunkWithIdCaption.GetId()] = chunk
			outputChunks = append(outputChunks, chunk)
			continue
		}
		if keywordChunk.Keyword == CaptionKeyword {
			captionChunks = append(captionChunks, chunk)
			continue
		}
		outputChunks = append(outputChunks, chunk)
	}

	for _, captionChunk := range captionChunks {
		id := captionChunk.(*KeywordChunk).Children[0]
		caption := captionChunk.(*KeywordChunk).Children[1]
		chunk, ok := idToChunk[id.GetValue()]
		if !ok {
			log.Println("caption ", caption, "has not found Id", id.GetValue())
			return nil, errExpectChunkWithId
		}
		chunk.(*KeywordChunk).Children[0].(WithIdCaption).SetCaption(caption.GetValue())
	}
	return outputChunks, nil
}

//IncludeChunkHandle filter the Include chunk, and import contents of the file it refers to
func IncludeChunkHandle(inputChunks []Chunk) ([]Chunk, error) {
	outputChunks := []Chunk{}

	for _, chunk := range inputChunks {
		keywordChunk, ok := chunk.(*KeywordChunk)
		if !ok {
			outputChunks = append(outputChunks, chunk)
			continue
		}

		if keywordChunk.Keyword != IncludeKeyword {
			outputChunks = append(outputChunks, chunk)
			continue
		}

		//now it is include keyword, we first figure out the path of the included file.
		//it is relative to the current file to be compiled or absolute path.
		//Note: If the included file itself contains include keyword,
		//it is still relative to the current file to be compiled(not the included file)
		//Note: the implementation of include keyword has limitations.
		//It is better the included content does not rely on chunks in other files. Otherwise, surprise may happens.
		parentDir := filepath.Dir(gDoc.FilePath)
		includedFileName := strings.Trim(keywordChunk.GetValue(), BlankChars)
		absolutePath := false
		if strings.HasPrefix(includedFileName, "/") || strings.HasPrefix(includedFileName, "\\") {
			absolutePath = true
		}
		var parts []string
		var includedFilePath string
		if strings.Contains(includedFileName, "/") {
			parts = strings.Split(includedFileName, "/")
			if absolutePath {
				includedFilePath = filepath.Join("/", filepath.Join(parts...))
			} else {
				includedFilePath = filepath.Join(parentDir, filepath.Join(parts...))
			}
		} else {
			parts = strings.Split(includedFileName, "\\")
			if absolutePath {
				includedFilePath = filepath.Join("\\", filepath.Join(parts...))
			} else {
				includedFilePath = filepath.Join(parentDir, filepath.Join(parts...))
			}
		}
		//included file to chunks, there are re-cursive calls inside
		includedChunks, err := fileToChunks(includedFilePath)
		if err != nil {
			log.Println(err)
			return nil, err
		}
		outputChunks = append(outputChunks, includedChunks...)
	}

	return outputChunks, nil
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
