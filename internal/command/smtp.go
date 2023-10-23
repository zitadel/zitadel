package command

import (
	"context"
	"net"
	"strings"

	"github.com/zitadel/zitadel/internal/api/authz"
	"github.com/zitadel/zitadel/internal/command/preparation"
	"github.com/zitadel/zitadel/internal/crypto"
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/errors"
	caos_errs "github.com/zitadel/zitadel/internal/errors"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/notification/channels/smtp"
	"github.com/zitadel/zitadel/internal/repository/instance"
	admin_pb "github.com/zitadel/zitadel/pkg/grpc/admin"
)

func (c *Commands) AddSMTPConfig(ctx context.Context, config *smtp.Config) (*domain.ObjectDetails, error) {
	instanceAgg := instance.NewAggregate(authz.GetInstance(ctx).InstanceID())
	validation := c.prepareAddSMTPConfig(instanceAgg, config.From, config.FromName, config.ReplyToAddress, config.SMTP.Host, config.SMTP.User, []byte(config.SMTP.Password), config.Tls, config.SMTP.IsActive, config.SMTP.ProviderType)
	cmds, err := preparation.PrepareCommands(ctx, c.eventstore.Filter, validation)
	if err != nil {
		return nil, err
	}
	events, err := c.eventstore.Push(ctx, cmds...)
	if err != nil {
		return nil, err
	}
	return &domain.ObjectDetails{
		Sequence:      events[len(events)-1].Sequence(),
		EventDate:     events[len(events)-1].CreationDate(),
		ResourceOwner: events[len(events)-1].Aggregate().InstanceID,
	}, nil
}

func (c *Commands) ChangeSMTPConfig(ctx context.Context, config *smtp.Config) (*domain.ObjectDetails, error) {
	instanceAgg := instance.NewAggregate(authz.GetInstance(ctx).InstanceID())
	validation := c.prepareChangeSMTPConfig(instanceAgg, config.From, config.FromName, config.ReplyToAddress, config.SMTP.Host, config.SMTP.User, config.Tls, config.SMTP.IsActive, config.SMTP.ProviderType)
	cmds, err := preparation.PrepareCommands(ctx, c.eventstore.Filter, validation)
	if err != nil {
		return nil, err
	}
	events, err := c.eventstore.Push(ctx, cmds...)
	if err != nil {
		return nil, err
	}
	return &domain.ObjectDetails{
		Sequence:      events[len(events)-1].Sequence(),
		EventDate:     events[len(events)-1].CreationDate(),
		ResourceOwner: events[len(events)-1].Aggregate().InstanceID,
	}, nil
}

func (c *Commands) ChangeSMTPConfigPassword(ctx context.Context, password string) (*domain.ObjectDetails, error) {
	instanceAgg := instance.NewAggregate(authz.GetInstance(ctx).InstanceID())
	// TODO @n40lab empty string ID or not?
	smtpConfigWriteModel, err := getSMTPConfigWriteModel(ctx, c.eventstore.Filter, "", "")
	if err != nil {
		return nil, err
	}
	if smtpConfigWriteModel.State != domain.SMTPConfigStateActive {
		return nil, errors.ThrowNotFound(nil, "COMMAND-3n9ls", "Errors.SMTPConfig.NotFound")
	}
	var smtpPassword *crypto.CryptoValue
	if password != "" {
		smtpPassword, err = crypto.Encrypt([]byte(password), c.smtpEncryption)
		if err != nil {
			return nil, err
		}
	}
	events, err := c.eventstore.Push(ctx, instance.NewSMTPConfigPasswordChangedEvent(
		ctx,
		&instanceAgg.Aggregate,
		smtpPassword))
	if err != nil {
		return nil, err
	}
	return &domain.ObjectDetails{
		Sequence:      events[len(events)-1].Sequence(),
		EventDate:     events[len(events)-1].CreationDate(),
		ResourceOwner: events[len(events)-1].Aggregate().InstanceID,
	}, nil
}

