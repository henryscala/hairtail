package main

type Doc struct {
	Title, SubTitle, CreateDate, ModifyDate, Keywords string
	Chunks                                            []Chunk
}
