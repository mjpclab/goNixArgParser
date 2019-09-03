package goNixArgParser

import (
	"os"
	"path"
)

var CommandLine *Command

func init() {
	var commandName string
	if len(os.Args) > 0 {
		commandName = path.Base(os.Args[0])
	}

	CommandLine = NewSimpleCommand(commandName, "")
}

func Append(opt *Option) error {
	return CommandLine.OptionSet.Append(opt)
}

func AddFlag(key, flag, summary string) error {
	return CommandLine.OptionSet.AddFlag(key, flag, summary)
}

func AddFlags(key string, flags []string, summary string) error {
	return CommandLine.OptionSet.AddFlags(key, flags, summary)
}

func AddFlagValue(key, flag, defaultValue, summary string) error {
	return CommandLine.OptionSet.AddFlagValue(key, flag, defaultValue, summary)
}

func AddFlagValues(key, flag string, defaultValues []string, summary string) error {
	return CommandLine.OptionSet.AddFlagValues(key, flag, defaultValues, summary)
}

func AddFlagsValue(key string, flags []string, defaultValue, summary string) error {
	return CommandLine.OptionSet.AddFlagsValue(key, flags, defaultValue, summary)
}

func AddFlagsValues(key string, flags, defaultValues []string, summary string) error {
	return CommandLine.OptionSet.AddFlagsValues(key, flags, defaultValues, summary)
}

func PrintHelp() {
	CommandLine.PrintHelp()
}

func Parse() *ParseResult {
	return CommandLine.Parse(os.Args)
}
