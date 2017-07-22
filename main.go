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

func main() {
	log.SetFlags(log.Lshortfile)
	flag.Parse()

	inputContent, err := ioutil.ReadFile(*gInputFile)
	if err != nil {
		log.Println(err)
		return
	}
	chunks, err := ParseChunks(string(inputContent))
	if err != nil {
		log.Println(err)
		return
	}

	//TODO first render inline format(May combine plain text chunks), and then section format

	outputContent, err := ChunkListRender(chunks)
	if err != nil {
		log.Println(err)
		return
	}
	err = ioutil.WriteFile(*gOutputFile, []byte(outputContent), 0666)
	if err != nil {
		log.Println(err)
	}

	return
}
