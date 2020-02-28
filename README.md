# goNixArgParser - Unix/Linux style cli args parser for Go

## Concepts
Command line arguments may contains several kinds of parts:
- Command and Sub Command
- Option
    - Flag
    - Value
- Rest

## Example
Here is an example for command `git remote`:
```
| Command | Sub Command | Sub Command |    Option     | Option |  Rest   |                Rest                 |
|         |             |             | Flag | Value  |  Flag  |         |                                     |
----------------------------------------------------------------------------------------------------------------
    git       remote          add        -t    master     -f      origin    https://repo.server.com/project.git
----------------------------------------------------------------------------------------------------------------
```

# Step 1 - Define Schema
## Root Command
```go
cmdGit := goNixArgParser.NewSimpleCommand("git", "A version control tool")
```

## Sub Commands
```go
cmdRemote := cmdGit.NewSimpleSubCommand("remote", "manage remotes")
cmdAdd := cmdRemote.NewSimpleSubCommand("add", "add a remote repository")
```

Command alias names can be specified as additional arguments:
```go
// "co" and "ckout" are alias names of "checkout" command
cmdCheckout := cmdGit.NewSimpleSubCommand("checkout", "checkout branches or files", "co", "ckout")
```

## Options
Here we want to define options on `cmdAdd` for `git remote add`. But if no sub command is needed, then just define them on root Command.

Available methods on `*Command.Options`:
- `AddFlag(key, flag, envVar, summary string)`  // single flag, without values
- `AddFlags(key string, flags []string, envVar, summary string)`  // multiple flag, without values
- `AddFlagValue(key, flag, envVar, defaultValue, summary string)`  // single flag, single value
- `AddFlagValues(key, flag, envVar string, defaultValues []string, summary string)`  // single flag, multiple values
- `AddFlagsValue(key string, flags []string, envVar, defaultValue, summary string)` // multiple flags, single value
- `AddFlagsValues(key string, flags []string, envVar string, defaultValues []string, summary string)`  / multiple flags, multiple values

`key` is a unique option name in `Options`.

```go
cmdAdd.Options().AddFlagValue("track", "-t", "", "", "only track specified branch")
cmdAdd.Options().AddFlag("fetch", "-f", "", "fetch after added")
```

# Step 2 - Parse
```go
// os.Args == []string{"git", "remote", "add", "-t", "master", "-f", "origin", "https://repo.server.com/project.git"}
results := cmdGit.Parse(os.Args, nil)
```
# Step 3 - Get Results
There are several methods on parsed result to get final values:
- `HasKey(key string) bool`
- `GetString(key string) (value string, found bool)`
- `GetBool(key string) (value bool, found bool)`
- `GetInt(key string) (value int, found bool)`
- `GetInt64(key string) (value int64, found bool)`
- `GetUint64(key string) (value uint64, found bool)`
- `GetFloat64(key string) (value float64, found bool)`
- `GetStrings(key string) (values []string, found bool)`
- `GetBools(key string) (values []bool, found bool)`
- `GetInts(key string) (values []int, found bool)`
- `GetInt64s(key string) (values []int64, found bool)`
- `GetUint64s(key string) (values []uint64, found bool)`
- `GetFloat64s(key string) (values []float64, found bool)`
- `GetRests() (rests []string)`
- `HasUndef() bool`
- `GetUndefs() []string`

Getting value for the example above:
```go
cmdPath := results.GetCommands()
if cmdPath[1] == "remote" && cmdPath[2] == "add" {
    if track, _ := results.GetString("track"); track != "" {
      // "-t" is supplied
    }
    if results.HasKey("fetch") {
      // "-f" is supplied
    }
}
```

# Configs
One application may have external config file. When application starts, it reads both command line args and config file.
Generally, the command line args is prior than config file.
The parser provides a quick and easy way to deal with this situation on parsing, if the contents of config are a list of arguments,
similar to the form of command line args.
The root command of config args can be omitted.

```go
cliArgs := []string{"cmd", "subCmd", "subSubCmd", "--option1", "value1"}
configArgs := []string{"cmd", "subCmd", "subSubCmd", "--option1", "value1FromConfig", "--option2", "value2FromConfig"}
// configArgs = configArgs[1:]  // omit root command
result := cmd.Parse(cliArgs, configArgs)

option1 := result.GetString("option1")  // "value1"
option2 := result.GetString("option2")  // "value2FromConfig"
```

# Env Var & Default Value
An option value can be set by Env var if it is not specified by other ways.
An option's related Env var can be specified when defining schema.

An option value can be set by default value if it is not specified by other ways.
An option's related default value can be specified when defining schema.

# Priority
The priority of getting an option's value is:
- input arg
- Env var
- config item
- default value

# Arg Groups
Sometimes a command may do tasks for multiple targets of a kind, e.g. start multiple spare services with different options.
By default, the parser treat `,,` as the separator of arg groups. Use `ParseGroups` instead of `Parse`,
Returns slice of parsed result for each group.
```go
cliArgs := []string{"cmd", "subCmd", "subSubCmd", "--option", "value1", ",,", "--option","value2"}
results := cmd.ParseGroups(cliArgs, nil)

service0option := results[0].GetString("optionKey")  // "value1"
service1option := results[1].GetString("optionKey")  // "value2"
```

