package email

import (
	"encoding/base64"
	"fmt"
	"github.com/caos/logging"
	"github.com/jaytaylor/html2text"
	"regexp"
	"strings"
)

var (
	isHTMLRgx = regexp.MustCompile(`.*<html.*>.*`)
)

type EmailMessage struct {
	Recipients []string
	BCC        []string
	CC         []string
	Sender     string
	Subject    string
	Content    string
}

func (msg *EmailMessage) GetContent() string {
	plainContent := toPlain(msg.Content)

	// Setup headers
	headers := make(map[string]string)
	headers["From"] = msg.Sender
	headers["To"] = strings.Join(msg.Recipients, ", ")
	headers["Cc"] = strings.Join(msg.CC, ", ")

	// Setup message
	message := ""
	mime := ""
	for k, v := range headers {
		message += fmt.Sprintf("%s: %s\r\n", k, v)
	}

	sDec, err := base64.StdEncoding.DecodeString(msg.Content)
	if err != nil {
		plainContent = msg.Content
		mime = "MIME-version: 1.0;\nContent-Type: text/plain; charset=\"UTF-8\";\n\n"
	} else {
		plainContent = string(sDec)
		mime = "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n"
	}
	subject := "Subject: " + msg.Subject + "\n"
	message += subject + mime + "\r\n" + plainContent

	return message
}

func toPlain(content string) string {
	if !isHTML(content) {
		return content
	}

	content, err := html2text.FromString(content, html2text.Options{PrettyTables: true})
	logging.Log("EMAIL-2ks94").OnError(err).Warn("could not get htmltext")

	return content
}

func isHTML(input string) bool {
	return isHTMLRgx.MatchString(input)
}
