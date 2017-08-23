package main

type Doc struct {
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
	CodeIndex string

	//Chunks                                            []Chunk
}

//to store meta data of a document
var gDoc Doc
