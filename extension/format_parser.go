package extension

import (
	"github.com/apollogoClient/v1/constant"
	parse "github.com/apollogoClient/v1/utils/parser"
)

var formatParser = make(map[constant.ConfigFileFormat]parse.ContentParser, 0)

// AddFormatParser 设置 formatParser
func AddFormatParser(key constant.ConfigFileFormat, contentParser parse.ContentParser) {
	formatParser[key] = contentParser
}

// GetFormatParser 获取 formatParser
func GetFormatParser(key constant.ConfigFileFormat) parse.ContentParser {
	return formatParser[key]
}
