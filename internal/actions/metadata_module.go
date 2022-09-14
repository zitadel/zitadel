package actions

import (
	"context"
	"encoding/json"
	"time"

	"github.com/dop251/goja"
	"github.com/zitadel/logging"

	"github.com/zitadel/zitadel/internal/command"
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/query"
)

func WithUserMetadata(ctx context.Context, q *query.Queries, cmd *command.Commands, userID, resourceOwner string) Option {
	return func(c *runConfig) {
		c.modules["zitadel/metadata/user"] = func(runtime *goja.Runtime, module *goja.Object) {
			config := &userMetadataRuntime{
				runtime:    runtime,
				module:     module,
				query:      q,
				command:    cmd,
				maxEndTime: c.end,
			}
			userMetadataModule(ctx, userID, resourceOwner, config)
		}
	}
}

type userMetadataRuntime struct {
	runtime *goja.Runtime
	module  *goja.Object

	query   *query.Queries
	command *command.Commands

	maxEndTime time.Time
}

func userMetadataModule(ctx context.Context, userID, resourceOwner string, c *userMetadataRuntime) {
	o := c.module.Get("exports").(*goja.Object)
	logging.OnError(o.Set("get", c.getFn(ctx, userID, resourceOwner))).Warn("unable to set module")
	logging.OnError(o.Set("set", c.setFn(ctx, userID, resourceOwner))).Warn("unable to set module")
}

func (c *userMetadataRuntime) userMetadataListFromQuery(metadata *query.UserMetadataList) *userMetadataList {
	result := &userMetadataList{
		Count:     metadata.Count,
		Sequence:  metadata.Sequence,
		Timestamp: metadata.Timestamp,
		Metadata:  make([]*userMetadata, len(metadata.Metadata)),
	}

	for i, md := range metadata.Metadata {
		var value interface{}
		err := json.Unmarshal(md.Value, &value)
		if err != nil {
			logging.WithError(err).Debug("unable to unmarshal into map")
			panic(err)
		}
		result.Metadata[i] = &userMetadata{
			CreationDate:  md.CreationDate,
			ChangeDate:    md.ChangeDate,
			ResourceOwner: md.ResourceOwner,
			Sequence:      md.Sequence,
			Key:           md.Key,
			Value:         c.runtime.ToValue(value),
		}
	}

	return result
}

type userMetadataList struct {
	Count     uint64          `json:"count"`
	Sequence  uint64          `json:"sequence"`
	Timestamp time.Time       `json:"timestamp"`
	Metadata  []*userMetadata `json:"metadata"`
}

type userMetadata struct {
	CreationDate  time.Time  `json:"creationDate"`
	ChangeDate    time.Time  `json:"changeDate"`
	ResourceOwner string     `json:"resourceOwner"`
	Sequence      uint64     `json:"sequence"`
	Key           string     `json:"key"`
	Value         goja.Value `json:"value"`
}

func (c *userMetadataRuntime) getFn(ctx context.Context, userID, resourceOwner string) func(call goja.FunctionCall) goja.Value {
	return func(call goja.FunctionCall) goja.Value {
		resourceOwnerQuery, err := query.NewUserMetadataResourceOwnerSearchQuery(resourceOwner)
		if err != nil {
			logging.WithError(err).Debug("unable to create search query")
			panic(err)
		}
		metadata, err := c.query.SearchUserMetadata(
			ctx,
			true,
			userID,
			&query.UserMetadataSearchQueries{Queries: []query.SearchQuery{resourceOwnerQuery}},
		)
		if err != nil {
			logging.WithError(err).Info("unable to get md in action")
			panic(err)
		}
		return c.runtime.ToValue(c.userMetadataListFromQuery(metadata))
	}
}

func (c *userMetadataRuntime) setFn(ctx context.Context, userID, resourceOwner string) func(call goja.FunctionCall) goja.Value {
	return func(call goja.FunctionCall) goja.Value {
		if len(call.Arguments) != 2 {
			panic("exactly 2 (key, value) arguments expected")
		}
		key := call.Arguments[0].Export().(string)
		val := call.Arguments[1].Export()

		value, err := json.Marshal(val)
		if err != nil {
			logging.WithError(err).Debug("unable to marshal")
			panic(err)
		}

		metadata := &domain.Metadata{
			Key:   key,
			Value: value,
		}
		if _, err = c.command.SetUserMetadata(ctx, metadata, userID, resourceOwner); err != nil {
			logging.WithError(err).Info("unable to set md in action")
			panic(err)
		}
		return nil
	}
}
