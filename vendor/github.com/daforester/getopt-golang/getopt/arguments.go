package getopt

import (
	"strings"
)

type Arguments struct {
	arguments []string
}

func NewArguments(arguments []string) *Arguments {
	a := new(Arguments)
	a.arguments = arguments

	return a
}

func (a *Arguments) Process(getopt *GetOpt, setOption setOptionFunc, setCommand setCommandFunc, addOperand addOperandFunc) (bool, error) {
	for len(a.arguments) > 0 {
		arg := a.arguments[0]
		a.arguments = a.arguments[1:]

		if a.isMeta(arg) {
			for _, argument := range a.arguments {
				err := addOperand(argument)
				if err != nil {
					return false, err
				}
			}
		}

		if a.isValue(arg) {
			operands := getopt.GetOperands()
			command := getopt.GetCommand(arg)
			if len(operands) == 0 && command != nil {
				err := setCommand(command)
				if err != nil {
					return false, err
				}
			} else {
				err := addOperand(arg)
				if err != nil {
					return false, err
				}
			}
		}

		if a.isLongOption(arg) {
			err := setOption(a.longName(arg), func() string {
				return a.Value(arg, "")
			})

			if err != nil {
				return false, err
			}
			continue
		}

		for _, name := range a.shortNames(arg) {
			requestedValue := false
			err := setOption(string(name), func() string {
				requestedValue = true
				return a.Value(arg, string(name))
			})

			if err != nil {
				return false, err
			}

			if requestedValue {
				break;
			}
		}
	}

	return true, nil
}

func (a *Arguments) isOption(arg string) bool {
	return !a.isValue(arg) && !a.isMeta(arg)
}

func (a *Arguments) isValue(arg string) bool {
	return len(arg) == 0 || arg == "=" || arg[0] != '-'
}

func (a *Arguments) isMeta(arg string) bool {
	return len(arg) > 0 && arg == "--"
}

func (a *Arguments) isLongOption(arg string) bool {
	return a.isOption(arg) && len(arg) > 1 && arg[1] == '-'
}

func (a *Arguments) longName(arg string) string {
	name := arg[2:]
	p := strings.IndexRune(name, '=')

	if p > 0 {
		return name[0:p]
	}

	return name
}

func (a *Arguments) shortNames(arg string) []rune {
	if !a.isOption(arg) || a.isLongOption(arg) {
		return []rune("")
	}

	return []rune(arg[1:])
}

func (a *Arguments) Value(arg string, name string) string {
	var p int

	if a.isLongOption(arg) {
		p = strings.IndexRune(arg, '=')
	} else {
		p = strings.Index(arg, name)
	}

	if a.isLongOption(arg) && p > 0 || !a.isLongOption(arg) && p < len(arg) - 1 {
		return arg[p+1:]
	}

	if len(a.arguments) > 0 && a.isValue(a.arguments[0]) {
		v := a.arguments[0]
		a.arguments = a.arguments[1:]

		return v
	}

	return ""
}

func (a *Arguments) FromString(argsString string) *Arguments {
	argv := make([]string, 0)
	argsString = strings.Trim(argsString, " \n\r\t")
	argc := 0

	if len(argsString) == 0 {
		return new(Arguments)
	}

	state := 'n'

	for i := 0; i < len(argsString); i ++ {
		char := argsString[i]

		switch state {
		case 'n':
			if char == '\'' {
				state = 's'
			} else if char == '"' {
				state = 'd'
			} else if char == '\n' || char == '\t' || char == ' ' {
				argc ++
				argv[argc] = ""
			} else {
				argv[argc] += string(char)
			}
		case 's':
			if char == '\'' {
				state = 'n'
			} else if char == '\\' {
				i ++
				argv[argc] += string(argsString[i])
			} else {
				argv[argc] += string(char)
			}
		case 'd':
			if char == '"' {
				state = 'n'
			} else if char == '\\' {
				i ++
				argv[argc] += string(argsString[i])
			} else {
				argv[argc] += string(char)
			}
		}
	}

	return NewArguments(argv)
}