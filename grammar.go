package main

import (
	"regexp"
)

const (
	BlankChars string = " \r\n\t"
	LineFeed   string = "\n"
)

//meta characters
const (
	EscapeChar     string = "\\"
	LeftBraceChar  string = "{"
	RightBraceChar string = "}"
	FillerChar     string = "#"
)

//supported grammar elements
const (
	RawTextChar string = "r"

	//inline element
	EmphasisFormat = "e"
	StrongFormat   = "s"
	HyperLink      = "w"
	InlineCode     = "c"
	AnchorBlock    = "a"
	ReferToBlock   = "k"
	ImageKeyword   = "image"
	CaptionKeyword = "caption"
	InlineTex      = "t"
	CommentKeyword = "--" //it is able to comment out other grammar elements

	//section
	SectionHeader  = "h"
	SectionHeader1 = "h1"
	SectionHeader2 = "h2"
	SectionHeader3 = "h3"
	SectionHeader4 = "h4"
	SectionHeader5 = "h5"
	SectionHeader6 = "h6"

	//Sections that may have caption and may be shown in specific index
	BlockTex     = "tex"
	BlockCode    = "code"
	OrderList    = "ol"
	BulletList   = "ul"
	ListItemMark = "-"
	TableKeyword = "table"
	//sub element of Table
	TableCellDelimiterKeyword = "d"

	//meta
	TitleKeyword      = "title"
	SubTitleKeyword   = "sub-title"
	AuthorKeyword     = "author"
	CreateDateKeyword = "create-date"
	ModifyDateKeyword = "modify-date"
	KeywordsKeyword   = "keywords"

	//
	IncludeKeyword string = "include" //to include other document

	//index
	SectionIndexKeyword    = "toc"
	ImageIndexKeyword      = "image-index"
	TableIndexKeyword      = "table-index"
	OrderListIndexKeyword  = "order-list-index"
	BulletListIndexKeyword = "bullet-list-index"
	CodeIndexKeyword       = "code-index"
	MathIndexKeyword       = "math-index"
)

var (
	gTokenPattern     = regexp.MustCompile(`[a-zA-Z-_\.][a-zA-Z0-9-_\.]*`)
	gParagraphDivider = regexp.MustCompile(`(\n\s*\n)|(\r\n\s*\r\n)`)
)

var (
	//MetaChars contains all meta characters in slice
	MetaChars = []string{EscapeChar, LeftBraceChar, RightBraceChar, FillerChar}
	//MetaCharMap contains all meta characters in map to be lookup
	MetaCharMap       = make(map[string]bool)
	gInlineFormatMap  = make(map[string]bool)
	gInlineFormatList = []string{
		EmphasisFormat, StrongFormat, HyperLink, InlineCode, CommentKeyword, AnchorBlock, ReferToBlock, InlineTex,
	}

	gChunkWithCaptionList = []string{
		OrderList, BulletList, TableKeyword, BlockCode, ImageKeyword, BlockTex,
	}
	gChunkWithCaptionMap = make(map[string]bool)

	gSectionLevel = map[string]int{SectionHeader: 1, SectionHeader1: 1, SectionHeader2: 2,
		SectionHeader3: 3, SectionHeader4: 4, SectionHeader5: 5, SectionHeader6: 6}

	gMetaInfoKeywordMap  = make(map[string]bool)
	gMetaInfoKeywordList = []string{
		TitleKeyword,
		SubTitleKeyword,
		AuthorKeyword,
		CreateDateKeyword,
		ModifyDateKeyword,
		KeywordsKeyword,
	}
)

func init() {
	for _, f := range gInlineFormatList {
		gInlineFormatMap[f] = true
	}
	for _, f := range gChunkWithCaptionList {
		gChunkWithCaptionMap[f] = true
	}
	for _, f := range gMetaInfoKeywordList {
		gMetaInfoKeywordMap[f] = true
	}
	for _, c := range MetaChars {
		MetaCharMap[c] = true
	}
}
