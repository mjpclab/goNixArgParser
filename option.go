package goNixArgParser

import (
	"bytes"
	"strings"
)

func (opt *Option) isDelimiter(r rune) bool {
	return strings.IndexRune(opt.Delimiters, r) >= 0
}

func (opt *Option) String() string {
	buffer := &bytes.Buffer{}

	for i, flag := range opt.Flags {
		if i > 0 {
			buffer.WriteString("|")
		}
		buffer.WriteString(flag.Name)
	}

	if opt.AcceptValue {
		buffer.WriteString(" <value>")
		if opt.MultiValues {
			buffer.WriteString(", ...")
		}
	}

	if len(opt.Summary) > 0 {
		buffer.WriteByte('\n')
		buffer.WriteString(opt.Summary)
	}

	if len(opt.Description) > 0 {
		buffer.WriteByte('\n')
		buffer.WriteString(opt.Description)
	}

	dftBuffer := &bytes.Buffer{}
	for _, d := range opt.DefaultValues {
		if len(d) > 0 {
			if dftBuffer.Len() > 0 {
				dftBuffer.WriteString(", ")
			}
			dftBuffer.WriteString(d)
		}
	}
	if dftBuffer.Len() > 0 {
		buffer.WriteByte('\n')
		buffer.WriteString("Default: ")
		buffer.WriteString(dftBuffer.String())
	}

	buffer.WriteByte('\n')

	return buffer.String()
}
