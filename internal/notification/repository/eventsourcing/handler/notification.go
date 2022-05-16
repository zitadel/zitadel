package handler

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/zitadel/logging"

	"github.com/zitadel/zitadel/internal/api/authz"
	http_utils "github.com/zitadel/zitadel/internal/api/http"
	"github.com/zitadel/zitadel/internal/command"
	"github.com/zitadel/zitadel/internal/crypto"
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/errors"
	"github.com/zitadel/zitadel/internal/eventstore"
	v1 "github.com/zitadel/zitadel/internal/eventstore/v1"
	"github.com/zitadel/zitadel/internal/eventstore/v1/models"
	queryv1 "github.com/zitadel/zitadel/internal/eventstore/v1/query"
	"github.com/zitadel/zitadel/internal/eventstore/v1/spooler"
	"github.com/zitadel/zitadel/internal/i18n"
	"github.com/zitadel/zitadel/internal/notification/channels/fs"
	"github.com/zitadel/zitadel/internal/notification/channels/log"
	"github.com/zitadel/zitadel/internal/notification/channels/smtp"
	"github.com/zitadel/zitadel/internal/notification/channels/twilio"
	"github.com/zitadel/zitadel/internal/notification/types"
	"github.com/zitadel/zitadel/internal/query"
	user_repo "github.com/zitadel/zitadel/internal/repository/user"
	es_model "github.com/zitadel/zitadel/internal/user/repository/eventsourcing/model"
	"github.com/zitadel/zitadel/internal/user/repository/view"
	"github.com/zitadel/zitadel/internal/user/repository/view/model"
)

const (
	notificationTable = "notification.notifications"
	NotifyUserID      = "NOTIFICATION"
)

type Notification struct {
	handler
	command            *command.Commands
	fileSystemPath     string
	statikDir          http.FileSystem
	subscription       *v1.Subscription
	assetsPrefix       string
	queries            *query.Queries
	userDataCrypto     crypto.EncryptionAlgorithm
	smtpPasswordCrypto crypto.EncryptionAlgorithm
	smsTokenCrypto     crypto.EncryptionAlgorithm
	externalPort       uint16
	externalSecure     bool
}

func newNotification(
	handler handler,
	command *command.Commands,
	query *query.Queries,
	externalPort uint16,
	externalSecure bool,
	statikDir http.FileSystem,
	assetsPrefix,
	fileSystemPath string,
	userEncryption crypto.EncryptionAlgorithm,
	smtpEncryption crypto.EncryptionAlgorithm,
	smsEncryption crypto.EncryptionAlgorithm,
) *Notification {
	h := &Notification{
		handler:            handler,
		command:            command,
		statikDir:          statikDir,
		assetsPrefix:       assetsPrefix,
		queries:            query,
		userDataCrypto:     userEncryption,
		smtpPasswordCrypto: smtpEncryption,
		smsTokenCrypto:     smsEncryption,
		externalSecure:     externalSecure,
		externalPort:       externalPort,
		fileSystemPath:     fileSystemPath,
	}

	h.subscribe()

	return h
}

func (k *Notification) subscribe() {
	k.subscription = k.es.Subscribe(k.AggregateTypes()...)
	go func() {
		for event := range k.subscription.Events {
			queryv1.ReduceEvent(k, event)
		}
	}()
}

func (n *Notification) ViewModel() string {
	return notificationTable
}

func (n *Notification) Subscription() *v1.Subscription {
	return n.subscription
}

func (_ *Notification) AggregateTypes() []models.AggregateType {
	return []models.AggregateType{user_repo.AggregateType}
}

func (n *Notification) CurrentSequence(instanceID string) (uint64, error) {
	sequence, err := n.view.GetLatestNotificationSequence(instanceID)
	if err != nil {
		return 0, err
	}
	return sequence.CurrentSequence, nil
}

func (n *Notification) EventQuery() (*models.SearchQuery, error) {
	sequences, err := n.view.GetLatestNotificationSequences()
	if err != nil {
		return nil, err
	}
	query := models.NewSearchQuery()
	instances := make([]string, 0)
	for _, sequence := range sequences {
		for _, instance := range instances {
			if sequence.InstanceID == instance {
				break
			}
		}
		instances = append(instances, sequence.InstanceID)
		query.AddQuery().
			AggregateTypeFilter(n.AggregateTypes()...).
			LatestSequenceFilter(sequence.CurrentSequence).
			InstanceIDFilter(sequence.InstanceID)
	}
	return query.AddQuery().
		AggregateTypeFilter(n.AggregateTypes()...).
		LatestSequenceFilter(0).
		ExcludedInstanceIDsFilter(instances...).
		SearchQuery(), nil
}