func (c *Commands) DeactivateSMTPConfig(ctx context.Context, req *admin_pb.DeactivateSMTPConfigRequest) (*domain.ObjectDetails, error) {
	instanceAgg := instance.NewAggregate(authz.GetInstance(ctx).InstanceID())
	validation := c.prepareDeactivateSMTPConfig(instanceAgg, req.Id)
	cmds, err := preparation.PrepareCommands(ctx, c.eventstore.Filter, validation)
	if err != nil {
		return nil, err
	}
	events, err := c.eventstore.Push(ctx, cmds...)
	if err != nil {
		return nil, err
	}
	return &domain.ObjectDetails{
		Sequence:      events[len(events)-1].Sequence(),
		EventDate:     events[len(events)-1].CreationDate(),
		ResourceOwner: events[len(events)-1].Aggregate().InstanceID,
	}, nil
}

func (c *Commands) RemoveSMTPConfig(ctx context.Context, req *admin_pb.RemoveSMTPConfigRequest) (*domain.ObjectDetails, error) {
	instanceAgg := instance.NewAggregate(authz.GetInstance(ctx).InstanceID())
	validation := c.prepareRemoveSMTPConfig(instanceAgg, req.Id)
	cmds, err := preparation.PrepareCommands(ctx, c.eventstore.Filter, validation)
	if err != nil {
		return nil, err
	}
	events, err := c.eventstore.Push(ctx, cmds...)
	if err != nil {
		return nil, err
	}
	return &domain.ObjectDetails{
		Sequence:      events[len(events)-1].Sequence(),
		EventDate:     events[len(events)-1].CreationDate(),
		ResourceOwner: events[len(events)-1].Aggregate().InstanceID,
	}, nil
}

func (c *Commands) prepareAddSMTPConfig(a *instance.Aggregate, from, name, replyTo, hostAndPort, user string, password []byte, tls, isActive bool, providerType uint32) preparation.Validation {
	return func() (preparation.CreateCommands, error) {
		if from = strings.TrimSpace(from); from == "" {
			return nil, errors.ThrowInvalidArgument(nil, "INST-mruNY", "Errors.Invalid.Argument")
		}

		replyTo = strings.TrimSpace(replyTo)
		hostAndPort = strings.TrimSpace(hostAndPort)

		if _, _, err := net.SplitHostPort(hostAndPort); err != nil {
			return nil, errors.ThrowInvalidArgument(nil, "INST-9JdRe", "Errors.Invalid.Argument")
		}
		return func(ctx context.Context, filter preparation.FilterToQueryReducer) ([]eventstore.Command, error) {
			fromSplitted := strings.Split(from, "@")
			senderDomain := fromSplitted[len(fromSplitted)-1]
			// TODO @n40lab empty string ID or not?
			writeModel, err := getSMTPConfigWriteModel(ctx, filter, "", senderDomain)
			if err != nil {
				return nil, err
			}

			err = checkSenderAddress(writeModel)
			if err != nil {
				return nil, err
			}
			var smtpPassword *crypto.CryptoValue
			if password != nil {
				smtpPassword, err = crypto.Encrypt(password, c.smtpEncryption)
				if err != nil {
					return nil, err
				}
			}

			configID, err := c.idGenerator.Next()
			if err != nil {
				return nil, err
			}

			return []eventstore.Command{
				instance.NewSMTPConfigAddedEvent(
					ctx,
					&a.Aggregate,
					configID,
					tls,
					from,
					name,
					replyTo,
					hostAndPort,
					user,
					smtpPassword,
					isActive,
					providerType,
				),
			}, nil
		}, nil
	}
}

