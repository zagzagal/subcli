package subcli

import (
	"flag"
	"io"
	"sync"
)

// Program is the overhead info on the software
type Program struct {
	Name    string
	Version string
	SDesc   string
}

// SubCommand is the description of a command
type SubCommand struct {
	Command string
	SDesc   string
	Help    string
	Cmd     func([]string)
	Flags   *flag.FlagSet
}

// HelpTopic is for additional help topics
type HelpTopic struct {
	Command string
	SDesc   string
	Help    string
}

// Action is the interface that contains the executed command
type Action interface {
	// Exec should proform the function taking all arguments after the
	// Sub command. Does no type checking, or varification
	Exec()
}

// SubCli is the housing struct for subcli
type SubCli struct {
	Program
	subs   []SubCommand
	help   []SubCommand
	Output io.Writer
	Args   []string
	data   *sync.Map
	parsed bool
}
