package main

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/zitadel/logging"
)

var idRegexp = regexp.MustCompile("[0-9]{16}")

func ids(ctx context.Context, cfg *E2EConfig, dbClient *sql.DB) (string, string, error) {
	zitadelProjectResourceID := strings.TrimPrefix(cfg.ZitadelProjectResourceID, "bignumber-")
	instanceID := strings.TrimPrefix(cfg.InstanceID, "bignumber-")

	if idRegexp.MatchString(zitadelProjectResourceID) && idRegexp.MatchString(instanceID) {
		return zitadelProjectResourceID, instanceID, nil
	}

	projCtx, projCancel := context.WithTimeout(ctx, time.Minute)
	defer projCancel()
	zitadelProjectResourceID, err := querySingleString(projCtx, dbClient, `select aggregate_id from eventstore.events where event_type = 'project.added' and event_data = '{"name": "ZITADEL"}'`)
	if err != nil {
		return "", "", err
	}

	instCtx, instCancel := context.WithTimeout(ctx, time.Minute)
	defer instCancel()
	instanceID, err = querySingleString(instCtx, dbClient, `select aggregate_id from eventstore.events where event_type = 'instance.added' and event_data = '{"name": "Localhost"}'`)
	return instanceID, zitadelProjectResourceID, err
}

func querySingleString(ctx context.Context, dbClient *sql.DB, query string) (_ string, err error) {
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
		select {
		case <-ctx.Done():
			return "", ctx.Err()
		default:
			logging.Warningf("no results for query yet. retrying in a second. query: %s", query)
			time.Sleep(time.Second)
			return querySingleString(ctx, dbClient, query)
		}
	}

	if *id == "" {
		return "", errors.New("could not parse result")
	}

	return *id, nil
}
