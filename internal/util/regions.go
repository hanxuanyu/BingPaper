package util

import "golang.org/x/text/language"

type Region struct {
	Value string `json:"value"`
	Label string `json:"label"`
}

// IsValidRegion 校验是否为标准的地区编码 (BCP 47)
func IsValidRegion(mkt string) bool {
	if mkt == "" {
		return false
	}
	_, err := language.Parse(mkt)
	return err == nil
}

var AllRegions = []Region{
	{Value: "zh-CN", Label: "中国"},
	{Value: "en-US", Label: "美国"},
	{Value: "ja-JP", Label: "日本"},
	{Value: "en-AU", Label: "澳大利亚"},
	{Value: "en-GB", Label: "英国"},
	{Value: "de-DE", Label: "德国"},
	{Value: "en-NZ", Label: "新西兰"},
	{Value: "en-CA", Label: "加拿大"},
	{Value: "fr-FR", Label: "法国"},
	{Value: "it-IT", Label: "意大利"},
	{Value: "es-ES", Label: "西班牙"},
	{Value: "pt-BR", Label: "巴西"},
	{Value: "ko-KR", Label: "韩国"},
	{Value: "en-IN", Label: "印度"},
	{Value: "ru-RU", Label: "俄罗斯"},
}
