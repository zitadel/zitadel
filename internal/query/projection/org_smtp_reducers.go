package projection

import (
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/eventstore/handler/v2"
	"github.com/zitadel/zitadel/internal/repository/org"
	"github.com/zitadel/zitadel/internal/zerrors"
)

func (p *orgSMTPConfigProjection) reduceOrgSMTPConfigAdded(event eventstore.Event) (*handler.Statement, error) {
	e, err := assertEvent[*org.OrgSMTPConfigAddedEvent](event)
	if err != nil {
		return nil, err
	}

	columns := []handler.Column{
		handler.NewCol(OrgSMTPConfigSMTPColumnInstanceID, e.Aggregate().InstanceID),
		handler.NewCol(OrgSMTPConfigSMTPColumnID, e.ID),
		handler.NewCol(OrgSMTPConfigSMTPColumnTLS, e.TLS),
		handler.NewCol(OrgSMTPConfigSMTPColumnSenderAddress, e.SenderAddress),
		handler.NewCol(OrgSMTPConfigSMTPColumnSenderName, e.SenderName),
		handler.NewCol(OrgSMTPConfigSMTPColumnReplyToAddress, e.ReplyToAddress),
		handler.NewCol(OrgSMTPConfigSMTPColumnHost, e.Host),
		handler.NewCol(OrgSMTPConfigSMTPColumnUser, e.User),
	}

	if e.PlainAuth != nil {
		columns = append(columns, removeXoauth()...)
		columns = append(columns, handler.NewCol(OrgSMTPConfigSMTPColumnPlainAuthPassword, e.PlainAuth.Password))
	} else if e.Password != nil {
		columns = append(columns, removeXoauth()...)
		columns = append(columns, handler.NewCol(OrgSMTPConfigSMTPColumnPlainAuthPassword, e.Password))
	}
	if e.XOAuth2Auth != nil {
		columns = append(columns, removePlainAuth()...)
		columns = append(columns,
			handler.NewCol(OrgSMTPConfigSMTPColumnXOAuth2AuthTokenEndpoint, e.XOAuth2Auth.TokenEndpoint),
			handler.NewCol(OrgSMTPConfigSMTPColumnXOAuth2AuthScope, e.XOAuth2Auth.Scopes),
		)
		if e.XOAuth2Auth.ClientCredentials != nil {
			columns = append(columns,
				handler.NewCol(OrgSMTPConfigSMTPColumnXOAuth2AuthClientCredentialsClientId, e.XOAuth2Auth.ClientCredentials.ClientId),
				handler.NewCol(OrgSMTPConfigSMTPColumnXOAuth2AuthClientCredentialsClientSecret, e.XOAuth2Auth.ClientCredentials.ClientSecret),
			)
		}
	}

	return handler.NewMultiStatement(
		e,
		handler.AddCreateStatement(
			[]handler.Column{
				handler.NewCol(OrgSMTPConfigColumnCreationDate, e.CreationDate()),
				handler.NewCol(OrgSMTPConfigColumnChangeDate, e.CreationDate()),
				handler.NewCol(OrgSMTPConfigColumnInstanceID, e.Aggregate().InstanceID),
				handler.NewCol(OrgSMTPConfigColumnResourceOwner, e.Aggregate().ResourceOwner),
				handler.NewCol(OrgSMTPConfigColumnAggregateID, e.Aggregate().ID),
				handler.NewCol(OrgSMTPConfigColumnID, e.ID),
				handler.NewCol(OrgSMTPConfigColumnSequence, e.Sequence()),
				handler.NewCol(OrgSMTPConfigColumnState, domain.SMTPConfigStateInactive),
				handler.NewCol(OrgSMTPConfigColumnDescription, e.Description),
			},
		),
		handler.AddCreateStatement(
			columns,
			handler.WithTableSuffix(smtpConfigSMTPTableSuffix),
		),
	), nil
}

