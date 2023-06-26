package actions

import (
	"context"

	"github.com/dop251/goja"
	"github.com/zitadel/logging"

	"github.com/google/uuid"
)

func WithUUID(ctx context.Context) Option {
	return func(c *runConfig) {
		c.modules["zitadel/uuid"] = func(runtime *goja.Runtime, module *goja.Object) {
			requireUUID(ctx, runtime, module)
		}
	}
}

func requireUUID(_ context.Context, runtime *goja.Runtime, module *goja.Object) {
	o := module.Get("exports").(*goja.Object)
	logging.OnError(o.Set("v1", inRuntime(uuid.NewUUID, runtime))).Warn("unable to set module")
	logging.OnError(o.Set("v3", inRuntimeHash(uuid.NewMD5, runtime))).Warn("unable to set module")
	logging.OnError(o.Set("v4", inRuntime(uuid.NewRandom, runtime))).Warn("unable to set module")
	logging.OnError(o.Set("v5", inRuntimeHash(uuid.NewSHA1, runtime))).Warn("unable to set module")
}

func inRuntime(function func() (uuid.UUID, error), runtime *goja.Runtime) func(call goja.FunctionCall) goja.Value {
  return func(call goja.FunctionCall) goja.Value {
    if len(call.Arguments) != 0 {
      panic("invalid arg count")
    }

    uuid, err := function();
    if err != nil {
      logging.WithError(err)
      panic(err)
    }

    return runtime.ToValue(uuid.String())
  }
}

func inRuntimeHash(function func(uuid.UUID, []byte) uuid.UUID, runtime *goja.Runtime) func(call goja.FunctionCall) goja.Value {
	return func(call goja.FunctionCall) goja.Value {
		if len(call.Arguments) != 2 {
			logging.WithFields("count", len(call.Arguments)).Debug("other than 2 args provided")
			panic("invalid arg count")
		}

		space, err := uuid.Parse(call.Arguments[0].Export().(string))
		if err != nil {
			logging.WithError(err).Debug("space failed parsing as UUID")
			panic(err)
		}

		return runtime.ToValue(function(space, call.Arguments[1].Export().([]byte)).String())
	}
}