func (n *Notification) Reduce(event *models.Event) (err error) {
	switch eventstore.EventType(event.Type) {
	case user_repo.UserV1InitialCodeAddedType,
		user_repo.HumanInitialCodeAddedType:
		err = n.handleInitUserCode(event)
	case user_repo.UserV1EmailCodeAddedType,
		user_repo.HumanEmailCodeAddedType:
		err = n.handleEmailVerificationCode(event)
	case user_repo.UserV1PhoneCodeAddedType,
		user_repo.HumanPhoneCodeAddedType:
		err = n.handlePhoneVerificationCode(event)
	case user_repo.UserV1PasswordCodeAddedType,
		user_repo.HumanPasswordCodeAddedType:
		err = n.handlePasswordCode(event)
	case user_repo.UserDomainClaimedType:
		err = n.handleDomainClaimed(event)
	case user_repo.HumanPasswordlessInitCodeRequestedType:
		err = n.handlePasswordlessRegistrationLink(event)
	}
	if err != nil {
		return err
	}
	return n.view.ProcessedNotificationSequence(event)
}

func (n *Notification) handleInitUserCode(event *models.Event) (err error) {
	initCode := new(es_model.InitUserCode)
	if err := initCode.SetData(event); err != nil {
		return err
	}
	ctx := getSetNotifyContextData(event.InstanceID, event.ResourceOwner)
	alreadyHandled, err := n.checkIfCodeAlreadyHandledOrExpired(ctx, event, initCode.Expiry,
		user_repo.UserV1InitialCodeAddedType, user_repo.UserV1InitialCodeSentType,
		user_repo.HumanInitialCodeAddedType, user_repo.HumanInitialCodeSentType)
	if err != nil || alreadyHandled {
		return err
	}
	colors, err := n.getLabelPolicy(ctx)
	if err != nil {
		return err
	}

	template, err := n.getMailTemplate(ctx)
	if err != nil {
		return err
	}

	user, err := n.getUserByID(event.AggregateID, event.InstanceID)
	if err != nil {
		return err
	}

	translator, err := n.getTranslatorWithOrgTexts(ctx, user.ResourceOwner, domain.InitCodeMessageType)
	if err != nil {
		return err
	}

	origin, err := n.origin(ctx)
	if err != nil {
		return err
	}
	err = types.SendUserInitCode(ctx, string(template.Template), translator, user, initCode, n.getSMTPConfig, n.getFileSystemProvider, n.getLogProvider, n.userDataCrypto, colors, n.assetsPrefix, origin)
	if err != nil {
		return err
	}
	return n.command.HumanInitCodeSent(ctx, event.ResourceOwner, event.AggregateID)
}

func (n *Notification) handlePasswordCode(event *models.Event) (err error) {
	pwCode := new(es_model.PasswordCode)
	if err := pwCode.SetData(event); err != nil {
		return err
	}
	ctx := getSetNotifyContextData(event.InstanceID, event.ResourceOwner)
	alreadyHandled, err := n.checkIfCodeAlreadyHandledOrExpired(ctx, event, pwCode.Expiry,
		user_repo.UserV1PasswordCodeAddedType, user_repo.UserV1PasswordCodeSentType,
		user_repo.HumanPasswordCodeAddedType, user_repo.HumanPasswordCodeSentType)
	if err != nil || alreadyHandled {
		return err
	}
	colors, err := n.getLabelPolicy(ctx)
	if err != nil {
		return err
	}

	template, err := n.getMailTemplate(ctx)
	if err != nil {
		return err
	}

	user, err := n.getUserByID(event.AggregateID, event.InstanceID)
	if err != nil {
		return err
	}

	translator, err := n.getTranslatorWithOrgTexts(ctx, user.ResourceOwner, domain.PasswordResetMessageType)
	if err != nil {
		return err
	}

	origin, err := n.origin(ctx)
	if err != nil {
		return err
	}
	err = types.SendPasswordCode(ctx, string(template.Template), translator, user, pwCode, n.getSMTPConfig, n.getTwilioConfig, n.getFileSystemProvider, n.getLogProvider, n.userDataCrypto, colors, n.assetsPrefix, origin)
	if err != nil {
		return err
	}
	return n.command.PasswordCodeSent(ctx, event.ResourceOwner, event.AggregateID)
}

