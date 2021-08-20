package auto_for_docker

import (
	"fmt"
	"github.com/urfave/cli/v2"
	"os"
)

var commands = []*cli.Command{
	{
		Name:    "build",
		Aliases: []string{"b"},
		Usage:   "build project",
		Action:  cmdBuild,
	},
	{
		Name:    "upload",
		Aliases: []string{"u"},
		Usage:   "upload image",
		Action:  cmdUpload,
	},
}

func cmdBuild(c *cli.Context) error {
	return nil
}

func cmdUpload(c *cli.Context) error {

	return nil
}

func run() int {
	app := cli.NewApp()
	app.Name = "AutoBuildImage"
	app.Usage = "Auto build image for frontend using Docker"

	app.Flags = []cli.Flag{
		&cli.StringFlag{
			Name:    "config",
			Aliases: []string{"c"},
			Usage:   "Load configuration from `FILE`",
		},
	}
	app.Commands = commands

	return msg(app.Run(os.Args))
}

func msg(err error) int {
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s: %v\n", os.Args[0], err)
		return 1
	}
	return 0
}

func main() {
	os.Exit(run())
}
