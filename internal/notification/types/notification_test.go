package types

import (
	"cmp"
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
	noChannel bool
	// noChain simulates a failed connection to the provider
	noChain     bool
	emailConfig *email.Config
	rule        smtp.Rule

	ruleRequested bool
	// chainsCreated counts the connections to the provider
	chainsCreated int
	messages      []zchannels.Message
}

func (c *testChannels) chain() *senders.Chain {
	return senders.ChainChannels(zchannels.HandleMessageFunc(func(message zchannels.Message) error {
		c.messages = append(c.messages, message)
		return nil
	}))
}

func (c *testChannels) EmailConfig(context.Context) (*email.Config, error) {
	if c.noChannel {
		return nil, errors.New("no provider")
	}
	return c.emailConfig, nil
}

func (c *testChannels) Email(context.Context, *email.Config) (*senders.Chain, error) {
	c.chainsCreated++
	if c.noChain {
		return senders.ChainChannels(), nil
	}
	return c.chain(), nil
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

func (c *testChannels) SMTPRule(context.Context, *smtp.Config, string) smtp.Rule {
	c.ruleRequested = true
	return c.rule
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
	expectMessage := func(content string, headers map[string]string) *messages.Email {
		return &messages.Email{
			Recipients:          []string{"user@example.com"},
			Subject:             "Invitation to App",
			Content:             content,
			Headers:             headers,
			TriggeringEventType: eventType,
		}
	}

	tests := []struct {
		name              string
		channels          *testChannels
		text              string
		displayName       string
		urlTemplate       string
		wantRuleRequested bool
		wantErr           func(t *testing.T, err error)
		wantChainsCreated int
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
			name:              "no connection to the provider, canceled",
			channels:          &testChannels{emailConfig: smtpConfig, noChain: true},
			text:              linkText,
			wantRuleRequested: true,
			wantErr: func(t *testing.T, err error) {
				assert.ErrorIs(t, err, new(zchannels.CancelError))
			},
			wantChainsCreated: 1,
		},
		{
			name:              "invalid url template, no connection to the provider",
			channels:          &testChannels{emailConfig: smtpConfig},
			text:              linkText,
			urlTemplate:       "{{.Missing",
			wantRuleRequested: true,
			wantErr: func(t *testing.T, err error) {
				assert.Error(t, err)
				assert.NotErrorIs(t, err, new(zchannels.CancelError))
			},
		},
		{
			name:              "no rule, HTML of custom texts is rendered",
			channels:          &testChannels{emailConfig: smtpConfig},
			text:              linkText,
			displayName:       "O'Brien <b>",
			wantRuleRequested: true,
			wantMessage: expectMessage(
				`<html><h1>Hello O&#39;Brien &lt;b&gt;,</h1>`+
					`<p>Click <a href="https://other.example">here</a> to join App</p>`+
					`<a href="https://login.example.com/invite?userID=user1&code=code1">Accept</a></html>`,
				nil,
			),
		},
		{
			name: "HTML disabled, link is removed, arguments are still escaped",
			channels: &testChannels{
				emailConfig: smtpConfig,
				rule:        smtp.Rule{RestrictCustomHTML: true},
			},
			text:              linkText,
			displayName:       "O'Brien <b>",
			wantRuleRequested: true,
			wantMessage: expectMessage(
				`<html><h1>Hello O&#39;Brien &lt;b&gt;,</h1>`+
					`<p>Click here to join App</p>`+
					`<a href="https://login.example.com/invite?userID=user1&code=code1">Accept</a></html>`,
				nil,
			),
		},
		{
			name: "HTML disabled, simple formatting is kept",
			channels: &testChannels{
				emailConfig: smtpConfig,
				rule:        smtp.Rule{RestrictCustomHTML: true},
			},
			text:              formattedText,
			displayName:       "Bob",
			wantRuleRequested: true,
			wantMessage: expectMessage(
				`<html><h1>Hello Bob,</h1>`+
					`<p>Use <strong>code1</strong><br>to join App</p>`+
					`<a href="https://login.example.com/invite?userID=user1&code=code1">Accept</a></html>`,
				nil,
			),
		},
		{
			name: "HTML disabled, display name with allowed tag is not rendered",
			channels: &testChannels{
				emailConfig: smtpConfig,
				rule:        smtp.Rule{RestrictCustomHTML: true},
			},
			text:              formattedText,
			displayName:       "<b>Bob</b><br>",
			wantRuleRequested: true,
			wantMessage: expectMessage(
				`<html><h1>Hello &lt;b&gt;Bob&lt;/b&gt;&lt;br&gt;,</h1>`+
					`<p>Use <strong>code1</strong><br>to join App</p>`+
					`<a href="https://login.example.com/invite?userID=user1&code=code1">Accept</a></html>`,
				nil,
			),
		},
		{
			name:              "no rule, German default text",
			channels:          &testChannels{emailConfig: smtpConfig},
			text:              germanText,
			displayName:       "Bob",
			wantRuleRequested: true,
			wantMessage: expectMessage(
				`<html><h1>Hello Bob,</h1>`+
					`<p>Dieser Benutzer wurde soeben im Zitadel erstellt. Mit dem Benutzernamen <br><strong>bob@example.com</strong><br> kannst du dich anmelden. (Code <strong>code1</strong>)</p>`+
					`<a href="https://login.example.com/invite?userID=user1&code=code1">Accept</a></html>`,
				nil,
			),
		},
		{
			name: "HTML disabled, German default text is identical",
			channels: &testChannels{
				emailConfig: smtpConfig,
				rule:        smtp.Rule{RestrictCustomHTML: true},
			},
			text:              germanText,
			displayName:       "Bob",
			wantRuleRequested: true,
			wantMessage: expectMessage(
				`<html><h1>Hello Bob,</h1>`+
					`<p>Dieser Benutzer wurde soeben im Zitadel erstellt. Mit dem Benutzernamen <br><strong>bob@example.com</strong><br> kannst du dich anmelden. (Code <strong>code1</strong>)</p>`+
					`<a href="https://login.example.com/invite?userID=user1&code=code1">Accept</a></html>`,
				nil,
			),
		},
		{
			name: "headers of the rule are passed to the message",
			channels: &testChannels{
				emailConfig: smtpConfig,
				rule: smtp.Rule{
					Headers: map[string]string{"X-Instance-ID": "instance1"},
				},
			},
			text:              linkText,
			displayName:       "Bob",
			wantRuleRequested: true,
			wantMessage: expectMessage(
				`<html><h1>Hello Bob,</h1>`+
					`<p>Click <a href="https://other.example">here</a> to join App</p>`+
					`<a href="https://login.example.com/invite?userID=user1&code=code1">Accept</a></html>`,
				map[string]string{"X-Instance-ID": "instance1"},
			),
		},
		{
			name: "webhook provider, no rule requested",
			channels: &testChannels{
				emailConfig: webhookConfig,
				rule:        smtp.Rule{RestrictCustomHTML: true},
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
				cmp.Or(tt.urlTemplate, urlTemplate),
				map[string]any{"Code": "code1", "ApplicationName": "App"},
				domain.InviteUserMessageType,
				true,
			)

			assert.Equal(t, tt.wantRuleRequested, tt.channels.ruleRequested)
			if tt.wantErr != nil {
				tt.wantErr(t, err)
				assert.Empty(t, tt.channels.messages)
				// the channels connect to the provider, they must not be created if nothing is sent
				assert.Equal(t, tt.wantChainsCreated, tt.channels.chainsCreated)
				return
			}
			require.NoError(t, err)
			require.Len(t, tt.channels.messages, 1)
			assert.Equal(t, 1, tt.channels.chainsCreated)
			if tt.wantMessage != nil {
				assert.Equal(t, tt.wantMessage, tt.channels.messages[0])
			}
		})
	}
}
