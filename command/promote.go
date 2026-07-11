package command

import (
	"context"
	"fmt"
	"io"
	"os"
	"slices"
	"time"

	"giiena.me/changelog/keepachangelog"
	"giiena.me/changelog/markdown"
	"github.com/urfave/cli/v3"
)

func init() {
	cmd := cli.Command{
		Name:                   "promote",
		Usage:                  "promote unreleased",
		ArgsUsage:              "<version>",
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
		Action: CommandPromote,
	}

	Register(&cmd)
}

func CommandPromote(ctx context.Context, cmd *cli.Command) error {
	if !cmd.Args().Present() {
		return cli.ShowSubcommandHelp(cmd)
	}

	target := normalizeVersion(cmd.Args().First())

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

	if len(cl.Versions) < 1 {
		return cli.Exit("Cannot promote non-existing draft!", 4)
	}

	var unreleased *keepachangelog.Version
	unreleasedIndex := make([]int, 2)
	for i, ver := range cl.Versions {
		if ver.Unreleased {
			unreleased = &ver
			unreleasedIndex[0] = i
			unreleasedIndex[1] = i + 1
		}
	}

	if unreleased == nil {
		return cli.Exit("Cannot promote non-existing draft!", 4)
	}

	unreleased.ID = target
	unreleased.ReleaseDate = time.Now()
	unreleased.Unreleased = false

	cl.Versions = slices.Insert(cl.Versions, 0, *unreleased)
	cl.Versions = slices.Delete(cl.Versions, unreleasedIndex[0]+1, unreleasedIndex[1]+1)

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
