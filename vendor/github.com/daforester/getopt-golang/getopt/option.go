package getopt

import (
	"github.com/daforester/getopt-golang/getopt/errors"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
)

type Option struct {
	className string

	short rune
	long string
	mode string
	description string
	argument *Argument
}

func NewOption(short rune, long string, mode string, description ...string) (*Option, error) {
	o := new(Option)
	if short == '\x00' && (long == "" || len(long) == 0) {
		return nil, errors.New("The short and long name may not both be empty")
	}

	_, err := o.SetShort(short)
	if err != nil {
		return nil, err
	}

	if long != "" {
		_, err = o.SetLong(long)
		if err != nil {
			return nil, err
		}
	}

	if mode == "" {
		_, err = o.SetMode(NO_ARGUMENT)
		if err != nil {
			return nil, err
		}
	} else {
		_, err = o.SetMode(mode)
		if err != nil {
			return nil, err
		}
	}

	o.argument = NewArgument(nil, nil, nil)
	o.argument.Multiple(o.mode == MULTIPLE_ARGUMENT)
	o.argument.SetOption(o)

	if description != nil && len(description) > 0 {
		o.SetDescription(description[0])
	}

	return o, nil
}

func (o *Option) SetDescription(description string) *Option {
	o.description = description
	return o
}

func (o *Option) GetDescription() string {
	return o.description
}

func (o *Option) SetDefaultValue(value string) *Option {
	o.argument.SetDefaultValue(value)
	return o
}

func (o *Option) SetValidation(validation validationFunc, message *string) *Option {
	o.argument.SetValidation(validation, message)
	return o
}

func (o *Option) SetArgumentName(name string) *Option {
	o.argument.SetName(name)
	return o
}

func (o *Option) SetArgument(arg *Argument) (*Option, error) {
	if o.mode == NO_ARGUMENT {
		return nil, errors.New("Option should not have any argument")
	}
	if o.argument == nil {
		o.argument = new(Argument)
	}

	*o.argument = *arg
	o.argument.Multiple(o.mode == MULTIPLE_ARGUMENT)
	o.argument.SetOption(o)

	return o, nil
}

func (o *Option) SetShort(short rune) (*Option, error) {

	if short != '\x00' {
		match, _ := regexp.MatchString("^[a-zA-Z0-9?!§$%#]$", string(short))

		if !match {
			return nil, errors.New(fmt.Sprintf("Short option must be null or one of [a-zA-Z0-9?!§$%%#], found '%s'", string(short)))
		} else {
			o.short = short
		}
	} else {
		o.short = short
	}

	return o, nil
}

func (o *Option) GetShort() rune {
	return o.short
}

func (o *Option) GetName() string {
	if len(o.GetLong()) > 0 {
		return o.GetLong()
	} else {
		return string(o.GetShort())
	}
}

func (o *Option) Short() rune {
	return o.short
}

func (o *Option) SetLong(long string) (*Option, error) {
	if long != "" {
		match, _ := regexp.MatchString("^[a-zA-Z0-9][a-zA-Z0-9_-]{1,}$", string(long))

		if !match {
			return nil, errors.New(fmt.Sprintf("Long option must be null or alphanumeric string, found '%s'", long))
		} else {
			o.long = long
		}
	} else {
		o.long = long
	}

	return o, nil
}

func (o *Option) GetLong() string {
	return o.long
}

func (o *Option) Long() string {
	return o.long
}

func (o *Option) SetMode(mode string) (*Option, error) {
	if mode != NO_ARGUMENT && mode != OPTIONAL_ARGUMENT && mode != REQUIRED_ARGUMENT && mode != MULTIPLE_ARGUMENT {
		return nil, errors.New(fmt.Sprintf("Option mode must be one of %s, %s, %s and %s", NO_ARGUMENT, OPTIONAL_ARGUMENT, REQUIRED_ARGUMENT, MULTIPLE_ARGUMENT))
	}

	o.mode = mode

	return o, nil
}

func (o *Option) GetMode() string {
	return o.mode
}

func (o *Option) Mode() string {
	return o.mode
}

func (o *Option) GetArgument() *Argument {
	return o.argument
}

func (o *Option) SetValue(value ...string) (*Option, error) {
	if value == nil || len(value) == 0 {

		if o.mode == REQUIRED_ARGUMENT || o.mode == MULTIPLE_ARGUMENT {
			return nil, errors.New(fmt.Sprintf(Translate("option-argument-missing"), o.GetName()))
		}

		if o.argument.GetValue() != nil && len(o.argument.GetValue()) > 0 {
			if len(o.argument.GetValue()) > 0 {
				i, err := strconv.Atoi(o.argument.GetValue()[0])
				if err != nil {
					value = append(value, strconv.Itoa(i + 1))
					_, err := o.argument.SetValue(value...)
					if err != nil {
						return nil, err
					}
				} else {
					return nil, errors.New(fmt.Sprintf("Unable to convert string to number: %s", o.argument.GetValue()[0]))
				}
			} else {
				_, err := o.argument.SetValue("1")
				if err != nil {
					return nil, err
				}
			}
		} else {
			_, err := o.argument.SetValue("1")
			if err != nil {
				return nil, err
			}
		}
	} else {
		_, err := o.argument.SetValue(value...)
		if err != nil {
			return nil, err
		}
	}

	return o, nil
}

func (o *Option) GetValue() []string {
	value := o.argument.GetValue()

	if len(value) == 0 {
		return o.argument.GetDefaultValue()
	}

	return value
}

func (o *Option) Value() []string {
	return o.GetValue()
}

func (o *Option) String() string {
	value := o.GetValue()
	if len(value) == 0 {
		return ""
	}

	if len(value) == 1 {
		return value[0]
	}

	return strings.Join(value, ",")
}

func (o *Option) Describe() string {
	return fmt.Sprintf("%s '%s'", Translate("option"), o.GetName())
}

func FileIsReadable(args ...string) bool {
	if len(args) == 0 {
		return false
	}

	for _, filename := range args {
		_, err := os.Stat(filename)
		if err != nil {
			return false
		}

		_, err = os.Open(filename)

		if err != nil {
			return false
		}
	}

	return true
}