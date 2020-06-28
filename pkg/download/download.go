package download

import (
	"errors"
	"io/ioutil"
	"os"
	"regexp"

	"github.com/timoknapp/soundcloud-cli/pkg/soundcloud"
	"github.com/timoknapp/soundcloud-cli/pkg/tag"
)

// Start will start the download process
func Start(downloadPath, format string, track soundcloud.Track) error {

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
