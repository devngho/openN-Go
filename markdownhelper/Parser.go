package markdownhelper


//This is example markdown parser. using html-to-markdown, markdown.
// Markdown to your parser name
type Markdown struct {}

func init() {
	register("example", Markdown{})
}

func (t Markdown) ToHTML(markdown string) string {
	//Markup to html code here
	return markdown
}

func (t Markdown) ToMarkdown(html string) string {
	//Html to markup code here
	return html
}