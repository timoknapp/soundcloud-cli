package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/urfave/cli/v2"

	"github.com/timoknapp/soundcloud-cli/pkg/download"
	"github.com/timoknapp/soundcloud-cli/pkg/flag"
	"github.com/timoknapp/soundcloud-cli/pkg/soundcloud"
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
			Action: func(c *cli.Context) error {
				dl := "download some stuff"
				fmt.Println(dl)
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
	// info()
	// commands()

	// err := app.Run(os.Args)
	// if err != nil {
	// 	log.Fatal(err)
	// }

	flags, err := flag.Read()
	if err != nil {
		log.Fatal(err)
	}
	if flags.Version {
		log.Println(BuildVersion)
		os.Exit(0)
	}
	if flags.TrackID == "" && flags.TrackURL == "" {
		log.Fatal("trackID or trackURL must be provided")
	}
	trackID := flags.TrackID
	if flags.TrackURL != "" {
		trackID, err = soundcloud.GetTrackIDByURL(flags.TrackURL)
	}
	track, err := soundcloud.GetTrack(trackID, flags.DownloadQuality)
	if err != nil {
		log.Fatal(err)
	}
	downloadStartTime := time.Now()
	absoluteDownloadPath, err := filepath.Abs(flags.DownloadPath)
	if err != nil {
		log.Printf("could not get absolute downloadPath")
		absoluteDownloadPath = "/"
	}
	log.Printf("start downloading track to %v", absoluteDownloadPath)
	err = download.Start(flags.DownloadPath, flags.DownloadQuality, track)
	if err != nil {
		log.Println(err)
	} else {
		downloadProgress := (float64(1) / float64(1)) * 100
		log.Printf("%3.f%% finished downloading: %v", downloadProgress, track.Title)
	}
	downloadDuration := time.Since(downloadStartTime)
	log.Printf("downloaded track in %s to %v", downloadDuration.Round(time.Second), absoluteDownloadPath)
}
