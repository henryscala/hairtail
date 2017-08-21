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

var gLanguageCaptionPrefix = map[string]map[string]string{
	"cn": gCaptionPrefixCn,
	"en": gCaptionPrefixEn,
}

var gCaptionPrefixCn = map[string]string{
	OrderList:    "有序列表",
	BulletList:   "无序列表",
	TableKeyword: "表格",
	BlockCode:    "代码",
	ImageKeyword: "图",
}
var gCaptionPrefixEn = map[string]string{
	OrderList:    "Ordered-List",
	BulletList:   "Bullet-List",
	TableKeyword: "Table",
	BlockCode:    "Code",
	ImageKeyword: "Figure",
}

func getCaptionPrefix(keyword string) string {
	prefixMap, ok := gLanguageCaptionPrefix[gConfig.Language]
	if !ok {
		log.Fatal("should not reach here ")
	}
	prefix, ok := prefixMap[keyword]
	if !ok {
		log.Fatal("should not reach here")
	}
	return prefix
}
