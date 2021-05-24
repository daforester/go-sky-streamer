package getopt

import (
	"github.com/daforester/getopt-golang/getopt/errors"
	"fmt"
	"regexp"
)

var defaultOptionParserMode = NO_ARGUMENT

type OptionParser struct {}

func NewOptionParser(mode string) *OptionParser {
	o := new (OptionParser)
	defaultOptionParserMode = mode
	return o
}

func (o *OptionParser) ParseString(s string) ([]*Option, error) {
	if len(s) == 0 {
		return nil, errors.New("Option string must not be empty")
	}

	options := make([]*Option, 0)
	eol := len(s) - 1
	nextCanBeColon := false

	for i := 0; i <= eol; i++ {
		ch := rune(s[i])
		match, _ := regexp.MatchString("^[A-Za-z0-9]$", string(ch))
		if !match {
			colon := ""
			if nextCanBeColon {
				colon = " or ':'"
			}

			return nil, errors.New(fmt.Sprintf("Option string is not well formed: " + "expected a letter%s, found '%s' at position %d", colon, string(ch), i + 1))
		}
		if i == eol || s[i + 1] != ':' {
			argType := NO_ARGUMENT
			option, err := NewOption(ch, "", argType)
			if err != nil {
				return nil, err
			}

			options = append(options, option)
			nextCanBeColon = true
		} else if (i < eol - 1 && s[i + 2] == ':') {
			argType := OPTIONAL_ARGUMENT
			option, err := NewOption(ch, "", argType)
			if err != nil {
				return nil, err
			}
			options = append(options, option)
			i += 2
			nextCanBeColon = false
		} else {
			argType := REQUIRED_ARGUMENT
			option, err := NewOption(ch, "", argType)
			if err != nil {
				return nil, err
			}
			options = append(options, option)
			i++
			nextCanBeColon = true
		}
	}

	return options, nil
}

func (o *OptionParser) ParseArray(p []string) (*Option, error) {
	if len(p) == 0 {
		return nil, errors.New("Invalid option array (at least a name has to be given)")
	}

	rowSize := len(p)
	if rowSize < 3 {
		p = o.CompleteOptionArray(p)
	}

	option, err := NewOption(rune(p[0][0]), p[1], p[2])

	if err != nil {
		return nil, err
	}

	if rowSize >= 4 {
		option.SetDescription(p[3])
	}

	if rowSize >= 5 && p[2] != NO_ARGUMENT {
		a := NewArgument(&p[4], nil, nil)
		_, err := option.SetArgument(a)

		if err != nil {
			return nil, err
		}
	}

	return option, nil
}

func (o *OptionParser) CompleteOptionArray(row []string) []string {
	short := ""

	if len(row[0]) == 1 {
		short = string(row[0][0])
	}

	long := ""

	if len(short) == 0 {
		long = row[0]
	} else if (len(row) > 1 && row[1][0] != ':') {
		long = row[1]
	}

	mode := defaultOptionParserMode

	if len(row) == 2 && row[1][0] == ':' {
		mode = row[1]
	}

	return []string{short, long, mode}

}

