package main

import (
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"time"

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
				trackID := input
				// check if input is trackID (int) otherwise expect url
				if _, err := strconv.Atoi(input); err != nil {
					trackID, err = soundcloud.GetTrackIDByURL(input)
					if err != nil {
						log.Fatal(errors.New("Track could not be found"))
					}
				}
				track, err := soundcloud.GetTrack(trackID, quality)
				if err != nil {
					log.Fatal(err)
				}
				downloadTrack(track, path, quality)
				return nil
			},
		},
		{
			Name:    "meta",
			Aliases: []string{"m"},
			Usage:   "Show metadata for a track",
			Action: func(c *cli.Context) error {
				m := "meta meta meta"
				fmt.Println(m)
				return nil
			},
		},
	}
	return cmds
}

func main() {
	// CLI based solution
	info()
	commands()

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}

	// Flag based solution
	// flags, err := flag.Read()
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// if flags.Version {
	// 	log.Println(BuildVersion)
	// 	os.Exit(0)
	// }
	// if flags.TrackID == "" && flags.TrackURL == "" {
	// 	log.Fatal("trackID or trackURL must be provided")
	// }
	// trackID := flags.TrackID
	// if flags.TrackURL != "" {
	// 	trackID, err = soundcloud.GetTrackIDByURL(flags.TrackURL)
	// }
	// track, err := soundcloud.GetTrack(trackID, flags.DownloadQuality)
	// if err != nil {
	// 	log.Fatal(err)
	// }
}

func downloadTrack(track soundcloud.Track, dlPath string, quality string) {
	downloadStartTime := time.Now()
	absoluteDownloadPath, err := filepath.Abs(dlPath)
	if err != nil {
		log.Printf("could not get absolute downloadPath")
		absoluteDownloadPath = "/"
	}
	log.Printf("start downloading track to %v", absoluteDownloadPath)
	err = download.Start(dlPath, quality, track)
	if err != nil {
		log.Println(err)
	} else {
		downloadProgress := (float64(1) / float64(1)) * 100
		log.Printf("%3.f%% finished downloading: %v", downloadProgress, track.Title)
	}
	downloadDuration := time.Since(downloadStartTime)
	log.Printf("downloaded track in %s to %v", downloadDuration.Round(time.Second), absoluteDownloadPath)
}
