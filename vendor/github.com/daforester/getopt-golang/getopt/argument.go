package getopt

import (
	"github.com/daforester/getopt-golang/getopt/errors"
	"fmt"
)

type Argument struct {
	className string

	adefault string
	validation validationFunc
	name string
	multiple bool
	value []string
	option *Option
	validationMessage string
}

func NewArgument(defaultValue *string, callable *validationFunc, nameArg *string) *Argument {
	a := new(Argument)

	if defaultValue != nil {
		a.SetDefaultValue(*defaultValue)
	}

	if callable != nil {
		a.SetValidation(*callable, nil)
	}

	if nameArg != nil {
		a.name = *nameArg
	} else {
		a.name = "arg"
	}

	return a
}

func (a *Argument) SetDefaultValue(value string) *Argument {
	a.adefault = value
	return a
}

func (a *Argument) SetValidation(callable validationFunc, message *string) *Argument {
	a.validation = callable

	if message != nil {
		a.validationMessage = *message
	}

	return a
}

func (a *Argument) SetName(name string) *Argument {
	a.name = name
	return a
}

func (a *Argument) GetValidationMessage(value ...string) string {
	if len(a.validationMessage) > 0 {
		fmt.Println("Fuck")
		return fmt.Sprintf(a.validationMessage, a.Describe(), value)
	} else {
		return fmt.Sprintf(Translate("value-invalid"), a.Describe(), value)
	}
}

func (a *Argument) IsMultiple() bool {
	return a.multiple
}

func(a *Argument) Multiple(m bool) *Argument {
	a.multiple = m
	return a
}

func (a *Argument) SetOption(o *Option) *Argument {
	a.option = o
	return a
}

func (a *Argument) SetValue(value ...string) (*Argument, error) {
	if a.validation != nil && !a.Validates(value...) {
		return nil, errors.New(a.GetValidationMessage(value...))
	}

	if a.IsMultiple() {
		if a.value == nil || len(a.value) == 0 {
			a.value = value
		} else {
			a.value = append(a.value, value...)
		}
	} else {
		a.value = value
	}

	return a, nil
}

func (a *Argument) GetValue() []string {
	/*
	if (a.value == nil || len(a.value) == 0) && a.IsMultiple() {
		return []string{}
	}
 	*/

	return a.value
}

func (a *Argument) Validates(args ...string) bool {
	return a.validation(args...)
}

func (a *Argument) HasValidation() bool {
	return a.validation != nil
}

func (a *Argument) HasDefaultValue() bool {
	return a.adefault != ""
}

func (a *Argument) GetDefaultValue() []string {
	if a.HasDefaultValue() {
		return []string{a.adefault}
	}

	return []string{}
}

func (a *Argument) GetName() string {
	return a.name
}

func (a *Argument) Describe() string {
	if a.option != nil {
		return a.option.Describe()
	}

	return fmt.Sprintf("%s '%s'", Translate(TRANSLATION_KEY), a.GetName())
}
