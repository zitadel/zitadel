package smtp

import (
	"crypto/tls"
	"errors"
	"net"
	"net/smtp"

	"github.com/zitadel/logging"

	"github.com/zitadel/zitadel/internal/notification/channels"
	"github.com/zitadel/zitadel/internal/notification/messages"
	"github.com/zitadel/zitadel/internal/zerrors"
)

var _ channels.NotificationChannel = (*Email)(nil)

type Email struct {
	smtpClient     *smtp.Client
	senderAddress  string
	senderName     string
	replyToAddress string
}

func InitChannel(cfg *Config) (*Email, error) {
	client, err := cfg.SMTP.connectToSMTP(cfg.Tls)
	if err != nil {
		logging.New().WithError(err).Error("could not connect to smtp")
		return nil, err
	}
	logging.New().Debug("successfully initialized smtp email channel")
	return &Email{
		smtpClient:     client,
		senderName:     cfg.FromName,
		senderAddress:  cfg.From,
		replyToAddress: cfg.ReplyToAddress,
	}, nil
}

func (email *Email) HandleMessage(message channels.Message) error {
	defer email.smtpClient.Close()
	emailMsg, ok := message.(*messages.Email)
	if !ok {
		return zerrors.ThrowInternal(nil, "EMAIL-s8JLs", "Errors.SMTP.NotEmailMessage")
	}

	if emailMsg.Content == "" || emailMsg.Subject == "" || len(emailMsg.Recipients) == 0 {
		return zerrors.ThrowInternal(nil, "EMAIL-zGemZ", "Errors.SMTP.RequiredAttributes")
	}
	emailMsg.SenderEmail = email.senderAddress
	emailMsg.SenderName = email.senderName
	emailMsg.ReplyToAddress = email.replyToAddress
	// To && From
	if err := email.smtpClient.Mail(emailMsg.SenderEmail); err != nil {
		return zerrors.ThrowInternal(err, "EMAIL-s3is3", "Errors.SMTP.CouldNotSetSender")
	}
	for _, recp := range append(append(emailMsg.Recipients, emailMsg.CC...), emailMsg.BCC...) {
		if err := email.smtpClient.Rcpt(recp); err != nil {
			return zerrors.ThrowInternal(err, "EMAIL-s4is4", "Errors.SMTP.CouldNotSetRecipient")
		}
	}

	// Data
	w, err := email.smtpClient.Data()
	if err != nil {
		return err
	}

	content, err := emailMsg.GetContent()
	if err != nil {
		return err
	}

	_, err = w.Write([]byte(content))
	if err != nil {
		return err
	}

	err = w.Close()
	if err != nil {
		return err
	}

	return email.smtpClient.Quit()
}

func (smtpConfig SMTP) connectToSMTP(tlsRequired bool) (client *smtp.Client, err error) {
	host, _, err := net.SplitHostPort(smtpConfig.Host)
	if err != nil {
		return nil, zerrors.ThrowInternal(err, "EMAIL-spR56", "Errors.SMTP.CouldNotSplit")
	}

	if !tlsRequired {
		client, err = smtpConfig.getSMTPClient()
	} else {
		client, err = smtpConfig.getSMTPClientWithTls(host)
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

func (smtpConfig SMTP) getSMTPClient() (*smtp.Client, error) {
	client, err := smtp.Dial(smtpConfig.Host)
	if err != nil {
		return nil, zerrors.ThrowInternal(err, "EMAIL-skwos", "Errors.SMTP.CouldNotDial")
	}
	return client, nil
}

func (smtpConfig SMTP) getSMTPClientWithTls(host string) (*smtp.Client, error) {
	conn, err := tls.Dial("tcp", smtpConfig.Host, &tls.Config{})

	if errors.As(err, &tls.RecordHeaderError{}) {
		logging.OnError(err).Warn("could not connect using normal tls. trying starttls instead...")
		return smtpConfig.getSMTPClientWithStartTls(host)
	}

	if err != nil {
		return nil, zerrors.ThrowInternal(err, "EMAIL-sl39s", "Errors.SMTP.CouldNotDialTLS")
	}

	client, err := smtp.NewClient(conn, host)
	if err != nil {
		return nil, zerrors.ThrowInternal(err, "EMAIL-skwi4", "Errors.SMTP.CouldNotCreateClient")
	}
	return client, err
}

func (smtpConfig SMTP) getSMTPClientWithStartTls(host string) (*smtp.Client, error) {
	client, err := smtpConfig.getSMTPClient()
	if err != nil {
		return nil, err
	}

	if err := client.StartTLS(&tls.Config{
		ServerName: host,
	}); err != nil {
		return nil, zerrors.ThrowInternal(err, "EMAIL-guvsQ", "Errors.SMTP.CouldNotStartTLS")
	}
	return client, nil
}

func (smtpConfig SMTP) smtpAuth(client *smtp.Client, host string) error {
	if !smtpConfig.HasAuth() {
		return nil
	}
	// Auth
	err := client.Auth(PlainOrLoginAuth(smtpConfig.User, smtpConfig.Password, host))
	if err != nil {
		return zerrors.ThrowInternal(err, "EMAIL-s9kfs", "Errors.SMTP.CouldNotAuth")
	}
	return nil
}

func TestConfiguration(cfg *Config, testEmail string) error {
	client, err := cfg.SMTP.connectToSMTP(cfg.Tls)
	if err != nil {
		return err
	}

	defer client.Close()

	message := &messages.Email{
		Recipients:  []string{testEmail},
		Subject:     "Test email",
		Content:     "This is a test email to check if your SMTP provider works fine",
		SenderEmail: cfg.From,
		SenderName:  cfg.FromName,
	}

	if err := client.Mail(cfg.From); err != nil {
		return zerrors.ThrowInternal(err, "EMAIL-s3is3", "Errors.SMTP.CouldNotSetSender")
	}

	if err := client.Rcpt(testEmail); err != nil {
		return zerrors.ThrowInternal(err, "EMAIL-s4is4", "Errors.SMTP.CouldNotSetRecipient")
	}

	// Open data connection
	w, err := client.Data()
	if err != nil {
		return err
	}

	// Send content
	content, err := message.GetContent()
	if err != nil {
		return err
	}
	_, err = w.Write([]byte(content))
	if err != nil {
		return err
	}

	// Close IO and quit smtp connection
	err = w.Close()
	if err != nil {
		return err
	}

	return client.Quit()
}
