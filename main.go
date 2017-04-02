package main

import (
	"fmt"
	"os"

	"github.com/urfave/cli"
)

var providers = []provider{
	swiftProvider{},
	s3Provider{},
}

func main() {
	app := cli.NewApp()
	app.Name = "clide"
	app.Version = version
	app.Author = "François SAMIN"
	app.Usage = "Cloud Storage CLI"
	app.Description = "Cloud Storage files management CLI"

	for _, p := range providers {
		cmd := cli.Command{
			Name:        p.name(),
			Description: p.description(),
		}

		if pu, ok := p.(uploader); ok {
			cmd.Subcommands = append(cmd.Subcommands, cli.Command{
				Name:        "upload",
				ShortName:   "up",
				Usage:       "clide " + p.name() + " upload <file 0> [file 1] ... [file n] <destination>",
				Description: "Upload files",
				Flags:       p.flags(),
				Action:      pu.uploadFiles,
			})
		}

		if pu, ok := p.(downloader); ok {
			cmd.Subcommands = append(cmd.Subcommands, cli.Command{
				Name:        "download",
				ShortName:   "dl",
				Usage:       "clide " + p.name() + " download <container 0> [container 1] ... [container n] <destination>",
				Description: "Download all files from several containers",
				Flags:       p.flags(),
				Action:      pu.downloadFiles,
			})
		}

		app.Commands = append(app.Commands, cmd)
	}

	if err := app.Run(os.Args); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
