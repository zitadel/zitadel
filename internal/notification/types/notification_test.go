package types

import (
	"context"
	"errors"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/text/language"

	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/i18n"
	zchannels "github.com/zitadel/zitadel/internal/notification/channels"
	"github.com/zitadel/zitadel/internal/notification/channels/email"
	"github.com/zitadel/zitadel/internal/notification/channels/set"
	"github.com/zitadel/zitadel/internal/notification/channels/sms"
	"github.com/zitadel/zitadel/internal/notification/channels/smtp"
	"github.com/zitadel/zitadel/internal/notification/channels/webhook"
	"github.com/zitadel/zitadel/internal/notification/messages"
	"github.com/zitadel/zitadel/internal/notification/senders"
	"github.com/zitadel/zitadel/internal/query"
)

func TestMain(m *testing.M) {
	i18n.SupportLanguages(language.English)
	os.Exit(m.Run())
}

var _ ChannelChains = (*testChannels)(nil)

// testChannels records the messages passed to the channels.
type testChannels struct {
	noChannel   bool
	emailConfig *email.Config
	rule        smtp.Rule
	ruleErr     error

	ruleRequested bool
	messages      []zchannels.Message
}

func (c *testChannels) chain() *senders.Chain {
	return senders.ChainChannels(zchannels.HandleMessageFunc(func(message zchannels.Message) error {
		c.messages = append(c.messages, message)
		return nil
	}))
}

func (c *testChannels) Email(context.Context) (*senders.Chain, *email.Config, error) {
	if c.noChannel {
		return nil, nil, errors.New("no provider")
	}
	return c.chain(), c.emailConfig, nil
}

func (c *testChannels) SMS(context.Context) (*senders.Chain, *sms.Config, error) {
	return c.chain(), nil, nil
}

func (c *testChannels) Webhook(context.Context, webhook.Config) (*senders.Chain, error) {
	return c.chain(), nil
}

func (c *testChannels) SecurityTokenEvent(context.Context, set.Config) (*senders.Chain, error) {
	return c.chain(), nil
}

func (c *testChannels) SMTPRule(context.Context, *smtp.Config, string) (smtp.Rule, error) {
	c.ruleRequested = true
	return c.rule, c.ruleErr
}

