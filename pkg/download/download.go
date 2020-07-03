package download

import (
	"errors"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"time"

	"github.com/timoknapp/soundcloud-cli/pkg/soundcloud"
	"github.com/timoknapp/soundcloud-cli/pkg/tag"
)

func start(downloadPath, format string, track soundcloud.Track) error {

	if _, err := os.Stat(downloadPath); os.IsNotExist(err) {
		os.Mkdir(downloadPath, os.ModePerm)
	}

	specialChars, err := regexp.Compile("[^a-zA-Z0-9-\\. ()äöüß,&]+")
	if err != nil {
		return err
	}
	fileName := specialChars.ReplaceAllString(track.Title+"."+format, "")
	var fileToWrite string
	fileToWrite = downloadPath + string(os.PathSeparator) + fileName
	if _, err := os.Stat(fileToWrite); err == nil {
		return errors.New("soundcloud-cli: track already exists")
	}

	downloadedTrack, err := downloadFile(track.MediaURL)
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(fileToWrite, downloadedTrack, 0644)
	if err != nil {
		return err
	}
	//todo add more informations
	err = tag.AddTags(fileToWrite, fileName, track)
	if err != nil {
		return err
	}
	return nil
}

// Prepare use track and prepare the download
func Prepare(track soundcloud.Track, dlPath string, quality string) {
	downloadStartTime := time.Now()
	absoluteDownloadPath, err := filepath.Abs(dlPath)
	if err != nil {
		log.Printf("could not get absolute downloadPath")
		absoluteDownloadPath = "/"
	}
	log.Printf("start downloading track to %v", absoluteDownloadPath)
	err = start(dlPath, quality, track)
	if err != nil {
		log.Println(err)
	} else {
		downloadProgress := (float64(1) / float64(1)) * 100
		log.Printf("%3.f%% finished downloading: %v", downloadProgress, track.Title)
	}
	downloadDuration := time.Since(downloadStartTime)
	log.Printf("downloaded track in %s to %v", downloadDuration.Round(time.Second), absoluteDownloadPath)
}
