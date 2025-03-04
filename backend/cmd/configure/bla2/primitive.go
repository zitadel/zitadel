package bla2

import (
	"fmt"
	"reflect"

	"github.com/manifoldco/promptui"
)

type primitive struct {
	typ reflect.Type

	tag fieldTag
}

func (p *primitive) defaultValue() any {
	if p.tag.currentValue != nil {
		return p.tag.currentValue
	}
	return reflect.Zero(p.typ).Interface()
}

func (p *primitive) label() string {
	if p.tag.description == "" {
		return p.tag.fieldName
	}
	return fmt.Sprintf("%s (%s)", p.tag.fieldName, p.tag.description)
}

func (p *primitive) toPrompt() prompt {
	return &linePrompt{
		Prompt: promptui.Prompt{
			Label:     p.label(),
			Default:   fmt.Sprintf("%v", p.defaultValue()),
			Validate:  p.validateInput,
			IsConfirm: p.typ.Kind() == reflect.Bool,
		},
	}
}

func promptFromPrimitive(p *primitive) prompt {
	return p.toPrompt()
}

func (p *primitive) validateInput(s string) error {
	_, err := mapValue(p.typ, s)
	return err
}
