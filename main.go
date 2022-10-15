package main

import (
	"os"
	"time"

	"github.com/shufo/ecs-fargate-oneshot/logs"
	"github.com/shufo/ecs-fargate-oneshot/tasks"

	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
)

// Version provides ecs-fargate version
var Version = "default"

func main() {
	log.SetOutput(os.Stdout)
	log.SetLevel(log.WarnLevel)

	err := run(os.Args)

	if err != nil {
		log.Fatal(err)
	}
}

func run(args []string) error {
	cli.VersionFlag = &cli.BoolFlag{
		Name: "version", Aliases: []string{"V"},
		Usage: "print only the version",
	}

	app := &cli.App{
		Name:     "ecs-fargate-oneshot",
		Version:  Version,
		Compiled: time.Now(),
		Usage:    "run oneshot task on ecs (fargate) with passed parameter",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     "cluster",
				Value:    "",
				Aliases:  []string{"c"},
				Usage:    "cluster name which task executes for",
				Required: true,
			},
			&cli.StringFlag{
				Name:     "service, s",
				Aliases:  []string{"s"},
				Usage:    "service where task executed in",
				Required: true,
			},
			&cli.BoolFlag{
				Name:     "verbose",
				Aliases:  []string{"v"},
				Usage:    "Show verbose logs",
				Required: false,
			},
			&cli.BoolFlag{
				Name:     "progress",
				Aliases:  []string{"p"},
				Usage:    "Show progress spinner",
				Required: false,
			},
		},
		Commands: []*cli.Command{
			{
				Name:  "run",
				Usage: "run the task with given args",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:     "task-definition",
						Aliases:  []string{"t"},
						Required: true,
					},
					&cli.Int64Flag{
						Name:     "cpu",
						Aliases:  []string{"C"},
						Required: false,
					},
					&cli.Int64Flag{
						Name:     "memory",
						Aliases:  []string{"m"},
						Required: false,
					},
					&cli.StringFlag{
						Name:     "container",
						Aliases:  []string{"n"},
						Required: true,
					},
					&cli.BoolFlag{
						Name:    "show-cloudwatch-logs",
						Aliases: []string{"l"},
						Value:   false,
					},
				},
				Action: func(c *cli.Context) error {
					if c.NArg() == 0 {
						log.Fatalln("ecs-fargate-oneshot: require commands to execute")
						log.Fatalln("e.g. ecs-fargate-oneshot run [options] echo 1")
						os.Exit(1)
					}

					if c.Bool("verbose") {
						log.SetLevel(log.InfoLevel)
					}

					tasks.RunTask(c)
					return nil
				},
			},
			{
				Name:  "logs",
				Usage: "show logs for tasks",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:     "container",
						Aliases:  []string{"n"},
						Value:    "",
						Required: true,
					},
					&cli.StringFlag{
						Name:     "task-id",
						Aliases:  []string{"t"},
						Value:    "",
						Required: false,
					},
				},
				Action: func(c *cli.Context) error {
					if c.Bool("verbose") {
						log.SetLevel(log.InfoLevel)
					}

					logs.RunShowLogsWithTaskId(c)
					return nil
				},
			},
		},
	}

	return app.Run(args)
}
