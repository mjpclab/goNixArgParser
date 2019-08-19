package goNixArgParser

import "os"

type Option struct {
	Key          string
	Summary      string
	Description  string
	Flags        []string
	AcceptValue  bool
	MultiValues  bool
	Delimiter    string
	DefaultValue []string
}

type OptionSet struct {
	canMergeOption    bool
	mergeOptionPrefix string
	options           []*Option
	keyOptionMap      map[string]*Option
	flagOptionMap     map[string]*Option
	keyDefaultMap     map[string][]string
}

type ParseResult struct {
	params   map[string][]string
	defaults map[string][]string
	rests    []string
}

var CommandLine *OptionSet = NewOptionSet(true, "-")

func Append(opt *Option) error {
	return CommandLine.Append(opt)
}

func Flag(key, flag, summary string) error {
	return CommandLine.Flag(key, flag, summary)
}

func Flags(key string, flags []string, summary string) error {
	return CommandLine.Flags(key, flags, summary)
}

func FlagValue(key, flag, defaultValue, summary string) error {
	return CommandLine.FlagValue(key, flag, defaultValue, summary)
}

func FlagValues(key, flag string, defaultValues []string, summary string) error {
	return CommandLine.FlagValues(key, flag, defaultValues, summary)
}

func FlagsValue(key string, flags []string, defaultValue, summary string) error {
	return CommandLine.FlagsValue(key, flags, defaultValue, summary)
}

func FlagsValues(key string, flags, defaultValues []string, summary string) error {
	return CommandLine.FlagsValues(key, flags, defaultValues, summary)
}

func PrintHelp() {
	CommandLine.PrintHelp()
}

func Parse() *ParseResult {
	return CommandLine.Parse(os.Args[1:])
}
