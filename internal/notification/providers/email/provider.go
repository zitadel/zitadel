package email

import (
	"crypto/tls"
	"net"
	"net/smtp"

	"github.com/caos/logging"
	caos_errs "github.com/caos/zitadel/internal/errors"
	"github.com/caos/zitadel/internal/notification/providers"
	"github.com/pkg/errors"
)

type Email struct {
	smtpClient *smtp.Client
}

func InitEmailProvider(config EmailConfig) (*Email, error) {
	client, err := config.SMTP.connectToSMTP(config.Tls)
	if err != nil {
		return nil, err
	}
	return &Email{
		smtpClient: client,
	}, nil
}

func (email *Email) CanHandleMessage(message providers.Message) bool {
	msg, ok := message.(*EmailMessage)
	if !ok {
		return false
	}
	return msg.Content != "" && msg.Subject != "" && len(msg.Recipients) > 0
}

func (email *Email) HandleMessage(message providers.Message) error {
	defer email.smtpClient.Close()
	emailMsg, ok := message.(*EmailMessage)
	if !ok {
		return caos_errs.ThrowInternal(nil, "EMAIL-s8JLs", "message is not EmailMessage")
	}
	// To && From
	if err := email.smtpClient.Mail(emailMsg.SenderEmail); err != nil {
		return caos_errs.ThrowInternalf(err, "EMAIL-s3is3", "could not set sender: %v", emailMsg.SenderEmail)
	}

	for _, recp := range append(append(emailMsg.Recipients, emailMsg.CC...), emailMsg.BCC...) {
		if err := email.smtpClient.Rcpt(recp); err != nil {
			return caos_errs.ThrowInternalf(err, "EMAIL-s4is4", "could not set recipient: %v", recp)
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

func (smtpConfig SMTP) connectToSMTP(tlsRequired bool) (client *smtp.Client, err error) {
	host, _, err := net.SplitHostPort(smtpConfig.Host)
	if err != nil {
		return nil, caos_errs.ThrowInternal(err, "EMAIL-spR56", "could not split host and port for connect to smtp")
	}

	if !tlsRequired {
		client, err = smtpConfig.getSMPTClient()
	} else {
		client, err = smtpConfig.getSMPTClientWithTls(host)
		if errors.As(err, &tls.RecordHeaderError{}) {
			logging.Log("EMAIL-xKIzT").OnError(err).Warn("could not connect using normal tls. Trying starttls instead...")
			client, err = smtpConfig.getSMPTClientWithStartTls(host)
		}
	}
	if err != nil {
		return nil, err
	}

	err = smtpConfig.smtpAuth(client, host)
	if err != nil {
		return nil, err
	}
	return client, nil
}

func (smtpConfig SMTP) getSMPTClient() (*smtp.Client, error) {
	client, err := smtp.Dial(smtpConfig.Host)
	if err != nil {
		return nil, caos_errs.ThrowInternal(err, "EMAIL-skwos", "Could not make smtp dial")
	}
	return client, nil
}

func (smtpConfig SMTP) getSMPTClientWithTls(host string) (*smtp.Client, error) {
	conn, err := tls.Dial("tcp", smtpConfig.Host, &tls.Config{})
	if err != nil {
		return nil, caos_errs.ThrowInternal(err, "EMAIL-sl39s", "Could not make tls dial")
	}

	client, err := smtp.NewClient(conn, host)
	if err != nil {
		return nil, caos_errs.ThrowInternal(err, "EMAIL-skwi4", "Could not create smtp client")
	}
	return client, err
}

func (smtpConfig SMTP) getSMPTClientWithStartTls(host string) (*smtp.Client, error) {
	client, err := smtpConfig.getSMPTClient()
	if err != nil {
		return nil, err
	}

	if err := client.StartTLS(&tls.Config{
		ServerName: host,
	}); err != nil {
		return nil, caos_errs.ThrowInternal(err, "EMAIL-guvsQ", "Could not start tls")
	}
	return client, nil
}

func (smtpConfig SMTP) smtpAuth(client *smtp.Client, host string) error {
	if !smtpConfig.HasAuth() {
		return nil
	}
	// Auth
	auth := smtp.PlainAuth("", smtpConfig.User, smtpConfig.Password, host)
	err := client.Auth(auth)
	logging.Log("EMAIL-s9kfs").WithField("smtp user", smtpConfig.User).OnError(err).Debug("Could not add smtp auth")
	return err
}
