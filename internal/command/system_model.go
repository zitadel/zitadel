package command

import (
	"strings"

	"github.com/zitadel/zitadel/internal/command/preparation"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/repository/instance"
	"github.com/zitadel/zitadel/internal/repository/project"
)

type SystemConfigWriteModel struct {
	eventstore.WriteModel

	Instances         map[string]*systemConfigChangesInstanceModel
	externalDomain    string
	externalPort      uint16
	externalSecure    bool
	newExternalDomain string
	newExternalPort   uint16
	newExternalSecure bool
}

type systemConfigChangesInstanceModel struct {
	Domains                []string
	GeneratedDomain        string
	ProjectID              string
	ConsoleAppID           string
	RedirectUris           []string
	PostLogoutRedirectUris []string
}

func NewSystemConfigWriteModel(externalDomain, newExternalDomain string, externalPort, newExternalPort uint16, externalSecure, newExternalSecure bool) *SystemConfigWriteModel {
	return &SystemConfigWriteModel{
		WriteModel:        eventstore.WriteModel{},
		Instances:         make(map[string]*systemConfigChangesInstanceModel),
		externalDomain:    externalDomain,
		externalPort:      externalPort,
		externalSecure:    externalSecure,
		newExternalDomain: newExternalDomain,
		newExternalPort:   newExternalPort,
		newExternalSecure: newExternalSecure,
	}
}

func (wm *SystemConfigWriteModel) Reduce() error {
	for _, event := range wm.Events {
		switch e := event.(type) {
		case *instance.InstanceAddedEvent:
			wm.Instances[e.Aggregate().InstanceID] = &systemConfigChangesInstanceModel{}
		case *instance.InstanceRemovedEvent:
			delete(wm.Instances, e.Aggregate().InstanceID)
		case *instance.DomainAddedEvent:
			if !e.Generated && e.Domain != wm.externalDomain && e.Domain != wm.newExternalDomain {
				continue
			}
			if e.Generated && !strings.HasSuffix(e.Domain, wm.externalDomain) && !strings.HasSuffix(e.Domain, wm.newExternalDomain) {
				continue
			}
			wm.Instances[e.Aggregate().InstanceID].Domains = append(wm.Instances[e.Aggregate().InstanceID].Domains, e.Domain)
		case *instance.DomainRemovedEvent:
			instance, ok := wm.Instances[e.Aggregate().InstanceID]
			if !ok {
				continue
			}

			for i, domain := range instance.Domains {
				if domain == e.Domain {
					instance.Domains[i] = instance.Domains[len(instance.Domains)-1]
					instance.Domains[len(instance.Domains)-1] = ""
					wm.Instances[e.Aggregate().InstanceID].Domains = instance.Domains[:len(instance.Domains)-1]
					break
				}
			}
		case *instance.ProjectSetEvent:
			wm.Instances[e.Aggregate().InstanceID].ProjectID = e.ProjectID
		case *instance.ConsoleSetEvent:
			wm.Instances[e.Aggregate().InstanceID].ConsoleAppID = e.AppID
		case *project.OIDCConfigAddedEvent:
			if wm.Instances[e.Aggregate().InstanceID].ConsoleAppID != e.AppID {
				continue
			}
			wm.Instances[e.Aggregate().InstanceID].RedirectUris = e.RedirectUris
			wm.Instances[e.Aggregate().InstanceID].PostLogoutRedirectUris = e.PostLogoutRedirectUris
		case *project.OIDCConfigChangedEvent:
			if wm.Instances[e.Aggregate().InstanceID].ConsoleAppID != e.AppID {
				continue
			}
			if e.RedirectUris != nil {
				wm.Instances[e.Aggregate().InstanceID].RedirectUris = *e.RedirectUris
			}
			if e.PostLogoutRedirectUris != nil {
				wm.Instances[e.Aggregate().InstanceID].PostLogoutRedirectUris = *e.PostLogoutRedirectUris
			}
		}
	}
	return nil
}

