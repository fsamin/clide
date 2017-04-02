package main

import (
	"github.com/graymeta/stow"
	"github.com/graymeta/stow/s3"
	"github.com/urfave/cli"
)

type s3Provider struct{}

func (s3Provider) name() string {
	return "s3"
}

func (s3Provider) description() string {
	return "Amazon S3 Cloud Storage"
}

func (s3Provider) flags() []cli.Flag {
	return []cli.Flag{
		cli.StringFlag{
			Name:   "key",
			EnvVar: "AWS_ACCESS_KEY_ID",
		},
		cli.StringFlag{
			Name:   "secret",
			EnvVar: "AWS_SECRET_ACCESS_KEY",
		},
		cli.StringFlag{
			Name:   "region",
			EnvVar: "AWS_DEFAULT_REGION",
		},
	}
}

func (s3Provider) location(c *cli.Context) (stow.Location, error) {
	secret := c.String("secret")
	key := c.String("key")
	region := c.String("region")

	cfg := stow.ConfigMap{
		s3.ConfigAccessKeyID: key,
		s3.ConfigSecretKey:   secret,
		s3.ConfigRegion:      region,
	}

	return stow.Dial(s3.Kind, cfg)
}

func (p s3Provider) uploadFiles(c *cli.Context) error {
	l, err := p.location(c)
	if err != nil {
		return err
	}
	return uploadFile(l, c)
}

func (p s3Provider) downloadFiles(c *cli.Context) error {
	l, err := p.location(c)
	if err != nil {
		return err
	}
	return downloadFiles(l, c)
}
