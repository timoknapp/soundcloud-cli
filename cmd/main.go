package main

import (
	"log"
	"os"

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
	if flags.TrackID == "" {
		log.Fatal("trackID must be provided")
	}
	track, err := soundcloud.GetStreamURL(flags.TrackID)
	if err != nil {
		log.Fatal(err)
	}
	log.Println(track.Title)
}
