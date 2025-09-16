package domain

import "errors"

var (
	ErrNoAdminSpecified = errors.New("at least one admin must be specified")
	ErrNoOrgIdSpecified = errors.New("organization id must be specified")
)
