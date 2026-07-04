// Command yup-uniq is the CLI wrapper around github.com/gloo-foo/cmd-uniq.
package main

import (
	clix "github.com/gloo-foo/cli"
	command "github.com/gloo-foo/cmd-uniq"
	urf "github.com/urfave/cli/v3"
)

// version is the build version. It defaults to "dev" for local builds and is
// overridden at release time via the linker: -ldflags "-X main.version=<v>".
var version = "dev"

const (
	name           = "uniq"
	flagCount      = "count"
	flagRepeated   = "repeated"
	flagUnique     = "unique"
	flagIgnoreCase = "ignore-case"
)

// synopsis is the multi-line --help usage block; urfave/cli indents it three
// spaces, so the lines stay flush-left.
const synopsis = `uniq [OPTIONS] [FILE...]

Filter adjacent matching lines from FILE (or standard input),
writing one copy of each group to standard output.`

// spec declares the uniq wrapper: a file-or-stdin filter with uniq's flags.
var spec = clix.Spec{
	Name:     name,
	Summary:  "report or omit repeated lines",
	Synopsis: synopsis,
	Build:    build,
	Flags: []urf.Flag{
		&urf.BoolFlag{Name: flagCount, Aliases: []string{"c"}, Usage: "prefix lines by the number of occurrences"},
		&urf.BoolFlag{
			Name:    flagRepeated,
			Aliases: []string{"d"},
			Usage:   "only print duplicate lines, one for each group",
		},
		&urf.BoolFlag{Name: flagUnique, Aliases: []string{"u"}, Usage: "only print unique lines"},
		&urf.BoolFlag{
			Name:    flagIgnoreCase,
			Aliases: []string{"i"},
			Usage:   "ignore differences in case when comparing",
		},
	},
}

// build maps the invocation to uniq's pipeline: a file-or-stdin source into the
// uniq command configured by the flags.
func build(inv clix.Invocation) (clix.Source, clix.Command, error) {
	return clix.OperandsOrStdin(inv), command.Uniq(options(inv.Args)...), nil
}

// options folds the parsed flags into uniq's option values.
func options(c *urf.Command) []any {
	var opts []any
	if c.Bool(flagCount) {
		opts = append(opts, command.UniqCount)
	}
	if c.Bool(flagRepeated) {
		opts = append(opts, command.UniqDuplicatesOnly)
	}
	if c.Bool(flagUnique) {
		opts = append(opts, command.UniqUniqueOnly)
	}
	if c.Bool(flagIgnoreCase) {
		opts = append(opts, command.UniqIgnoreCase)
	}
	return opts
}

// runMain is an indirection seam so main's wiring is testable without spawning
// the process; a test swaps it and restores it.
var runMain = clix.Main

func main() { runMain(spec, version) }