func TestSendEmail(t *testing.T) {
	const (
		mailTemplate = `<html><h1>{{.Greeting}}</h1><p>{{.Text}}</p><a href="{{.URL}}">{{.ButtonText}}</a></html>`
		urlTemplate  = "https://login.example.com/invite?userID={{.UserID}}&code={{.Code}}"
		eventType    = eventstore.EventType("user.human.invite.code.added")

		linkText      = `Click <a href="https://other.example">here</a> to join {{.ApplicationName}}`
		formattedText = `Use <strong>{{.Code}}</strong><br>to join {{.ApplicationName}}`
		// as seeded by DefaultInstance.MessageTexts
		germanText = `Dieser Benutzer wurde soeben im Zitadel erstellt. Mit dem Benutzernamen <br><strong>{{.PreferredLoginName}}</strong><br> kannst du dich anmelden. (Code <strong>{{.Code}}</strong>)`
	)
	smtpConfig := &email.Config{
		SMTPConfig: &smtp.Config{
			SMTP: smtp.SMTP{Host: "smtp.example.com:587"},
			From: "noreply@example.com",
		},
	}
	webhookConfig := &email.Config{
		WebhookConfig: &webhook.Config{CallURL: "https://relay.example.com"},
	}
	ruleErr := errors.New("rule error")

	expectMessage := func(content string) *messages.Email {
		return &messages.Email{
			Recipients:          []string{"user@example.com"},
			Subject:             "Invitation to App",
			Content:             content,
			TriggeringEventType: eventType,
		}
	}

	tests := []struct {
		name              string
		channels          *testChannels
		text              string
		displayName       string
		wantRuleRequested bool
		wantErr           func(t *testing.T, err error)
		wantMessage       zchannels.Message
	}{
		{
			name:     "no provider, canceled",
			channels: &testChannels{noChannel: true},
			wantErr: func(t *testing.T, err error) {
				assert.ErrorIs(t, err, new(zchannels.CancelError))
			},
		},
		{
			name: "rule error",
			channels: &testChannels{
				emailConfig: smtpConfig,
				ruleErr:     ruleErr,
			},
			wantRuleRequested: true,
			wantErr: func(t *testing.T, err error) {
				assert.ErrorIs(t, err, ruleErr)
			},
		},
		{
			name:              "no rule, HTML of custom texts is rendered",
			channels:          &testChannels{emailConfig: smtpConfig},
			text:              linkText,
			displayName:       "O'Brien <b>",
			wantRuleRequested: true,
			wantMessage: expectMessage(
				`<html><h1>Hello O&#39;Brien &lt;b&gt;,</h1>` +
					`<p>Click <a href="https://other.example">here</a> to join App</p>` +
					`<a href="https://login.example.com/invite?userID=user1&code=code1">Accept</a></html>`,
			),
		},
		{
			name: "HTML disabled, link is removed, arguments are still escaped",
			channels: &testChannels{
				emailConfig: smtpConfig,
				rule:        smtp.Rule{DisableCustomHTML: true},
			},
			text:              linkText,
			displayName:       "O'Brien <b>",
			wantRuleRequested: true,
			wantMessage: expectMessage(
				`<html><h1>Hello O&#39;Brien &lt;b&gt;,</h1>` +
					`<p>Click here to join App</p>` +
					`<a href="https://login.example.com/invite?userID=user1&code=code1">Accept</a></html>`,
			),
		},
		{
			name: "HTML disabled, simple formatting is kept",
			channels: &testChannels{
				emailConfig: smtpConfig,
				rule:        smtp.Rule{DisableCustomHTML: true},
			},
			text:              formattedText,
			displayName:       "Bob",
			wantRuleRequested: true,
			wantMessage: expectMessage(
				`<html><h1>Hello Bob,</h1>` +
					`<p>Use <strong>code1</strong><br>to join App</p>` +
					`<a href="https://login.example.com/invite?userID=user1&code=code1">Accept</a></html>`,
			),
		},
		{
			name: "HTML disabled, display name with allowed tag is not rendered",
			channels: &testChannels{
				emailConfig: smtpConfig,
				rule:        smtp.Rule{DisableCustomHTML: true},
			},
			text:              formattedText,
			displayName:       "<b>Bob</b><br>",
			wantRuleRequested: true,
			wantMessage: expectMessage(
				`<html><h1>Hello &lt;b&gt;Bob&lt;/b&gt;&lt;br&gt;,</h1>` +
					`<p>Use <strong>code1</strong><br>to join App</p>` +
					`<a href="https://login.example.com/invite?userID=user1&code=code1">Accept</a></html>`,
			),
		},
		{
			name:              "no rule, German default text",
			channels:          &testChannels{emailConfig: smtpConfig},
			text:              germanText,
			displayName:       "Bob",
			wantRuleRequested: true,
			wantMessage: expectMessage(
				`<html><h1>Hello Bob,</h1>` +
					`<p>Dieser Benutzer wurde soeben im Zitadel erstellt. Mit dem Benutzernamen <br><strong>bob@example.com</strong><br> kannst du dich anmelden. (Code <strong>code1</strong>)</p>` +
					`<a href="https://login.example.com/invite?userID=user1&code=code1">Accept</a></html>`,
			),
		},
		{
			name: "HTML disabled, German default text is identical",
			channels: &testChannels{
				emailConfig: smtpConfig,
				rule:        smtp.Rule{DisableCustomHTML: true},
			},
			text:              germanText,
			displayName:       "Bob",
			wantRuleRequested: true,
			wantMessage: expectMessage(
				`<html><h1>Hello Bob,</h1>` +
					`<p>Dieser Benutzer wurde soeben im Zitadel erstellt. Mit dem Benutzernamen <br><strong>bob@example.com</strong><br> kannst du dich anmelden. (Code <strong>code1</strong>)</p>` +
					`<a href="https://login.example.com/invite?userID=user1&code=code1">Accept</a></html>`,
			),
		},
		{
			name: "webhook provider, no rule requested",
			channels: &testChannels{
				emailConfig: webhookConfig,
				rule:        smtp.Rule{DisableCustomHTML: true},
			},
			text:              linkText,
			displayName:       "Bob",
			wantRuleRequested: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			translator := i18n.NewNotificationTranslator(language.English, nil)
			err := translator.AddMessages(language.English,
				i18n.Message{ID: "InviteUser.Subject", Text: "Invitation to {{.ApplicationName}}"},
				i18n.Message{ID: "InviteUser.Greeting", Text: "Hello {{.DisplayName}},"},
				i18n.Message{ID: "InviteUser.Text", Text: tt.text},
				i18n.Message{ID: "InviteUser.ButtonText", Text: "Accept"},
			)
			require.NoError(t, err)
			user := &query.NotifyUser{
				ID:                 "user1",
				ResourceOwner:      "org1",
				DisplayName:        tt.displayName,
				PreferredLoginName: "bob@example.com",
				LastEmail:          "user@example.com",
				PreferredLanguage:  language.English,
			}

			notify := SendEmail(t.Context(), tt.channels, mailTemplate, translator, user, &query.LabelPolicy{}, eventType)
			err = notify(
				urlTemplate,
				map[string]any{"Code": "code1", "ApplicationName": "App"},
				domain.InviteUserMessageType,
				true,
			)

			assert.Equal(t, tt.wantRuleRequested, tt.channels.ruleRequested)
			if tt.wantErr != nil {
				tt.wantErr(t, err)
				assert.Empty(t, tt.channels.messages)
				return
			}
			require.NoError(t, err)
			require.Len(t, tt.channels.messages, 1)
			if tt.wantMessage != nil {
				assert.Equal(t, tt.wantMessage, tt.channels.messages[0])
			}
		})
	}
}
