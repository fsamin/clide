package main

import (
	"fmt"
	"os"

	clide "github.com/fsamin/clide/lib"
	"github.com/graymeta/stow"
	"github.com/urfave/cli"
)

func downloadFiles(location stow.Location, c *cli.Context) error {
	if c.NArg() < 2 {
		return fmt.Errorf("Invalid usage")
	}

	containers := append([]string{c.Args().First()}, c.Args().Tail()...)
	dest := containers[len(containers)-1:][0]
	containers = containers[:len(containers)-1]

	downloadFiles, err := clide.DownloadFiles(location, dest, containers, fmt.Printf)
	if err != nil {
		fmt.Printf("Error: %s\n", err)
		os.Exit(1)
	}

	for _, f := range downloadFiles {
		fmt.Printf("%s\t%s\n", f.URL, f.Filename)
	}
	return nil
}
