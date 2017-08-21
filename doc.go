package main

type Doc struct {
	Title, SubTitle, CreateDate, ModifyDate, Keywords string
	//Chunks                                            []Chunk
	SectionIndex,
	ImageIndex,
	TableIndex,
	OrderListIndex,
	BulletListIndex,
	CodeIndex string
}

//to store meta data of a document
var gDoc Doc
