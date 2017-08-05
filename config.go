package main

import (
	"log"
)

type Config struct {
	Language string //en,cn
}

var GConfig = Config{
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
	prefixMap, ok := gLanguageCaptionPrefix[GConfig.Language]
	if !ok {
		log.Fatal("should not reach here ")
	}
	prefix, ok := prefixMap[keyword]
	if !ok {
		log.Fatal("should not reach here")
	}
	return prefix
}
