package main

import (
	"flag"
	"io/ioutil"
	"log"
)

//command flags
var (
	gInputFile  = flag.String("i", "input.txt", "input file to process")
	gOutputFile = flag.String("o", "output.html", "out file to put result")
)

func Compile(inputFile, outputFile string) error {
	inputContent, err := ioutil.ReadFile(inputFile)
	if err != nil {
		log.Println(err)
		return err
	}
	chunks, err := ParseChunks(string(inputContent))
	if err != nil {
		log.Println(err)
		return err
	}
	log.Println("###########################################")
	log.Println("rendering", inputFile)
	log.Println("###########################################")
	log.Println("intermediate result is :")
	log.Println(chunks)
	gInlineRenderMode = false
	outputContent, err := ChunkListRender(chunks)
	if err != nil {
		log.Println(err)
		return err
	}
	err = ioutil.WriteFile(outputFile, []byte(outputContent), 0666)
	if err != nil {
		log.Println(err)
	}
	return nil
}

func init() {
	log.SetFlags(log.Lshortfile)
}

func main() {

	flag.Parse()

	err := Compile(*gInputFile, *gOutputFile)

	if err != nil {
		log.Println(err)
	}

	return
}
