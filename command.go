package goNixArgParser

import (
	"bytes"
	"os"
	"path"
)

func NewCommand(
	names []string,
	summary, mergeFlagPrefix string,
	restsSigns, groupSeps, assignSigns, undefFlagPrefixes []string,
) *Command {
	return &Command{
		names:       names,
		summary:     summary,
		options:     NewOptionSet(mergeFlagPrefix, restsSigns, groupSeps, assignSigns, undefFlagPrefixes),
		subCommands: []*Command{},
	}
}

func NewSimpleCommand(name, summary string, aliasNames ...string) *Command {
	names := make([]string, 1+len(aliasNames))
	names[0] = name
	copy(names[1:], aliasNames)

	return &Command{
		names:       names,
		summary:     summary,
		options:     NewSimpleOptionSet(),
		subCommands: []*Command{},
	}
}

func (c *Command) NewSubCommand(
	names []string,
	summary, mergeFlagPrefix string,
	restsSigns, groupSeps, assignSigns, undefFlagPrefixes []string,
) *Command {
	subCommand := NewCommand(names, summary, mergeFlagPrefix, restsSigns, groupSeps, assignSigns, undefFlagPrefixes)
	c.subCommands = append(c.subCommands, subCommand)
	return subCommand
}

func (c *Command) NewSimpleSubCommand(name, summary string, aliasNames ...string) *Command {
	subCommand := NewSimpleCommand(name, summary, aliasNames...)
	c.subCommands = append(c.subCommands, subCommand)
	return subCommand
}

func (c *Command) hasName(name string) bool {
	for _, n := range c.names {
		if n == name {
			return true
		}
	}
	return false
}

func (c *Command) GetSubCommand(name string) *Command {
	for _, cmd := range c.subCommands {
		if cmd.hasName(name) {
			return cmd
		}
	}
	return nil
}

func (c *Command) Name() (name string) {
	if len(c.names) > 0 {
		name = c.names[0]
	}

	return
}

func (c *Command) Names() []string {
	names := make([]string, len(c.names))
	copy(names, c.names)
	return names
}

func (c *Command) Summary() string {
	return c.summary
}

func (c *Command) Options() *OptionSet {
	return c.options
}

func (c *Command) SubCommands() []*Command {
	return c.subCommands
}

func (c *Command) getLeafCmd(args []string) (explicitCmd *Command, inferredCmd *Command, cmdPaths []string) {
	inferredCmd = c

	if len(args) == 0 {
		return explicitCmd, inferredCmd, []string{}
	}

	for i, arg := range args {
		if i == 0 && inferredCmd.hasName(arg) {
			explicitCmd = c
			cmdPaths = append(cmdPaths, inferredCmd.Name())
		} else if subCmd := inferredCmd.GetSubCommand(arg); subCmd != nil {
			explicitCmd = subCmd
			inferredCmd = subCmd
			cmdPaths = append(cmdPaths, inferredCmd.Name())
		} else {
			break
		}
	}

	return
}

func (c *Command) extractCmdOptionArgs(specifiedArgs, configArgs []string) (
	specifiedCmd *Command,
	cmdPaths, specifiedOptionArgs, configOptionArgs []string,
) {
	_, specifiedCmd, specifiedCmdPaths := c.getLeafCmd(specifiedArgs)
	explicitConfigCmd, configCmd, configCmdPaths := c.getLeafCmd(configArgs)

	cmdPaths = specifiedCmdPaths

	specifiedOptionArgs = specifiedArgs[len(specifiedCmdPaths):]

	if specifiedCmd == configCmd {
		configOptionArgs = configArgs[len(configCmdPaths):]
	} else if explicitConfigCmd == nil {
		configOptionArgs = configArgs
	} else {
		configOptionArgs = []string{}
	}

	return
}

func (c *Command) Parse(specifiedArgs, configArgs []string) *ParseResult {
	cmd, cmdPaths, specifiedOptionArgs, configOptionArgs := c.extractCmdOptionArgs(specifiedArgs, configArgs)
	result := cmd.options.Parse(specifiedOptionArgs, configOptionArgs)
	result.commands = cmdPaths

	return result
}

func (c *Command) ParseGroups(specifiedArgs, configArgs []string) (results []*ParseResult) {
	cmd, cmdPaths, specifiedOptionArgs, configOptionArgs := c.extractCmdOptionArgs(specifiedArgs, configArgs)

	if len(specifiedOptionArgs) == 0 && len(configOptionArgs) == 0 {
		result := cmd.options.Parse(specifiedOptionArgs, configOptionArgs)
		results = append(results, result)
	} else {
		results = cmd.options.ParseGroups(specifiedOptionArgs, configOptionArgs)
	}

	for _, result := range results {
		result.commands = cmdPaths
	}

	return results
}

func (c *Command) GetHelp() []byte {
	buffer := &bytes.Buffer{}

	name := c.Name()
	if len(name) > 0 {
		buffer.WriteString(path.Base(name))
		buffer.WriteString(": ")
	}
	if len(c.summary) > 0 {
		buffer.WriteString(c.summary)
	}
	if buffer.Len() > 0 {
		buffer.WriteByte('\n')
	} else {
		buffer.WriteString("Usage:\n")
	}

	optionsHelp := c.options.GetHelp()
	if len(optionsHelp) > 0 {
		buffer.WriteString("\nOptions:\n\n")
		buffer.Write(optionsHelp)
	}

	if len(c.subCommands) > 0 {
		buffer.WriteString("\nSub commands:\n\n")
		for _, cmd := range c.subCommands {
			buffer.WriteString(cmd.Name())
			buffer.WriteByte('\n')
			if len(cmd.summary) > 0 {
				buffer.WriteString(cmd.summary)
				buffer.WriteByte('\n')
			}
			buffer.WriteByte('\n')
		}
	}

	return buffer.Bytes()
}

func (c *Command) PrintHelp() {
	os.Stdout.Write(c.GetHelp())
}
