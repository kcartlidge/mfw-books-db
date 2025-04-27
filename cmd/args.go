package main

import (
	"fmt"
	"strings"
)

// ArgsParser represents a command line parser
type ArgsParser struct {
	Arguments []CommandArgument
	Flags     []CommandFlag
	Errors    []string

	providedArguments              map[string]string
	providedArgumentsWithoutValues map[string]bool
	providedFlags                  map[string]bool
	providedPlainText              map[string]bool
}

// CommandArgument represents an argument in the command line
type CommandArgument struct {
	Name         string
	Description  string
	DefaultValue string
	Provided     bool
	Value        string
	IsRequired   bool
}

// CommandFlag represents a flag in the command line
type CommandFlag struct {
	Name        string
	Description string
	Provided    bool
}

// NewArgsParser creates a new command line parser
func NewArgsParser() *ArgsParser {
	return &ArgsParser{
		Arguments:                      []CommandArgument{},
		Flags:                          []CommandFlag{},
		Errors:                         []string{},
		providedArguments:              make(map[string]string),
		providedArgumentsWithoutValues: make(map[string]bool),
		providedFlags:                  make(map[string]bool),
		providedPlainText:              make(map[string]bool),
	}
}

// AddArgument adds an argument to the command line parser
func (parser *ArgsParser) AddArgument(name string, description string, defaultValue string, isRequired bool) {
	parser.Arguments = append(parser.Arguments, CommandArgument{
		Name:         strings.ToLower(name),
		Description:  description,
		DefaultValue: defaultValue,
		Provided:     false,
		Value:        "",
		IsRequired:   isRequired,
	})
}

// AddFlag adds a flag to the command line parser
func (parser *ArgsParser) AddFlag(name string, description string) {
	parser.Flags = append(parser.Flags, CommandFlag{Name: strings.ToLower(name), Description: description, Provided: false})
}

// ShowUsage shows the usage of the command line arguments and flags
func (parser *ArgsParser) ShowUsage() {
	maxNameLength := 0
	for _, argument := range parser.Arguments {
		if len(argument.Name)+8 > maxNameLength {
			maxNameLength = len(argument.Name) + 8
		}
	}
	for _, flag := range parser.Flags {
		if len(flag.Name) > maxNameLength {
			maxNameLength = len(flag.Name)
		}
	}
	maxNameLength += 2

	fmt.Println()
	fmt.Println("Usage:")
	hasRequired := false
	for _, argument := range parser.Arguments {
		name := argument.Name + " <value>"
		spaces := strings.Repeat(" ", maxNameLength-len(name)+1)
		desc := argument.Description
		if argument.IsRequired {
			desc = "* " + desc
			hasRequired = true
		} else {
			desc = "  " + desc
		}
		fmt.Printf("  -%s%s %s\n", name, spaces, desc)
	}
	for _, flag := range parser.Flags {
		spaces := strings.Repeat(" ", maxNameLength-len(flag.Name))
		fmt.Printf("  --%s%s   %s\n", flag.Name, spaces, flag.Description)
	}
	if hasRequired {
		fmt.Println()
		fmt.Printf("  %s\n", "(arguments marked with * are required)")
	}
}

