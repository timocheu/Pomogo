package main

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/schollz/progressbar/v3"
	"github.com/urfave/cli/v3"
)

// The Duration is in "Minutes"
func NewSession(duration int) *progressbar.ProgressBar {
	// Convert duration into seconds
	// We want to update the bar by every second, so we have to
	//   multiply it by second
	session := progressbar.NewOptions(duration*60,
		progressbar.OptionUseANSICodes(false),
		progressbar.OptionEnableColorCodes(true),
		progressbar.OptionSetWidth(30),
		progressbar.OptionSetPredictTime(false),
		progressbar.OptionSetTheme(progressbar.Theme{
			Saucer:         "[red]\uEE04[reset]", // 
			SaucerHead:     "[red]\uEE04[reset]", // 
			SaucerPadding:  "[red]\uEE01[reset]", // 
			BarStart:       "[red]\uEE00[reset]", // 
			BarStartFilled: "[red]\uEE03[reset]", // 
			BarEnd:         "[red]\uEE02[reset]", // 
			BarEndFilled:   "[red]\uEE05[reset]", // 
		}),
		progressbar.OptionOnCompletion(func() {
			fmt.Println()
		}),
	)

	return session
}

// The Duration is in "Minutes"
func NewRest(duration int) *progressbar.ProgressBar {
	// Convert duration into seconds
	// We want to update the bar by every second, so we have to
	//   multiply it by second
	session := progressbar.NewOptions(duration*60,
		progressbar.OptionUseANSICodes(false),
		progressbar.OptionEnableColorCodes(true),
		progressbar.OptionSetWidth(15),
		progressbar.OptionSetPredictTime(false),
		progressbar.OptionSetTheme(progressbar.Theme{
			Saucer:         "[blue]\uEE04[reset]", // 
			SaucerHead:     "[blue]\uEE04[reset]", // 
			SaucerPadding:  "[blue]\uEE01[reset]", // 
			BarStart:       "[blue]\uEE00[reset]", // 
			BarStartFilled: "[blue]\uEE03[reset]", // 
			BarEnd:         "[blue]\uEE02[reset]", // 
			BarEndFilled:   "[blue]\uEE05[reset]", // 
		}),
		progressbar.OptionOnCompletion(func() {
			fmt.Println()
		}),
	)

	return session
}

// Duration is in "MINUTES"
func Play(bar *progressbar.ProgressBar, duration int, color string) {
	option := make(chan string)
	playing := make(chan bool)

	go func() {
		// blank time 00:00:00
		timeSession := time.Time{}

		for i := 0; i < duration*60; i++ {
			select {
			case cmd := <-option:
				switch cmd {
				case "r":
					bar.Reset()
					timeSession = time.Time{}
					i = 0
					fmt.Println("Reset successfully...")
				case "c":
					fmt.Println("Canceled ...")
					return
				case "p":
					// Display Option to resume
					fmt.Println("\n▄▄▄ [PAUSED] ▄▄▄")
					fmt.Println("\033[31m[Y]\033[0m - Resume")

					// Block the ticker
					// Ask for input
					for {
						scan := bufio.NewScanner(os.Stdin)
						if scan.Scan() {
							cmd := strings.TrimSpace(strings.ToLower(scan.Text()))
							if cmd == "y" {
								break
							} else {
								fmt.Println("Error: Invalid command")
							}
						}
					}
				}
			default:
				// Add second to the timer
				timeSession = timeSession.Add(time.Second)

				m := int(timeSession.Minute())
				s := int(timeSession.Second())
				bar.Describe(fmt.Sprintf("[[%s]%02dm, %02ds[reset]] Session", color, m, s))
				bar.Add(1)
			}
			time.Sleep(time.Second)
		}
	}()

	go func() {
		scanner := bufio.NewScanner(os.Stdin)
		for {
			if scanner.Scan() {
				cmd := strings.TrimSpace(strings.ToLower(scanner.Text()))
				option <- cmd
			}
		}
	}()
	<-playing
}

func main() {
	cmd := &cli.Command{
		Name:  "pomogo",
		Usage: "Pomodoro cli-app made in go.",
		Flags: []cli.Flag{
			&cli.IntFlag{
				Name:    "session",
				Aliases: []string{"s"},
				// Minutes
				Value: 25,
				Usage: "Time duration(mins) for each session",
			},
			&cli.IntFlag{
				Name:    "rest",
				Aliases: []string{"b"},
				// Minutes
				Value: 5,
				Usage: "Time duration(mins) for rest time",
			},
		},
		Action: func(context.Context, *cli.Command) error {
			fmt.Println("▗▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄")
			fmt.Println("▐▌ ▗▄▄▖  ▗▄▖ ▗▖  ▗▖ ▗▄▖  ▗▄▄▖ ▗▄▖  ")
			fmt.Println("▐▌ ▐▌ ▐▌▐▌ ▐▌▐▛▚▞▜▌▐▌ ▐▌▐▌   ▐▌ ▐▌ ")
			fmt.Println("▐▌ ▐▛▀▘ ▐▌ ▐▌▐▌  ▐▌▐▌ ▐▌▐▌▝▜▌▐▌ ▐▌ ")
			fmt.Println("▐▌ ▐▌   ▝▚▄▞▘▐▌  ▐▌▝▚▄▞▘▝▚▄▞▘▝▚▄▞▘ ")
			fmt.Println("▗▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄\n")

			fmt.Println("\033[31m[P]\033[0m - Pause")
			fmt.Println("\033[34m[C]\033[0m - Cancel")
			fmt.Println("[R] - Reset")
			fmt.Println("--- --- --- --- --- --- --- ")

			return nil
		},
	}

	if err := cmd.Run(context.Background(), os.Args); err != nil {
		log.Fatal(err)
	}

	// Stop if the program if its only asking help
	if len(os.Args) > 1 {
		arg := os.Args[1]
		if arg == "-h" || arg == "--help" {
			return
		}
	}

	// Create a new session && rest
	session := NewSession(int(cmd.Int("session")))
	rest := NewRest(int(cmd.Int("rest")))

	// Play
	Play(session, int(cmd.Int("session")), "red")
	Play(rest, int(cmd.Int("rest")), "blue")
}
