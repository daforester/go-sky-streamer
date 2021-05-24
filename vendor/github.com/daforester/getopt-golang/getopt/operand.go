package getopt

import "strings"

const TRANSLATION_KEY = "operand"
const OPTIONAL = 0
const REQUIRED = 1
const MULTIPLE = 2

type Operand struct {
	Argument
	required bool
	multiple bool

	description string
}

func NewOperand(name string, mode int) *Operand {
	o := new(Operand)

	o.required = (mode & REQUIRED > 0)
	o.multiple = (mode & MULTIPLE > 0)

	return o
}

func (o *Operand) IsRequired() bool {
	return o.required
}

func (o *Operand) Required(required bool) *Operand {
	o.required = required
	return o
}

func (o *Operand) SetValue(value ...string) (*Operand, error) {
	_, err := o.Argument.SetValue(value...)
	if err != nil {
		return nil, err
	}

	return o, nil
}

func (o *Operand) GetValue() []string {
	value := o.Argument.GetValue()

	if value == nil || len(value) == 0 {
		return o.Argument.GetDefaultValue()
	}

	return value
}

func (o *Operand) GetDescription() string {
	return o.description
}

func (o *Operand) SetDescription(description string) *Operand {
	o.description = description
	return o
}

func (o *Operand) Value() []string {
	return o.GetValue()
}

func (o *Operand) String() string {
	value := o.GetValue()
	if len(value) == 0 {
		return ""
	}
	if len(value) == 1 {
		return value[1]
	}
	return strings.Join(value, ",")
}