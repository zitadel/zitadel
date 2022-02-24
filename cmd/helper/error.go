package helper

import (
	"errors"
	"fmt"
	"reflect"
)

type commandError struct {
	err       string
	parent    error
	userError bool
}

func (c *commandError) Error() string {
	message := c.err
	if c.parent != nil {
		message += "parrent: " + c.parent.Error()
	}
	return message
}

func (c *commandError) isUserError() bool {
	return c.userError
}

func (c *commandError) WithParent(err error) *commandError {
	c.parent = err
	return c
}

func (c *commandError) Is(target error) bool {
	t, ok := target.(*commandError)
	if !ok {
		return false
	}
	return (c.userError == t.userError) &&
		(c.err == t.err || t.err == "") &&
		(t.parent == nil || errors.Is(t.parent, c.parent))
}

func (c *commandError) As(target interface{}) bool {
	_, ok := target.(**commandError)
	if !ok {
		return false
	}
	reflect.Indirect(reflect.ValueOf(target)).Set(reflect.ValueOf(c))
	return true
}

func NewUserError(error ...string) *commandError {
	return &commandError{err: fmt.Sprintln(error), userError: true}
}

func NewUserErrorf(format string, a ...interface{}) *commandError {
	return &commandError{err: fmt.Sprintf(format, a...), userError: true}
}

func NewSystemError(error ...string) commandError {
	return commandError{err: fmt.Sprintln(error), userError: false}
}

func NewSystemErrorF(format string, a ...interface{}) commandError {
	return commandError{err: fmt.Sprintf(format, a...), userError: false}
}

//func isUserError(err error) bool {
//	if cErr, ok := err.(commandError); ok && cErr.isUserError() {
//		return true
//	}
//
//	return userErrorRegexp.MatchString(err.Error())
//}
