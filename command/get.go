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

func init() {
	cmd := cli.Command{
		Name:                   "get",
		Usage:                  "show changes for a specific version(s)",
		ArgsUsage:              "<version> [version...]",
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
			&cli.BoolFlag{
				Name:    "hide-id",
				Aliases: []string{"v"},
				Value:   false,
				Usage:   "hide version IDs (i.e. aggregate the sections)",
			},
			&cli.BoolFlag{
				Name:    "single-line",
				Aliases: []string{"s"},
				Value:   false,
				Usage:   "do not break sentences into their own lines",
			},
		},
		Action: CommandGet,
		ShellComplete: func(ctx context.Context, cmd *cli.Command) {
			if shellCompleteFlag(ctx, cmd) {
				return
			}

			targets := make(map[string]struct{})
			for _, target := range cmd.Args().Slice() {
				target = normalizeVersion(strings.ToLower(target))
				targets[target] = struct{}{}
			}

			f, err := os.ReadFile(cmd.String("file"))
			if err != nil {
				return
			}

			var cl keepachangelog.Changelog
			if cmd.Bool("old-parser") {
				err = markdown.Unmarshal(f, &cl)
			} else {
				cl, err = keepachangelog.Parse(f)
			}
			if err != nil {
				fmt.Println(err)
				return
			}

			for _, version := range cl.Versions {
				if _, present := targets[strings.ToLower(version.ID)]; present {
					continue
				}
				fmt.Println(version.ID)
			}
		},
	}

	Register(&cmd)
}

func CommandGet(ctx context.Context, cmd *cli.Command) error {
	if !cmd.Args().Present() {
		return cli.ShowSubcommandHelp(cmd)
	}

	targets := []string{strings.ToLower(cmd.Args().First())}
	for _, version := range cmd.Args().Tail() {
		targets = append(targets, strings.ToLower(version))
	}

	normalizeVersions(targets)

	f, err := os.OpenFile(cmd.String("file"), os.O_RDONLY, 0644)
	if err != nil {
		return cli.Exit(fmt.Sprintf("Cannot open changelog file %s!", cmd.String("file")), 1)
	}
	defer f.Close()

	data, err := io.ReadAll(f)
	if err != nil {
		return cli.Exit(fmt.Sprintf("Cannot read changelog file %s!", cmd.String("file")), 1)
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

	if !cmd.Bool("hide-id") {
		for _, ver := range cl.Versions {
			if slices.Contains(targets, strings.ToLower(ver.ID)) {
				if !cmd.Bool("single-line") {
					fmt.Printf("%s", ver)
				} else {
					fmt.Printf("%s", ver.SingleLineString())
				}
			}
		}
	} else {
		sections := make(map[string]*keepachangelog.Section)

		for _, ver := range cl.Versions {
			if slices.Contains(targets, strings.ToLower(ver.ID)) {
				for _, section := range ver.Sections {
					if sec, ok := sections[section.Heading]; ok {
						sec.Changes = append(sec.Changes, section.Changes...)
					} else {
						sections[section.Heading] = &section
					}
				}
			}
		}

		for _, section := range sections {
			if !cmd.Bool("single-line") {
				fmt.Printf("%s", section)
			} else {
				fmt.Printf("%s", section.SingleLineString())
			}
		}
	}

	return nil
}
