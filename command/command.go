package command

import (
	"context"
	"fmt"
	"os"
	"slices"
	"strings"

	"github.com/urfave/cli/v3"
)

var cmds []*cli.Command

func init() {
	cmds = make([]*cli.Command, 0)
}

func Register(cmd *cli.Command) {
	cmds = append(cmds, cmd)
}

func Retrieve() []*cli.Command {
	return cmds
}

// normalizeVersion normalizes all elements of versions by removing leading `v` from each element.
//
// Changelog version heading should be without leading `v` (remember, `v` is short for `version`).
// But sometimes, release tags do have leading `v`.
// It may be difficult to strip the leading `v` when using this tool in CI context (e.g. in
// GitHub Actions), so we normalize version targets.
func normalizeVersions(versions []string) {
	for i, version := range versions {
		versions[i] = normalizeVersion(version)
	}
}

// normalizeVersion normalizes version by removing leading `v` from it.
//
// Changelog version heading should be without leading `v` (remember, `v` is short for `version`).
// But sometimes, release tags do have leading `v`.
// It may be difficult to strip the leading `v` when using this tool in CI context (e.g. in
// GitHub Actions), so we normalize version targets.
func normalizeVersion(version string) string {
	if version[0:1] == "v" {
		return version[1:]
	}
	return version
}

// shellCompleteFlag performs shell completions for flags in cmd and its parents.
// shellCompleteFlag returns true if it offers flag candidates.
func shellCompleteFlag(_ context.Context, cmd *cli.Command) bool {
	args := os.Args
	if len(args) > 0 && args[len(args)-1] == "--generate-shell-completion" {
		args = args[:len(args)-1]
	}

	if len(args) > 0 && strings.HasPrefix(args[len(args)-1], "-") {
		if slices.ContainsFunc(cmd.FlagNames(), func(flag string) bool {
			return strings.EqualFold(flag, strings.TrimLeft(args[len(args)-1], "-"))
		}) {
			return false
		}

		for _, flag := range getFlags(cmd) {
			if !slices.ContainsFunc(flag.Names(), func(name string) bool {
				return strings.EqualFold(name, strings.TrimLeft(args[len(args)-1], "-"))
			}) {
				continue
			}

			if flag.IsSet() {
				return false
			}

			switch flag := flag.(type) {
			case *cli.StringFlag:
				return flag.TakesFile
			default:
				return false
			}
		}

		for _, name := range getFlagNames(cmd) {
			if strings.HasPrefix(args[len(args)-1], "--") {
				if len(name) > 1 {
					fmt.Printf("--%s\n", name)
				} else {
					continue
				}
			} else if len(name) == 1 {
				fmt.Printf("-%s\n", name)
			}
		}
		return true
	}
	return false
}

// getFlagNames returns flag names in cmd and its parents.
// Flags are grouped by one-letter alias and full name, then lexicographically sorted.
func getFlagNames(cmd *cli.Command) []string {
	var flagNames []string

	for _, flag := range getFlags(cmd) {
		for _, name := range flag.Names() {
			flagNames = append(flagNames, name)
		}
	}

	slices.SortFunc(flagNames, func(a, b string) int {
		if len(a) == 1 && len(b) > 1 {
			return -1
		} else if len(a) > 1 && len(b) == 1 {
			return 1
		}
		return strings.Compare(a, b)
	})
	flagNames = slices.Compact(flagNames)

	return flagNames
}

func getFlags(cmd *cli.Command) []cli.Flag {
	var flags []cli.Flag
	commands := cmd.Lineage()
	slices.Reverse(commands)
	for _, cmd := range commands {
		flags = append(flags, cmd.Flags...)
	}
	return flags
}