func (n *Notification) handleEmailVerificationCode(event *models.Event) (err error) {
	emailCode := new(es_model.EmailCode)
	if err := emailCode.SetData(event); err != nil {
		return err
	}
	ctx := getSetNotifyContextData(event.InstanceID, event.ResourceOwner)
	alreadyHandled, err := n.checkIfCodeAlreadyHandledOrExpired(ctx, event, emailCode.Expiry,
		user_repo.UserV1EmailCodeAddedType, user_repo.UserV1EmailCodeSentType,
		user_repo.HumanEmailCodeAddedType, user_repo.HumanEmailCodeSentType)
	if err != nil || alreadyHandled {
		return nil
	}
	colors, err := n.getLabelPolicy(ctx)
	if err != nil {
		return err
	}

	template, err := n.getMailTemplate(ctx)
	if err != nil {
		return err
	}

	user, err := n.getUserByID(event.AggregateID, event.InstanceID)
	if err != nil {
		return err
	}

	translator, err := n.getTranslatorWithOrgTexts(ctx, user.ResourceOwner, domain.VerifyEmailMessageType)
	if err != nil {
		return err
	}

	origin, err := n.origin(ctx)
	if err != nil {
		return err
	}
	err = types.SendEmailVerificationCode(ctx, string(template.Template), translator, user, emailCode, n.getSMTPConfig, n.getFileSystemProvider, n.getLogProvider, n.userDataCrypto, colors, n.assetsPrefix, origin)
	if err != nil {
		return err
	}
	return n.command.HumanEmailVerificationCodeSent(ctx, event.ResourceOwner, event.AggregateID)
}

func (n *Notification) handlePhoneVerificationCode(event *models.Event) (err error) {
	phoneCode := new(es_model.PhoneCode)
	if err := phoneCode.SetData(event); err != nil {
		return err
	}
	ctx := getSetNotifyContextData(event.InstanceID, event.ResourceOwner)
	alreadyHandled, err := n.checkIfCodeAlreadyHandledOrExpired(ctx, event, phoneCode.Expiry,
		user_repo.UserV1PhoneCodeAddedType, user_repo.UserV1PhoneCodeSentType,
		user_repo.HumanPhoneCodeAddedType, user_repo.HumanPhoneCodeSentType)
	if err != nil || alreadyHandled {
		return nil
	}
	user, err := n.getUserByID(event.AggregateID, event.InstanceID)
	if err != nil {
		return err
	}
	translator, err := n.getTranslatorWithOrgTexts(ctx, user.ResourceOwner, domain.VerifyPhoneMessageType)
	if err != nil {
		return err
	}
	err = types.SendPhoneVerificationCode(ctx, translator, user, phoneCode, n.getTwilioConfig, n.getFileSystemProvider, n.getLogProvider, n.userDataCrypto)
	if err != nil {
		return err
	}
	return n.command.HumanPhoneVerificationCodeSent(ctx, event.ResourceOwner, event.AggregateID)
}

func (n *Notification) handleDomainClaimed(event *models.Event) (err error) {
	ctx := getSetNotifyContextData(event.InstanceID, event.ResourceOwner)
	alreadyHandled, err := n.checkIfAlreadyHandled(ctx, event.AggregateID, event.InstanceID, event.Sequence, user_repo.UserDomainClaimedType, user_repo.UserDomainClaimedSentType)
	if err != nil || alreadyHandled {
		return nil
	}
	data := make(map[string]string)
	if err := json.Unmarshal(event.Data, &data); err != nil {
		logging.Log("HANDLE-Gghq2").WithError(err).Error("could not unmarshal event data")
		return errors.ThrowInternal(err, "HANDLE-7hgj3", "could not unmarshal event")
	}
	user, err := n.getUserByID(event.AggregateID, event.InstanceID)
	if err != nil {
		return err
	}
	if user.LastEmail == "" {
		return nil
	}
	colors, err := n.getLabelPolicy(ctx)
	if err != nil {
		return err
	}

	template, err := n.getMailTemplate(ctx)
	if err != nil {
		return err
	}

	translator, err := n.getTranslatorWithOrgTexts(ctx, user.ResourceOwner, domain.DomainClaimedMessageType)
	if err != nil {
		return err
	}

	origin, err := n.origin(ctx)
	if err != nil {
		return err
	}
	err = types.SendDomainClaimed(ctx, string(template.Template), translator, user, data["userName"], n.getSMTPConfig, n.getFileSystemProvider, n.getLogProvider, colors, n.assetsPrefix, origin)
	if err != nil {
		return err
	}
	return n.command.UserDomainClaimedSent(ctx, event.ResourceOwner, event.AggregateID)
}