// Parse parses the arguments and flags from the command line
func (parser *ArgsParser) Parse(args []string) {

	parser.Errors = []string{}

	// Parse the provided arguments and flags
	for i := 0; i < len(args); i++ {
		// Normalise and check the dashes
		arg := args[i]
		plainArg := strings.ToLower(arg)
		for strings.HasPrefix(plainArg, "-") {
			plainArg = strings.TrimPrefix(plainArg, "-")
		}
		dashes := len(arg) - len(plainArg)

		if dashes == 2 { // Flag (--name)
			parser.providedFlags[plainArg] = true
		} else if dashes == 1 { // Argument (-name value)
			if i+1 >= len(args) {
				parser.providedArgumentsWithoutValues[plainArg] = true
				continue
			}
			value := args[i+1]
			if strings.HasPrefix(value, "-") {
				parser.providedArgumentsWithoutValues[plainArg] = true
				continue
			}
			parser.providedArguments[plainArg] = value
			i++
		} else { // Plain text
			parser.providedPlainText[arg] = true
		}
	}

	// Add in any default values that were not provided
	for i := range parser.Arguments {
		if parser.Arguments[i].DefaultValue != "" {
			if _, ok := parser.providedArguments[parser.Arguments[i].Name]; !ok {
				parser.providedArguments[parser.Arguments[i].Name] = parser.Arguments[i].DefaultValue
			}
		}
	}

	// Run through, allocating to expected flags
	for i := range parser.Flags {
		// Matched flags are tracked then removed from the provided list
		if _, ok := parser.providedFlags[parser.Flags[i].Name]; ok {
			parser.Flags[i].Provided = true
			delete(parser.providedFlags, parser.Flags[i].Name)
		}
	}

	// Run through, allocating to expected arguments
	for i := range parser.Arguments {
		// Matched arguments with values are tracked then removed from the provided list
		if _, ok := parser.providedArguments[parser.Arguments[i].Name]; ok {
			parser.Arguments[i].Provided = true
			parser.Arguments[i].Value = parser.providedArguments[parser.Arguments[i].Name]
			delete(parser.providedArguments, parser.Arguments[i].Name)
		}
	}

	// Check for required arguments that were not provided
	for i := range parser.Arguments {
		if parser.Arguments[i].IsRequired && !parser.Arguments[i].Provided {
			if _, ok := parser.providedArgumentsWithoutValues[parser.Arguments[i].Name]; ok {
				parser.Errors = append(parser.Errors, fmt.Sprintf("Argument expected a value:  -%s", parser.Arguments[i].Name))
				delete(parser.providedArgumentsWithoutValues, parser.Arguments[i].Name)
			} else {
				parser.Errors = append(parser.Errors, fmt.Sprintf("Argument is required:  -%s", parser.Arguments[i].Name))
			}
		}
	}

	// Add errors for unexpected arguments that included a value
	for k := range parser.providedArguments {
		parser.Errors = append(parser.Errors, fmt.Sprintf("Unknown argument:  -%s", k))
		delete(parser.providedArguments, k)
	}

	// Add errors for arguments that were missing a value
	for k := range parser.providedArgumentsWithoutValues {
		found := false
		for _, argument := range parser.Arguments {
			if argument.Name == k {
				parser.Errors = append(parser.Errors, fmt.Sprintf("Missing value for argument:  -%s", k))
				delete(parser.providedArgumentsWithoutValues, k)
				found = true
			}
		}
		if !found {
			parser.Errors = append(parser.Errors, fmt.Sprintf("Unknown argument (and missing value):  -%s", k))
			delete(parser.providedArgumentsWithoutValues, k)
		}
	}

	// Add errors for unexpected flags
	for k := range parser.providedFlags {
		parser.Errors = append(parser.Errors, fmt.Sprintf("Unknown flag:  --%s", k))
		delete(parser.providedFlags, k)
	}

	// Add errors for unexpected plain text
	for k := range parser.providedPlainText {
		parser.Errors = append(parser.Errors, fmt.Sprintf("Unexpected plain text:  %s", k))
		delete(parser.providedPlainText, k)
	}
}

// HasArgument returns true if an argument was provided, otherwise it returns false
func (parser *ArgsParser) HasArgument(name string) bool {
	for _, argument := range parser.Arguments {
		if strings.EqualFold(argument.Name, name) {
			return argument.Provided
		}
	}
	return false
}

// GetArgument returns the value of an argument (or an empty string)
// You can also look at the Provided flag to see if it was provided
func (parser *ArgsParser) GetArgument(name string) string {
	for _, argument := range parser.Arguments {
		if strings.EqualFold(argument.Name, name) {
			return argument.Value
		}
	}
	return ""
}

// GetFlag returns true if a flag was provided, otherwise it returns false
func (parser *ArgsParser) GetFlag(name string) bool {
	for _, flag := range parser.Flags {
		if strings.EqualFold(flag.Name, name) {
			return flag.Provided
		}
	}
	return false
}

// HasErrors returns true if there are any errors in the parser
func (parser *ArgsParser) HasErrors() bool {
	return len(parser.Errors) > 0
}

// PrintErrors prints all collected errors
func (parser *ArgsParser) PrintErrors() {
	if len(parser.Errors) > 0 {
		fmt.Println()
		fmt.Println("Errors:")
		for _, err := range parser.Errors {
			fmt.Printf("  %s\n", err)
		}
	}
}

// ShowProvided shows which arguments and flags were provided
func (parser *ArgsParser) ShowProvided() {
	fmt.Println()
	fmt.Println("Provided:")

	// Show provided arguments
	for _, argument := range parser.Arguments {
		if argument.Provided {
			fmt.Printf("  -%s %s\n", argument.Name, argument.Value)
		}
	}

	// Show provided flags
	for _, flag := range parser.Flags {
		if flag.Provided {
			fmt.Printf("  --%s\n", flag.Name)
		}
	}
}
