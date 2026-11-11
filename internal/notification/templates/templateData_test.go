package templates

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/text/language"

	"github.com/zitadel/zitadel/internal/i18n"
)

const testMsgType = "TestMessage"

func newTestTranslator(t *testing.T, messages ...i18n.Message) *i18n.Translator {
	t.Helper()
	translator := i18n.NewNotificationTranslator(language.English, []language.Tag{language.English})
	require.NoError(t, translator.AddMessages(language.English, messages...))
	return translator
}

func TestTemplateData_Translate(t *testing.T) {
	tests := []struct {
		name             string
		greeting         string
		buttonText       string
		wantGreeting     bool
		wantButton       bool
		wantGreetingText string
		wantButtonText   string
	}{
		{
			name:             "greeting and button text set",
			greeting:         "Hello",
			buttonText:       "Click me",
			wantGreeting:     true,
			wantButton:       true,
			wantGreetingText: "Hello",
			wantButtonText:   "Click me",
		},
		{
			name:             "blank button text hides the button",
			greeting:         "Hello",
			buttonText:       " ",
			wantGreeting:     true,
			wantButton:       false,
			wantGreetingText: "Hello",
			wantButtonText:   " ",
		},
		{
			name:             "blank greeting hides the greeting",
			greeting:         " \t ",
			buttonText:       "Click me",
			wantGreeting:     false,
			wantButton:       true,
			wantGreetingText: " \t ",
			wantButtonText:   "Click me",
		},
		{
			name:             "blank greeting and button text hide both",
			greeting:         " ",
			buttonText:       " ",
			wantGreeting:     false,
			wantButton:       false,
			wantGreetingText: " ",
			wantButtonText:   " ",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			translator := newTestTranslator(t,
				i18n.Message{ID: testMsgType + ".Greeting", Text: tt.greeting},
				i18n.Message{ID: testMsgType + ".ButtonText", Text: tt.buttonText},
			)
			data := &TemplateData{}
			data.Translate(translator, testMsgType, nil, "en")
			assert.Equal(t, tt.wantGreeting, data.IncludeGreeting)
			assert.Equal(t, tt.wantButton, data.IncludeButton)
			assert.Equal(t, tt.wantGreetingText, data.Greeting)
			assert.Equal(t, tt.wantButtonText, data.ButtonText)
		})
	}
}

func TestGetParsedTemplate_defaultTemplate(t *testing.T) {
	content, err := os.ReadFile(filepath.Join("..", "static", "templates", "template.html"))
	require.NoError(t, err)

	base := TemplateData{
		Greeting:   "Hello Max",
		Text:       "Please verify",
		ButtonText: "Verify now",
		URL:        "https://example.com/verify",
	}
	tests := []struct {
		name            string
		includeGreeting bool
		includeButton   bool
	}{
		{name: "greeting and button", includeGreeting: true, includeButton: true},
		{name: "no button", includeGreeting: true},
		{name: "no greeting", includeButton: true},
		{name: "no greeting and no button"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			data := base
			data.IncludeGreeting = tt.includeGreeting
			data.IncludeButton = tt.includeButton
			got, err := GetParsedTemplate(string(content), data)
			require.NoError(t, err)
			assert.Equal(t, tt.includeGreeting, strings.Contains(got, "Hello Max"))
			assert.Equal(t, tt.includeButton, strings.Contains(got, "Verify now"))
			assert.Equal(t, tt.includeButton, strings.Contains(got, "https://example.com/verify"))
			assert.True(t, strings.Contains(got, "Please verify"))
		})
	}
}
