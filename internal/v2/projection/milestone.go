package projection

import (
	"context"
	"time"

	"github.com/zitadel/zitadel/internal/v2/eventstore"
)

type milestone struct {
	reachedAt time.Time

	position    float64
	eventFilter []eventstore.EventFilterOpt
}

type instanceCreateMilestone struct {
	milestone
}

func (p *instanceCreateMilestone) Filter(ctx context.Context) *eventstore.Filter {
	return eventstore.NewFilter(
		ctx,
	)
}

/*
filter:
	instance: string
	owner: string
	pagination
	aggregates: []aggregateFilter

eventFilter:
	type textFilter[string]
	revision numberFilter[uint16]
	createdAt: timeFilter
	sequence: numberFilter[uint32]
	creator: textFilter[string] => not system
	pagination

aggregateFilter:
	type string
	id ListFilter
	events []eventFilter
	pagination

pagination:
	limit uint32
	offset uint32
	positionFilter

positionFilter:
	position float64
	in_tx_order uint32
*/

/*
milestonesFilter:
  instance: asdf
  aggregates:
    - type: instance
      events:
  	    - type: instance.added
		  pagination:
		    limit: 1
  	    - type: instance.removed
		  pagination:
		    limit: 1
  	    - type: instance.domain.primary.set
  	      creator: NotInList("", "system")
		  pagination:
		    limit: 1
    - type: project
      id: NotInFilter: (zitadel)
      events:
  	    - type: project.added
		  pagination:
		    limit: 1
  	    - type: project.application.added
		  pagination:
		    limit: 1
    - type: user
      events:
	    - user.token.added # query as long as e.applicationId is not empty and not console
		  pagination:
		    position:
			  position >= {{current position}}
			  in_tx_order > {{current in_tx_order}}
	- type: instance
	  pagination:
	    limit: 1
	  events:
	    - type: instance.idp.config.added
	    - type: instance.idp.oauth.added
	    - type: instance.idp.oidc.added
	    - type: instance.idp.jwt.added
	    - type: instance.idp.azure.added
	    - type: instance.idp.github.added
	    - type: instance.idp.github.enterprise.added
	    - type: instance.idp.gitlab.added
	    - type: instance.idp.gitlab.selfhosted.added
	    - type: instance.idp.google.added
	    - type: instance.idp.ldap.added
	    - type: instance.idp.config.apple.added
	    - type: instance.idp.saml.added
	- type: org
	  pagination:
	    limit: 1
	  events:
	    - type: org.idp.config.added
	    - type: org.idp.oauth.added
	    - type: org.idp.oidc.added
	    - type: org.idp.jwt.added
	    - type: org.idp.azure.added
	    - type: org.idp.github.added
	    - type: org.idp.github.enterprise.added
	    - type: org.idp.gitlab.added
	    - type: org.idp.gitlab.selfhosted.added
	    - type: org.idp.google.added
	    - type: org.idp.ldap.added
	    - type: org.idp.config.apple.added
	    - type: org.idp.saml.added
	- type: instance
	  pagination:
	    limit: 1
	  events:
	    - type: instance.login.policy.idp.added
	- type: org
	  pagination:
	    limit: 1
	  events:
	    - type: org.login.policy.idp.added
	- type: instance
	  events:
	  - type: instance.label.policy.added
	    pagination:
	      limit: 1
	  - type: instance.label.policy.activated
	    pagination:
	      limit: 1
	- type: org
	  events:
	  - type: org.label.policy.added
	    pagination:
	      limit: 1
	  - type: org.label.policy.activated
	    pagination:
	      limit: 1
	- type: instance
	  pagination:
	    limit: 1
	  events:
	    - type: instance.smtp.config.added
*/
