package actions

import (
	"context"

	"github.com/dop251/goja"
	"github.com/google/uuid"
	"github.com/zitadel/logging"
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
	logging.OnError(o.Set("namespaceDNS", uuid.NameSpaceDNS)).Warn("unable to set namespace")
	logging.OnError(o.Set("namespaceURL", uuid.NameSpaceURL)).Warn("unable to set namespace")
	logging.OnError(o.Set("namespaceOID", uuid.NameSpaceOID)).Warn("unable to set namespace")
	logging.OnError(o.Set("namespaceX500", uuid.NameSpaceX500)).Warn("unable to set namespace")
}

func inRuntime(function func() (uuid.UUID, error), runtime *goja.Runtime) func(call goja.FunctionCall) goja.Value {
	return func(call goja.FunctionCall) goja.Value {
		if len(call.Arguments) != 0 {
			panic("invalid arg count")
		}

		uuid, err := function()
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

		var err error
		var namespace uuid.UUID
		switch n := call.Arguments[0].Export().(type) {
		case string:
			namespace, err = uuid.Parse(n)
			if err != nil {
				logging.WithError(err).Debug("namespace failed parsing as UUID")
				panic(err)
			}
		case uuid.UUID:
			namespace = n
		default:
			logging.WithError(err).Debug("invalid type for namespace")
			panic(err)
		}

		var data []byte
		switch d := call.Arguments[1].Export().(type) {
		case string:
			data = []byte(d)
		case []byte:
			data = d
		default:
			logging.WithError(err).Debug("invalid type for data")
			panic(err)
		}

		return runtime.ToValue(function(namespace, data).String())
	}
}
