package messages

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/caos/zitadel/internal/notification/channels"
)

var (
	isHTMLRgx = regexp.MustCompile(`.*<html.*>.*`)
	lineBreak = "\r\n"
)

var _ channels.Message = (*Email)(nil)

type Email struct {
	Recipients  []string
	BCC         []string
	CC          []string
	SenderEmail string
	SenderName  string
	Subject     string
	Content     string
}

func (msg *Email) GetContent() string {
	headers := make(map[string]string)
	from := msg.SenderEmail
	if msg.SenderName != "" {
		from = msg.SenderName + " " + from
	}
	headers["From"] = from
	headers["To"] = strings.Join(msg.Recipients, ", ")
	headers["Cc"] = strings.Join(msg.CC, ", ")

	message := ""
	for k, v := range headers {
		message += fmt.Sprintf("%s: %s"+lineBreak, k, v)
	}

	//default mime-type is html
	mime := "MIME-version: 1.0;" + lineBreak + "Content-Type: text/html; charset=\"UTF-8\";" + lineBreak + lineBreak
	if !isHTML(msg.Content) {
		mime = "MIME-version: 1.0;" + lineBreak + "Content-Type: text/plain; charset=\"UTF-8\";" + lineBreak + lineBreak
	}
	subject := "Subject: " + msg.Subject + lineBreak
	message += subject + mime + lineBreak + msg.Content

	return message
}

func isHTML(input string) bool {
	return isHTMLRgx.MatchString(input)
}
