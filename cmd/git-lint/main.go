package main

import (
	"fmt"
	"github.com/troykinsella/git-lint"
	"github.com/urfave/cli"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
)

func loadConfig(path string) (*git_lint.Config, error) {

	var config git_lint.Config
	configData, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	err = yaml.Unmarshal(configData, &config)
	if err != nil {
		return nil, err
	}

	return &config, nil
}

func run(c *cli.Context) error {

	configFilePath := c.String("c")

	config, err := loadConfig(configFilePath)
	if err != nil {
		return err
	}

	gl, err := git_lint.New(config)
	if err != nil {
		return err
	}

	passed, err := gl.Run()
	if err != nil {
		return err
	}

	if !passed {
		os.Exit(1)
	}

	return nil
}

func newCliApp() *cli.App {

	app := cli.NewApp()
	app.Name = "git-lint"
	app.Version = git_lint.AppVersion
	//app.Usage = ""
	app.Author = "Troy Kinsella"

	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "c, config",
			Usage: "Yaml config file path",
			Value: ".git-lint.yml",
		},
	}

	app.Action = run

	return app
}

func main() {
	app := newCliApp()
	err := app.Run(os.Args)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
