package main

import (
	"fmt"
	"os"

	"github.com/Sirupsen/logrus"
	"github.com/joho/godotenv"
	"github.com/urfave/cli"
)

var build = "0" // build number set at compile-time

func main() {
	app := cli.NewApp()
	app.Name = "sftp cache plugin"
	app.Usage = "sftp cache plugin"
	app.Action = run
	app.Version = fmt.Sprintf("1.0.%s", build)
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:   "repo.name",
			Usage:  "repository full name",
			EnvVar: "DRONE_REPO",
		},
		cli.StringFlag{
			Name:   "repo.branch",
			Usage:  "repository default branch",
			EnvVar: "DRONE_REPO_BRANCH",
		},
		cli.StringFlag{
			Name:   "commit.branch",
			Value:  "master",
			Usage:  "repository branch",
			EnvVar: "DRONE_COMMIT_BRANCH",
		},
		cli.StringSliceFlag{
			Name:   "mount",
			Usage:  "cache directories",
			EnvVar: "PLUGIN_MOUNT",
		},
		cli.BoolFlag{
			Name:   "rebuild",
			Usage:  "rebuild the cache directories",
			EnvVar: "PLUGIN_REBUILD",
		},
		cli.BoolFlag{
			Name:   "restore",
			Usage:  "restore the cache directories",
			EnvVar: "PLUGIN_RESTORE",
		},
		cli.StringFlag{
			Name:   "sftp",
			Usage:  "sftp configuration",
			EnvVar: "PLUGIN_SFTP",
		},
		cli.StringFlag{
			Name:   "s3",
			Usage:  "s3 configuration",
			EnvVar: "PLUGIN_S3",
		},
		cli.StringFlag{
			Name:  "env-file",
			Usage: "source env file",
		},
	}

	if err := app.Run(os.Args); err != nil {
		logrus.Fatal(err)
	}
}

func run(c *cli.Context) error {
	if c.String("env-file") != "" {
		_ = godotenv.Load(c.String("env-file"))
	}

	plugin := Plugin{
		Rebuild:  c.Bool("rebuild"),
		Restore:  c.Bool("restore"),
		Mount:    c.StringSlice("mount"),
		Repo:     c.String("repo.name"),
		Default:  c.String("repo.branch"),
		Branch:   c.String("commit.branch"),
		Path:     c.String("path"),
		SFTP:     c.String("sftp"),
		S3:       c.String("s3"),
	}

	return plugin.Exec()
}
