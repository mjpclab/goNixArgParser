package goNixArgParser

import (
	"fmt"
	"testing"
)

func TestOptionSet(t *testing.T) {
	var err error

	s := NewOptionSet(true, "")

	err = s.Append(&Option{
		Key:          "deft",
		Flags:        []string{"-df", "--default"},
		AcceptValue:  true,
		DefaultValue: []string{"myDefault"},
	})
	if err != nil {
		t.Error(err)
	}

	err = s.Flag("flag", "-flag", "flag option")
	if err != nil {
		t.Error(err)
	}

	err = s.Flag("p", "p", "flag p")
	if err != nil {
		t.Error(err)
	}

	err = s.Flag("q", "q", "flag q")
	if err != nil {
		t.Error(err)
	}

	err = s.Flags("flags", []string{"-flags", "--flags"}, "flags option")
	if err != nil {
		t.Error(err)
	}

	err = s.FlagValue("port", "--port", "21", "port to listen")
	if err != nil {
		t.Error(err)
	}

	err = s.FlagValues("ports", "--ports", []string{"80", "8080"}, "ports to listen for http")
	if err != nil {
		t.Error(err)
	}

	err = s.FlagsValue("file", []string{"-f", "--files"}, "", "file to open")
	if err != nil {
		t.Error(err)
	}

	err = s.FlagsValues("props", []string{"-p", "--props"}, []string{}, "properties")
	if err != nil {
		t.Error(err)
	}

	args := []string{
		"-flag",
		"pq",
		"--flags",
		"--port", "22",
		"--ports", "443", "4443",
		"--ports", "4444",
		"-f", "file1", "file2",
		"--props", "aa", "bb", "cc", "dd", "ee",
	}

	parsed := s.Parse(args)
	if parsed.GetValue("deft") != "myDefault" {
		t.Error("deft")
	}
	if !parsed.Contains("flag") {
		t.Error("flag")
	}
	if !parsed.Contains("flags") {
		t.Error("flags")
	}
	if !parsed.Contains("p") {
		t.Error("p")
	}
	if !parsed.Contains("q") {
		t.Error("q")
	}
	if parsed.GetValue("port") != "22" {
		t.Error("port")
	}

	ports := parsed.GetValues("ports")
	fmt.Println("ports:", ports)
	if len(ports) != 3 {
		t.Error("ports")
	}

	if parsed.GetValue("file") != "file1" {
		t.Error("file")
	}

	props := parsed.GetValues("props")
	fmt.Println("props:", props)
	if len(props) != 5 {
		t.Error("props")
	}

	fmt.Println("rests:", parsed.GetRests())
	if len(parsed.GetRests()) != 1 {
		t.Error("rests")
	}
}