func (p *orgSMTPConfigProjection) reduceOrgSMTPConfigHTTPAdded(event eventstore.Event) (*handler.Statement, error) {
	e, err := assertEvent[*org.OrgSMTPConfigHTTPAddedEvent](event)
	if err != nil {
		return nil, err
	}

	return handler.NewMultiStatement(
		e,
		handler.AddCreateStatement(
			[]handler.Column{
				handler.NewCol(OrgSMTPConfigColumnCreationDate, e.CreationDate()),
				handler.NewCol(OrgSMTPConfigColumnChangeDate, e.CreationDate()),
				handler.NewCol(OrgSMTPConfigColumnInstanceID, e.Aggregate().InstanceID),
				handler.NewCol(OrgSMTPConfigColumnResourceOwner, e.Aggregate().ResourceOwner),
				handler.NewCol(OrgSMTPConfigColumnAggregateID, e.Aggregate().ID),
				handler.NewCol(OrgSMTPConfigColumnID, e.ID),
				handler.NewCol(OrgSMTPConfigColumnSequence, e.Sequence()),
				handler.NewCol(OrgSMTPConfigColumnState, domain.SMTPConfigStateInactive),
				handler.NewCol(OrgSMTPConfigColumnDescription, e.Description),
			},
		),
		handler.AddCreateStatement(
			[]handler.Column{
				handler.NewCol(OrgSMTPConfigHTTPColumnInstanceID, e.Aggregate().InstanceID),
				handler.NewCol(OrgSMTPConfigHTTPColumnID, e.ID),
				handler.NewCol(OrgSMTPConfigHTTPColumnEndpoint, e.Endpoint),
				handler.NewCol(OrgSMTPConfigHTTPColumnSigningKey, e.SigningKey),
			},
			handler.WithTableSuffix(smtpConfigHTTPTableSuffix),
		),
	), nil
}

func (p *orgSMTPConfigProjection) reduceOrgSMTPConfigHTTPChanged(event eventstore.Event) (*handler.Statement, error) {
	e, err := assertEvent[*org.OrgSMTPConfigHTTPChangedEvent](event)
	if err != nil {
		return nil, err
	}

	stmts := make([]func(eventstore.Event) handler.Exec, 0, 3)
	columns := []handler.Column{
		handler.NewCol(OrgSMTPConfigColumnChangeDate, e.CreationDate()),
		handler.NewCol(OrgSMTPConfigColumnSequence, e.Sequence()),
	}
	if e.Description != nil {
		columns = append(columns, handler.NewCol(OrgSMTPConfigColumnDescription, *e.Description))
	}
	stmts = append(stmts, handler.AddUpdateStatement(
		columns,
		[]handler.Condition{
			handler.NewCond(OrgSMTPConfigColumnID, e.ID),
			handler.NewCond(OrgSMTPConfigColumnInstanceID, e.Aggregate().InstanceID),
		},
	))

	httpColumns := make([]handler.Column, 0, 2)
	if e.Endpoint != nil {
		httpColumns = append(httpColumns, handler.NewCol(OrgSMTPConfigHTTPColumnEndpoint, *e.Endpoint))
	}
	if e.SigningKey != nil {
		httpColumns = append(httpColumns, handler.NewCol(OrgSMTPConfigHTTPColumnSigningKey, e.SigningKey))
	}
	if len(httpColumns) > 0 {
		stmts = append(stmts, handler.AddUpdateStatement(
			httpColumns,
			[]handler.Condition{
				handler.NewCond(OrgSMTPConfigHTTPColumnID, e.ID),
				handler.NewCond(OrgSMTPConfigHTTPColumnInstanceID, e.Aggregate().InstanceID),
			},
			handler.WithTableSuffix(smtpConfigHTTPTableSuffix),
		))
	}

	return handler.NewMultiStatement(e, stmts...), nil
}

