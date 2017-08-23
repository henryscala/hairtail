package main

import (
	"bytes"
	"log"
	"text/template"
)

var (
	gSectionTemplate      *template.Template
	gParagraphTemplate    *template.Template
	gEmphasisTemplate     *template.Template
	gStrongTemplate       *template.Template
	gHyperLinkTemplate    *template.Template
	gInlineTexTemplate    *template.Template
	gInlineCodeTemplate   *template.Template
	gBlockTexTemplate     *template.Template
	gBlockCodeTemplate    *template.Template
	gListTemplate         *template.Template
	gListItemTemplate     *template.Template
	gAnchorTemplate       *template.Template
	gReferToTemplate      *template.Template
	gTableTemplate        *template.Template
	gTableRowTemplate     *template.Template
	gTableCellTemplate    *template.Template
	gImageTemplate        *template.Template
	gSectionIndexTemplate *template.Template
	gGlobalIndexTemplate  *template.Template //to generate index for entities other than section
	gTitleTemplate        *template.Template
	gMetaDataTemplate     *template.Template

	gInlineRenderMode bool
)

func init() {
	gSectionTemplate, _ = template.New("Section").Parse(`<h{{.Level}} id="{{.Id}}">{{.Numbering}} {{.Caption}}</h{{.Level}}>` + "\n")
	gParagraphTemplate, _ = template.New("Paragraph").Parse(`<p>{{.}}</p>` + "\n")
	gEmphasisTemplate, _ = template.New("Emphasis").Parse(`<em>{{.}}</em>`)
	gStrongTemplate, _ = template.New("Strong").Parse(`<strong>{{.}}</strong>`)
	gHyperLinkTemplate, _ = template.New("HyperLink").Parse(`<a href="{{.Url}}">{{.Text}}</a>`)
	gInlineTexTemplate, _ = template.New("InlineTex").Parse(`<span class="inline-tex">\({{.}}\)</span>`) //need mathjax to support this
	gInlineCodeTemplate, _ = template.New("InlineCode").Parse(`<code>{{.}}</code>`)
	gBlockTexTemplate, _ = template.New("BlockTex").Parse(`<p><a id="{{.Id}}" class="caption">{{.Numbering}} {{.Caption}}</a></p><div class="math">\[{{.Value}}\]</div>` + "\n")
	gBlockCodeTemplate, _ = template.New("BlockCode").Parse(`<p><a id="{{.Id}}" class="caption">{{.Numbering}} {{.Caption}}</a></p><pre>{{.Value}}</pre>` + "\n")
	gListTemplate, _ = template.New("List").Parse(`<p><a id="{{.Id}}" class="caption">{{.Numbering}} {{.Caption}}</a></p><{{.ListType}}>{{.Value}}</{{.ListType}}>` + "\n")
	gListItemTemplate, _ = template.New("ListItem").Parse(`<li>{{.}}</li>` + "\n")
	gAnchorTemplate, _ = template.New("Anchor").Parse(`<a id="{{.Id}}" class="anchor">{{.Value}}</a>`)
	gReferToTemplate, _ = template.New("ReferTo").Parse(`<a class="referto" href="#{{.Id}}">{{.Id}}</a>`)
	gTableTemplate, _ = template.New("Table").Parse(`<p><a id="{{.Id}}" class="caption">{{.Numbering}} {{.Caption}}</a></p><table>{{.Content}}</table>` + "\n")
	gTableRowTemplate, _ = template.New("TableRow").Parse(`<tr>{{.}}</tr>` + "\n")
	gTableCellTemplate, _ = template.New("TableCell").Parse(`<td>{{.}}</td>`)
	gImageTemplate, _ = template.New("Image").Parse(`<p><a id="{{.Id}}" class="caption">{{.Numbering}} {{.Caption}}</a></p><img src="{{.Src}}" alt="{{.Caption}}">`)
	gGlobalIndexTemplate, _ = template.New("GlobalIndex").Parse(`<p><a href="#{{.Id}}">{{.Numbering}} {{.Caption}}</a></p>` + "\n")
	gSectionIndexTemplate, _ = template.New("SectionIndex").Parse(`<p><a href="#{{.Id}}">{{.Numbering}} {{.Caption}}</a></p>` + "\n")
	gTitleTemplate, _ = template.New("Title").Parse(`<h{{.Level}} class="title{{.Level}}">{{.Title}}</h{{.Level}}>` + "\n")
	gMetaDataTemplate, _ = template.New("MetaData").Parse(`<span class="meta-data-name"><strong>{{.Name}}:</strong></span> <span class="meta-data-value">{{.Value}}</span>` + "\n")
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
	case ImageKeyword:
		imageChunk := keywordChunk.Children[0].(*ImageChunk)
		err = gImageTemplate.Execute(&buf, imageChunk)
		if err != nil {
			log.Println(err)
			return text, err
		}
	case InlineTex:
		text, err = ChunkRender(keywordChunk.Children[0])
		if err != nil {
			log.Println(err)
			return text, err
		}
		err = gInlineTexTemplate.Execute(&buf, text)
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
	case BlockTex:
		blockTexChunk := keywordChunk.Children[0].(*BlockTexChunk)
		err = gBlockTexTemplate.Execute(&buf, blockTexChunk)
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
		err = gListTemplate.Execute(&buf, struct{ Id, Caption, Numbering, ListType, Value string }{listChunk.Id, listChunk.Caption, listChunk.Numbering, listChunk.ListType, tmpBuf.String()})
		if err != nil {
			log.Println(err)
			return text, err
		}
	case TableKeyword:
		tableChunk := keywordChunk.Children[0].(*TableChunk)
		var rowBuf bytes.Buffer

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
		err := gTableTemplate.Execute(&buf, struct{ Id, Caption, Numbering, Content string }{tableChunk.Id, tableChunk.Caption, tableChunk.Numbering, rowBuf.String()})
		if err != nil {
			log.Println(err)
			return text, err
		}

	//output different kind of index
	case SectionIndexKeyword:
		return gDoc.SectionIndex, nil
	case ImageIndexKeyword:
		return gDoc.ImageIndex, nil
	case TableIndexKeyword:
		return gDoc.TableIndex, nil
	case OrderListIndexKeyword:
		return gDoc.OrderListIndex, nil
	case BulletListIndexKeyword:
		return gDoc.BulletListIndex, nil
	case CodeIndexKeyword:
		return gDoc.CodeIndex, nil
	case MathIndexKeyword:
		return gDoc.MathIndex, nil

	//different kind of meta data handling
	case TitleKeyword:
		err = gTitleTemplate.Execute(&buf, struct {
			Level int
			Title string
		}{1, gDoc.Title})
		if err != nil {
			return text, err
		}

	case SubTitleKeyword:
		err = gTitleTemplate.Execute(&buf, struct {
			Level int
			Title string
		}{2, gDoc.SubTitle})
		if err != nil {
			return text, err
		}

	case AuthorKeyword:
		err = gMetaDataTemplate.Execute(&buf, struct{ Name, Value string }{getKeywordName(AuthorKeyword), gDoc.Author})
		if err != nil {
			return text, err
		}

	case CreateDateKeyword:
		err = gMetaDataTemplate.Execute(&buf, struct{ Name, Value string }{getKeywordName(CreateDateKeyword), gDoc.CreateDate})
		if err != nil {
			return text, err
		}
	case ModifyDateKeyword:
		err = gMetaDataTemplate.Execute(&buf, struct{ Name, Value string }{getKeywordName(ModifyDateKeyword), gDoc.ModifyDate})
		if err != nil {
			return text, err
		}
	case KeywordsKeyword:
		err = gMetaDataTemplate.Execute(&buf, struct{ Name, Value string }{getKeywordName(KeywordsKeyword), gDoc.Keywords})
		if err != nil {
			return text, err
		}

	default:
		log.Fatal(errNotImplemented)
		panic(errNotImplemented)
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
