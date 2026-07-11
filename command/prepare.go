package command

import (
	"context"
	"fmt"
	"io"
	"os"
	"slices"
	"strings"

	"giiena.me/changelog/keepachangelog"
	"giiena.me/changelog/markdown"
	"github.com/urfave/cli/v3"
)

var defaultSections = []string{
	"Added",
	"Changed",
	"Deprecated",
	"Removed",
	"Fixed",
	"Security",
}

func init() {
	cmd := cli.Command{
		Name:                   "prepare",
		Usage:                  "prepare changelog for next cycle",
		ArgsUsage:              "[section...]",
		HideHelp:               true,
		UseShortOptionHandling: true,
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:      "file",
				TakesFile: true,
				Aliases:   []string{"f"},
				Value:     "CHANGELOG.md",
				Usage:     "changelog file",
			},
		},
		Action: CommandPrepare,
	}

	Register(&cmd)
}

func CommandPrepare(ctx context.Context, cmd *cli.Command) error {
	var targets []string

	if cmd.Args().Present() {
		targets = []string{cmd.Args().First()}
		targets = append(targets, cmd.Args().Tail()...)
	} else {
		targets = defaultSections
	}

	for i, target := range targets {
		if len(target) > 1 {
			targets[i] = strings.ToUpper(target[:1]) + strings.ToLower(target[1:])
		} else {
			targets[i] = strings.ToUpper(target)
		}
	}

	f, err := os.Open(cmd.String("file"))
	if err != nil {
		return cli.Exit(fmt.Sprintf("Cannot open changelog file %s!", cmd.String("file")), 1)
	}

	data, err := io.ReadAll(f)
	if err != nil {
		return cli.Exit(fmt.Sprintf("Cannot read changelog file %s!", cmd.String("file")), 1)
	}

	err = f.Close()
	if err != nil {
		return cli.Exit(fmt.Sprintf("Cannot close changelog file %s!", cmd.String("file")), 1)
	}

	var cl keepachangelog.Changelog
	if cmd.Bool("old-parser") {
		err = markdown.Unmarshal(data, &cl)
	} else {
		cl, err = keepachangelog.Parse(data)
	}
	if err != nil {
		return cli.Exit(fmt.Sprintf("Cannot parse changelog file %s!\n%v", cmd.String("file"), err), 2)
	}

	var unreleased *keepachangelog.Version
	unreleasedIndex := []int{-1, -1}

	for i, ver := range cl.Versions {
		if ver.Unreleased {
			unreleased = &ver
			unreleasedIndex[0] = i
			unreleasedIndex[1] = i + 1
		}
	}

	if unreleased == nil {
		unreleased = &keepachangelog.Version{
			ID:         "UNRELEASED",
			Unreleased: true,
			Sections:   make([]keepachangelog.Section, 0),
		}
	}

	for _, heading := range targets {
		if !slices.ContainsFunc(unreleased.Sections, func(sec keepachangelog.Section) bool {
			return strings.EqualFold(sec.Heading, heading)
		}) {
			unreleased.Sections = append(unreleased.Sections, keepachangelog.Section{
				Heading: heading,
			})
		}
	}

	cl.Versions = slices.Insert(cl.Versions, 0, *unreleased)
	if unreleasedIndex[0] > -1 {
		cl.Versions = slices.Delete(cl.Versions, unreleasedIndex[0]+1, unreleasedIndex[1]+1)
	}

	f, err = os.OpenFile(cmd.String("file"), os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0644)
	if err != nil {
		return cli.Exit(fmt.Sprintf("Cannot open changelog file %s!", cmd.String("file")), 1)
	}
	defer f.Close()

	md, err := markdown.Marshal(cl)
	if err != nil {
		return cli.Exit(fmt.Sprintf("Cannot encode changelog: %v", err), 3)
	}

	_, err = f.Write(md)
	if err != nil {
		return cli.Exit(fmt.Sprintf("Cannot write changelog: %v", err), 3)
	}

	return nil
}
