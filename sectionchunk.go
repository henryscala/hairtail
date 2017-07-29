package main

import (
	"fmt"
)

// SectionChunk denotes Sections in article. It is nested structure.
type SectionChunk struct {
	Position int
	Level    int //1 2 .. 6
	Id       string
	Caption  string
	Children []Chunk
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

	return p.Caption

}

/*
// SectionChunkHandle
func SectionChunkHandle(inputChunks []Chunk) ([]Chunk, error) {
	var (
		outputChunks []Chunk

		sectionExist bool
	)

	//first pass, the SectionChunk has no children

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
*/
