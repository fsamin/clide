package main

import (
	"github.com/graymeta/stow"
	"github.com/urfave/cli"
)

type uploader interface {
	uploadFiles(c *cli.Context) error
}

type downloader interface {
	downloadFiles(c *cli.Context) error
}

type provider interface {
	name() string
	description() string
	flags() []cli.Flag
	location(c *cli.Context) (stow.Location, error)
}
