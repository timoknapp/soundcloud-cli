package main

import (
	"errors"
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/timoknapp/soundcloud-cli/pkg/download"
	"github.com/timoknapp/soundcloud-cli/pkg/soundcloud"
	"github.com/urfave/cli/v2"
)

// BuildVersion contains SemVer of latest master merge
var BuildVersion = "development"

var app = cli.NewApp()

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
			Action: func(c *cli.Context) error {
				if c.Args().Len() == 0 {
					log.Fatal(errors.New("No track provided"))
				}
				path := c.String("path")
				quality := c.String("quality")
				input := c.Args().First()
				track, err := getTrackByInput(input, quality)
				if err != nil {
					log.Fatal(err)
				}
				download.Prepare(track, path, quality)
				return nil
			},
		},
		//TODO retrieve simple info for track + add search
		{
			Name:    "meta",
			Aliases: []string{"m"},
			Usage:   "Show metadata for a track",
			Action: func(c *cli.Context) error {
				if c.Args().Len() == 0 {
					log.Fatal(errors.New("No track provided"))
				}
				input := c.Args().First()
				track, err := getTrackByInput(input, "mp3")
				if err != nil {
					log.Fatal(err)
				}
				fmt.Printf("ID:\t\t" + strconv.Itoa(track.ID) + "\nArtist:\t\t" + track.Artist.FullName + "\nTitle:\t\t" + track.Title + "\nArtwork-URL:\t" + track.ArtworkURL)
				return nil
			},
		},
	}
	return cmds
}

func getTrackByInput(input string, quality string) (soundcloud.Track, error) {
	fmt.Println("Fetching meta data ...")
	trackID := input
	// check if input is trackID (int) otherwise expect url
	if _, err := strconv.Atoi(input); err != nil {
		trackID, err = soundcloud.GetTrackIDByURL(input)
		if err != nil {
			return soundcloud.Track{}, errors.New("Track could not be found")
		}
	}
	track, err := soundcloud.GetTrack(trackID, quality)
	if err != nil {
		return soundcloud.Track{}, err
	}
	return track, nil
}

func main() {
	info()
	commands()

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