func (p *orgSMTPConfigProjection) reduceOrgSMTPConfigChanged(event eventstore.Event) (*handler.Statement, error) {
	e, err := assertEvent[*org.OrgSMTPConfigChangedEvent](event)
	if err != nil {
		return nil, err
	}

	stmts := make([]func(eventstore.Event) handler.Exec, 0, 3)
	columns := []handler.Column{
		handler.NewCol(OrgSMTPConfigColumnChangeDate, e.CreationDate()),
		handler.NewCol(OrgSMTPConfigColumnSequence, e.Sequence()),
	}
	if e.Description != nil {
		columns = append(columns, handler.NewCol(OrgSMTPConfigColumnDescription, *e.Description))
	}
	if len(columns) > 0 {
		stmts = append(stmts, handler.AddUpdateStatement(
			columns,
			[]handler.Condition{
				handler.NewCond(OrgSMTPConfigColumnID, e.ID),
				handler.NewCond(OrgSMTPConfigColumnInstanceID, e.Aggregate().InstanceID),
			},
		))
	}

	smtpColumns := make([]handler.Column, 0, 7)
	if e.TLS != nil {
		smtpColumns = append(smtpColumns, handler.NewCol(OrgSMTPConfigSMTPColumnTLS, *e.TLS))
	}
	if e.FromAddress != nil {
		smtpColumns = append(smtpColumns, handler.NewCol(OrgSMTPConfigSMTPColumnSenderAddress, *e.FromAddress))
	}
	if e.FromName != nil {
		smtpColumns = append(smtpColumns, handler.NewCol(OrgSMTPConfigSMTPColumnSenderName, *e.FromName))
	}
	if e.ReplyToAddress != nil {
		smtpColumns = append(smtpColumns, handler.NewCol(OrgSMTPConfigSMTPColumnReplyToAddress, *e.ReplyToAddress))
	}
	if e.Host != nil {
		smtpColumns = append(smtpColumns, handler.NewCol(OrgSMTPConfigSMTPColumnHost, *e.Host))
	}
	if e.User != nil {
		smtpColumns = append(smtpColumns, handler.NewCol(OrgSMTPConfigSMTPColumnUser, *e.User))
	}
	if e.Password != nil {
		smtpColumns = append(smtpColumns, handler.NewCol(OrgSMTPConfigSMTPColumnPlainAuthPassword, *e.Password))
	}
	if !e.PlainAuth.IsEmpty() {
		smtpColumns = append(smtpColumns,
			handler.NewCol(OrgSMTPConfigSMTPColumnPlainAuthPassword, e.PlainAuth.Password),
			// Clear XOAuth2 columns when switching to plain auth
			handler.NewCol(OrgSMTPConfigSMTPColumnXOAuth2AuthTokenEndpoint, ""),
			handler.NewCol(OrgSMTPConfigSMTPColumnXOAuth2AuthScope, nil),
			handler.NewCol(OrgSMTPConfigSMTPColumnXOAuth2AuthClientCredentialsClientId, ""),
			handler.NewCol(OrgSMTPConfigSMTPColumnXOAuth2AuthClientCredentialsClientSecret, nil),
		)
	}
	if !e.XOAuth2Auth.IsEmpty() {
		smtpColumns = append(smtpColumns,
			handler.NewCol(OrgSMTPConfigSMTPColumnXOAuth2AuthTokenEndpoint, e.XOAuth2Auth.TokenEndpoint),
			handler.NewCol(OrgSMTPConfigSMTPColumnXOAuth2AuthScope, e.XOAuth2Auth.Scopes),
			// Clear plain auth columns when switching to XOAuth2
			handler.NewCol(OrgSMTPConfigSMTPColumnPlainAuthPassword, nil),
		)
		if !e.XOAuth2Auth.ClientCredentials.IsEmpty() {
			smtpColumns = append(smtpColumns,
				handler.NewCol(OrgSMTPConfigSMTPColumnXOAuth2AuthClientCredentialsClientId, e.XOAuth2Auth.ClientCredentials.ClientId),
				handler.NewCol(OrgSMTPConfigSMTPColumnXOAuth2AuthClientCredentialsClientSecret, e.XOAuth2Auth.ClientCredentials.ClientSecret),
			)
		}
	}
	if len(smtpColumns) > 0 {
		stmts = append(stmts, handler.AddUpdateStatement(
			smtpColumns,
			[]handler.Condition{
				handler.NewCond(OrgSMTPConfigSMTPColumnID, e.ID),
				handler.NewCond(OrgSMTPConfigSMTPColumnInstanceID, e.Aggregate().InstanceID),
			},
			handler.WithTableSuffix(smtpConfigSMTPTableSuffix),
		))
	}

	return handler.NewMultiStatement(e, stmts...), nil
}

func (p *orgSMTPConfigProjection) reduceOrgSMTPConfigPasswordChanged(event eventstore.Event) (*handler.Statement, error) {
	e, err := assertEvent[*org.OrgSMTPConfigPasswordChangedEvent](event)
	if err != nil {
		return nil, err
	}

	return handler.NewMultiStatement(
		e,
		handler.AddUpdateStatement(
			[]handler.Column{
				handler.NewCol(OrgSMTPConfigSMTPColumnPlainAuthPassword, e.Password),
			},
			[]handler.Condition{
				handler.NewCond(OrgSMTPConfigSMTPColumnID, e.ID),
				handler.NewCond(OrgSMTPConfigSMTPColumnInstanceID, e.Aggregate().InstanceID),
			},
			handler.WithTableSuffix(smtpConfigSMTPTableSuffix),
		),
		handler.AddUpdateStatement(
			[]handler.Column{
				handler.NewCol(OrgSMTPConfigColumnChangeDate, e.CreationDate()),
				handler.NewCol(OrgSMTPConfigColumnSequence, e.Sequence()),
			},
			[]handler.Condition{
				handler.NewCond(OrgSMTPConfigColumnID, e.ID),
				handler.NewCond(OrgSMTPConfigColumnInstanceID, e.Aggregate().InstanceID),
			},
		),
	), nil
}

