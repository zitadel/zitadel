package model

import (
	b64 "encoding/base64"
	"encoding/json"

	"github.com/caos/zitadel/internal/errors"
	"github.com/caos/zitadel/internal/eventstore/models"
	es_models "github.com/caos/zitadel/internal/eventstore/models"
	iam_model "github.com/caos/zitadel/internal/iam/model"
)

type MailTemplate struct {
	models.ObjectRoot
	State    int32 `json:"-"`
	Template []byte
}

func MailTemplateToModel(policy *MailTemplate) *iam_model.MailTemplate {
	return &iam_model.MailTemplate{
		ObjectRoot: policy.ObjectRoot,
		State:      iam_model.PolicyState(policy.State),
		Template:   policy.Template,
	}
}

func MailTemplateFromModel(policy *iam_model.MailTemplate) *MailTemplate {
	return &MailTemplate{
		ObjectRoot: policy.ObjectRoot,
		State:      int32(policy.State),
		Template:   policy.Template,
	}
}

func (p *MailTemplate) Changes(changed *MailTemplate) map[string]interface{} {
	changes := make(map[string]interface{}, 2)

	if b64.StdEncoding.EncodeToString(changed.Template) != b64.StdEncoding.EncodeToString(p.Template) {
		changes["mailTemplate"] = changed.Template
	}

	return changes
}

func (i *IAM) appendAddMailTemplateEvent(event *es_models.Event) error {
	i.DefaultMailTemplate = new(MailTemplate)
	err := i.DefaultMailTemplate.SetDataLabel(event)
	if err != nil {
		return err
	}
	i.DefaultMailTemplate.ObjectRoot.CreationDate = event.CreationDate
	return nil
}

func (i *IAM) appendChangeMailTemplateEvent(event *es_models.Event) error {
	return i.DefaultMailTemplate.SetDataLabel(event)
}

func (p *MailTemplate) SetDataLabel(event *es_models.Event) error {
	err := json.Unmarshal(event.Data, p)
	if err != nil {
		return errors.ThrowInternal(err, "MODEL-ikjhf", "unable to unmarshal data")
	}
	return nil
}
