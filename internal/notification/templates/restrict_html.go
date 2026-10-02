package templates

import "github.com/microcosm-cc/bluemonday"

// restrictedHTMLPolicy keeps simple formatting of custom message texts.
// Every other element is removed, its text content is kept.
// The allowed elements cannot carry attributes, so no links, images or styles can be used.
var restrictedHTMLPolicy = bluemonday.NewPolicy().
	AllowElements("br", "strong", "b", "em", "i", "u", "p")

// RestrictHTML removes all HTML but simple formatting from the texts of the message.
// It is used if the operator disabled custom HTML for the email provider.
func (data *TemplateData) RestrictHTML() {
	data.Title = restrictHTML(data.Title)
	data.PreHeader = restrictHTML(data.PreHeader)
	data.Greeting = restrictHTML(data.Greeting)
	data.Text = restrictHTML(data.Text)
	data.ButtonText = restrictHTML(data.ButtonText)
	data.FooterText = restrictHTML(data.FooterText)
	// a footer consisting of removed elements only is not rendered
	data.IncludeFooter = data.IncludeFooter && data.FooterText != ""
}

func restrictHTML(text string) string {
	return restrictedHTMLPolicy.Sanitize(text)
}
