package tag

import (
	"io/ioutil"
	"net/http"

	"github.com/bogem/id3v2"
	"github.com/timoknapp/soundcloud-cli/pkg/soundcloud"
)

// AddTags adds some metadata to the file
func AddTags(file string, fileName string, track soundcloud.Track) error {
	var trackImage []byte
	if track.ArtworkURL != "" {
		client := http.Client{}
		req, err := http.NewRequest("GET", track.ArtworkURL, nil)
		if err != nil {
			return err
		}
		resp, err := client.Do(req)
		if err != nil {
			return err
		}
		trackImage, err = ioutil.ReadAll(resp.Body)
		if err != nil {
			return err
		}
	}
	tag, err := id3v2.Open(file, id3v2.Options{Parse: true})
	if err != nil {
		return err
	}
	defer tag.Close()
	tag.SetTitle(track.Title)
	tag.SetArtist(track.Artist.FullName)

	if trackImage != nil {
		tag.AddAttachedPicture(id3v2.PictureFrame{Encoding: id3v2.EncodingUTF8, MimeType: "image/jpeg", Picture: trackImage, PictureType: byte(3)})
	}
	if err = tag.Save(); err != nil {
		return err
	}
	return nil
}
