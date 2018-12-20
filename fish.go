package subcli

import (
	"flag"
	"fmt"
	"io"
	"os"
	"strings"
)

func buildSub(rootCmd string, s SubCommand) string {
	var sb strings.Builder
	cmd := s.Command
	eFprintf(&sb, "# %s\n", cmd)
	eFprintf(&sb, `complete -c %s -n "__fish_use_subcommand`, rootCmd)
	eFprintf(&sb, " -a %s -d '%s'\n", cmd, s.SDesc)
	s.Flags.VisitAll(func(fl *flag.Flag) {
		eFprintf(&sb, `complete -c %s -n '__fish_seen_subcommand_from %s'`,
			rootCmd,
			cmd)
		eFprintf(&sb, " -a %s -d '%s'\n", fl.Name, fl.Usage)
	})
	return sb.String()
}

func eFprintf(w io.Writer, format string, a ...interface{}) {
	_, err := fmt.Fprintf(w, format, a...)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
