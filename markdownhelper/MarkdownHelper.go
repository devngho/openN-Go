package markdownhelper

import (
	"errors"
	"github.com/devngho/openN-Go/iohelper"
	"github.com/devngho/openN-Go/settinghelper"
	"github.com/devngho/openN-Go/types"
)

type RawReturn struct{}

func (t RawReturn) ToHTML(markdown string) string {
	return markdown
}

func (t RawReturn) ToMarkdown(html string) string {
	return html
}

var parsers = make(map[string]types.MarkdownParser)
var parser types.MarkdownParser

func register(name string, markdownParser types.MarkdownParser) {
	parsers[name] = markdownParser
}

func SetParser() {
	use, err := settinghelper.ReadSetting("wiki", "use_markdown").Bool()
	iohelper.ErrLog(err)
	if use {
		name := settinghelper.ReadSetting("wiki", "markdown").String()
		if parsers[name] == nil {
			iohelper.ErrFatal(errors.New("can't found parser " + name))
		} else {
			parser = parsers[name]
		}
	} else {
		parser = RawReturn{}
	}
}

func ToHTML(markdown string) string {
	return parser.ToHTML(markdown)
}

func ToMarkdown(html string) string {
	return parser.ToMarkdown(html)
}
