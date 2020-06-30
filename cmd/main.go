package main

import (
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/timoknapp/soundcloud-cli/pkg/download"
	"github.com/timoknapp/soundcloud-cli/pkg/flag"
	"github.com/timoknapp/soundcloud-cli/pkg/soundcloud"
)

// BuildVersion contains SemVer of latest master merge
var BuildVersion = "development"

func main() {
	flags, err := flag.Read()
	if err != nil {
		log.Fatal(err)
	}
	if flags.Version {
		log.Println(BuildVersion)
		os.Exit(0)
	}
	if flags.TrackID == "" && flags.TrackPath == "" {
		log.Fatal("trackID or trackPath must be provided")
	}
	trackID := flags.TrackID
	if flags.TrackPath != "" {
		trackID, err = soundcloud.GetTrackIDByPath(flags.TrackPath)
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
