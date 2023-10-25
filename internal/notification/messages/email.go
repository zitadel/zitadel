package messages

import (
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/notification/channels"
	"github.com/zitadel/zitadel/internal/notification/templates"
)

var (
	isHTMLRgx = regexp.MustCompile(`.*<html.*>.*`)
	lineBreak = "\r\n"
)

var _ channels.Message = (*Email)(nil)

type Email struct {
	Recipients      []string               `json:"recipients,omitempty"`
	BCC             []string               `json:"bcc,omitempty"`
	CC              []string               `json:"cc,omitempty"`
	SenderEmail     string                 `json:"senderEmail,omitempty"`
	SenderName      string                 `json:"senderName,omitempty"`
	ReplyToAddress  string                 `json:"replyToAddress,omitempty"`
	Subject         string                 `json:"subject,omitempty"`
	Content         string                 `json:"content,omitempty"`
	SMTPMessage     string                 `json:"smtpMessage,omitempty"`
	TemplateData    templates.TemplateData `json:"templateData,omitempty"`
	TriggeringEvent eventstore.Event       `json:"-"`
}

func (msg *Email) GetContent() (string, error) {
	headers := make(map[string]string)
	from := msg.SenderEmail
	if msg.SenderName != "" {
		from = fmt.Sprintf("%s <%s>", msg.SenderName, msg.SenderEmail)
	}
	headers["From"] = from
	if msg.ReplyToAddress != "" {
		headers["Reply-to"] = msg.ReplyToAddress
	}
	headers["Return-Path"] = msg.SenderEmail
	headers["To"] = strings.Join(msg.Recipients, ", ")
	headers["Cc"] = strings.Join(msg.CC, ", ")
	headers["Date"] = time.Now().Format(time.RFC1123Z)

	message := ""
	for k, v := range headers {
		message += fmt.Sprintf("%s: %s"+lineBreak, k, v)
	}

	//default mime-type is html
	mime := "MIME-version: 1.0;" + lineBreak + "Content-Type: text/html; charset=\"UTF-8\";" + lineBreak + lineBreak
	if !isHTML(msg.Content) {
		mime = "MIME-version: 1.0;" + lineBreak + "Content-Type: text/plain; charset=\"UTF-8\";" + lineBreak + lineBreak
	}
	subject := "Subject: " + qEncodeSubject(msg.Subject) + lineBreak
	message += subject + mime + lineBreak + msg.Content

	return message, nil
}

func (msg *Email) GetTriggeringEvent() eventstore.Event {
	return msg.TriggeringEvent
}

func (msg *Email) ToJSON(includeContent, includeSMTPMessage bool) (json *JSON, err error) {
	webhookEmail := *msg
	if !includeContent {
		webhookEmail.Content = ""
	}
	if includeSMTPMessage {
		webhookEmail.SMTPMessage, err = webhookEmail.GetContent()
	}
	return &JSON{
		Serializable:    webhookEmail,
		TriggeringEvent: msg.TriggeringEvent,
	}, err
}

func isHTML(input string) bool {
	return isHTMLRgx.MatchString(input)
}

// returns a RFC1342 "Q" encoded string to allow non-ascii characters
func qEncodeSubject(subject string) string {
	return "=?utf-8?q?" + subject + "?="
}
