package main

import (
	"flag"
	"fmt"
	"io/ioutil"
)

//command flags
var (
	gInputFile  = flag.String("i", "input.txt", "input file to process")
	gOutputFile = flag.String("o", "output.html", "out file to put result")
)

func main() {
	flag.Parse()

	inputContent, err := ioutil.ReadFile(*gInputFile)
	if err != nil {
		fmt.Println(err)
		return
	}
	chunks, err := ParseChunks(string(inputContent))
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("debug chunks: ", chunks)
	outputContent, err := ChunkListRender(chunks)
	if err != nil {
		fmt.Println(err)
		return
	}
	err = ioutil.WriteFile(*gOutputFile, []byte(outputContent), 0666)
	if err != nil {
		fmt.Println(err)
	}
	return
}