func (n *Notification) handlePasswordlessRegistrationLink(event *models.Event) (err error) {
	addedEvent := new(user_repo.HumanPasswordlessInitCodeRequestedEvent)
	if err := json.Unmarshal(event.Data, addedEvent); err != nil {
		return err
	}
	ctx := getSetNotifyContextData(event.InstanceID, event.ResourceOwner)
	events, err := n.getUserEvents(ctx, event.AggregateID, event.InstanceID, event.Sequence)
	if err != nil {
		return err
	}
	for _, e := range events {
		if eventstore.EventType(e.Type) == user_repo.HumanPasswordlessInitCodeSentType {
			sentEvent := new(user_repo.HumanPasswordlessInitCodeSentEvent)
			if err := json.Unmarshal(e.Data, sentEvent); err != nil {
				return err
			}
			if sentEvent.ID == addedEvent.ID {
				return nil
			}
		}
	}
	user, err := n.getUserByID(event.AggregateID, event.InstanceID)
	if err != nil {
		return err
	}
	colors, err := n.getLabelPolicy(ctx)
	if err != nil {
		return err
	}

	template, err := n.getMailTemplate(ctx)
	if err != nil {
		return err
	}

	translator, err := n.getTranslatorWithOrgTexts(ctx, user.ResourceOwner, domain.PasswordlessRegistrationMessageType)
	if err != nil {
		return err
	}

	origin, err := n.origin(ctx)
	if err != nil {
		return err
	}
	err = types.SendPasswordlessRegistrationLink(ctx, string(template.Template), translator, user, addedEvent, n.getSMTPConfig, n.getFileSystemProvider, n.getLogProvider, n.userDataCrypto, colors, n.assetsPrefix, origin)
	if err != nil {
		return err
	}
	return n.command.HumanPasswordlessInitCodeSent(ctx, event.AggregateID, event.ResourceOwner, addedEvent.ID)
}

func (n *Notification) checkIfCodeAlreadyHandledOrExpired(ctx context.Context, event *models.Event, expiry time.Duration, eventTypes ...eventstore.EventType) (bool, error) {
	if event.CreationDate.Add(expiry).Before(time.Now().UTC()) {
		return true, nil
	}
	return n.checkIfAlreadyHandled(ctx, event.AggregateID, event.InstanceID, event.Sequence, eventTypes...)
}

func (n *Notification) checkIfAlreadyHandled(ctx context.Context, userID, instanceID string, sequence uint64, eventTypes ...eventstore.EventType) (bool, error) {
	events, err := n.getUserEvents(ctx, userID, instanceID, sequence)
	if err != nil {
		return false, err
	}
	for _, event := range events {
		for _, eventType := range eventTypes {
			if eventstore.EventType(event.Type) == eventType {
				return true, nil
			}
		}
	}
	return false, nil
}

func (n *Notification) getUserEvents(ctx context.Context, userID, instanceID string, sequence uint64) ([]*models.Event, error) {
	query, err := view.UserByIDQuery(userID, instanceID, sequence)
	if err != nil {
		return nil, err
	}

	return n.es.FilterEvents(ctx, query)
}

func (n *Notification) OnError(event *models.Event, err error) error {
	logging.WithFields("id", event.AggregateID, "sequence", event.Sequence).WithError(err).Warn("something went wrong in notification handler")
	return spooler.HandleError(event, err, n.view.GetLatestNotificationFailedEvent, n.view.ProcessedNotificationFailedEvent, n.view.ProcessedNotificationSequence, n.errorCountUntilSkip)
}

func (n *Notification) OnSuccess() error {
	return spooler.HandleSuccess(n.view.UpdateNotificationSpoolerRunTimestamp)
}

func getSetNotifyContextData(instanceID, orgID string) context.Context {
	ctx := authz.WithInstanceID(context.Background(), instanceID)
	return authz.SetCtxData(ctx, authz.CtxData{UserID: NotifyUserID, OrgID: orgID})
}

// Read organization specific colors
func (n *Notification) getLabelPolicy(ctx context.Context) (*query.LabelPolicy, error) {
	return n.queries.ActiveLabelPolicyByOrg(ctx, authz.GetCtxData(ctx).OrgID)
}

