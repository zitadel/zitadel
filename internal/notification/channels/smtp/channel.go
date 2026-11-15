package smtp

import (
	"crypto/tls"
	"errors"
	"net"
	"net/smtp"
	"strings"

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
	client, err := cfg.SMTP.connectToSMTP(cfg.Tls, ehloName(cfg.From))
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

// ehloName returns the domain of the sender address, which the client uses in EHLO/HELO.
// Without it net/smtp greets the server as "localhost", which some servers reject (e.g. Google's SMTP relay).
func ehloName(senderAddress string) string {
	at := strings.LastIndex(senderAddress, "@")
	if at < 0 {
		return ""
	}
	return strings.Trim(senderAddress[at+1:], " >")
}

func (smtpConfig SMTP) connectToSMTP(tlsRequired bool, ehloName string) (client *smtp.Client, err error) {
	host, _, err := net.SplitHostPort(smtpConfig.Host)
	if err != nil {
		return nil, zerrors.ThrowInternal(err, "EMAIL-spR56", "Errors.SMTP.CouldNotSplit")
	}

	if !tlsRequired {
		client, err = smtpConfig.getSMTPClient(ehloName)
	} else {
		client, err = smtpConfig.getSMTPClientWithTls(host, ehloName)
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

func (smtpConfig SMTP) getSMTPClient(ehloName string) (*smtp.Client, error) {
	client, err := smtp.Dial(smtpConfig.Host)
	if err != nil {
		return nil, zerrors.ThrowInternal(err, "EMAIL-skwos", "Errors.SMTP.CouldNotDial")
	}
	return hello(client, ehloName)
}

// hello must run before any other command, as net/smtp sends its own EHLO with the first one.
func hello(client *smtp.Client, ehloName string) (*smtp.Client, error) {
	if ehloName == "" {
		return client, nil
	}
	if err := client.Hello(ehloName); err != nil {
		client.Close()
		return nil, zerrors.ThrowInternal(err, "EMAIL-Hel0n", "Errors.SMTP.CouldNotCreateClient")
	}
	return client, nil
}

func (smtpConfig SMTP) getSMTPClientWithTls(host, ehloName string) (*smtp.Client, error) {
	conn, err := tls.Dial("tcp", smtpConfig.Host, &tls.Config{})

	if errors.As(err, &tls.RecordHeaderError{}) {
		logging.OnError(err).Warn("could not connect using normal tls. trying starttls instead...")
		return smtpConfig.getSMTPClientWithStartTls(host, ehloName)
	}

	if err != nil {
		return nil, zerrors.ThrowInternal(err, "EMAIL-sl39s", "Errors.SMTP.CouldNotDialTLS")
	}

	client, err := smtp.NewClient(conn, host)
	if err != nil {
		return nil, zerrors.ThrowInternal(err, "EMAIL-skwi4", "Errors.SMTP.CouldNotCreateClient")
	}
	return hello(client, ehloName)
}

func (smtpConfig SMTP) getSMTPClientWithStartTls(host, ehloName string) (*smtp.Client, error) {
	client, err := smtpConfig.getSMTPClient(ehloName)
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
	if smtpConfig.XOAuth2Auth != nil {
		err := client.Auth(XOAuth2Auth(*smtpConfig.XOAuth2Auth, host))
		if err != nil {
			return zerrors.ThrowInternal(err, "EMAIL-av88SB", "Errors.SMTP.CouldNotXOAuth2")
		}
		return nil
	}

	if smtpConfig.PlainAuth != nil {
		err := client.Auth(PlainOrLoginAuth(*smtpConfig.PlainAuth, host))
		if err != nil {
			return zerrors.ThrowInternal(err, "EMAIL-s9kfs", "Errors.SMTP.CouldNotAuth")
		}
		return nil
	}

	return nil
}

func TestConfiguration(cfg *Config, testEmail string) error {
	client, err := cfg.SMTP.connectToSMTP(cfg.Tls, ehloName(cfg.From))
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
		return zerrors.ThrowInternal(err, "EMAIL-s5is5", "Errors.SMTP.CouldNotSetSender")
	}

	if err := client.Rcpt(testEmail); err != nil {
		return zerrors.ThrowInternal(err, "EMAIL-s6is6", "Errors.SMTP.CouldNotSetRecipient")
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
