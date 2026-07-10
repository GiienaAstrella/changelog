package command

import (
	"context"
	"fmt"
	"strings"

	"github.com/urfave/cli/v3"
)

func init() {
	cmd := cli.Command{
		Name:   "version",
		Usage:  "print the app version",
		Action: CommandVersion,
	}

	Register(&cmd)
}

func CommandVersion(ctx context.Context, cmd *cli.Command) error {
	ShowBanner(ctx, cmd)

	return nil
}

func ShowBanner(ctx context.Context, cmd *cli.Command) {
	name := strings.ToUpper(cmd.Root().Name[0:1]) + cmd.Root().Name[1:]
	fmt.Printf("%s v%s\n%s\n\n", name, cmd.Root().Version, cmd.Root().Copyright)
}
