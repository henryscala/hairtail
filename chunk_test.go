package main

import (
	"fmt"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
)

func TestParseChunks(t *testing.T) {

	inputFiles := []string{
		"inline.txt",
		"list.txt",
		"nestlist.txt",
		"section.txt",
		"table.txt",
		"index.txt",
	}
	gConfig.TemplateFile = "template.html"
	for _, file := range inputFiles {
		outputFilePath := filepath.Join("test", fmt.Sprintf("%s%s", strings.TrimSuffix(file, filepath.Ext(file)), ".html"))
		inputFilePath := filepath.Join("test", file)
		err := CompileFile(inputFilePath, outputFilePath)
		if err != nil {
			t.Fatal(err)
		}
	}

}

func TestRawTextChunkHandle(t *testing.T) {
	var text string
	var err error
	var chunkList []Chunk
	var chunk Chunk

	//test one chunk of raw text
	text = `abc\r{raw}def`
	chunkList, err = RawTextChunkHandle(text)
	if err != nil {
		t.FailNow()
	}
	if len(chunkList) != 3 {
		t.FailNow()
	}

	values := []string{"abc", "raw", "def"}
	positions := []int{0, 6, 10}
	types := []string{
		"*main.PlainTextChunk",
		"*main.RawTextChunk",
		"*main.PlainTextChunk",
	}
	for i := 0; i < 3; i++ {
		chunk = chunkList[i]
		if chunk.GetValue() != values[i] {
			t.FailNow()
		}
		if chunk.GetPosition() != positions[i] {
			t.FailNow()
		}
		if reflect.TypeOf(chunk).String() != types[i] {
			t.FailNow()
		}
	}

	//test one chunk of raw text that contains meta chars
	text = `\r##{\{}#}##`
	chunkList, err = RawTextChunkHandle(text)
	if err != nil {
		t.Fatal(err)
	}
	if len(chunkList) != 1 {
		t.FailNow()
	}
	chunk = chunkList[0]
	if chunk.GetValue() != `\{}#` {
		t.Fatal(chunk)
	}

	//test multiple chunks of raw text
	text = `1\r#{aaa}#2\r#{bbb}#3`
	chunkList, err = RawTextChunkHandle(text)
	if err != nil {
		t.FailNow()
	}
	values = []string{"1", "aaa", "2", "bbb", "3"}
	positions = []int{0, 5, 9, 15, 19}
	types = []string{
		"*main.PlainTextChunk",
		"*main.RawTextChunk",
		"*main.PlainTextChunk",
		"*main.RawTextChunk",
		"*main.PlainTextChunk",
	}
	for i := 0; i < 5; i++ {
		chunk = chunkList[i]
		if chunk.GetValue() != values[i] {
			t.FailNow()
		}
		if chunk.GetPosition() != positions[i] {
			t.FailNow()
		}
		if reflect.TypeOf(chunk).String() != types[i] {
			t.FailNow()
		}
	}

	//test multiple line raw text
	text = `
	//followed are computer code 
	\r##{
		int main () {
			return 0; 
		}
	}##
	`
	chunkList, err = RawTextChunkHandle(text)
	if err != nil {
		t.FailNow()
	}
	if len(chunkList) != 3 {
		t.FailNow()
	}
	chunk = chunkList[1]
	if _, ok := chunk.(*RawTextChunk); !ok {
		t.FailNow()
	}

}

func TestMetaCharChunkHandle(t *testing.T) {
	var text string
	var err error
	var chunkList []Chunk
	var chunk Chunk

	// test plain text without meta char
	text = `abc`
	chunkList, err = metaCharChunkHandle(text)
	if err != nil {
		t.FailNow()
	}

	if len(chunkList) != 1 {
		t.FailNow()
	}

	chunk = chunkList[0]

	if plainTextChunk, ok := chunk.(*PlainTextChunk); ok {
		if plainTextChunk.GetPosition() != 0 {
			t.FailNow()
		}

		if plainTextChunk.GetValue() != text {
			t.FailNow()
		}
	} else {
		t.FailNow()
	}

	// test plain text with escape char
	text = `aa \emphasis{param} content`
	chunkList, err = metaCharChunkHandle(text)
	if err != nil {
		t.FailNow()
	}

	if len(chunkList) != 7 {
		t.FailNow()
	}
	// check chunk 0
	chunk = chunkList[0]
	if _, ok := chunk.(*PlainTextChunk); !ok {
		t.FailNow()
	}
	if chunk.GetPosition() != 0 {
		t.FailNow()
	}
	// check values of chunks
	values := []string{"aa ", `\`, "emphasis", "{", "param", "}", " content"}
	positions := []int{0, 3, 4, 12, 13, 18, 19}
	types := []string{
		"*main.PlainTextChunk",
		"*main.MetaCharChunk",
		"*main.PlainTextChunk",
		"*main.MetaCharChunk",
		"*main.PlainTextChunk",
		"*main.MetaCharChunk",
		"*main.PlainTextChunk",
	}
	for i := 0; i < 7; i++ {
		chunk = chunkList[i]
		if chunk.GetValue() != values[i] {
			t.FailNow()
		}
		if chunk.GetPosition() != positions[i] {
			t.FailNow()
		}
		if reflect.TypeOf(chunk).String() != types[i] {
			t.FailNow()
		}
	}

	//test escaped meta chars
	text = `a\\\{\}\#b`
	chunkList, err = metaCharChunkHandle(text)
	if err != nil {
		t.FailNow()
	}
	if len(chunkList) != 1 {
		t.FailNow()
	}
	if chunkList[0].GetValue() != `a\{}#b` {
		t.FailNow()
	}
}
