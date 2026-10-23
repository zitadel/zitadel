package templates

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_restrictHTML(t *testing.T) {
	tests := []struct {
		name string
		text string
		want string
	}{
		{
			name: "plain text",
			text: "Hello, please click the button below.",
			want: "Hello, please click the button below.",
		},
		{
			name: "allowed formatting is kept",
			text: "Use <strong>code1</strong><br>to <b>continue</b> <em>now</em> <i>or</i> <u>later</u><p>Thanks</p>",
			want: "Use <strong>code1</strong><br>to <b>continue</b> <em>now</em> <i>or</i> <u>later</u><p>Thanks</p>",
		},
		{
			name: "any case and self closing break",
			text: "<BR><Strong>x</STRONG><br/><br />",
			want: "<br><strong>x</strong><br/><br/>",
		},
		{
			name: "attributes are removed",
			text: `<strong style="color:red" onclick="x()">Urgent</strong><br class="x">`,
			want: "<strong>Urgent</strong><br>",
		},
		{
			name: "link is removed, text is kept",
			text: `Click <a href="https://other.example">here</a> to join`,
			want: "Click here to join",
		},
		{
			name: "image, script and style are removed",
			text: `<img src="https://other.example/x.png"><script>alert(1)</script><style>p{}</style>after`,
			want: "after",
		},
		{
			name: "layout elements are removed, text is kept",
			text: `<div><h1>Title</h1><table><tr><td>cell</td></tr></table><span>span</span></div>`,
			want: "Titlecellspan",
		},
		{
			name: "escaped user data stays escaped",
			text: "Hello O&#39;Brien &lt;b&gt;Bob &amp; Co,",
			want: "Hello O&#39;Brien &lt;b&gt;Bob &amp; Co,",
		},
		{
			name: "lone angle brackets are escaped",
			text: "1 < 2 > 0",
			want: "1 &lt; 2 &gt; 0",
		},
		{
			name: "unknown element names are removed",
			text: "<bold>x</bold><brx>",
			want: "x",
		},
		{
			name: "German default text",
			text: "Mit dem Benutzernamen <br><strong>name</strong><br> kannst du dich anmelden (Code <strong>code1</strong>).",
			want: "Mit dem Benutzernamen <br><strong>name</strong><br> kannst du dich anmelden (Code <strong>code1</strong>).",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, restrictHTML(tt.text))
		})
	}
}

func TestTemplateData_RestrictHTML(t *testing.T) {
	data := &TemplateData{
		Title:      `<a href="x">Title</a>`,
		PreHeader:  `<a href="x">PreHeader</a>`,
		Subject:    `<a href="x">Subject</a>`,
		Greeting:   `<a href="x">Greeting</a>`,
		Text:       `<a href="x">Text</a>`,
		URL:        `https://login.example.com/?a=<b>`,
		ButtonText: `<a href="x">ButtonText</a>`,
		FooterText: `<a href="x">FooterText</a>`,
	}
	data.RestrictHTML()
	assert.Equal(t, &TemplateData{
		Title:      "Title",
		PreHeader:  "PreHeader",
		Subject:    `<a href="x">Subject</a>`,
		Greeting:   "Greeting",
		Text:       "Text",
		URL:        `https://login.example.com/?a=<b>`,
		ButtonText: "ButtonText",
		FooterText: "FooterText",
	}, data)
}

func TestTemplateData_RestrictHTML_footer(t *testing.T) {
	tests := []struct {
		name          string
		footer        string
		includeFooter bool
		want          string
		wantInclude   bool
	}{
		{name: "text footer stays", footer: "Contact <b>us</b>", includeFooter: true, want: "Contact <b>us</b>", wantInclude: true},
		{name: "footer of removed elements only is dropped", footer: `<img src="https://other.example/x.png">`, includeFooter: true, want: "", wantInclude: false},
		{name: "not included stays not included", footer: "text", includeFooter: false, want: "text", wantInclude: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			data := &TemplateData{FooterText: tt.footer, IncludeFooter: tt.includeFooter}
			data.RestrictHTML()
			assert.Equal(t, tt.want, data.FooterText)
			assert.Equal(t, tt.wantInclude, data.IncludeFooter)
		})
	}
}