Similar to `Parse`, The second parameter of `ParseGroups` is groups of config args, each group is related to its input arg groups by index,
and root command can be omitted.
```go
cliArgs := []string{"cmd", "subCmd", "subSubCmd", "--optionX", "valueX1", ",,", "--optionY","valueY2"}
configArgs := []string{"cmd", "subCmd", "subSubCmd", "--optionX", "valueX1FromConfig", ",,", "--optionX", "valueX2FromConfig"}
// configArgs = configArgs[1:]  // omit root command
results := cmd.ParseGroups(cliArgs, nil)

service0option := results[0].GetString("optionXKey")  // "valueX1"
service1option := results[1].GetString("optionXKey")  // "valueX2FromConfig"
```

if arg group separator is the last arg, then there is an empty option set follows.

# Control the Detail
When defining schemas, methods like `NewSimpleXXX` on command, or `AddXXX` on options,
are shortcuts that hides detail of bottom layer.
If you want to control or customize on more detail level, then the following part explains.
Use `NewCommand` to create command schema manually.
Use `*OptionSet.Append(*Option)` to add option schema manually.
use `NewFlag` to create flag schema manually, and append to `*Option.Flags`.

## Command struct
Both root command and sub commands are of type `*Command`. Initial parameters:

### `name`
Command or sub command name used for parsing

### `summary`
Summary about the command, which will be shown when invoking `GetHelp()` method

## OptionSet struct
One OptionSet manages all options supported by its command or sub command. Initial parameters:

### `mergeFlagPrefix`
Specify the prefix of flag that can be merged together. For example, following commands are equal:
```bash
ls -a -l
ls -al
```
The `-a` and `-l` are merged with the same prefix `-`.
Flag names which has only 1 suffix character can be merged.

### `restsSigns`
Sometimes we want to specify rest args explicitly, e.g. for "git checkout":
```bash
git checkout -- file1 file2 file3
```
Here `--` is a rests sign. It can be specified by other values when initializing an OptionSet.

### `groupSeps`
`groupSeps` is the separator to split "Arg Groups". Can be customized when initializing an OptionSet.

### `undefFlagPrefixes`
If an argument is not a flag, and begin with one of the `undefFlagPrefixes`, treat it as an undefined flag.
Otherwise treat this argument as previous flag's value or rests value.

## Option struct
`Option` represents an individual option. Some initial parameter:

### `AcceptValue`
Specifies if this option is flag only or can receive values.

### `MultiValues`
For option that can receive values, specify if it accepts multiple values.

### `OverridePrev`
For option that accepts values, when it is supplied by multiple times, specify if the later one will override the previous one.
For multiple-value option, if `OverridePrev` is `false`, then later items will be appended to previous.

### `Delimiters`
A multiple-value option's values can be supplied as a string separated by `Delimiters`.
Following args have the same effect if delimiter is `,`:
```bash
cmd subCmd --option value1,value2,value3
cmd subCmd --option value1 --option value2 --option value3
```

### `UniqueValues`
Remove duplicated values for parsed result automatically if true.

### `EnvVars`
Env var names for the option as fallback if option is not supplied.
Will look into it one by one util find a non-empty value.
Multiple-value option's value should be separated by `Delimiters`.

### `DefaultValues`
Default values for the option as fallback if option is not supplied.
For option that only accepts single value, only first element is valid.

## Flag struct
`Flag` represents a flag of option. Some initial parameter:

### `Name`
Name of the flag, e.g. `--option-name`.

### `canMerge`
Specify if this flag can be merged with others when the name starts with option's `mergeFlagPrefix` and suffix has only 1 character.

### `canFollowAssign`
Specify if treat args that follow after a flag as option values. Example for follow assign:
```bash
ls --hide '*.go'
```

### `assignSigns`
Specify symbols (e.g `=`) as assign symbols, separate value to its flag. Example for `=` assign:
```bash
ls --hide='*.go'
```

### `canConcatAssign`
Specify if option value can be concatenated after flag name, like `mysql` client tool:
```bash
mysql -uroot
# same as:
mysql -u root
```

## Creating Custom Command Schema
Use `NewCommand` to create a customized Command, instead of `NewSimpleCommand`:
```go
func NewCommand(
	names []string,
	summary, mergeFlagPrefix string,
	restsSigns, groupSeps, undefFlagPrefixes []string,
) *Command
```

Similarly, use `NewSubCommand` instead of `NewSimpleSubCommand` to create a sub command:
```go
func (c *Command) NewSubCommand(
	names []string,
	summary, mergeFlagPrefix string,
	restsSigns, groupSeps, undefFlagPrefixes []string,
) *Command
```

## Creating Custom Option Schema
Option is defined on `*OptionSet`, which can be got by `*Command.Options()`.
Use `*OptionSet.Append` to create a customized option, instead of the `AddXXX` method:
```go
func (s *OptionSet) Append(opt *Option) error

type Option struct {
	Key           string
	Summary       string
	Description   string
	Flags         []*Flag
	AcceptValue   bool
	MultiValues   bool
	OverridePrev  bool
	Delimiters    []rune
	UniqueValues  bool
	EnvVars       []string
	DefaultValues []string
}
```

## Creating Custom Flag Schema
`Option.Flags` is a slice of `*Flag`, Some useful functions to create them:
```go
func NewFlag(name string, canMerge, canFollowAssign, canConcatAssign bool, assignSigns []string) *Flag
func NewSimpleFlag(name string) *Flag
func NewSimpleFlags(names []string) []*Flag
```
