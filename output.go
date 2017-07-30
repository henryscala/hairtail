package main

import (
	"bytes"
	"log"
	"text/template"
)

var (
	gSectionTemplate    *template.Template
	gParagraphTemplate  *template.Template
	gEmphasisTemplate   *template.Template
	gStrongTemplate     *template.Template
	gHyperLinkTemplate  *template.Template
	gInlineCodeTemplate *template.Template
	gBlockCodeTemplate  *template.Template
	gListTemplate       *template.Template
	gListItemTemplate   *template.Template
	gAnchorTemplate     *template.Template
	gReferToTemplate    *template.Template
	gTableTemplate      *template.Template
	gTableRowTemplate   *template.Template
	gTableCellTemplate  *template.Template

	gInlineRenderMode bool
)

func init() {
	gSectionTemplate, _ = template.New("Section").Parse(`<h{{.Level}} id="{{.Id}}">{{.Caption}}</h{{.Level}}>` + "\n")
	gParagraphTemplate, _ = template.New("Paragraph").Parse(`<p>{{.}}</p>` + "\n")
	gEmphasisTemplate, _ = template.New("Emphasis").Parse(`<em>{{.}}</em>`)
	gStrongTemplate, _ = template.New("Strong").Parse(`<strong>{{.}}</strong>`)
	gHyperLinkTemplate, _ = template.New("HyperLink").Parse(`<a href="{{.Url}}">{{.Text}}</a>`)
	gInlineCodeTemplate, _ = template.New("InlineCode").Parse(`<code>{{.}}</code>`)
	gBlockCodeTemplate, _ = template.New("BlockCode").Parse(`<p><a id="{{.Id}}" class="caption">{{.Caption}}</a></p><pre>{{.Value}}</pre>` + "\n")
	gListTemplate, _ = template.New("List").Parse(`<p><a id="{{.Id}}" class="caption">{{.Caption}}</a></p><{{.ListType}}>{{.Value}}</{{.ListType}}>` + "\n")
	gListItemTemplate, _ = template.New("ListItem").Parse(`<li>{{.}}</li>` + "\n")
	gAnchorTemplate, _ = template.New("Anchor").Parse(`<a id="{{.Id}}" class="anchor">{{.Value}}</a>`)
	gReferToTemplate, _ = template.New("ReferTo").Parse(`<a class="referto" href="#{{.Id}}">{{.Id}}</a>`)
	gTableTemplate, _ = template.New("Table").Parse(`<p><a id="{{.Id}}" class="caption">{{.Caption}}</a></p><table>{{.Content}}</table>` + "\n")
	gTableRowTemplate, _ = template.New("TableRow").Parse(`<tr>{{.}}</tr>` + "\n")
	gTableCellTemplate, _ = template.New("TableCell").Parse(`<td>{{.}}</td>`)
}

func InlineChunkListRender(chunkList []Chunk) ([]Chunk, error) {
	var (
		curr         Chunk
		outputChunks []Chunk
	)

	pushToOutputChunks := func(chunk Chunk) {
		num := len(outputChunks)
		if num == 0 {
			outputChunks = append(outputChunks, chunk)
			return
		}
		topChunk := outputChunks[num-1]
		topPlainTextChunk, topIsPlainText := topChunk.(*PlainTextChunk)
		if !topIsPlainText {
			outputChunks = append(outputChunks, chunk)
			return
		}

		_, currIsPlainText := chunk.(*PlainTextChunk)
		if !currIsPlainText {
			outputChunks = append(outputChunks, chunk)
			return
		}
		//merge the two PlainTextChunk
		topPlainTextChunk.Value += chunk.GetValue()
	}

	for i := 0; i < len(chunkList); i++ {

		curr = chunkList[i]

		if keyword, isKeyword := curr.(*KeywordChunk); isKeyword {
			if gInlineFormatMap[keyword.Keyword] {
				str, err := KeywordChunkRender(curr)
				if err != nil {
					log.Fatalln(err)
					return outputChunks, err
				}
				pushToOutputChunks(&PlainTextChunk{Position: curr.GetPosition(), Value: str})
				continue
			}
		}

		pushToOutputChunks(curr)

	}
	return outputChunks, nil
}

func ChunkRender(chunk Chunk) (string, error) {
	switch chunk.(type) {
	case *KeywordChunk:
		return KeywordChunkRender(chunk)
	case *PlainTextChunk:
		return PlainTextChunkRender(chunk)
	case *RawTextChunk:
		return RawTextChunkRender(chunk)
	default:
		log.Fatalln("not implemented")
		panic("not implemented")
	}
}

func ChunkListRender(chunkList []Chunk) (string, error) {
	var buf bytes.Buffer
	for _, chunk := range chunkList {
		text, err := ChunkRender(chunk)
		if err != nil {
			return buf.String(), err
		}
		buf.WriteString(text)
	}
	return buf.String(), nil
}

