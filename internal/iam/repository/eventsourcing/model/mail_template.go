package model

import (
	b64 "encoding/base64"
	"encoding/json"

	"github.com/caos/zitadel/internal/errors"
	es_models "github.com/caos/zitadel/internal/eventstore/v1/models"
	iam_model "github.com/caos/zitadel/internal/iam/model"
)

type MailTemplate struct {
	es_models.ObjectRoot
	State    int32 `json:"-"`
	Template []byte
}

func MailTemplateToModel(template *MailTemplate) *iam_model.MailTemplate {
	return &iam_model.MailTemplate{
		ObjectRoot: template.ObjectRoot,
		State:      iam_model.PolicyState(template.State),
		Template:   template.Template,
	}
}

func (p *MailTemplate) Changes(changed *MailTemplate) map[string]interface{} {
	changes := make(map[string]interface{}, 1)
	if b64.StdEncoding.EncodeToString(changed.Template) != b64.StdEncoding.EncodeToString(p.Template) {
		changes["template"] = b64.StdEncoding.EncodeToString(changed.Template)
	}

	return changes
}

func (p *MailTemplate) SetDataLabel(event *es_models.Event) error {
	err := json.Unmarshal(event.Data, p)
	if err != nil {
		return errors.ThrowInternal(err, "MODEL-ikjhf", "unable to unmarshal data")
	}
	return nil
}
