package util

type Region struct {
	Value string `json:"value"`
	Label string `json:"label"`
}

var AllRegions = []Region{
	{Value: "zh-CN", Label: "中国 (zh-CN)"},
	{Value: "en-US", Label: "美国 (en-US)"},
	{Value: "ja-JP", Label: "日本 (ja-JP)"},
	{Value: "en-AU", Label: "澳大利亚 (en-AU)"},
	{Value: "en-GB", Label: "英国 (en-GB)"},
	{Value: "de-DE", Label: "德国 (de-DE)"},
	{Value: "en-NZ", Label: "新西兰 (en-NZ)"},
	{Value: "en-CA", Label: "加拿大 (en-CA)"},
	{Value: "fr-FR", Label: "法国 (fr-FR)"},
	{Value: "it-IT", Label: "意大利 (it-IT)"},
	{Value: "es-ES", Label: "西班牙 (es-ES)"},
	{Value: "pt-BR", Label: "巴西 (pt-BR)"},
	{Value: "ko-KR", Label: "韩国 (ko-KR)"},
	{Value: "en-IN", Label: "印度 (en-IN)"},
	{Value: "ru-RU", Label: "俄罗斯 (ru-RU)"},
	{Value: "zh-HK", Label: "中国香港 (zh-HK)"},
	{Value: "zh-TW", Label: "中国台湾 (zh-TW)"},
}
