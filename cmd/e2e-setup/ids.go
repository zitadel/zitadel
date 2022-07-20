package main

import (
	"database/sql"
	"errors"
	"fmt"
	"strings"
)

func ids(cfg *E2EConfig, dbClient *sql.DB) (string, string, error) {
	zitadelProjectResourceID := strings.TrimPrefix(cfg.ZitadelProjectResourceID, "bignumber-")
	instanceID := strings.TrimPrefix(cfg.InstanceID, "bignumber-")

	if zitadelProjectResourceID != "" && instanceID != "" {
		return zitadelProjectResourceID, instanceID, nil
	}

	zitadelProjectResourceID, err := querySingleString(dbClient, `select aggregate_id from eventstore.events where event_type = 'project.added' and event_data = '{\"name\": \"ZITADEL\"}'`)
	if err != nil {
		return "", "", err
	}

	instanceID, err = querySingleString(dbClient, `select aggregate_id from eventstore.events where event_type = 'instance.added' and event_data = '{\"name\": \"Localhost\"}'`)
	return instanceID, zitadelProjectResourceID, err
}

func querySingleString(dbClient *sql.DB, query string) (_ string, err error) {
	defer func() {
		if err != nil {
			err = fmt.Errorf("getting single string failed for query %s: %w", query, err)
		}
	}()

	rows, err := dbClient.Query(query)
	if err != nil {
		return "", err
	}

	var read bool
	id := new(string)
	for rows.Next() {
		if read {
			return "", errors.New("read more than one row")
		}
		read = true
		if err := rows.Scan(id); err != nil {
			return "", err
		}
	}
	if !read {
		return "", errors.New("no result")
	}

	if *id == "" {
		return "", errors.New("could not parse result")
	}
	return *id, nil
}
