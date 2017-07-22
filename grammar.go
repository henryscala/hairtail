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
	HyperLink      string = "w"
)

var ()

func init() {

}
