package main

import (
	"fmt"
	"os"

	"github.com/itsubaki/prstats/cmd/actions/jobs"
	"github.com/itsubaki/prstats/cmd/actions/runs"
	"github.com/urfave/cli/v2"
)

var date, hash, goversion string

func New(version string) *cli.App {
	app := cli.NewApp()
	app.Name = "prstats"
	app.Usage = "Github Productivity Stats"
	app.Version = version

	dir := cli.StringFlag{
		Name:    "dir",
		Aliases: []string{"d"},
		Value:   fmt.Sprintf("/var/tmp/%v", app.Name),
	}

	own := cli.StringFlag{
		Name:    "owner",
		Aliases: []string{"o"},
		Value:   "itsubaki",
	}

	repo := cli.StringFlag{
		Name:    "repo",
		Aliases: []string{"r"},
		Value:   "q",
	}

	format := cli.StringFlag{
		Name:    "format",
		Aliases: []string{"f"},
		Value:   "json",
		Usage:   "json, csv",
	}

	pat := cli.StringFlag{
		Name:    "pat",
		Aliases: []string{"t"},
		EnvVars: []string{"PAT"},
		Usage:   "Personal Access Token",
	}

	page := cli.IntFlag{
		Name:  "page",
		Value: 0,
	}

	perpage := cli.IntFlag{
		Name:  "perpage",
		Value: 1000,
	}

	runs := cli.Command{
		Name:    "runs",
		Aliases: []string{"r"},
		Subcommands: []*cli.Command{
			{
				Name:    "fetch",
				Aliases: []string{"f"},
				Action:  runs.Fetch,
				Flags: []cli.Flag{
					&dir,
					&own,
					&repo,
					&pat,
					&page,
					&perpage,
				},
			},
			{
				Name:    "list",
				Aliases: []string{"l"},
				Action:  runs.List,
				Flags: []cli.Flag{
					&dir,
					&own,
					&repo,
					&format,
				},
			},
			{
				Name:    "analyze",
				Aliases: []string{"a"},
				Action:  runs.Analyze,
				Flags: []cli.Flag{
					&dir,
					&own,
					&repo,
					&format,
					&cli.Int64Flag{
						Name:  "weeks",
						Value: 52,
					},
					&cli.BoolFlag{
						Name:  "excluding_weekends",
						Value: false,
					},
				},
			},
		},
	}

	jobs := cli.Command{
		Name:    "jobs",
		Aliases: []string{"j"},
		Subcommands: []*cli.Command{
			{
				Name:    "fetch",
				Aliases: []string{"f"},
				Action:  jobs.Fetch,
				Flags: []cli.Flag{
					&dir,
					&own,
					&repo,
					&pat,
					&page,
					&perpage,
				},
			},
			{
				Name:    "list",
				Aliases: []string{"l"},
				Action:  jobs.List,
				Flags: []cli.Flag{
					&dir,
					&own,
					&repo,
					&format,
				},
			},
			{
				Name:    "analyze",
				Aliases: []string{"a"},
				Action:  jobs.Analyze,
				Flags: []cli.Flag{
					&dir,
					&own,
					&repo,
					&format,
					&cli.Int64Flag{
						Name:  "weeks",
						Value: 52,
					},
					&cli.BoolFlag{
						Name:  "excluding_weekends",
						Value: false,
					},
				},
			},
		},
	}

	actions := cli.Command{
		Name:    "actions",
		Aliases: []string{"a"},
		Subcommands: []*cli.Command{
			&jobs,
			&runs,
		},
	}

	app.Commands = []*cli.Command{
		&actions,
	}

	return app
}

func main() {
	v := fmt.Sprintf("%s %s %s", date, hash, goversion)
	if err := New(v).Run(os.Args); err != nil {
		panic(err)
	}
}
