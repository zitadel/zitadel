package smtp

import (
	"bytes"
	"crypto/tls"
	"encoding/base64"
	"encoding/json"
	"net"
	"net/http"
	"net/smtp"
	"os"

	"github.com/pkg/errors"
	"github.com/zitadel/logging"

	caos_errs "github.com/zitadel/zitadel/internal/errors"
	"github.com/zitadel/zitadel/internal/notification/channels"
	"github.com/zitadel/zitadel/internal/notification/messages"
)

var _ channels.NotificationChannel = (*Email)(nil)

type Email struct {
	smtpClient      *smtp.Client
	httpRelayClient *HttpRelayClient
	senderAddress   string
	senderName      string
	replyToAddress  string
}

type HttpRelayClient struct {
	User     string
	Password string
	Host     string
	Protocol string
}

func InitChannel(cfg *Config) (*Email, error) {
	if os.Getenv("ZITADEL_SEND_SMTP_EMAIL_TO_HTTP_RELAY") == "true" {
		logging.New().Debug("successfully initialized email http relay channel")
		protocol := "http://"
		if cfg.Tls {
			protocol = "https://"
		}
		return &Email{
			smtpClient: nil,
			httpRelayClient: &HttpRelayClient{
				Host:     cfg.SMTP.Host,
				User:     cfg.SMTP.User,
				Password: cfg.SMTP.Password,
				Protocol: protocol,
			},
			senderName:     cfg.FromName,
			senderAddress:  cfg.From,
			replyToAddress: cfg.ReplyToAddress,
		}, nil
	}
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
	if email.httpRelayClient != nil {
		emailMsg, ok := message.(*messages.Email)
		if !ok {
			return caos_errs.ThrowInternal(nil, "EMAIL-proxy1", "message is not EmailMessage")
		}
		if os.Getenv("ZITADEL_HTTP_RELAY_DATA_ONLY") == "true" {
			emailMsg.Content = ""
		}
		type EmailProxyPost struct {
			EmailMsg *messages.Email
			User     string
			Password string
		}
		emailMsg.SenderEmail = email.senderAddress
		emailMsg.SenderName = email.senderName
		emailMsg.ReplyToAddress = email.replyToAddress
		emailProxyPost := EmailProxyPost{EmailMsg: emailMsg}
		if os.Getenv("ZITADEL_HTTP_RELAY_INCLUDE_CREDENTIALS") == "true" {
			emailProxyPost.User = email.httpRelayClient.User
			emailProxyPost.Password = email.httpRelayClient.Password
		}
		jsonStr, err := json.Marshal(emailProxyPost)
		if err != nil {
			return caos_errs.ThrowInternalf(nil, "EMAIL-proxy2", "could not marshal json")
		}
		req, err := http.NewRequest(http.MethodPost, email.httpRelayClient.Protocol+email.httpRelayClient.Host, bytes.NewBuffer(jsonStr))
		if err != nil {
			return caos_errs.ThrowInternalf(nil, "EMAIL-proxy3", "could not create request: "+err.Error())
		}
		if os.Getenv("ZITADEL_HTTP_RELAY_BASIC_AUTH") == "true" {
			auth := base64.StdEncoding.EncodeToString([]byte(email.httpRelayClient.User + ":" + email.httpRelayClient.Password))
			req.Header.Add("Authorization", "Basic "+auth)
		}
		_, err = http.DefaultClient.Do(req)
		if err != nil {
			return caos_errs.ThrowInternalf(nil, "EMAIL-proxy4", "error making request: "+err.Error())
		}
		return nil
	}
	defer email.smtpClient.Close()
	emailMsg, ok := message.(*messages.Email)
	if !ok {
		return caos_errs.ThrowInternal(nil, "EMAIL-s8JLs", "message is not EmailMessage")
	}

	if emailMsg.Content == "" || emailMsg.Subject == "" || len(emailMsg.Recipients) == 0 {
		return caos_errs.ThrowInternalf(nil, "EMAIL-zGemZ", "subject, recipients and content must be set but got subject %s, recipients length %d and content length %d", emailMsg.Subject, len(emailMsg.Recipients), len(emailMsg.Content))
	}
	emailMsg.SenderEmail = email.senderAddress
	emailMsg.SenderName = email.senderName
	emailMsg.ReplyToAddress = email.replyToAddress

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
		return nil, caos_errs.ThrowInternal(err, "EMAIL-spR56", "could not split host and port for connect to smtp")
	}

	if !tlsRequired {
		client, err = smtpConfig.getSMPTClient()
	} else {
		client, err = smtpConfig.getSMPTClientWithTls(host)
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
		return nil, caos_errs.ThrowInternal(err, "EMAIL-skwos", "could not make smtp dial")
	}
	return client, nil
}

func (smtpConfig SMTP) getSMPTClientWithTls(host string) (*smtp.Client, error) {
	conn, err := tls.Dial("tcp", smtpConfig.Host, &tls.Config{})

	if errors.As(err, &tls.RecordHeaderError{}) {
		logging.Log("MAIN-xKIzT").OnError(err).Warn("could not connect using normal tls. trying starttls instead...")
		return smtpConfig.getSMPTClientWithStartTls(host)
	}

	if err != nil {
		return nil, caos_errs.ThrowInternal(err, "EMAIL-sl39s", "could not make tls dial")
	}

	client, err := smtp.NewClient(conn, host)
	if err != nil {
		return nil, caos_errs.ThrowInternal(err, "EMAIL-skwi4", "could not create smtp client")
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
		return nil, caos_errs.ThrowInternal(err, "EMAIL-guvsQ", "could not start tls")
	}
	return client, nil
}

func (smtpConfig SMTP) smtpAuth(client *smtp.Client, host string) error {
	if !smtpConfig.HasAuth() {
		return nil
	}
	// Auth
	auth := unencryptedAuth{
		smtp.PlainAuth("", smtpConfig.User, smtpConfig.Password, host),
	}
	err := client.Auth(auth)
	if err != nil {
		return caos_errs.ThrowInternalf(err, "EMAIL-s9kfs", "could not add smtp auth for user %s", smtpConfig.User)
	}
	return nil
}
