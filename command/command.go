package command

import "github.com/urfave/cli/v3"

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
