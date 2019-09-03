package goNixArgParser

import (
	"bytes"
	"io"
)

func (opt *Option) isDelimiter(r rune) bool {
	for _, delimiter := range opt.Delimiters {
		if r == delimiter {
			return true
		}
	}
	return false
}

func (opt *Option) GetHelp() []byte {
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
			buffer.WriteString(" ...")
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
		io.Copy(buffer, dftBuffer)
	}

	if buffer.Len() > 0 {
		buffer.WriteByte('\n')
	}

	return buffer.Bytes()
}
