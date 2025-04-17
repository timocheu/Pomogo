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
func NewBreak(duration int) *progressbar.ProgressBar {
	// Convert duration into seconds
	// We want to update the bar by every second, so we have to
	//   multiply it by second
	session := progressbar.NewOptions(duration*60,
		progressbar.OptionUseANSICodes(false),
		progressbar.OptionEnableColorCodes(true),
		progressbar.OptionSetWidth(10),
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
	// blank time 00:00:00
	timeSession := time.Time{}
	option := make(chan string)

	go func() {
		scanner := bufio.NewScanner(os.Stdin)
		for {
			if scanner.Scan() {
				cmd := strings.TrimSpace(strings.ToLower(scanner.Text()))
				option <- cmd
			}
		}
	}()

	for i := 0; i < duration*60; i++ {
		select {
		case cmd := <-option:
			switch cmd {
			case "q":
				// Display Option to resume
				fmt.Println("\n▄▄▄ [PAUSED] ▄▄▄")
				fmt.Println("\033[31m[R]\033[0m - Resume")

				for {
					if <-option == "r" {
						for i := 1; i <= 3; i++ {
							EraseLines(1)
							fmt.Println("Resuming in... ", i)

							time.Sleep(time.Second)
						}
						EraseLines(5)
						break
					} else {
						EraseLines(1)
					}
				}
			case "w":
				fmt.Println("Canceled...")
				return
			case "e":
				bar.Reset()
				timeSession = time.Time{}
				i = 0
				fmt.Println("Reset successful...")
			}
		default:
			// Add second to the timer
			timeSession = timeSession.Add(time.Second)

			minutes := int(timeSession.Minute())
			seconds := int(timeSession.Second())
			bar.Describe(fmt.Sprintf("[[%s]%02dm, %02ds[reset]] %s", color, minutes, seconds))
			bar.Add(1)
			time.Sleep(time.Second)
		}
	}
}

// Erase the lines in the terminal in upward motion
func EraseLines(n int) {
	for i := 0; i < n; i++ {
		fmt.Print("\033[A\033[2K")
	}
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

			fmt.Println("\033[31m[Q]\033[0m - Pause")
			fmt.Println("\033[34m[W]\033[0m - Cancel")
			fmt.Println("[E] - Reset")
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
	rest := NewBreak(int(cmd.Int("rest")))

	// Play
	Play(session, int(cmd.Int("session")), "red")
	Play(rest, int(cmd.Int("rest")), "blue")
}
