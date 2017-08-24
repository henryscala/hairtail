package main

import (
	"bytes"
	"flag"
	"html/template"
	"io/ioutil"
	"log"
)

var (
	gTemplate *template.Template
)

//command flags
var (
	gInputFile    = flag.String("i", "input.txt", "input file to process")
	gOutputFile   = flag.String("o", "output.html", "out file to put result")
	gTemplateFile = flag.String("t", "template.html", "template file with hole to be filled in")
	gLanguage     = flag.String("language", "cn", "language of the output (cn|en)")
)

func fileToChunks(inputFile string) ([]Chunk, error) {
	inputContent, err := ioutil.ReadFile(inputFile)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	chunks, err := ParseChunks(string(inputContent))
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return chunks, nil
}

func CompileFile(inputFile, outputFile string) error {
	gDoc.FilePath = inputFile
	chunks, err := fileToChunks(inputFile)
	if err != nil {
		log.Println(err)
		return err
	}
	//	log.Println("###########################################")
	//	log.Println("rendering", inputFile)
	//	log.Println("###########################################")
	//	log.Println("intermediate result is :")
	//	log.Println(chunks)
	gInlineRenderMode = false
	outputContent, err := ChunkListRender(chunks)
	if err != nil {
		log.Println(err)
		return err
	}

	if gConfig.TemplateFile != "" {
		templateFileContent, err := ioutil.ReadFile(gConfig.TemplateFile)
		if err != nil {
			return err
		}

		gTemplate, err = template.New("global").Parse(string(templateFileContent))
		if err != nil {
			return err
		}
		var buf bytes.Buffer
		err = gTemplate.Execute(&buf, template.HTML(outputContent))
		if err != nil {
			return err
		}

		err = ioutil.WriteFile(outputFile, buf.Bytes(), 0666)
		if err != nil {
			log.Println(err)
		}
	} else {
		err = ioutil.WriteFile(outputFile, []byte(outputContent), 0666)
		if err != nil {
			log.Println(err)
		}
	}
	return nil
}

func init() {
	log.SetFlags(log.Lshortfile)
}

func handleArguments() error {
	flag.Parse()
	gConfig.Language = *gLanguage
	gConfig.TemplateFile = *gTemplateFile

	return nil
}

func main() {

	err := handleArguments()
	if err != nil {
		log.Fatalln(err)
	}

	err = CompileFile(*gInputFile, *gOutputFile)

	if err != nil {
		log.Fatalln(err)
	}

	return
}
