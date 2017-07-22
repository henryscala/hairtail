package main

import (
	"regexp"
)

var (
	gTokenPattern        = regexp.MustCompile(`[a-zA-Z-_\.][a-zA-Z0-9-_\.]*`)
	gSectionTokenPattern = regexp.MustCompile(`h\d?`) //h followed by 0 or one digit
	gLineDivider         = regexp.MustCompile(`\n|(\r\n)`)
	gParagraphDivider    = regexp.MustCompile(`(\n\n)|(\r\n\r\n)`)
)

const (
	EmphasisFormat string = "e"
	StrongFormat   string = "s"
	HyperLink string = "w"
)

type InlineFormatDescription struct {
	Keyword string
	NumEmbracedBlock int //number of EmbracedBlocks following the keyword
}

var (
	gInlineFormatKeywordMap map[string]*InlineFormatDescription = make(map[string]*InlineFormatDescription)
	gInlineFormatDescriptions = []*InlineFormatDescription {
		&InlineFormatDescription{
			Keyword: EmphasisFormat,
			NumEmbracedBlock: 1,
		},
		&InlineFormatDescription{
			Keyword: StrongFormat,
			NumEmbracedBlock: 1,
		},
		&InlineFormatDescription{
			Keyword: HyperLink,
			NumEmbracedBlock: 2,
		},
	}
)






func init() {
	for _, f := range gInlineFormatDescriptions {
		gInlineFormatKeywordMap[f.Keyword] = f
	}
}