func KeywordChunkRender(chunk Chunk) (string, error) {
	keywordChunk := chunk.(*KeywordChunk)
	var err error
	var text string
	var buf bytes.Buffer
	switch keywordChunk.Keyword {
	case EmphasisFormat:
		text, err = ChunkRender(keywordChunk.Children[0]) //only care one child
		if err != nil {
			log.Println(err)
			return text, err
		}
		err = gEmphasisTemplate.Execute(&buf, text)
		if err != nil {
			log.Println(err)
			return text, err
		}
	case StrongFormat:
		text, err = ChunkRender(keywordChunk.Children[0]) //only care one child
		if err != nil {
			log.Println(err)
			return text, err
		}
		err = gStrongTemplate.Execute(&buf, text)
		if err != nil {
			log.Println(err)
			return text, err
		}
	case HyperLink:
		url, err := ChunkRender(keywordChunk.Children[0])
		if err != nil {
			log.Println(err)
			return text, err
		}
		content, err := ChunkRender(keywordChunk.Children[1])
		if err != nil {
			log.Println(err)
			return text, err
		}
		err = gHyperLinkTemplate.Execute(&buf, struct{ Url, Text string }{url, content})
		if err != nil {
			log.Println(err)
			return text, err
		}
	case InlineCode:
		text, err = ChunkRender(keywordChunk.Children[0])
		if err != nil {
			log.Println(err)
			return text, err
		}
		err = gInlineCodeTemplate.Execute(&buf, text)
		if err != nil {
			log.Println(err)
			return text, err
		}
	case BlockCode:
		blockCodeChunk := keywordChunk.Children[0].(*BlockCodeChunk)
		err = gBlockCodeTemplate.Execute(&buf, blockCodeChunk)
		if err != nil {
			log.Println(err)
			return text, err
		}
	case AnchorBlock:
		anchorChunk := keywordChunk.Children[0].(*AnchorChunk)
		err = gAnchorTemplate.Execute(&buf, anchorChunk)
		if err != nil {
			log.Println(err)
			return text, err
		}
	case ReferToBlock:
		referToChunk := keywordChunk.Children[0].(*ReferToChunk)
		err = gReferToTemplate.Execute(&buf, referToChunk)
		if err != nil {
			log.Println(err)
			return text, err
		}
	case SectionHeader, SectionHeader1, SectionHeader2, SectionHeader3,
		SectionHeader4, SectionHeader5, SectionHeader6:
		sectionChunk := keywordChunk.Children[0].(*SectionChunk)
		err = gSectionTemplate.Execute(&buf, sectionChunk)
		if err != nil {
			log.Println(err)
			return text, err
		}
	case OrderList, BulletList:
		listChunk := keywordChunk.Children[0].(*ListChunk)
		var tmpBuf bytes.Buffer
		for _, item := range listChunk.Items {
			itemText, err := ChunkListRender(item.Value)
			if err != nil {
				return text, err
			}
			err = gListItemTemplate.Execute(&tmpBuf, itemText)
			if err != nil {
				log.Fatalln(err)
				return text, err
			}
		}
		err = gListTemplate.Execute(&buf, struct{ Id, Caption, ListType, Value string }{listChunk.Id, listChunk.Caption, listChunk.ListType, tmpBuf.String()})
		if err != nil {
			log.Println(err)
			return text, err
		}
	case TableKeyword:
		tableChunk := keywordChunk.Children[0].(*TableChunk)
		var rowBuf bytes.Buffer
		var tableBuf bytes.Buffer
		for row := 0; row < len(tableChunk.Cells); row++ {
			rowChunks := tableChunk.Cells[row]
			var cellBuf bytes.Buffer
			for col := 0; col < len(rowChunks); col++ {
				cellChunk := rowChunks[col]

				err := gTableCellTemplate.Execute(&cellBuf, cellChunk.GetValue())
				if err != nil {
					return text, err
				}
			}

			err := gTableRowTemplate.Execute(&rowBuf, cellBuf.String())
			if err != nil {
				return text, err
			}

		}
		err := gTableTemplate.Execute(&tableBuf, struct{ Id, Caption, Content string }{tableChunk.Id, tableChunk.Caption, rowBuf.String()})
		if err != nil {
			log.Println(err)
			return text, err
		}
		return tableBuf.String(), nil

	default:
		panic("not implemented")
	}
	return buf.String(), nil
}

func RawTextChunkRender(chunk Chunk) (string, error) {
	return chunk.GetValue(), nil
}

func PlainTextChunkRender(chunk Chunk) (string, error) {
	if gInlineRenderMode {
		return chunk.GetValue(), nil
	}

	var buf bytes.Buffer
	plainTextChunk := chunk.(*PlainTextChunk)
	plist := plainTextChunk.ToParagraphList()
	for _, paragraph := range plist {
		err := gParagraphTemplate.Execute(&buf, paragraph)
		if err != nil {
			return buf.String(), err
		}
	}
	return buf.String(), nil
}