// Read organization specific template
func (n *Notification) getMailTemplate(ctx context.Context) (*query.MailTemplate, error) {
	return n.queries.MailTemplateByOrg(ctx, authz.GetCtxData(ctx).OrgID)
}

// Read iam smtp config
func (n *Notification) getSMTPConfig(ctx context.Context) (*smtp.EmailConfig, error) {
	config, err := n.queries.SMTPConfigByAggregateID(ctx, authz.GetInstance(ctx).InstanceID())
	if err != nil {
		return nil, err
	}
	password, err := crypto.Decrypt(config.Password, n.smtpPasswordCrypto)
	if err != nil {
		return nil, err
	}
	return &smtp.EmailConfig{
		From:     config.SenderAddress,
		FromName: config.SenderName,
		Tls:      config.TLS,
		SMTP: smtp.SMTP{
			Host:     config.Host,
			User:     config.User,
			Password: string(password),
		},
	}, nil
}

// Read iam twilio config
func (n *Notification) getTwilioConfig(ctx context.Context) (*twilio.TwilioConfig, error) {
	config, err := n.queries.SMSProviderConfigByID(ctx, authz.GetInstance(ctx).InstanceID())
	if err != nil {
		return nil, err
	}
	if config.TwilioConfig == nil {
		return nil, errors.ThrowNotFound(nil, "HANDLER-8nfow", "Errors.SMS.Twilio.NotFound")
	}
	token, err := crypto.Decrypt(config.TwilioConfig.Token, n.smtpPasswordCrypto)
	if err != nil {
		return nil, err
	}
	return &twilio.TwilioConfig{
		SID:          config.TwilioConfig.SID,
		Token:        string(token),
		SenderNumber: config.TwilioConfig.SenderNumber,
	}, nil
}

// Read iam filesystem provider config
func (n *Notification) getFileSystemProvider(ctx context.Context) (*fs.FSConfig, error) {
	config, err := n.queries.NotificationProviderByIDAndType(ctx, authz.GetInstance(ctx).InstanceID(), domain.NotificationProviderTypeFile)
	if err != nil {
		return nil, err
	}
	return &fs.FSConfig{
		Compact: config.Compact,
		Path:    n.fileSystemPath,
	}, nil
}

// Read iam log provider config
func (n *Notification) getLogProvider(ctx context.Context) (*log.LogConfig, error) {
	config, err := n.queries.NotificationProviderByIDAndType(ctx, authz.GetInstance(ctx).InstanceID(), domain.NotificationProviderTypeLog)
	if err != nil {
		return nil, err
	}
	return &log.LogConfig{
		Compact: config.Compact,
	}, nil
}

func (n *Notification) getTranslatorWithOrgTexts(ctx context.Context, orgID, textType string) (*i18n.Translator, error) {
	translator, err := i18n.NewTranslator(n.statikDir, n.queries.GetDefaultLanguage(ctx), "")
	if err != nil {
		return nil, err
	}

	allCustomTexts, err := n.queries.CustomTextListByTemplate(ctx, authz.GetInstance(ctx).InstanceID(), textType)
	if err != nil {
		return translator, nil
	}
	customTexts, err := n.queries.CustomTextListByTemplate(ctx, orgID, textType)
	if err != nil {
		return translator, nil
	}
	allCustomTexts.CustomTexts = append(allCustomTexts.CustomTexts, customTexts.CustomTexts...)

	for _, text := range allCustomTexts.CustomTexts {
		msg := i18n.Message{
			ID:   text.Template + "." + text.Key,
			Text: text.Text,
		}
		translator.AddMessages(text.Language, msg)
	}
	return translator, nil
}

func (n *Notification) getUserByID(userID, instanceID string) (*model.NotifyUser, error) {
	return n.view.NotifyUserByID(userID, instanceID)
}

func (n *Notification) origin(ctx context.Context) (string, error) {
	primary, err := query.NewInstanceDomainPrimarySearchQuery(true)
	domains, err := n.queries.SearchInstanceDomains(ctx, &query.InstanceDomainSearchQueries{
		Queries: []query.SearchQuery{primary},
	})
	if err != nil {
		return "", err
	}
	if len(domains.Domains) < 1 {
		return "", errors.ThrowInternal(nil, "NOTIF-Ef3r1", "Errors.Notification.NoDomain")
	}
	return http_utils.BuildHTTP(domains.Domains[0].Domain, n.externalPort, n.externalSecure), nil
}
