package flag

import (
	"flag"
)

// Flag represents the available flags
type Flag struct {
	TrackID         string
	DownloadPath    string
	DownloadQuality string
	Version         bool
}

// Read returns the set flags
func Read() (Flag, error) {
	trackID := flag.String("trackID", "", "the SoundCloud track ID you want to download")
	downloadPath := flag.String("downloadPath", "download", "the path you want to download the tracks to")
	downloadQuality := flag.String("downloadQuality", "progressive", "the quality of the music files (progressive/hls)")
	version := flag.Bool("version", false, "prints current version")

	flag.Parse()

	return Flag{
		TrackID:         *trackID,
		DownloadPath:    *downloadPath,
		DownloadQuality: *downloadQuality,
		Version:         *version,
	}, nil
}
