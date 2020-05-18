package model

import "github.com/caos/zitadel/internal/eventstore/models"

const (
	IamAggregate models.AggregateType = "iam"

	IamSetupStarted  models.EventType = "iam.setup.started"
	IamSetupDone     models.EventType = "iam.setup.done"
	GlobalOrgSet     models.EventType = "iam.global.org.set"
	IamProjectSet    models.EventType = "iam.project.iam.set"
	IamMemberAdded   models.EventType = "iam.member.added"
	IamMemberChanged models.EventType = "iam.member.changed"
	IamMemberRemoved models.EventType = "iam.member.removed"
)
