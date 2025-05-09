package bla4

import (
	"fmt"
	"log/slog"
	"reflect"

	"github.com/manifoldco/promptui"
)

type Iterator interface {
	// NextField returns the name of the next field to configure.
	// If it returns an empty string, the configuration is complete.
	NextField() string
	// FinishAllowed returns true if the configuration is complete and the user can skip the remaining fields.
	FinishAllowed() bool
}

func fieldFunctionByIterator(l *slog.Logger, f *field) Configure {
	if !f.value.Type().Implements(reflect.TypeFor[Iterator]()) {
		return nil
	}

	return func() (value any, err error) {
		iter := f.value.Interface().(Iterator)
		sub := f.sub()
		for {
			fieldName := iter.NextField()
			if fieldName == "" {
				break
			}
			if iter.FinishAllowed() {
				prompt := promptui.Prompt{
					Label:     fmt.Sprintf("Finish configuration of %s", f.info.Name),
					IsConfirm: true,
				}
				_, err := prompt.Run()
				if err == nil {
					break
				}
			}
			info, ok := f.value.Type().FieldByName(fieldName)
			if !ok {
				return nil, fmt.Errorf("field %s not found", fieldName)
			}
			tag, err := newTag(info, f.value.FieldByName(fieldName))
			if err != nil {
				return nil, err
			}
			if tag.skip {
				continue
			}

			res, err := fieldFunction(
				l.With(slog.String("field", fieldName)),
				&field{
					info:  info,
					viper: sub, //TODO: sub of sub
					tag:   tag,
				},
			)()
			if err != nil {
				return nil, err
			}
			reflect.Indirect(f.value).FieldByName(fieldName).Set(reflect.ValueOf(res))
		}
		f.viper.Set(f.tag.fieldName, sub.AllSettings())
		return value, err
	}
}