func (c *Commands) prepareChangeSMTPConfig(a *instance.Aggregate, from, name, replyTo, hostAndPort, user string, tls, isActive bool, providerType uint32) preparation.Validation {
	return func() (preparation.CreateCommands, error) {
		if from = strings.TrimSpace(from); from == "" {
			return nil, errors.ThrowInvalidArgument(nil, "INST-ASv2d", "Errors.Invalid.Argument")
		}

		replyTo = strings.TrimSpace(replyTo)
		hostAndPort = strings.TrimSpace(hostAndPort)
		if _, _, err := net.SplitHostPort(hostAndPort); err != nil {
			return nil, errors.ThrowInvalidArgument(nil, "INST-Kv875", "Errors.Invalid.Argument")
		}
		return func(ctx context.Context, filter preparation.FilterToQueryReducer) ([]eventstore.Command, error) {
			fromSplitted := strings.Split(from, "@")
			senderDomain := fromSplitted[len(fromSplitted)-1]
			// TODO @n40lab empty string ID or not?
			writeModel, err := getSMTPConfigWriteModel(ctx, filter, "", senderDomain)
			if err != nil {
				return nil, err
			}
			if writeModel.State != domain.SMTPConfigStateActive {
				return nil, errors.ThrowNotFound(nil, "INST-Svq1a", "Errors.SMTPConfig.NotFound")
			}
			err = checkSenderAddress(writeModel)
			if err != nil {
				return nil, err
			}
			changedEvent, hasChanged, err := writeModel.NewChangedEvent(
				ctx,
				&a.Aggregate,
				tls,
				from,
				name,
				replyTo,
				hostAndPort,
				user,
				isActive,
				providerType,
			)
			if err != nil {
				return nil, err
			}
			if !hasChanged {
				return nil, errors.ThrowPreconditionFailed(nil, "COMMAND-m0o3f", "Errors.NoChangesFound")
			}
			return []eventstore.Command{
				changedEvent,
			}, nil
		}, nil
	}
}

func (c *Commands) prepareDeactivateSMTPConfig(a *instance.Aggregate, id string) preparation.Validation {
	return func() (preparation.CreateCommands, error) {
		return func(ctx context.Context, filter preparation.FilterToQueryReducer) ([]eventstore.Command, error) {
			writeModel, err := getSMTPConfigWriteModel(ctx, filter, id, "")
			if err != nil {
				return nil, err
			}
			if !writeModel.State.Exists() {
				return nil, errors.ThrowNotFound(nil, "INST-Sfefg", "Errors.SMTPConfig.NotFound")
			}
			if writeModel.State == domain.SMTPConfigStateInactive {
				return nil, caos_errs.ThrowNotFound(nil, "COMMAND-xk8k0", "Errors.SMTPConfig.AlreadyDeactivated")
			}
			return []eventstore.Command{
				instance.NewSMTPConfigDeactivatedEvent(ctx, &a.Aggregate, id),
			}, nil
		}, nil
	}
}

func (c *Commands) prepareRemoveSMTPConfig(a *instance.Aggregate, id string) preparation.Validation {
	return func() (preparation.CreateCommands, error) {
		return func(ctx context.Context, filter preparation.FilterToQueryReducer) ([]eventstore.Command, error) {
			writeModel, err := getSMTPConfigWriteModel(ctx, filter, id, "")
			if err != nil {
				return nil, err
			}
			if !writeModel.State.Exists() {
				return nil, errors.ThrowNotFound(nil, "INST-Sfefg", "Errors.SMTPConfig.NotFound")
			}
			return []eventstore.Command{
				instance.NewSMTPConfigRemovedEvent(ctx, &a.Aggregate, id),
			}, nil
		}, nil
	}
}

func checkSenderAddress(writeModel *InstanceSMTPConfigWriteModel) error {
	if !writeModel.smtpSenderAddressMatchesInstanceDomain {
		return nil
	}
	if !writeModel.domainState.Exists() {
		return errors.ThrowInvalidArgument(nil, "INST-83nl8", "Errors.SMTPConfig.SenderAdressNotCustomDomain")
	}
	return nil
}

func getSMTPConfigWriteModel(ctx context.Context, filter preparation.FilterToQueryReducer, id, domain string) (_ *InstanceSMTPConfigWriteModel, err error) {
	writeModel := NewInstanceSMTPConfigWriteModel(authz.GetInstance(ctx).InstanceID(), id, domain)
	events, err := filter(ctx, writeModel.Query())
	if err != nil {
		return nil, err
	}
	if len(events) == 0 {
		return writeModel, nil
	}
	writeModel.AppendEvents(events...)
	err = writeModel.Reduce()
	return writeModel, err
}
