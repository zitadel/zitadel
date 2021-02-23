package model

import (
	"encoding/json"

	"github.com/caos/logging"

	es_models "github.com/caos/zitadel/internal/eventstore/v1/models"
	"github.com/caos/zitadel/internal/project/model"
)

type Application struct {
	es_models.ObjectRoot
	AppID      string      `json:"appId"`
	State      int32       `json:"-"`
	Name       string      `json:"name,omitempty"`
	Type       int32       `json:"appType,omitempty"`
	OIDCConfig *OIDCConfig `json:"-"`
	APIConfig  *APIConfig  `json:"-"`
}

type ApplicationID struct {
	es_models.ObjectRoot
	AppID string `json:"appId"`
}

func GetApplication(apps []*Application, id string) (int, *Application) {
	for i, a := range apps {
		if a.AppID == id {
			return i, a
		}
	}
	return -1, nil
}

func (a *Application) Changes(changed *Application) map[string]interface{} {
	changes := make(map[string]interface{}, 1)
	changes["appId"] = a.AppID
	if changed.Name != "" && a.Name != changed.Name {
		changes["name"] = changed.Name
	}
	return changes
}

func AppsToModel(apps []*Application) []*model.Application {
	convertedApps := make([]*model.Application, len(apps))
	for i, a := range apps {
		convertedApps[i] = AppToModel(a)
	}
	return convertedApps
}

func AppsFromModel(apps []*model.Application) []*Application {
	convertedApps := make([]*Application, len(apps))
	for i, a := range apps {
		convertedApps[i] = AppFromModel(a)
	}
	return convertedApps
}

func AppFromModel(app *model.Application) *Application {
	converted := &Application{
		ObjectRoot: app.ObjectRoot,
		AppID:      app.AppID,
		Name:       app.Name,
		State:      int32(app.State),
		Type:       int32(app.Type),
	}
	if app.OIDCConfig != nil {
		converted.OIDCConfig = OIDCConfigFromModel(app.OIDCConfig)
	}
	if app.APIConfig != nil {
		converted.APIConfig = APIConfigFromModel(app.APIConfig)
	}
	return converted
}

func AppToModel(app *Application) *model.Application {
	converted := &model.Application{
		ObjectRoot: app.ObjectRoot,
		AppID:      app.AppID,
		Name:       app.Name,
		State:      model.AppState(app.State),
		Type:       model.AppType(app.Type),
	}
	if app.OIDCConfig != nil {
		converted.OIDCConfig = OIDCConfigToModel(app.OIDCConfig)
	}
	if app.APIConfig != nil {
		converted.APIConfig = APIConfigToModel(app.APIConfig)
	}
	return converted
}

func (p *Project) appendAddAppEvent(event *es_models.Event) error {
	app := new(Application)
	err := app.setData(event)
	if err != nil {
		return err
	}
	app.ObjectRoot.CreationDate = event.CreationDate
	p.Applications = append(p.Applications, app)
	return nil
}

func (p *Project) appendChangeAppEvent(event *es_models.Event) error {
	app := new(Application)
	err := app.setData(event)
	if err != nil {
		return err
	}
	if i, a := GetApplication(p.Applications, app.AppID); a != nil {
		return p.Applications[i].setData(event)
	}
	return nil
}

func (p *Project) appendRemoveAppEvent(event *es_models.Event) error {
	app := new(Application)
	err := app.setData(event)
	if err != nil {
		return err
	}
	if i, a := GetApplication(p.Applications, app.AppID); a != nil {
		p.Applications[i] = p.Applications[len(p.Applications)-1]
		p.Applications[len(p.Applications)-1] = nil
		p.Applications = p.Applications[:len(p.Applications)-1]
	}
	return nil
}

func (p *Project) appendAppStateEvent(event *es_models.Event, state model.AppState) error {
	app := new(Application)
	err := app.setData(event)
	if err != nil {
		return err
	}

	if i, a := GetApplication(p.Applications, app.AppID); a != nil {
		a.State = int32(state)
		p.Applications[i] = a
	}
	return nil
}

func (a *Application) setData(event *es_models.Event) error {
	a.ObjectRoot.AppendEvent(event)
	if err := json.Unmarshal(event.Data, a); err != nil {
		logging.Log("EVEN-8die3").WithError(err).Error("could not unmarshal event data")
		return err
	}
	return nil
}
