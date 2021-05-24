package getopt

import (
	"fmt"
	"golang.org/x/crypto/ssh/terminal"
	"log"
	"os"
	"strconv"
	"strings"
)

const DESCRIPTION = "description"
const MAX_WIDTH = "maxWidth"
const HIDE_OPERANDS = "hideOperands"

type HelpInterface interface {
	Render(*GetOpt, map[string]string) string
}

type Help struct {
	usageTemplate string
	optionsTemplate string
	commandsTemplate string
	settings map[string]string
	texts map[string]string
	getOpt *GetOpt
	screenWidth int
}

func NewHelp(settings map[string]string) *Help {
	h := new(Help)

	h.settings = make(map[string]string)
	h.texts = make(map[string]string)

	h.settings[MAX_WIDTH] = "120"
	h.texts["placeholder"] = "<>"
	h.texts["optional"] = "[]"
	h.texts["multiple"] = "..."
	h.texts["options-listing"] = ", "

	for setting, value := range settings {
		h.Set(setting, value)
	}

	return h
}

func (h *Help) Set(setting string, value string) *Help {
	switch setting {
	case "optionsTemplate":
		h.SetOptionsTemplate(value)
	case "commandsTemplate":
		h.SetCommandsTemplate(value)
	case "usageTemplate":
		h.SetUsageTemplate(value)
	default:
		h.settings[setting] = value
	}

	return h
}

func (h *Help) SetTexts(texts map[string]string) *Help {
	for k, v := range texts {
		h.texts[k] = v
	}

	for k, v := range h.texts {
		h.texts[k] = strings.ReplaceAll(v, "\r\n", "\n")
	}

	return h
}

func (h *Help) GetText(key string) string {
	v, isset := h.texts[key]
	if !isset {
		return Translate(key)
	}

	return v
}

func (h *Help) Render(getOpt *GetOpt, data map[string]string) string {
	helpText := ""

	h.getOpt = getOpt

	for setting, value := range data {
		h.Set(setting, value)
	}

	if h.usageTemplate != "" {
		data := make(map[string][]string)
		data["command"] = append(data["command"], getOpt.GetCommand("").GetName())
		helpText = h.renderTemplate(h.usageTemplate, data)
	} else {
		helpText = h.renderUsage()
	}

	_, isset := h.settings[HIDE_OPERANDS]

	if getOpt.HasOperands() && isset {
		helpText += h.renderOperands()
	}

	if getOpt.HasOptions() {
		if h.optionsTemplate != "" {
			helpText += h.renderTemplate(h.optionsTemplate, getOpt.GetOptions())
		} else {
			helpText += h.renderOptions()
		}
	}

	if getOpt.GetCommand("") == nil && getOpt.HasCommands() {
		if h.commandsTemplate != "" {
			data := make(map[string][]string)
			for _, command := range getOpt.GetCommands() {
				data[command.GetName()] = append(data[command.GetName()], command.GetShortDescription())
			}

			helpText += h.renderTemplate(h.commandsTemplate, data)
		} else {
			helpText += h.renderCommands()
		}
	}

	return helpText
}

func (h *Help) surround(text string, with string) string {
	return string(with[0]) + text + with[len(with)-1:]
}

func (h *Help) renderUsage() string {
	return h.GetText("usage-title") +
		h.getOpt.Get(SETTING_COMMAND_NAME) +
		" " +
		h.renderUsageCommand() +
		h.renderUsageOptions() +
		h.renderUsageOperands() +
		"\n\n" +
		h.renderDescription()
}

func (h *Help) renderOperands() string {
	data := make([][]string, 0)
	definitionWidth := 0
	hasDescriptions := false

	for _, operand := range h.getOpt.WithOperands.GetOperands() {
		definition := h.surround(operand.GetName(), h.texts["placeholder"])
		if !operand.IsRequired() {
			definition = h.surround(definition, h.texts["optional"])
		}

		if len(definition) > definitionWidth {
			definitionWidth = len(definition)
		}

		if operand.GetDescription() != "" {
			hasDescriptions = true
		}

		data = append(data, []string{operand.GetDescription()})
	}

	if hasDescriptions {
		return ""
	}

	return h.GetText("operands-title") + h.renderColumns(definitionWidth, data) + "\n"
}

func (h *Help) renderOptions() string {
	data := make([][]string, 0)
	definitionWidth := 0

	for _, option := range h.getOpt.WithOptions.GetOptions() {
		optionStrings := make([]string,0)
		if option.GetShort() != '\x00' {
			optionStrings = append(optionStrings, "-" + string(option.GetShort()))
		}
		if option.GetLong() != "" {
			optionStrings = append(optionStrings, "--" + string(option.GetLong()))
		}

		definition := strings.Join(optionStrings,h.texts["options-listing"])

		if option.GetMode() != NO_ARGUMENT {
			argument := h.surround(option.GetArgument().GetName(), h.texts["placeholder"])
			if option.GetMode() == OPTIONAL_ARGUMENT {
				argument = h.surround(argument, h.texts["optional"])
			}

			definition += " " + argument
		}

		if len(definition) > definitionWidth {
			definitionWidth = len(definition)
		}

		data = append(data, []string{definition, option.GetDescription()})
	}

	return h.GetText("options-title") + h.renderColumns(definitionWidth, data) + "\n"
}

