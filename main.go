package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/graymeta/stow"
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

func main() {
	app := cli.NewApp()
	app.Name = "stowcp"
	app.Usage = "stowcp <file 0> [file 1] ... [file n] <swift container>"
	app.Commands = []cli.Command{
		cli.Command{
			Name:   "swift",
			Flags:  swiftFlags,
			Action: swiftAction,
		},
	}
	app.Run(os.Args)
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
	containers, _, err := location.Containers("dest", stow.CursorStart, 100)
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
		fmt.Printf("%s\t => %s\n", name, url)

	}

	return nil
}
