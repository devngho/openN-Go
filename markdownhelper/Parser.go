package markdownhelper

//This is example markdown parser.
// Markdown to your parser name
type Example struct{}

func init() {
	register("example", Example{})
}

func (t Example) ToHTML(markdown string) string {
	return markdown
}

func (t Example) ToMarkdown(html string) string {
	return html
}
