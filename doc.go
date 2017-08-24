package main

type Doc struct {
	FilePath, //the file path of the document to be compiled

	//meta data
	Title,
	SubTitle,
	Author,
	CreateDate,
	ModifyDate,
	Keywords string

	//index
	SectionIndex,
	ImageIndex,
	TableIndex,
	OrderListIndex,
	BulletListIndex,
	CodeIndex,
	MathIndex string

	//Chunks                                            []Chunk
}

//to store meta data of a document
var gDoc Doc