func (wm *SystemConfigWriteModel) Query() *eventstore.SearchQueryBuilder {
	return eventstore.NewSearchQueryBuilder(eventstore.ColumnsEvent).
		AddQuery().
		AggregateTypes(instance.AggregateType).
		EventTypes(
			instance.InstanceAddedEventType,
			instance.InstanceRemovedEventType,
			instance.InstanceDomainAddedEventType,
			instance.InstanceDomainRemovedEventType,
			instance.ProjectSetEventType,
			instance.ConsoleSetEventType,
		).
		Or().
		AggregateTypes(project.AggregateType).
		EventTypes(
			project.OIDCConfigAddedType,
			project.OIDCConfigChangedType,
		).
		Builder()
}

type SystemConfigChangesValidation struct {
	ProjectID    string
	ConsoleAppID string
	Validations  []preparation.Validation
	InstanceID   string
}

func (wm *SystemConfigWriteModel) NewChangedEvents(commands *Commands) map[string]*SystemConfigChangesValidation {
	var newCustomDomainExists, isInstanceOfCustomDomain bool
	var instanceOfCustomDomain string
	cmds := make(map[string]*SystemConfigChangesValidation)
	for i, inst := range wm.Instances {
		cmds[i] = &SystemConfigChangesValidation{
			InstanceID:   i,
			ProjectID:    inst.ProjectID,
			ConsoleAppID: inst.ConsoleAppID,
		}
		//check each instance separately for changes (using the generated domain) and check if there's an existing custom domain
		newCustomDomainExists, isInstanceOfCustomDomain = wm.changeConfig(cmds[i], inst, commands)
		if isInstanceOfCustomDomain || newCustomDomainExists {
			instanceOfCustomDomain = i
		}
	}
	//handle the custom domain at last
	if newCustomDomainExists {
		//if the domain itself already exists, then only check if the uris of the console app exist as well
		wm.changeURIs(cmds[instanceOfCustomDomain], wm.Instances[instanceOfCustomDomain], commands, wm.newExternalDomain)
		return cmds
	}
	//otherwise the add instance domain will take care of the uris
	cmds[instanceOfCustomDomain].Validations = append(cmds[instanceOfCustomDomain].Validations, commands.addInstanceDomain(instance.NewAggregate(instanceOfCustomDomain), wm.newExternalDomain, false))
	return cmds
}

func (wm *SystemConfigWriteModel) changeConfig(validation *SystemConfigChangesValidation, inst *systemConfigChangesInstanceModel, commands *Commands) (newCustomDomainExists, isInstanceOfCustomDomain bool) {
	var newGeneratedDomain string
	var newGeneratedDomainExists bool
	for _, domain := range inst.Domains {
		if domain == wm.newExternalDomain {
			newCustomDomainExists = true
			continue
		}
		if domain != wm.newExternalDomain && strings.HasSuffix(domain, wm.newExternalDomain) {
			newGeneratedDomainExists = true
		}
		if !newCustomDomainExists && domain == wm.externalDomain {
			isInstanceOfCustomDomain = true
		}
		if domain != wm.externalDomain && strings.HasSuffix(domain, wm.externalDomain) {
			newGeneratedDomain = strings.TrimSuffix(domain, wm.externalDomain) + wm.newExternalDomain
		}
	}
	if newGeneratedDomainExists {
		//if the domain itself already exists, then only check if the uris of the console app exist as well
		wm.changeURIs(validation, inst, commands, newGeneratedDomain)
		return newCustomDomainExists, isInstanceOfCustomDomain
	}
	//otherwise the add instance domain will take care of the uris
	validation.Validations = append(validation.Validations, commands.addInstanceDomain(instance.NewAggregate(validation.InstanceID), newGeneratedDomain, true))
	return newCustomDomainExists, isInstanceOfCustomDomain
}

func (wm *SystemConfigWriteModel) changeURIs(validation *SystemConfigChangesValidation, inst *systemConfigChangesInstanceModel, commands *Commands, domain string) {
	validation.Validations = append(validation.Validations, commands.prepareUpdateConsoleRedirectURIs(domain))
}
