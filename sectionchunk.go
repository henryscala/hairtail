package main

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

// SectionChunk denotes Sections in article. It is nested structure.
type SectionChunk struct {
	Position int
	Level    int //1 2 .. 6
	Id       string
	Caption  string
	Children []Chunk
}

func (p SectionChunk) IsTerminal() bool {
	return len(p.Children) == 0
}

// String implements the Stringer interface
func (p SectionChunk) String() string {
	return fmt.Sprintf("SectionChunk{Position: %d, Level: %d, Id: %v, Caption: %v, Children: %v}",
		p.GetPosition(), p.Level, p.Id, p.Caption, p.Children)
}

// GetPosition implements the Chunk interface
func (p *SectionChunk) GetPosition() int {
	return p.Position
}

// SetPosition implements the Chunk interface
func (p *SectionChunk) SetPosition(pos int) {
	p.Position = pos
}

// GetValue implements the Chunk interface
func (p *SectionChunk) GetValue() string {
	if p.IsTerminal() {
		return p.Caption
	}
	panic("should not directly use")
	return "not implemented"
}

// SectionChunkHandle
func SectionChunkHandle(inputChunks []Chunk) ([]Chunk, error) {
	var (
		outputChunks []Chunk
		err          error
		sectionExist bool
	)

	//first pass, the SectionChunk has no children
	for i := 0; i < len(inputChunks); i++ {

		inputChunk := inputChunks[i]

		keywordChunk, isKeywordChunk := inputChunk.(*KeywordChunk)
		if !isKeywordChunk {
			outputChunks = append(outputChunks, inputChunk)

			continue
		}
		sectionName := gSectionTokenPattern.FindString(keywordChunk.Keyword)
		if len(sectionName) == 0 || sectionName != strings.Trim(keywordChunk.Keyword, BlankChars) {
			outputChunks = append(outputChunks, inputChunk)

			continue
		}
		//from here, it is sure that the chunk is section chunk

		level := 1
		if len(sectionName) > 1 {
			level, err = strconv.Atoi(sectionName[1:2]) //only one digit is supported
			if err != nil {
				return outputChunks, err
			}
		}

		if i+2 >= len(inputChunks) {
			return outputChunks, errors.New("section chunk does not follow a ID chunk and a caption chunk")
		}
		nextChunk := inputChunks[i+1]
		nextNextChunk := inputChunks[i+2]
		embracedChunk, isEmbracedChunk := nextChunk.(*EmbracedChunk)
		if !isEmbracedChunk {
			return outputChunks, errors.New("section chunk shall contain an Embraced chunnk that stands for ID of the section")
		}

		if !embracedChunk.IsTerminal() {
			return outputChunks, errors.New("ID chunk shall be terminal not nested")

		}
		id := embracedChunk.GetValue()
		plainTextChunk, isPlainTextChunk := nextNextChunk.(*PlainTextChunk)
		if !isPlainTextChunk {
			return outputChunks, errors.New("captain chunk must be contained in PlainTextChunk")
		}
		firstLine, restLines, err := plainTextChunk.FirstLineRestLines()
		if err != nil {
			return outputChunks, err
		}

		sectionChunk := &SectionChunk{Position: keywordChunk.GetPosition(),
			Level:   level,
			Id:      id,
			Caption: firstLine.GetValue(),
		}
		outputChunks = append(outputChunks, sectionChunk)
		if restLines != nil {
			outputChunks = append(outputChunks, restLines)
		}
		sectionExist = true
		i += 2 // jump 2 chunks

	}

	if !sectionExist {
		return outputChunks, nil
	}

	var outputChunks2 []Chunk
	var lastLevel int = 0
	//second pass, put children in SectionChunk
	for i := 0; i < len(outputChunks); i++ {
		chunk := outputChunks[i]
		sectionChunk, isSectionChunk := chunk.(*SectionChunk)
		if !isSectionChunk {
			outputChunks2 = append(outputChunks2, chunk)
			continue
		}
		index, found := findPrevSectionChunk(outputChunks2, sectionChunk.Level)
		if !found {
			outputChunks2 = append(outputChunks2, sectionChunk)
			continue
		}
		prevSectionChunk := outputChunks2[index].(*SectionChunk)
		//put sub Chunks to the prev section chunk
		prevSectionChunk.Children = append([]Chunk{}, outputChunks2[index+1:]...)
		outputChunks2 = outputChunks2[:index+1]
		outputChunks2 = append(outputChunks2, sectionChunk)
		lastLevel = sectionChunk.Level
	}
	//special handling the last section
	if lastLevel > 0 {
		index, found := findPrevSectionChunk(outputChunks2, lastLevel)
		if !found {
			panic("should not happen")
		}
		prevSectionChunk := outputChunks2[index].(*SectionChunk)
		//put sub Chunks to the prev section chunk
		prevSectionChunk.Children = outputChunks2[index+1:]
		outputChunks2 = outputChunks2[:index+1]
	}

	//third pass, combine low level sections to higher level sections (h2 is lower than h1)

	outputChunks3 := putLowerLevelSectionsToHigherLevelSections(outputChunks2)

	return outputChunks3, nil
}

func findPrevSectionChunk(chunkList []Chunk, level int) (index int, found bool) {
	for index = len(chunkList) - 1; index >= 0; index-- {
		if sectionChunk, isSectionChunk := chunkList[index].(*SectionChunk); isSectionChunk && sectionChunk.Level == level {
			found = true
			return
		}
	}
	return
}

//return the index that the level is equal or less than lowerLimit. Or the end index.
func findNextSectionChunk(chunkList []Chunk, fromIndex, lowerLimit int) (index int) {
	for index = fromIndex; index < len(chunkList); index++ {
		inputChunk := chunkList[index]
		sectionChunk, isSectionChunk := inputChunk.(*SectionChunk)
		if isSectionChunk && sectionChunk.Level <= lowerLimit {
			return
		}
	}
	return
}

func putLowerLevelSectionsToHigherLevelSections(inputChunks []Chunk) []Chunk {
	var chunkList []Chunk

	for i := 0; i < len(inputChunks); {
		inputChunk := inputChunks[i]
		sectionChunk, isSectionChunk := inputChunk.(*SectionChunk)
		if !isSectionChunk {

			chunkList = append(chunkList, inputChunk)
			i++
			continue
		}

		index := findNextSectionChunk(inputChunks, i+1, sectionChunk.Level)
		sectionChunk.Children = append(sectionChunk.Children, inputChunks[i+1:index]...)
		sectionChunk.Children = putLowerLevelSectionsToHigherLevelSections(sectionChunk.Children) //recursive call
		chunkList = append(chunkList, sectionChunk)
		i = index
	}
	return chunkList
}
