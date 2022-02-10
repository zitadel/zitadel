package main

import (
	"context"
	"fmt"
	"github.com/caos/logging"
	"time"
)

func awaitConsistency(ctx context.Context, cfg E2EConfig, expectUsers []user) (err error) {

	retry := make(chan struct{})
	go func() {
		// trigger first check
		retry <- struct{}{}
	}()
	for {
		select {
		case <-retry:
			err = checkCondition(ctx, cfg, expectUsers)
			if err == nil {
				logging.Log("AWAIT-QIOOJ").Info("setup is consistent")
				return nil
			}
			logging.Log("AWAIT-VRk3Y").Info("setup is not consistent yet, retrying in a second: ", err)
			time.Sleep(time.Second)
			go func() {
				retry <- struct{}{}
			}()
		case <-ctx.Done():
			return fmt.Errorf("setup failed to come to a consistent state: %s: %w", ctx.Err(), err)
		}
	}
}

func checkCondition(ctx context.Context, cfg E2EConfig, expectUsers []user) error {
	token, err := newToken(cfg)
	if err != nil {
		return err
	}

	foundUsers, err := listUsers(ctx, cfg.APIURL, token)
	if err != nil {
		return err
	}

	var awaitingUsers []string
expectLoop:
	for i := range expectUsers {
		expectUser := expectUsers[i]
		for j := range foundUsers {
			foundUser := foundUsers[j]
			if expectUser.desc+"_user_name" == foundUser {
				continue expectLoop
			}
		}
		awaitingUsers = append(awaitingUsers, expectUser.desc)
	}

	if len(awaitingUsers) > 0 {
		return fmt.Errorf("users %v are not consistent yet", awaitingUsers)
	}
	return nil
}