func (h *Help) renderCommands() string {
	data := make([][]string, 0)
	nameWidth := 0

	for _, command := range h.getOpt.GetCommands() {
		if len(command.GetName()) > nameWidth {
			nameWidth = len(command.GetName())
		}

		data = append(data, []string{command.GetName(), command.GetShortDescription()})
	}

	return h.GetText("commands-title") + h.renderColumns(nameWidth, data) + "\n"
}

func (h *Help) renderUsageCommand() string {
	command := h.getOpt.GetCommand("")

	if command != nil {
		return command.GetName() + " "
	} else if h.getOpt.HasCommands() {
		return h.surround(h.GetText("usage-command"), h.texts["placeholder"]) + " "
	}

	return ""
}

func (h *Help) renderUsageOptions() string {
	if h.getOpt.HasOptions() || len(h.getOpt.Get(SETTING_STRICT_OPTIONS)) > 0 {
		return h.surround(h.GetText("usage-options"), h.texts["optional"]) + " "
	}

	return ""
}

func (h *Help) renderUsageOperands() string {
	usage := ""

	lastOperandMultiple := false

	if h.getOpt.HasOperands() {
		for _, operand := range h.getOpt.WithOperands.GetOperands() {
			name := h.surround(operand.GetName(), h.texts["placeholder"])
			if !operand.IsRequired() {
				name = h.surround(name, h.texts["optional"])
			}
			usage += name + " "
			if operand.IsMultiple() {
				usage += h.surround(
					h.surround(
						operand.GetName(),
						h.texts["placeholder"]) +
						h.texts["multiple"],
					h.texts["optional"])
				lastOperandMultiple = true
			}
		}
	}

	if !lastOperandMultiple && h.getOpt.Get(SETTING_STRICT_OPERANDS) == ""{
		usage += h.surround(h.GetText("usage-operands"), h.texts["optional"])
	}

	return usage
}

func (h *Help) renderDescription() string {
	command := h.getOpt.GetCommand("")
	if command != nil {
		return command.GetDescription() + "\n\n"
	} else if h.settings[DESCRIPTION] != "" {
		return h.settings[DESCRIPTION] + "\n\n"
	}

	return ""
}

func (h *Help) getScreenWidth() int {
	if h.screenWidth == 0 {
		columns := os.Getenv("COLUMNS")
		if len(columns) == 0 {
			w, _, err := terminal.GetSize(int(os.Stdout.Fd()))
			if err != nil {
				log.Println(err)
				return 0
			}
			h.screenWidth = w
		} else {
			w, err := strconv.Atoi(columns)

			if err != nil {
				log.Println(err)
				return 0
			}

			h.screenWidth = w
		}

		maxWidth, err := strconv.Atoi(h.settings[MAX_WIDTH])

		if err != nil {
			log.Println(err)
			return 0
		}

		if h.screenWidth > maxWidth {
			h.screenWidth = maxWidth
		}
	}

	return h.screenWidth
}

func (h *Help) renderColumns(columnWidth int, data [][]string) string {
	text := ""
	screenWidth := h.getScreenWidth()

	for _, dataRow := range data {
		row := fmt.Sprintf("  % -" + strconv.Itoa(columnWidth) + "s %s", dataRow[0], dataRow[1])
		for len(row) > screenWidth {

			p := strings.IndexRune(row[:screenWidth], ' ')

			if p < columnWidth + 4 {
				p = strings.IndexRune(row[:screenWidth], '-')
				if p < columnWidth + 4 {
					p = screenWidth - 1
				}
			}
			if p < 0 {
				p = 0
			}
			c := row[p:p+1]
			pc := ""
			if c != " " {
				pc = c
			}
			text += row[:p] + pc + "\n"
			row = fmt.Sprintf("  %s %s", strings.Repeat(" ", columnWidth), row[p+1:])
		}
		text += row + "\n"
	}

	return text
}

func (h *Help) renderTemplate(template string, data map[string][]string) string {
	return ""
}

func (h *Help) GetUsageTemplate() string {
	return h.usageTemplate
}

func (h *Help) SetUsageTemplate(usageTemplate string) *Help {
	h.usageTemplate = usageTemplate
	return h
}

func (h *Help) GetOptionsTemplate() string {
	return h.optionsTemplate
}

func (h *Help) SetOptionsTemplate(optionsTemplate string) *Help {
	h.optionsTemplate = optionsTemplate
	return h
}

func (h *Help) GetCommandsTemplate() string {
	return h.commandsTemplate
}

func (h *Help) SetCommandsTemplate(commandsTemplate string) *Help {
	h.commandsTemplate = commandsTemplate
	return h
}
