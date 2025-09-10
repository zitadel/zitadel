package domain

import (
	"context"

	legacy_es "github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/v2/org"
)

// AddOrgCommand adds a new organization.
// I'm unsure if we should add the Admins here or if this should be a separate command.
type AddOrgCommand struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	// Admins []*AddMemberCommand `json:"admins"`
}

// func NewAddOrgCommand(name string, admins ...*AddMemberCommand) *AddOrgCommand {
// 	return &AddOrgCommand{
// 		Name:   name,
// 		Admins: admins,
// 	}
// }

// // String implements [Commander].
// func (cmd *AddOrgCommand) String() string {
// 	return "AddOrgCommand"
// }

// // Execute implements Commander.
// func (cmd *AddOrgCommand) Execute(ctx context.Context, opts *CommandOpts) (err error) {
// 	if len(cmd.Admins) == 0 {
// 		return ErrNoAdminSpecified
// 	}
// 	if err = cmd.ensureID(); err != nil {
// 		return err
// 	}

// 	close, err := opts.EnsureTx(ctx)
// 	if err != nil {
// 		return err
// 	}
// 	defer func() { err = close(ctx, err) }()
// 	err = orgRepo(opts.DB).Create(ctx, &Org{
// 		ID:   cmd.ID,
// 		Name: cmd.Name,
// 	})
// 	if err != nil {
// 		return err
// 	}

// 	for _, admin := range cmd.Admins {
// 		admin.orgID = cmd.ID
// 		if err = opts.Invoke(ctx, admin); err != nil {
// 			return err
// 		}
// 	}

// 	orgCache.Set(ctx, &Org{
// 		ID:   cmd.ID,
// 		Name: cmd.Name,
// 	})

// 	return nil
// }

// Events implements [eventer].
func (cmd *AddOrgCommand) Events(ctx context.Context) []legacy_es.Command {
	command, err := org.NewAddedCommand(ctx, cmd.Name)
	if err != nil {
		return nil
	}
	return []legacy_es.Command{command}
}

// var (
// 	_ Commander = (*AddOrgCommand)(nil)
// 	_ eventer   = (*AddOrgCommand)(nil)
// )

// func (cmd *AddOrgCommand) ensureID() (err error) {
// 	if cmd.ID != "" {
// 		return nil
// 	}
// 	cmd.ID, err = generateID()
// 	return err
// }

// // AddMemberCommand adds a new member to an organization.
// // I'm not sure if we should make it more generic to also use it for instances.
// type AddMemberCommand struct {
// 	orgID  string
// 	UserID string   `json:"userId"`
// 	Roles  []string `json:"roles"`
// }

// func NewAddMemberCommand(userID string, roles ...string) *AddMemberCommand {
// 	return &AddMemberCommand{
// 		UserID: userID,
// 		Roles:  roles,
// 	}
// }

// // String implements [Commander].
// func (cmd *AddMemberCommand) String() string {
// 	return "AddMemberCommand"
// }

// // Execute implements Commander.
// func (a *AddMemberCommand) Execute(ctx context.Context, opts *CommandOpts) (err error) {
// 	close, err := opts.EnsureTx(ctx)
// 	if err != nil {
// 		return err
// 	}
// 	defer func() { err = close(ctx, err) }()

// 	return orgRepo(opts.DB).Member().AddMember(ctx, a.orgID, a.UserID, a.Roles)
// }

// // Events implements [eventer].
// func (a *AddMemberCommand) Events() []*eventstore.Event {
// 	return []*eventstore.Event{
// 		{
// 			AggregateType: "org",
// 			AggregateID:   a.UserID,
// 			Type:          "member.added",
// 			Payload:       a,
// 		},
// 	}
// }

// var (
// 	_ Commander = (*AddMemberCommand)(nil)
// 	_ eventer   = (*AddMemberCommand)(nil)
// )
