package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/trezorg/plato/pkg/logger"
	"github.com/urfave/cli/v2"
)

func prepareCliArgs(version string) []string {

	defaultCommand := func(c *cli.Context) error {
		if c.NArg() == 0 {
			err := cli.ShowAppHelp(c)
			if err != nil {
				return fmt.Errorf("there are no urls to fetch, %w", err)
			}
			return fmt.Errorf("there are no urls to fetch")
		}
		return nil
	}

	app := cli.NewApp()
	app.Version = version
	app.HideHelp = false
	app.HideVersion = false
	app.Authors = []*cli.Author{{
		Name:  "Igor Nemilentsev",
		Email: "trezorg@gmail.com",
	}}
	app.Usage = "Urls fetcher"
	app.EnableBashCompletion = true
	app.ArgsUsage = "Multiple urls can be supplied"
	app.Action = defaultCommand

	err := app.Run(os.Args)
	if err != nil {
		if strings.Contains(err.Error(), "help requested") {
			os.Exit(0)
		}
		logger.Fatal(err)
	}
	return os.Args[1:]

}
