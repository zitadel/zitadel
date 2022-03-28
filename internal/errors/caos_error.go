package errors

import (
	"errors"
	"fmt"
	"reflect"
)

var _ Error = (*CaosError)(nil)

type CaosError struct {
	Parent  error
	Message string
	ID      string
}

func ThrowError(parent error, id, message string) error {
	return CreateCaosError(parent, id, message)
}

func CreateCaosError(parent error, id, message string) *CaosError {
	return &CaosError{
		Parent:  parent,
		ID:      id,
		Message: message,
	}
}

func (err *CaosError) Error() string {
	if err.Parent != nil {
		return fmt.Sprintf("ID=%s Message=%s Parent=(%v)", err.ID, err.Message, err.Parent)
	}
	return fmt.Sprintf("ID=%s Message=%s", err.ID, err.Message)
}

func (err *CaosError) Unwrap() error {
	return err.GetParent()
}

func (err *CaosError) GetParent() error {
	return err.Parent
}

func (err *CaosError) GetMessage() string {
	return err.Message
}

func (err *CaosError) SetMessage(msg string) {
	err.Message = msg
}

func (err *CaosError) GetID() string {
	return err.ID
}

func (err *CaosError) Is(target error) bool {
	t, ok := target.(*CaosError)
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

func (err *CaosError) As(target interface{}) bool {
	_, ok := target.(**CaosError)
	if !ok {
		return false
	}
	reflect.Indirect(reflect.ValueOf(target)).Set(reflect.ValueOf(err))
	return true
}
