package main

import (
	"bytes"
	"log"
	"text/template"
)

type OutputRenderFunc func(chunk Chunk) (string, error)

var (
	gSectionTemplate   *template.Template
	gParagraphTemplate *template.Template
	gEmphasisTemplate  *template.Template
	gStrongTemplate    *template.Template
	gHyperLinkTemplate *template.Template
)

func init() {
	gSectionTemplate, _ = template.New("Section").Parse(`<h{{.Level}} name="{{.Id}}">{{.Caption}}</h{{.Level}}>`)
	gParagraphTemplate, _ = template.New("Paragraph").Parse(`<p>{{.}}</p>`)
	gEmphasisTemplate, _ = template.New("Emphasis").Parse(`<em>{{.}}</em>`)
	gStrongTemplate, _ = template.New("Strong").Parse(`<strong>{{.}}</strong>`)
	gHyperLinkTemplate, _ = template.New("HyperLink").Parse(`<a href="{{.Url}}">{{.Text}}</a>`)
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

		if _, isInlineChunk := curr.(*InlineChunk); isInlineChunk {
			str, err := InlineChunkRender(curr)
			if err != nil {
				return outputChunks, err
			}
			pushToOutputChunks(&PlainTextChunk{Position: curr.GetPosition(), Value: str})
			continue
		}
		if embracedChunk, isEmbracedChunk := curr.(*EmbracedChunk); isEmbracedChunk {
			var err error
			embracedChunk.Children, err = InlineChunkListRender(embracedChunk.Children) //recursive call
			if err != nil {
				return outputChunks, err
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
	case *InlineChunk:
		return InlineChunkRender(chunk)
	case *EmbracedChunk:
		embracedChunk := chunk.(*EmbracedChunk)
		return ChunkListRender(embracedChunk.Children)
	default:
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
	default:
		panic("not implemented")
	}
	return buf.String(), nil
}

func InlineChunkRender(chunk Chunk) (string, error) {
	inlineChunk := chunk.(*InlineChunk)

	var err error
	var text string
	var buf bytes.Buffer

	switch inlineChunk.Keyword {
	case EmphasisFormat:
		text, err = ChunkRender(inlineChunk.Children[0]) //only care one child
		if err != nil {
			return text, err
		}

		err = gEmphasisTemplate.Execute(&buf, text)
	case StrongFormat:
		text, err = ChunkRender(inlineChunk.Children[0]) //only care one child
		if err != nil {
			return text, err
		}

		err = gStrongTemplate.Execute(&buf, text)
	case HyperLink:
		text1 := inlineChunk.Children[0].(*EmbracedChunk).Children[0].GetValue()
		text2 := inlineChunk.Children[1].(*EmbracedChunk).Children[0].GetValue()

		err = gHyperLinkTemplate.Execute(&buf, struct{ Url, Text string }{text1, text2})
	default:
		panic("not implemented")
	}
	if err != nil {
		return text, err
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
