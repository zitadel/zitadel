package command

import (
	"context"

	"github.com/zitadel/zitadel/backend/command/receiver"
)

type createInstance struct {
	receiver receiver.InstanceManipulator
	*receiver.Instance
}

func CreateInstance(receiver receiver.InstanceManipulator, instance *receiver.Instance) *createInstance {
	return &createInstance{
		Instance: instance,
		receiver: receiver,
	}
}

func (c *createInstance) Execute(ctx context.Context) error {
	c.State = receiver.InstanceStateActive
	return c.receiver.Create(ctx, c.Instance)
}

func (c *createInstance) Name() string {
	return "CreateInstance"
}

type deleteInstance struct {
	receiver receiver.InstanceManipulator
	*receiver.Instance
}

func DeleteInstance(receiver receiver.InstanceManipulator, instance *receiver.Instance) *deleteInstance {
	return &deleteInstance{
		Instance: instance,
		receiver: receiver,
	}
}

func (d *deleteInstance) Execute(ctx context.Context) error {
	return d.receiver.Delete(ctx, d.Instance)
}

func (c *deleteInstance) Name() string {
	return "DeleteInstance"
}

type updateInstance struct {
	receiver receiver.InstanceManipulator

	*receiver.Instance

	name string
}

func UpdateInstance(receiver receiver.InstanceManipulator, instance *receiver.Instance, name string) *updateInstance {
	return &updateInstance{
		Instance: instance,
		receiver: receiver,
		name:     name,
	}
}

func (u *updateInstance) Execute(ctx context.Context) error {
	u.Instance.Name = u.name
	// return u.receiver.Update(ctx, u.Instance)
	return nil
}

func (c *updateInstance) Name() string {
	return "UpdateInstance"
}

type addDomain struct {
	receiver receiver.InstanceManipulator

	*receiver.Instance
	*receiver.Domain
}

func AddDomain(receiver receiver.InstanceManipulator, instance *receiver.Instance, domain *receiver.Domain) *addDomain {
	return &addDomain{
		Instance: instance,
		Domain:   domain,
		receiver: receiver,
	}
}

func (a *addDomain) Execute(ctx context.Context) error {
	return a.receiver.AddDomain(ctx, a.Instance, a.Domain)
}

func (c *addDomain) Name() string {
	return "AddDomain"
}
