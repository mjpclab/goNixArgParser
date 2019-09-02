package goNixArgParser

import "fmt"

func (c *Command) String() string {
	return ""
}

func (c *Command) PrintHelp() {
	fmt.Print(c.String())
}
