package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/graymeta/stow"
	"github.com/graymeta/stow/s3"
	"github.com/graymeta/stow/swift"
	"gopkg.in/urfave/cli.v1"
)

var swiftFlags = []cli.Flag{
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

var s3Flags = []cli.Flag{
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

func main() {
	app := cli.NewApp()
	app.Name = "stowcp"
	app.Usage = "stowcp <file 0> [file 1] ... [file n] <container/bucket/what else>"
	app.Commands = []cli.Command{
		cli.Command{
			Name:   "swift",
			Flags:  swiftFlags,
			Action: swiftAction,
		},
		cli.Command{
			Name:   "s3",
			Flags:  s3Flags,
			Action: s3Action,
		},
	}
	if err := app.Run(os.Args); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func s3Action(c *cli.Context) error {
	secret := c.String("secret")
	key := c.String("key")
	region := c.String("region")

	cfg := stow.ConfigMap{
		s3.ConfigAccessKeyID: key,
		s3.ConfigSecretKey:   secret,
		s3.ConfigRegion:      region,
	}

	location, err := stow.Dial(s3.Kind, cfg)
	if err != nil {
		fmt.Println("Error: ", err)
		return err
	}

	return uploadFile(location, c)
}

func swiftAction(c *cli.Context) error {
	username := c.String("username")
	key := c.String("key")
	tenantName := c.String("tenant")
	tenantAuthUTL := c.String("url")

	cfg := stow.ConfigMap{
		"username":        username,
		"key":             key,
		"tenant_name":     tenantName,
		"tenant_auth_url": tenantAuthUTL,
	}

	location, err := stow.Dial(swift.Kind, cfg)
	if err != nil {
		fmt.Println("Error: ", err)
		return err
	}

	return uploadFile(location, c)
}

func uploadFile(location stow.Location, c *cli.Context) error {
	if c.NArg() < 2 {
		return fmt.Errorf("Invalid usage")
	}

	files := append([]string{c.Args().First()}, c.Args().Tail()...)
	dest := files[len(files)-1:][0]
	files = files[:len(files)-1]

	fmt.Printf("Uploading %v to container %s\n", files, dest)

	var container stow.Container
	containers, _, err := location.Containers(dest, stow.CursorStart, 100)
	if err != nil {
		return err
	}
	for _, c := range containers {
		if c.Name() == dest {
			container = c
			break
		}
	}

	if container == nil {
		container, err = location.CreateContainer(dest)
		if err != nil {
			return err
		}
	}

	for i := range files {
		btes, err := ioutil.ReadFile(files[i])
		if err != nil {
			return err
		}
		name := filepath.Base(files[i])
		item, err := container.Put(name, bytes.NewBuffer(btes), int64(len(btes)), map[string]interface{}{})
		if err != nil {
			return err
		}

		url := strings.Replace(item.URL().String(), "swift://", "https://", 1)
		url = strings.Replace(url, "s3://", "", 1)
		fmt.Printf("%s\t => %s\n", name, url)

	}

	return nil
}
