package getopt

import (
	"github.com/daforester/getopt-golang/getopt/errors"
	"fmt"
)

type WithOperands struct {
	operands []*Operand
}

func (w *WithOperands) AddOperands(operands ...*Operand) (*WithOperands, error) {
	if operands != nil && len(operands) > 0 {
		for _, o := range operands {
			if o != nil {
				_, err := w.AddOperand(o)
				if err != nil {
					return nil, err
				}
			}
		}
	}

	return w, nil
}

func (w *WithOperands) AddOperand(operand *Operand) (*WithOperands, error) {
	if operand.IsRequired() {
		for _, p := range w.operands {
			p.Required(true)
		}
	}

	if w.HasOperands() {
		last := w.operands[len(w.operands) - 1]
		if last.IsMultiple() {
			return nil, errors.New(fmt.Sprintf("Operand %s is multiple - no more operands allowed.", last.GetName()))
		}
	}

	w.operands = append(w.operands, operand)

	return w, nil
}


func (w *WithOperands) GetOperands() []*Operand {
	return w.operands
}

func (w *WithOperands) GetOperand(index interface{}) *Operand {
	switch index.(type) {
	case string:
		name := index.(string)
		for _, o := range w.operands {
			if o.GetName() == name {
				return o
			}
		}
		return nil
	case int:
		if index.(int) >= len(w.operands) || index.(int) < 0 {
			return nil
		}
		return w.operands[index.(int)]
	}

	return nil
}

func (w *WithOperands) HasOperands() bool {
	return len(w.operands) > 0
}