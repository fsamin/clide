package main

import (
	"github.com/graymeta/stow"
	"github.com/graymeta/stow/swift"
	"github.com/urfave/cli"
)

type swiftProvider struct{}

func (swiftProvider) name() string {
	return "swift"
}

func (swiftProvider) description() string {
	return "Openstack Swift Cloud Storage"
}

func (swiftProvider) flags() []cli.Flag {
	return []cli.Flag{
		cli.StringFlag{
			Name:   "username",
			EnvVar: "OS_USERNAME",
		},
		cli.StringFlag{
			Name:   "key",
			EnvVar: "OS_PASSWORD",
		},
		cli.StringFlag{
			Name:   "tenant",
			EnvVar: "OS_TENANT_NAME",
		},
		cli.StringFlag{
			Name:   "url",
			EnvVar: "OS_AUTH_URL",
		},
	}
}

func (swiftProvider) location(c *cli.Context) (stow.Location, error) {
	username := c.String("username")
	key := c.String("key")
	tenantName := c.String("tenant")
	tenantAuthUTL := c.String("url")

	cfg := stow.ConfigMap{
		swift.ConfigUsername:      username,
		swift.ConfigKey:           key,
		swift.ConfigTenantName:    tenantName,
		swift.ConfigTenantAuthURL: tenantAuthUTL,
	}

	return stow.Dial(swift.Kind, cfg)
}

func (p swiftProvider) uploadFile(c *cli.Context) error {
	l, err := p.location(c)
	if err != nil {
		return err
	}
	return uploadFile(l, c)
}
