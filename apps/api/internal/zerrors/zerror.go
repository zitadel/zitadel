package zerrors

import (
	"errors"
	"fmt"
	"reflect"
)

var _ Error = (*ZitadelError)(nil)

type ZitadelError struct {
	Parent  error
	Message string
	ID      string
}

func ThrowError(parent error, id, message string) error {
	return CreateZitadelError(parent, id, message)
}

func CreateZitadelError(parent error, id, message string) *ZitadelError {
	return &ZitadelError{
		Parent:  parent,
		ID:      id,
		Message: message,
	}
}

func (err *ZitadelError) Error() string {
	if err.Parent != nil {
		return fmt.Sprintf("ID=%s Message=%s Parent=(%v)", err.ID, err.Message, err.Parent)
	}
	return fmt.Sprintf("ID=%s Message=%s", err.ID, err.Message)
}

func (err *ZitadelError) Unwrap() error {
	return err.GetParent()
}

func (err *ZitadelError) GetParent() error {
	return err.Parent
}

func (err *ZitadelError) GetMessage() string {
	return err.Message
}

func (err *ZitadelError) SetMessage(msg string) {
	err.Message = msg
}

func (err *ZitadelError) GetID() string {
	return err.ID
}

func (err *ZitadelError) Is(target error) bool {
	t, ok := target.(*ZitadelError)
	if !ok {
		return false
	}
	if t.ID != "" && t.ID != err.ID {
		return false
	}
	if t.Message != "" && t.Message != err.Message {
		return false
	}
	if t.Parent != nil && !errors.Is(err.Parent, t.Parent) {
		return false
	}

	return true
}

func (err *ZitadelError) As(target interface{}) bool {
	_, ok := target.(**ZitadelError)
	if !ok {
		return false
	}
	reflect.Indirect(reflect.ValueOf(target)).Set(reflect.ValueOf(err))
	return true
}

func IsZitadelError(err error) bool {
	zitadelErr := new(ZitadelError)
	return errors.As(err, &zitadelErr)
}
