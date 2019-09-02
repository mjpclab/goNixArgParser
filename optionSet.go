package goNixArgParser

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"path"
)

func NewOptionSet(mergeOptionPrefix string) *OptionSet {
	s := &OptionSet{
		mergeOptionPrefix: mergeOptionPrefix,
		options:           []*Option{},
		keyOptionMap:      map[string]*Option{},
		flagOptionMap:     map[string]*Option{},
		flagMap:           map[string]*Flag{},
		keyDefaultMap:     map[string][]string{},
	}
	return s
}

func NewSimpleOptionSet() *OptionSet {
	return NewOptionSet("-")
}

func (s *OptionSet) Append(opt *Option) error {
	if len(opt.Key) == 0 {
		return errors.New("key is empty")
	}
	if s.keyOptionMap[opt.Key] != nil {
		return errors.New("key '" + opt.Key + "' already exists")
	}

	if len(opt.Flags) == 0 {
		return errors.New("flag not found")
	}
	for _, flag := range opt.Flags {
		flagName := flag.Name
		if len(flagName) == 0 {
			return errors.New("flag name is empty")
		}
		if s.flagMap[flagName] != nil {
			return errors.New("flag '" + flagName + "' already exists")
		}
	}

	optCopied := *opt
	option := &optCopied

	s.options = append(s.options, option)
	s.keyOptionMap[option.Key] = option
	for _, flag := range option.Flags {
		flagName := flag.Name
		s.flagOptionMap[flagName] = option
		s.flagMap[flagName] = flag
	}
	if len(option.DefaultValues) > 0 {
		s.keyDefaultMap[option.Key] = option.DefaultValues
	}
	return nil
}

func (s *OptionSet) AddFlag(key, flag, summary string) error {
	return s.Append(&Option{
		Key:     key,
		Flags:   []*Flag{NewSimpleFlag(flag)},
		Summary: summary,
	})
}

func (s *OptionSet) AddFlags(key string, flags []string, summary string) error {
	return s.Append(&Option{
		Key:     key,
		Flags:   NewSimpleFlags(flags),
		Summary: summary,
	})
}

func (s *OptionSet) AddFlagValue(key, flag, defaultValue, summary string) error {
	return s.Append(&Option{
		Key:           key,
		Flags:         []*Flag{NewSimpleFlag(flag)},
		AcceptValue:   true,
		DefaultValues: []string{defaultValue},
		Summary:       summary,
	})
}

func (s *OptionSet) AddFlagValues(key, flag string, defaultValues []string, summary string) error {
	return s.Append(&Option{
		Key:           key,
		Flags:         []*Flag{NewSimpleFlag(flag)},
		AcceptValue:   true,
		MultiValues:   true,
		DefaultValues: defaultValues,
		Summary:       summary,
	})
}

func (s *OptionSet) AddFlagsValue(key string, flags []string, defaultValue, summary string) error {
	return s.Append(&Option{
		Key:           key,
		Flags:         NewSimpleFlags(flags),
		AcceptValue:   true,
		DefaultValues: []string{defaultValue},
		Summary:       summary,
	})
}

func (s *OptionSet) AddFlagsValues(key string, flags, defaultValues []string, summary string) error {
	return s.Append(&Option{
		Key:           key,
		Flags:         NewSimpleFlags(flags),
		AcceptValue:   true,
		MultiValues:   true,
		DefaultValues: defaultValues,
		Summary:       summary,
	})
}

func (s *OptionSet) String() string {
	buffer := &bytes.Buffer{}
	for _, opt := range s.options {
		buffer.WriteByte('\n')
		buffer.WriteString(opt.String())
	}
	buffer.WriteByte('\n')
	return buffer.String()
}

func (s *OptionSet) PrintHelp() {
	fmt.Print("Usage of " + path.Base(os.Args[0]) + ":\n")

	fmt.Print(s.String())
}
