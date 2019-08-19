package goNixArgParser

import (
	"fmt"
	"testing"
)

func TestParse(t *testing.T) {
	var err error

	s := NewOptionSet(true, "-");
	err = s.Append(&Option{
		Key:         "tag",
		Flags:       []string{"-t", "--tag"},
		AcceptValue: false,
	})
	if err != nil {
		t.Error(err)
	}

	err = s.Append(&Option{
		Key:         "single",
		Flags:       []string{"-s", "--single"},
		AcceptValue: true,
	})
	if err != nil {
		t.Error(err)
	}

	err = s.Append(&Option{
		Key:         "multi",
		Flags:       []string{"-m", "--multi"},
		AcceptValue: true,
		MultiValues: true,
		Delimiter:   ",",
	})
	if err != nil {
		t.Error(err)
	}

	err = s.Append(&Option{
		Key:          "deft",
		Flags:        []string{"-df", "--default"},
		AcceptValue:  true,
		DefaultValue: []string{"myDefault"},
	})
	if err != nil {
		t.Error(err)
	}

	err = s.Append(&Option{
		Key:         "singleMissingValue",
		Flags:       []string{"-sm", "--single-missing"},
		AcceptValue: true,
	})
	if err != nil {
		t.Error(err)
	}

	err = s.Append(&Option{
		Key:         "flagX",
		Flags:       []string{"-x"},
		AcceptValue: true,
	})
	if err != nil {
		t.Error(err)
	}

	err = s.Append(&Option{
		Key:         "flagY",
		Flags:       []string{"-y"},
		AcceptValue: true,
	})
	if err != nil {
		t.Error(err)
	}

	args := []string{
		"-t",
		"-un1", "val1",
		"--single", "singleval1",
		"xxx",
		"-m", "multival1", "multival2",
		"--multi", "multival3,multival4",
		"-sm",
		"-xy",
	}
	r := s.Parse(args)
	fmt.Printf("%+v\n", r)

	if r.Contains("deft") {
		t.Error("deft")
	}

	if r.GetValue("deft") != "myDefault" {
		t.Error("default")
	}

	single := r.GetValue("single")
	fmt.Println("single:", single)
	if single != "singleval1" {
		t.Error("single")
	}

	multi := r.GetValues("multi")
	fmt.Println("multi:", multi)
	if len(multi) != 4 {
		t.Error("multi should have 4 values")
	}

	if !r.Contains("flagX") {
		t.Error("flagX")
	}

	if !r.Contains("flagY") {
		t.Error("flagY")
	}

}
