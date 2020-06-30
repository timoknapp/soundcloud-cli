package flag

import (
	"flag"
)

// Flag represents the available flags
type Flag struct {
	TrackURL        string
	TrackID         string
	DownloadPath    string
	DownloadQuality string
	Version         bool
}

// Read returns the set flags
func Read() (Flag, error) {
	trackURL := flag.String("trackURL", "", "the SoundCloud track url, either with the hostname or without")
	trackID := flag.String("trackID", "", "the SoundCloud track ID")
	downloadPath := flag.String("downloadPath", "download", "the path you want to download the tracks to")
	downloadQuality := flag.String("downloadQuality", "mp3", "the quality of the music files (mp3/ogg)")
	version := flag.Bool("version", false, "prints current version")

	flag.Parse()

	return Flag{
		TrackURL:        *trackURL,
		TrackID:         *trackID,
		DownloadPath:    *downloadPath,
		DownloadQuality: *downloadQuality,
		Version:         *version,
	}, nil
}
