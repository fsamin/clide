package main

import (
	"github.com/graymeta/stow"
	"github.com/urfave/cli"
)

type uploader interface {
	uploadFile(c *cli.Context) error
}

type provider interface {
	name() string
	description() string
	flags() []cli.Flag
	location(c *cli.Context) (stow.Location, error)
}