func (p *orgSMTPConfigProjection) reduceOrgSMTPConfigActivated(event eventstore.Event) (*handler.Statement, error) {
	e, err := assertEvent[*org.OrgSMTPConfigActivatedEvent](event)
	if err != nil {
		return nil, err
	}

	return handler.NewMultiStatement(
		e,
		handler.AddUpdateStatement(
			[]handler.Column{
				handler.NewCol(OrgSMTPConfigColumnChangeDate, e.CreationDate()),
				handler.NewCol(OrgSMTPConfigColumnSequence, e.Sequence()),
				handler.NewCol(OrgSMTPConfigColumnState, domain.SMTPConfigStateInactive),
			},
			[]handler.Condition{
				handler.Not(handler.NewCond(OrgSMTPConfigColumnID, e.ID)),
				handler.NewCond(OrgSMTPConfigColumnState, domain.SMTPConfigStateActive),
				handler.NewCond(OrgSMTPConfigColumnInstanceID, e.Aggregate().InstanceID),
				handler.NewCond(OrgSMTPConfigColumnResourceOwner, e.Aggregate().ResourceOwner),
			},
		),
		handler.AddUpdateStatement(
			[]handler.Column{
				handler.NewCol(OrgSMTPConfigColumnChangeDate, e.CreationDate()),
				handler.NewCol(OrgSMTPConfigColumnSequence, e.Sequence()),
				handler.NewCol(OrgSMTPConfigColumnState, domain.SMTPConfigStateActive),
			},
			[]handler.Condition{
				handler.NewCond(OrgSMTPConfigColumnID, e.ID),
				handler.NewCond(OrgSMTPConfigColumnInstanceID, e.Aggregate().InstanceID),
			},
		),
	), nil
}

func (p *orgSMTPConfigProjection) reduceOrgSMTPConfigDeactivated(event eventstore.Event) (*handler.Statement, error) {
	e, err := assertEvent[*org.OrgSMTPConfigDeactivatedEvent](event)
	if err != nil {
		return nil, err
	}

	return handler.NewUpdateStatement(
		e,
		[]handler.Column{
			handler.NewCol(OrgSMTPConfigColumnChangeDate, e.CreationDate()),
			handler.NewCol(OrgSMTPConfigColumnSequence, e.Sequence()),
			handler.NewCol(OrgSMTPConfigColumnState, domain.SMTPConfigStateInactive),
		},
		[]handler.Condition{
			handler.NewCond(OrgSMTPConfigColumnID, e.ID),
			handler.NewCond(OrgSMTPConfigColumnInstanceID, e.Aggregate().InstanceID),
		},
	), nil
}

func (p *orgSMTPConfigProjection) reduceOrgSMTPConfigRemoved(event eventstore.Event) (*handler.Statement, error) {
	e, err := assertEvent[*org.OrgSMTPConfigRemovedEvent](event)
	if err != nil {
		return nil, err
	}

	return handler.NewDeleteStatement(
		e,
		[]handler.Condition{
			handler.NewCond(OrgSMTPConfigColumnID, e.ID),
			handler.NewCond(OrgSMTPConfigColumnInstanceID, e.Aggregate().InstanceID),
		},
	), nil
}

func (p *orgSMTPConfigProjection) reduceOrgSMTPOwnerRemoved(event eventstore.Event) (*handler.Statement, error) {
	e, ok := event.(*org.OrgRemovedEvent)
	if !ok {
		return nil, zerrors.ThrowInvalidArgumentf(nil, "PROJE-orgSmtpOwn", "reduce.wrong.event.type %s", org.OrgRemovedEventType)
	}

	return handler.NewDeleteStatement(
		e,
		[]handler.Condition{
			handler.NewCond(OrgSMTPConfigColumnInstanceID, e.Aggregate().InstanceID),
			handler.NewCond(OrgSMTPConfigColumnResourceOwner, e.Aggregate().ID),
		},
	), nil
}
