package main

import (
	"context"
	"log"
	"os"

	"giiena.me/changelog/command"
	"github.com/urfave/cli/v3"
)

const helpTemplate = `USAGE:
    {{.Name}}{{if .Commands}} command [command options]{{end}} {{if .ArgsUsage}}{{.ArgsUsage}}{{else}}[arguments...]{{end}}

COMMANDS:
{{range .Commands}}    {{join .Names ", "}}{{ "\t"}}{{.Usage}}{{ "\n" }}{{end}}
`

func main() {
	cmd := cli.Command{
		Name:                          "changelog",
		Version:                       "0.4.0",
		Copyright:                     "(c) 2026 Giiena Astrella",
		HideVersion:                   true,
		CustomRootCommandHelpTemplate: helpTemplate,
		Commands:                      command.Retrieve(),
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:    "old-parser",
				Aliases: []string{"x"},
				Value:   false,
				Usage:   "use the old text parser implementation",
			},
		},
	}

	if err := cmd.Run(context.Background(), os.Args); err != nil {
		log.Fatal(err)
	}
}
