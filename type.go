package goNixArgParser

type Command struct {
	Name        string
	OptionSet   *OptionSet
	SubCommands []*Command
}

type OptionSet struct {
	mergeOptionPrefix string
	options           []*Option

	keyOptionMap  map[string]*Option
	flagOptionMap map[string]*Option
	flagMap       map[string]*Flag
	keyDefaultMap map[string][]string
}

type Option struct {
	Key           string
	Summary       string
	Description   string
	Flags         []*Flag
	AcceptValue   bool
	MultiValues   bool
	OverridePrev  bool
	Delimiters    string
	DefaultValues []string
}
type Flag struct {
	Name            string
	canMerge        bool
	canEqualAssign  bool
	canConcatAssign bool
}

type ParseResult struct {
	inputs   []*Arg
	params   map[string][]string
	defaults map[string][]string
	rests    []string
}

type ArgType int

const (
	Unknown ArgType = iota
	SubCmd
	FlagName
	Value
	Rest
)

type Arg struct {
	Text string
	Type ArgType
}
