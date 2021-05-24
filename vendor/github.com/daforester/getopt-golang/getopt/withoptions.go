package getopt

import (
	"github.com/daforester/getopt-golang/getopt/errors"
	"fmt"
)

type WithOptions struct {
	options []*Option
	optionMapping map[string]*Option
}

func (w *WithOptions) AddOptionString(option string) (*WithOptions, error) {
	op := NewOptionParser(NO_ARGUMENT)
	options, err := op.ParseString(option)

	if err != nil {
		return nil, err
	}

	return w.AddOptions(options...)
}

func (w *WithOptions) AddOptionArray(s []string) (*WithOptions, error) {
	op := NewOptionParser(NO_ARGUMENT)
	option, err := op.ParseArray(s)

	if err != nil {
		return nil, err
	}

	return w.AddOption(option)
}

func (w *WithOptions) AddOptions(options ...*Option) (*WithOptions, error) {
	var err error

	if len(options) > 0 {
		for _, option := range options {
			if option != nil {
				_, err := w.AddOption(option)
				if err != nil {
					return nil, err
				}
			} else {
				err = errors.New("Nil option provided to AddOptions")
			}
		}
	}

	return w, err
}

func (w *WithOptions) AddOption(option *Option) (*WithOptions, error) {
	if w.Conflicts(option) {
		return nil, errors.New(fmt.Sprintf("%s's short and long name have to be unique", option))
	}

	w.options = append(w.options, option)

	short := option.GetShort()
	long := option.GetLong()

	if short != '\x00'{
		w.optionMapping[string(short)] = option
	}
	if long != "" {
		w.optionMapping[long] = option
	}

	return w, nil
}

func (w *WithOptions) Conflicts(option *Option) bool {
	short := option.GetShort()
	long := option.GetLong()

	_, hasShortMap := w.optionMapping[string(short)]
	_, hasLongMap := w.optionMapping[long]

	return (short != '\x00' && hasShortMap) || (long != "" && hasLongMap)
}

func (w *WithOptions) GetOptions() []*Option {
	return w.options
}

func (w *WithOptions) GetOption(name string) *Option {
	option, hasOption := w.optionMapping[name]

	if hasOption {
		return option
	}

	return nil
}

func (w *WithOptions) HasOptions() bool {
	return len(w.options) > 0
}