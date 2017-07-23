package main

import (
	"bytes"
	"log"
	"text/template"
)

type OutputRenderFunc func(chunk Chunk) (string, error)

var (
	gSectionTemplate    *template.Template
	gParagraphTemplate  *template.Template
	gEmphasisTemplate   *template.Template
	gStrongTemplate     *template.Template
	gHyperLinkTemplate  *template.Template
	gInlineCodeTemplate *template.Template
	gBlockCodeTemplate  *template.Template
)

func init() {
	gSectionTemplate, _ = template.New("Section").Parse(`<h{{.Level}} id="{{.Id}}">{{.Caption}}</h{{.Level}}>`)
	gParagraphTemplate, _ = template.New("Paragraph").Parse(`<p>{{.}}</p>`)
	gEmphasisTemplate, _ = template.New("Emphasis").Parse(`<em>{{.}}</em>`)
	gStrongTemplate, _ = template.New("Strong").Parse(`<strong>{{.}}</strong>`)
	gHyperLinkTemplate, _ = template.New("HyperLink").Parse(`<a href="{{.Url}}">{{.Text}}</a>`)
	gInlineCodeTemplate, _ = template.New("InlineCode").Parse(`<code>{{.}}</code>`)
	gBlockCodeTemplate, _ = template.New("BlockCode").Parse(`<p><a id="{{.Id}}" class="caption">{{.Caption}}</a></p><pre>{{.Value}}</pre>`)
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
	case *SectionChunk:
		return SectionChunkRender(chunk)
	case *RawTextChunk:
		rawTextChunk := chunk.(*RawTextChunk)
		return rawTextChunk.GetValue(), nil
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
	case SectionHeader, SectionHeader1, SectionHeader2, SectionHeader3,
		SectionHeader4, SectionHeader5, SectionHeader6:
		sectionChunk := keywordChunk.Children[0].(*SectionChunk)
		err = gSectionTemplate.Execute(&buf, sectionChunk)
		if err != nil {
			log.Println(err)
			return text, err
		}
	default:
		panic("not implemented")
	}
	return buf.String(), nil
}

func PlainTextChunkRender(chunk Chunk) (string, error) {
	return chunk.GetValue(), nil
}

func SectionChunkRender(chunk Chunk) (string, error) {
	sectionChunk := chunk.(*SectionChunk)
	var buf bytes.Buffer
	var err error
	err = gSectionTemplate.Execute(&buf, sectionChunk)
	if err != nil {
		return "", err
	}
	//no direct paragraphs followed the section caption
	if sectionChunk.IsTerminal() {

		return buf.String(), nil
	}
	content, err := ChunkListRender(sectionChunk.Children)
	if err != nil {
		return "", err
	}
	buf.WriteString(content)
	return buf.String(), nil
}
