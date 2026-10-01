// Package cli implements the tidyid command-line interface.
package cli

import (
	"fmt"
	"io"
	"strconv"

	"github.com/sheldonix/tidyid-go/v2"
)

const help = `Usage
  tidyid [options]

Options
  -s, --size <SIZE>       Generated ID size (3-256)
  -u, --allow-uppercase   Allow uppercase letters
  -v, --version           Show version number
  -h, --help              Show this help`

// Run executes the CLI with arguments that exclude the executable name.
func Run(args []string, stdout io.Writer) error {
	for _, argument := range args {
		if argument == "-v" || argument == "--version" {
			_, err := fmt.Fprintln(stdout, tidyid.Version)
			return err
		}
	}
	for _, argument := range args {
		if argument == "-h" || argument == "--help" {
			_, err := fmt.Fprintln(stdout, help)
			return err
		}
	}

	length := tidyid.DefaultLength
	allowUppercase := false
	for index := 0; index < len(args); index++ {
		switch argument := args[index]; argument {
		case "-u", "--allow-uppercase":
			allowUppercase = true
		case "-s", "--size":
			index++
			if index == len(args) {
				return invalidSizeError()
			}
			value, err := strconv.Atoi(args[index])
			if err != nil || value < tidyid.MinLength || value > tidyid.MaxLength {
				return invalidSizeError()
			}
			length = value
		default:
			return fmt.Errorf("Unknown argument %s", argument)
		}
	}

	value, err := tidyid.Generate(length, allowUppercase)
	if err != nil {
		return err
	}
	_, err = fmt.Fprintln(stdout, value)
	return err
}

func invalidSizeError() error {
	return fmt.Errorf(
		"Size must be an integer between %d and %d",
		tidyid.MinLength,
		tidyid.MaxLength,
	)
}
