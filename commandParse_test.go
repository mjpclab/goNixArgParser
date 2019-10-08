package goNixArgParser

import (
	"fmt"
	"testing"
)

func getGitCommand() *Command {
	cmdGit := NewSimpleCommand("git", "A version control tool")

	cmdGit.Options.AddFlag("version", "--version", "", "display version")
	cmdGit.Options.AddFlag("help", "--help", "", "show git help")

	cmdSetUrl := cmdGit.NewSimpleSubCommand("remote", "manage remotes").NewSimpleSubCommand("set-url", "set remote url")
	cmdSetUrl.Options.AddFlag("push", "--push", "", "")
	cmdSetUrl.Options.AddFlagValue("dummy", "--dummy", "", "", "dummy option")

	cmdReset := cmdGit.NewSimpleSubCommand("reset", "reset command")
	cmdReset.Options.AddFlag("hard", "--hard", "", "hard reset")
	cmdReset.Options.AddFlag("mixed", "--mixed", "", "mixed reset")
	cmdReset.Options.AddFlag("soft", "--soft", "", "soft reset")

	return cmdGit
}

func TestNormalizeCmdArgs(t *testing.T) {
	cmd := getGitCommand()
	args := []string{"git", "remote", "set-url", "--push", "origin", "https://github.com/mjpclab/goNixArgParser.git"}
	_, normalizedArgs := cmd.getNormalizedArgs(args)
	for i, arg := range normalizedArgs {
		fmt.Printf("%d %+v\n", i, arg)
	}
}

func TestParseCommand1(t *testing.T) {
	cmd := getGitCommand()
	args := []string{"git", "remote", "set-url", "--push", "origin", "https://github.com/mjpclab/goNixArgParser.git"}

	result := cmd.Parse(args, nil)
	if result.commands[0] != "git" ||
		result.commands[1] != "remote" ||
		result.commands[2] != "set-url" {
		t.Error("commands", result.commands)
	}

	if !result.HasFlagKey("push") {
		t.Error("push")
	}

	if result.argRests[0] != "origin" ||
		result.argRests[1] != "https://github.com/mjpclab/goNixArgParser.git" {
		t.Error("rests", result.argRests)
	}

	cmd.PrintHelp()
}

func TestParseCommand2(t *testing.T) {
	cmd := getGitCommand()
	args := []string{"git", "remote", "xxx", "set-url", "origin", "https://github.com/mjpclab/goNixArgParser.git"}

	result := cmd.Parse(args, nil)
	if result.commands[0] != "git" ||
		result.commands[1] != "remote" {
		t.Error("commands", result.commands)
	}

	if result.argRests[0] != "xxx" ||
		result.argRests[1] != "set-url" ||
		result.argRests[2] != "origin" ||
		result.argRests[3] != "https://github.com/mjpclab/goNixArgParser.git" {
		t.Error("rests", result.argRests)
	}
}

func TestParseCommand3(t *testing.T) {
	cmd := getGitCommand()
	args := []string{"git", "remote", "set-url", "origin", "https://github.com/mjpclab/goNixArgParser.git"}
	configArgs := []string{"git", "remote", "set-url", "--dummy", "dummyconfigvalue"}
	result := cmd.Parse(args, configArgs)

	dummy, _ := result.GetString("dummy")
	if dummy != "dummyconfigvalue" {
		fmt.Println("dummy:", dummy)
		t.Error("dummy config value error")
	}

	configArgs = configArgs[1:]
	result = cmd.Parse(args, configArgs)

	dummy, _ = result.GetString("dummy")
	if dummy != "dummyconfigvalue" {
		fmt.Println("dummy:", dummy)
		t.Error("dummy config value error")
	}
}

func TestParseCommand4(t *testing.T) {
	cmd := getGitCommand()
	args := []string{"git", "remote", "set-url", "--dummy", "dummy0", "github", "https://github.com/mjpclab/goNixArgParser.git", ",,", "--dummy", "dummy1", "bitbucket", "https://bitbucket.com/mjpclab/goNixArgParser.git"}
	results := cmd.ParseGroups(args, nil)

	dummy0, _ := results[0].GetString("dummy")
	if dummy0 != "dummy0" {
		t.Error(results[0].GetStrings("dummy"))
	}

	dummy1, _ := results[1].GetString("dummy")
	if dummy1 != "dummy1" {
		t.Error(results[1].GetStrings("dummy"))
	}
}

func TestParseCommand5(t *testing.T) {
	cmd := getGitCommand()
	args := []string{"git", "remote", "set-url", "github", "https://github.com/mjpclab/goNixArgParser.git", ",,", "bitbucket", "https://bitbucket.com/mjpclab/goNixArgParser.git"}
	configArgs := []string{"git", "remote", "set-url", "--dummy", "dummy0", ",,", "--dummy", "dummy1"}
	results := cmd.ParseGroups(args, configArgs)

	dummy0, _ := results[0].GetString("dummy")
	if dummy0 != "dummy0" {
		t.Error(results[0].GetStrings("dummy"))
	}

	dummy1, _ := results[1].GetString("dummy")
	if dummy1 != "dummy1" {
		t.Error(results[1].GetStrings("dummy"))
	}

	configArgs = configArgs[1:]
	results = cmd.ParseGroups(args, configArgs)

	dummy0, _ = results[0].GetString("dummy")
	if dummy0 != "dummy0" {
		t.Error(results[0].GetStrings("dummy"))
	}

	dummy1, _ = results[1].GetString("dummy")
	if dummy1 != "dummy1" {
		t.Error(results[1].GetStrings("dummy"))
	}
}
