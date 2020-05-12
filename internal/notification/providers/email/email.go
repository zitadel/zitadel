package email

import (
	"crypto/tls"
	"github.com/caos/logging"
	caos_errs "github.com/caos/zitadel/internal/errors"
	"github.com/caos/zitadel/internal/notification/providers"
	"net"
	"net/smtp"
)

type Email struct {
	smtpClient *smtp.Client
	host       string
}

func InitEmailProvider(config *EmailConfig) (*Email, error) {
	client, host, err := connectToSMTP(config.SMTP, config.Tls)
	if err != nil {
		return nil, err
	}
	return &Email{
		smtpClient: client,
		host:       host,
	}, nil
}

func (email *Email) CanHandleMessage(message providers.Message) bool {
	msg := message.(EmailMessage)
	return msg.GetContent() != "" && msg.Subject != "" && len(msg.Recipients) > 0
}

func (email *Email) HandleMessage(message providers.Message) error {
	emailMsg := message.(EmailMessage)

	// To && From
	if err := email.smtpClient.Mail(emailMsg.Sender); err != nil {
		return caos_errs.ThrowInternalf(err, "EMAIL-s3is3", "could not set sender: %v", emailMsg.Sender)
	}

	for _, recp := range append(append(emailMsg.Recipients, emailMsg.CC...), emailMsg.BCC...) {
		if err := email.smtpClient.Rcpt(recp); err != nil {
			return caos_errs.ThrowInternalf(err, "EMAIL-s3is3", "could not set recipient: %v", recp)
		}
	}

	// Data
	w, err := email.smtpClient.Data()
	if err != nil {
		return err
	}

	_, err = w.Write([]byte(emailMsg.GetContent()))
	if err != nil {
		return err
	}

	err = w.Close()
	if err != nil {
		return err
	}

	defer logging.LogWithFields("EMAI-a1c87ec8").Debug("email sent")
	return email.smtpClient.Quit()
}

func connectToSMTP(smtpConfig SMTP, tlsRequired bool) (*smtp.Client, string, error) {
	host, _, _ := net.SplitHostPort(smtpConfig.Host)
	tlsconfig := &tls.Config{}

	var client *smtp.Client
	if !tlsRequired {
		var err error
		client, err = smtp.Dial(smtpConfig.Host)
		if err != nil {
			return nil, "", caos_errs.ThrowInternal(err, "EMAIL-skwos", "Could not make smtp dial")
		}
		client.StartTLS(tlsconfig)
		return client, host, nil
	}

	conn, err := tls.Dial("tcp", smtpConfig.Host, tlsconfig)
	if err != nil {
		return nil, "", caos_errs.ThrowInternal(err, "EMAIL-sl39s", "Could not make tls dial")
	}

	client, err = smtp.NewClient(conn, host)
	if err != nil {
		return nil, "", caos_errs.ThrowInternal(err, "EMAIL-skwi4", "Could not create smtp client")
	}
	defer client.Close()

	// Auth
	if smtpConfig.User != "" && smtpConfig.Password != "" {
		auth := smtp.PlainAuth("", smtpConfig.User, smtpConfig.Password, host)
		if err := client.Auth(auth); err != nil {
			logging.Log("EMAIL-s9kfs").WithField("smtp user", smtpConfig.User).Debug("Could not add smtp auth")
			return nil, "", err
		}
	}

	return client, host, nil
}
