package email

import (
	"fmt"
	"regexp"
	"strings"
)

var (
	isHTMLRgx = regexp.MustCompile(`.*<html.*>.*`)
)

type EmailMessage struct {
	Recipients  []string
	BCC         []string
	CC          []string
	SenderEmail string
	//SenderDisplayName     string
	Subject string
	Content string
}

func (msg *EmailMessage) GetContent() string {
	headers := make(map[string]string)
	headers["From"] = msg.SenderEmail
	headers["To"] = strings.Join(msg.Recipients, ", ")
	headers["Cc"] = strings.Join(msg.CC, ", ")

	message := ""
	for k, v := range headers {
		message += fmt.Sprintf("%s: %s\r\n", k, v)
	}

	mime := ""
	if !isHTML(msg.Content) {
		mime = "MIME-version: 1.0;\nContent-Type: text/plain; charset=\"UTF-8\";\n\n"
	} else {
		mime = "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n"
	}
	subject := "Subject: " + msg.Subject + "\n"
	message += subject + mime + "\r\n" + msg.Content

	return message
}

//
//func (msg *EmailMessage) toHtml() bool {
//	if !isHTML(msg.Content) {
//		return false
//	}
//
//	content, err := html2text.FromString(msg.Content, html2text.Options{PrettyTables: true})
//	if err != nil {
//		logging.Log("EMAIL-2ks94").OnError(err).Warn("could not get htmltext")
//		return true
//	}
//	msg.Content = content
//	return true
//}

func isHTML(input string) bool {
	return isHTMLRgx.MatchString(input)
}
