package domain

import (
	"context"
)

type AddOrgCommand struct {
	ID     string            `json:"id"`
	Name   string            `json:"name"`
	Admins []AddAdminCommand `json:"admins"`
}

func NewAddOrgCommand(name string, admins ...AddAdminCommand) *AddOrgCommand {
	return &AddOrgCommand{
		Name:   name,
		Admins: admins,
	}
}

// Execute implements Commander.
func (cmd *AddOrgCommand) Execute(ctx context.Context, opts *CommandOpts) (err error) {
	if len(cmd.Admins) == 0 {
		return ErrNoAdminSpecified
	}
	if err = cmd.ensureID(); err != nil {
		return err
	}

	close, err := opts.EnsureTx(ctx)
	if err != nil {
		return err
	}
	defer func() { err = close(ctx, err) }()
	err = orgRepo(opts.DB).Create(ctx, &Org{
		ID:   cmd.ID,
		Name: cmd.Name,
	})
	if err != nil {
		return err
	}

	return nil
}

var (
	_ Commander = (*AddOrgCommand)(nil)
)

func (cmd *AddOrgCommand) ensureID() (err error) {
	if cmd.ID != "" {
		return nil
	}
	cmd.ID, err = generateID()
	return err
}

type AddAdminCommand struct {
	UserID string   `json:"userId"`
	Roles  []string `json:"roles"`
}

// Execute implements Commander.
func (a *AddAdminCommand) Execute(ctx context.Context, opts *CommandOpts) (err error) {
	close, err := opts.EnsureTx(ctx)
	if err != nil {
		return err
	}
	defer func() { err = close(ctx, err) }()
	return nil
}

var (
	_ Commander = (*AddAdminCommand)(nil)
)
