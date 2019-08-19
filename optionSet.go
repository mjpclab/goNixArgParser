package goNixArgParser

import (
	"errors"
	"fmt"
	"os"
)

func NewOptionSet(canMergeOption bool, mergeOptionPrefix string) *OptionSet {
	s := &OptionSet{
		canMergeOption:    canMergeOption,
		mergeOptionPrefix: mergeOptionPrefix,
		options:           []*Option{},
		keyOptionMap:      map[string]*Option{},
		flagOptionMap:     map[string]*Option{},
		keyDefaultMap:     map[string][]string{},
	}
	return s
}

func (s *OptionSet) Append(opt *Option) error {
	if len(opt.Key) == 0 {
		return errors.New("key not found")
	}
	if len(opt.Flags) == 0 {
		return errors.New("flag not found")
	}
	if s.keyOptionMap[opt.Key] != nil {
		return errors.New("key already exists")
	}
	for _, flag := range opt.Flags {
		if s.flagOptionMap[flag] != nil {
			return errors.New("flag '" + flag + "' already exists")
		}
	}

	optCopied := *opt
	option := &optCopied

	s.options = append(s.options, option)
	s.keyOptionMap[option.Key] = option
	for _, flag := range option.Flags {
		s.flagOptionMap[flag] = option
	}
	if len(option.DefaultValue) > 0 {
		s.keyDefaultMap[option.Key] = option.DefaultValue
	}
	return nil
}

func (s *OptionSet) Flag(key, flag, summary string) error {
	return s.Append(&Option{
		Key:     key,
		Flags:   []string{flag},
		Summary: summary,
	})
}

func (s *OptionSet) Flags(key string, flags []string, summary string) error {
	return s.Append(&Option{
		Key:     key,
		Flags:   flags,
		Summary: summary,
	})
}

func (s *OptionSet) FlagValue(key, flag, defaultValue, summary string) error {
	return s.Append(&Option{
		Key:          key,
		Flags:        []string{flag},
		AcceptValue:  true,
		DefaultValue: []string{defaultValue},
		Summary:      summary,
	})
}

func (s *OptionSet) FlagValues(key, flag string, defaultValues []string, summary string) error {
	return s.Append(&Option{
		Key:          key,
		Flags:        []string{flag},
		AcceptValue:  true,
		MultiValues:  true,
		DefaultValue: defaultValues,
		Summary:      summary,
	})
}

func (s *OptionSet) FlagsValue(key string, flags []string, defaultValue, summary string) error {
	return s.Append(&Option{
		Key:          key,
		Flags:        flags,
		AcceptValue:  true,
		DefaultValue: []string{defaultValue},
		Summary:      summary,
	})
}

func (s *OptionSet) FlagsValues(key string, flags, defaultValues []string, summary string) error {
	return s.Append(&Option{
		Key:          key,
		Flags:        flags,
		AcceptValue:  true,
		MultiValues:  true,
		DefaultValue: defaultValues,
		Summary:      summary,
	})
}

func (s *OptionSet) PrintHelp() {
	fmt.Println("Usage of " + os.Args[0] + ":")
	fmt.Println()

	for _, opt := range s.options {
		for i, flag := range opt.Flags {
			if i > 0 {
				fmt.Print(", ")
			}
			fmt.Print(flag)
		}
		fmt.Println()

		if len(opt.Summary) > 0 {
			fmt.Println(opt.Summary)
		}

		if len(opt.Description) > 0 {
			fmt.Println(opt.Description)
		}

		fmt.Println()
	}
}
