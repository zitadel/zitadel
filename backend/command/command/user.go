package command

import "github.com/zitadel/zitadel/backend/command/receiver"

type ChangeUsername struct {
	*receiver.User

	Username string
}

func (c *ChangeUsername) Execute() error {
	c.User.Username = c.Username
	return nil
}

func (c *ChangeUsername) Name() string {
	return "ChangeUsername"
}

type SetEmail struct {
	*receiver.User
	*receiver.Email
}

func (s *SetEmail) Execute() error {
	s.User.Email = s.Email
	return nil
}

func (s *SetEmail) Name() string {
	return "SetEmail"
}

type SetPhone struct {
	*receiver.User
	*receiver.Phone
}

func (s *SetPhone) Execute() error {
	s.User.Phone = s.Phone
	return nil
}

func (s *SetPhone) Name() string {
	return "SetPhone"
}
