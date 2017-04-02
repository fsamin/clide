package main

import (
	"fmt"
	"os"

	"github.com/graymeta/stow"
	"github.com/urfave/cli"

	"github.com/fsamin/clide/lib"
)

func uploadFile(location stow.Location, c *cli.Context) error {
	if c.NArg() < 2 {
		return fmt.Errorf("Invalid usage")
	}

	files := append([]string{c.Args().First()}, c.Args().Tail()...)
	dest := files[len(files)-1:][0]
	files = files[:len(files)-1]

	uploadedFiles, err := clide.UploadFiles(location, dest, files, fmt.Printf)
	if err != nil {
		fmt.Printf("Error: %s\n", err)
		os.Exit(1)
	}

	for _, f := range uploadedFiles {
		fmt.Printf("%s\t%s\n", f.Filename, f.URL)
	}

	return nil
}
