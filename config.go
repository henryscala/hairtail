package main

import (
	"log"
)

type Config struct {
	Language     string //en,cn
	TemplateFile string // the template file with whole to put render result in
	//GenerateTitle bool   //main title and sub title
	//GenerateMeta  bool   //create date, modify date, keywords
}

var gConfig = Config{
	Language: "cn",
}

var gLanguageKeywordName = map[string]map[string]string{
	"cn": gKeywordNameCn,
	"en": gKeywordNameEn,
}

var gKeywordNameCn = map[string]string{
	OrderList:    "有序列表",
	BulletList:   "无序列表",
	TableKeyword: "表格",
	BlockCode:    "代码",
	ImageKeyword: "图",

	AuthorKeyword:     "作者",
	CreateDateKeyword: "创建日期",
	ModifyDateKeyword: "修改日期",
	KeywordsKeyword:   "关键词",
}
var gKeywordNameEn = map[string]string{
	OrderList:    "Ordered-List",
	BulletList:   "Bullet-List",
	TableKeyword: "Table",
	BlockCode:    "Code",
	ImageKeyword: "Figure",

	AuthorKeyword:     "Author",
	CreateDateKeyword: "Create-Date",
	ModifyDateKeyword: "Modify-Date",
	KeywordsKeyword:   "Keywords",
}

func getKeywordName(keyword string) string {
	prefixMap, ok := gLanguageKeywordName[gConfig.Language]
	if !ok {
		log.Fatal("not supported language")
	}
	prefix, ok := prefixMap[keyword]
	if !ok {
		log.Fatal("not supported keyword")
	}
	return prefix
}
