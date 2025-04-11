package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/urfave/cli/v3"
)

func main() {
	cmd := &cli.Command{
		Name:  "pomogo",
		Usage: "Pomodoro cli-app made in go.",
		Flags: []cli.Flag{
			&cli.IntFlag{
				Name:    "session",
				Aliases: []string{"s"},
				Value:   25,
				Usage:   "Time duration(mins) for each session",
			},
			&cli.IntFlag{
				Name:    "break",
				Aliases: []string{"b"},
				Value:   5,
				Usage:   "Time duration(mins) for break time",
			},
		},
		Action: func(context.Context, *cli.Command) error {
			fmt.Println("HALLO")
			return nil
		},
	}

	if err := cmd.Run(context.Background(), os.Args); err != nil {
		log.Fatal(err)
	}
}
