// Package subcli is a command line parser that aims to replicate the interface
// of the go command.
package subcli

import (
	"flag"
	"fmt"
	"os"
	"sort"
)

// New is the SubCli constructor function
// Output defaults to os.Stderr
func New(p Program) *SubCli {
	s := &SubCli{
		Program: p,
		subs:    *new([]SubCommand),
		help:    *new([]SubCommand),
		Output:  os.Stderr,
		Args:    os.Args,
	}
	s.AddCmd(SubCommand{
		Command: "version",
		SDesc:   "prints " + s.Name + " version",
		Help: "usage: " + s.Name +
			" version\n\nPrint the version of " + s.Name,
		Cmd:   nil,
		Flags: flag.NewFlagSet("version", flag.ContinueOnError),
	})
	return s
}

type subCommands []SubCommand

func (s subCommands) Len() int {
	return len(s)
}

func (s subCommands) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

func (s subCommands) Less(i, j int) bool {
	return s[i].Command < s[j].Command
}

// AddCmd adds a new SubCommand
func (s *SubCli) AddCmd(c SubCommand) {
	c.Flags = flag.NewFlagSet(c.Command, flag.ExitOnError)
	c.Flags.SetOutput(s.Output)
	s.subs = append(s.subs, c)
	sort.Sort(subCommands(s.subs))
}

// AddHelp adds a help topic
func (s *SubCli) AddHelp(h HelpTopic) {
	c := SubCommand{
		Command: h.Command,
		Help:    h.Help,
		SDesc:   h.SDesc,
		Cmd:     func(s []string) {},
	}
	s.help = append(s.help, c)
	sort.Sort(subCommands(s.help))
}

// Parse parses the command-line from os.Args[1:]. It then forks to the
// Action.Exec of the relevent sub command. Must be called after all
// SubCommands and HelpTopics are defined. Help, and version sub commands are
// automatically added
func (s SubCli) Parse(arguments []string) {
	if len(arguments) == 0 {
		s.cmdTree("")
	} else {
		s.cmdTree(arguments[0])
	}
}

func (s SubCli) cmdTree(cmd string) {
	// version builtin
	if cmd == "version" {
		s.cmdVersion()
		return
	}

	// help builtin
	if cmd == "help" || cmd == "" {
		if len(os.Args) > 3 {
			s.cmdHelp(os.Args[3:])
		} else {
			s.cmdHelp(nil)
		}
		return
	}

	// Check all the subs
	for _, v := range s.subs {
		if v.Command == cmd {
			err := v.Flags.Parse(os.Args[2:])
			if err != nil {
				s.println(err)
			}
			s.Args = v.Flags.Args()
			v.Cmd(v.Flags.Args())
			return
		}
	}

	// Not found
	s.printf("%s: unknown subcommand \"%s\"\n", s.Name, cmd)
	s.printf("Run \"%s help\" for usage.\n", s.Name)
}

func (s SubCli) cmdHelp(a []string) {
	// help with no arguments
	if a == nil {
		s.printHelp()
		return
	}

	// help with more than 1 arugment
	if len(a) > 1 {
		s.printf("usage: %s help command\n\n", s.Name)
		s.println("Too many arguments given.")
		return
	}

	// Check all the subs
	for _, v := range s.subs {
		if v.Command == a[0] {
			s.println(helpformated(v.Help))
			v.Flags.PrintDefaults()
			return
		}
	}

	// Check all the Help Topics
	for _, v := range s.help {
		if v.Command == a[0] {
			s.println(helpformated(v.Help))
			return
		}
	}

	// Nope its not here
	s.printf("Unkown help topic \"%s\". Run \"%s help\".\n", a[0], s.Name)
}

func (s SubCli) cmdVersion() {
	s.printf("%s version %s\n", s.Name, s.Version)
}

func (s SubCli) printHelp() {
	s.println(s.SDesc)
	s.print("\n")
	s.printf("Usage:\n\t%s command [arguments]\n\n", s.Name)
	s.print("The commands are:\n\n")
	for _, c := range s.subs {
		s.printf("\t%s\t\t%s\n", c.Command, c.SDesc)
	}
	s.printf(
		"\nUse \"%s help [command]\" for more information about a command\n\n",
		s.Name,
	)
	if len(s.help) == 0 {
		return
	}
	s.printf("Additional help topics:\n\n")
	for _, h := range s.help {
		s.printf("\t%s\t\t%s\n", h.Command, h.SDesc)
	}
	s.printf(
		"\nUse \"%s help [topic]\" for more information about that topic\n\n",
		s.Name,
	)
}

// print is a private function for redirecting output
func (s SubCli) print(a ...interface{}) {
	_, err := fmt.Fprint(s.Output, a...)
	if err != nil {
		println(err)
	}
}

func (s SubCli) println(a ...interface{}) {
	_, err := fmt.Fprintln(s.Output, a...)

	if err != nil {
		println(err)
	}
}

func (s SubCli) printf(f string, a ...interface{}) {
	_, err := fmt.Fprintf(s.Output, f, a...)
	if err != nil {
		println(err)
	}
}

// helpformated is a private function to print help messages in case I deside to
// preprocess the data (like say allow markdown for color or what not)
func helpformated(s string) string {
	return s
}

var commandLine = New(Program{Name: os.Args[0], Version: "", SDesc: ""})

func SetProgram(p Program) {
	commandLine.Program = p
}

// AddCmd adds a subcommand to the system
func AddCmd(c SubCommand) {
	commandLine.AddCmd(c)
}

// AddHelp adds a help topic to the system
func AddHelp(h HelpTopic) {
	commandLine.AddHelp(h)
}

// Parse parses the command-line from os.Args[1:]. Must be called after all
// subcommands, help topics, and flags are defined and before they are accessed.
func Parse() {
	commandLine.Parse(os.Args[1:])
}

func Args() []string {
	return commandLine.Args
}
