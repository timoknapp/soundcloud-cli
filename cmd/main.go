package main

import (
	"context"
	"errors"
	"log"
	"os"
	"strconv"

	"github.com/timoknapp/soundcloud-cli/pkg/download"
	"github.com/timoknapp/soundcloud-cli/pkg/sccli"
	"github.com/timoknapp/soundcloud-cli/pkg/soundcloud"
	"github.com/urfave/cli/v3"
)

// BuildVersion contains SemVer of latest master merge
var BuildVersion = "development"

var app = &cli.Command{}

func info() {
	app.Name = "SoundCloud CLI"
	app.Usage = "A simple CLI to interact with tracks on SoundCloud"
	app.Version = "1.0.0"
	app.Commands = commands()
}

func commands() []*cli.Command {
	cmds := []*cli.Command{
		{
			Name:    "download",
			Aliases: []string{"dl"},
			Usage:   "Download a track",
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:  "path",
					Value: "download",
					Usage: "Path where the files will be stored",
				},
				&cli.StringFlag{
					Name:  "quality",
					Value: "mp3",
					Usage: "Quality of the track",
				},
			},
			Action: func(ctx context.Context, c *cli.Command) error {
				if c.Args().Len() == 0 {
					log.Fatal(errors.New("No track provided"))
				}
				path := c.String("path")
				quality := c.String("quality")
				input := c.Args().First()
				track, err := soundcloud.GetTrackByInput(input, quality)
				if err != nil {
					log.Fatal(err)
				}
				download.Prepare(track, path, quality)
				return nil
			},
		},
		{
			Name:    "meta",
			Aliases: []string{"m"},
			Usage:   "Show metadata for a track",
			Action: func(ctx context.Context, c *cli.Command) error {
				if c.Args().Len() == 0 {
					log.Fatal(errors.New("No track provided"))
				}
				input := c.Args().First()
				track, err := soundcloud.GetTrackByInput(input, "mp3")
				if err != nil {
					log.Fatal(err)
				}
				sccli.PrintTable([]soundcloud.Track{track})
				return nil
			},
		},
		{
			Name:    "search",
			Aliases: []string{"ls"},
			Usage:   "Search for tracks",
			Flags: []cli.Flag{
				&cli.IntFlag{
					Name:  "limit",
					Value: 5,
					Usage: "Amount of search results",
				},
			},
			Action: func(ctx context.Context, c *cli.Command) error {
				if c.Args().Len() == 0 {
					log.Fatal(errors.New("No Search term provided"))
				}
				limit := strconv.Itoa(c.Int("limit"))
				input := c.Args().First()
				tracks, err := soundcloud.SearchTracks(input, limit)
				if err != nil {
					log.Fatal(err)
				}
				sccli.PrintTable(tracks)
				return nil
			},
		},
	}
	return cmds
}

func main() {
	info()
	commands()

	err := app.Run(context.Background(), os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
