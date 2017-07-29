package main

import (
	"regexp"
)

var (
	gTokenPattern     = regexp.MustCompile(`[a-zA-Z-_\.][a-zA-Z0-9-_\.]*`)
	gParagraphDivider = regexp.MustCompile(`(\n\s*\n)|(\r\n\s*\r\n)`)
)

const (
	//inline
	EmphasisFormat string = "e"
	StrongFormat   string = "s"
	HyperLink      string = "w"
	InlineCode     string = "c"
	AnchorBlock    string = "a"
	ReferToBlock   string = "k"

	//section
	BlockCode      string = "code"
	SectionHeader  string = "h"
	SectionHeader1 string = "h1"
	SectionHeader2 string = "h2"
	SectionHeader3 string = "h3"
	SectionHeader4 string = "h4"
	SectionHeader5 string = "h5"
	SectionHeader6 string = "h6"

	OrderList    string = "ol"
	BulletList   string = "ul"
	ListItemMark string = "-"

	TableKeyword              = "table"
	TableCellDelimiterKeyword = "d"

	//meta
	TitleKeyword      = "title"
	SubTitleKeyword   = "sub-title"
	AuthorKeyword     = "author"
	CreateDateKeyword = "create-date"
	ModifyDateKeyword = "modify-date"
	KeywordsKeyword   = "keywords"

	//index
	ChapterIndexKeyword = "toc"
	FigureIndexKeyword  = "figure-index"
	TableIndexKeyword   = "table-index"
	ListIndexKeyword    = "list-index"
	CodeIndexKeyword    = "code-index"
)

var (
	gInlineFormatMap  = make(map[string]bool)
	gInlineFormatList = []string{
		EmphasisFormat, StrongFormat, HyperLink, InlineCode, AnchorBlock, ReferToBlock,
	}
	gSectionLevel = map[string]int{SectionHeader: 1, SectionHeader1: 1, SectionHeader2: 2,
		SectionHeader3: 3, SectionHeader4: 4, SectionHeader5: 5, SectionHeader6: 6}
)

func init() {
	for _, f := range gInlineFormatList {
		gInlineFormatMap[f] = true
	}
}
