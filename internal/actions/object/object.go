package object

import "github.com/dop251/goja"

func objectFromFirstArgument(call goja.FunctionCall, runtime *goja.Runtime) *goja.Object {
	if len(call.Arguments) != 1 {
		panic("exactly one argument expected")
	}
	object := call.Arguments[0].ToObject(runtime)
	if object == nil {
		panic("unable to unmarshal arg")
	}
	return object
}
