package soundcloud

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"regexp"
	"strings"

	"github.com/anaskhan96/soup"
)

var soundCloudHost = "https://soundcloud.com"
var soundCloudAPIHost = "https://api-v2.soundcloud.com"

// ClientID for SoundCloud
var ClientID = ""

// Track represents the track from SoundCloud
type Track struct {
	Title    string `json:"title"`
	Media    media  `json:"media"`
	MediaURL string
}

type media struct {
	Transcodings []transcoding `json:"transcodings"`
}

type transcoding struct {
	URL      string `json:"url"`
	Preset   string `json:"preset"`
	Duration int    `json:"duration"`
	Format   format `json:"format"`
	Quality  string `json:"quality"`
}

type format struct {
	Protocol string `json:"protocol"`
	MimeType string `json:"mime_type"`
}

type mediaURL struct {
	URL string `json:"url"`
}

// GetTrack returns the streamURL of SoundCloud track
func GetTrack(trackID string, quality string) (Track, error) {
	_, err := getClientID()
	if err != nil {
		return Track{}, err
	}
	body, err := fetchHTTPBody(soundCloudAPIHost + "/tracks/" + trackID + "?client_id=" + ClientID)
	if err != nil {
		return Track{}, err
	}
	var trackResponse Track
	if err = json.Unmarshal(body, &trackResponse); err != nil {
		return Track{}, err
	}
	mediaURL, err := getMediaURL(trackResponse, getQualityByFileFormat(quality))
	if err != nil {
		return Track{}, err
	}
	trackResponse.MediaURL = mediaURL
	return trackResponse, nil
}

func getTranscodingByQuality(track Track, quality string) (transcoding, error) {
	errorTranscodingDoesNotExist := errors.New("soundcloud: desired quality does not exist")
	for _, transcodingType := range track.Media.Transcodings {
		if transcodingType.Format.Protocol == quality {
			return transcodingType, nil
		}
	}
	return transcoding{}, errorTranscodingDoesNotExist
}

func getMediaURL(track Track, quality string) (string, error) {
	_, err := getClientID()
	if err != nil {
		return "", err
	}
	trackInQuality, err := getTranscodingByQuality(track, quality)
	if err != nil {
		return "", err
	}
	body, err := fetchHTTPBody(trackInQuality.URL + "?client_id=" + ClientID)
	if err != nil {
		return "", err
	}
	var mediaResponse mediaURL
	if err = json.Unmarshal(body, &mediaResponse); err != nil {
		return "", err
	}
	return mediaResponse.URL, nil
}

func getClientID() (string, error) {
	if ClientID != "" {
		return ClientID, nil
	}
	body, err := fetchHTTPBody(soundCloudHost + "/mt-marcy/cold-nights")
	if err != nil {
		return "", err
	}
	parsedHTML := soup.HTMLParse(string(body))
	scriptElements := parsedHTML.FindAll("script")
	if len(scriptElements) == 0 {
		return "", errors.New("soundcloud: clientID could not be parsed")
	}
	for _, element := range scriptElements {
		if val, exists := element.Attrs()["src"]; exists {
			script, err := fetchHTTPBody(val)
			if err != nil {
				return "", err
			}
			var clientID = regexp.MustCompile(`client_id:+\"[a-zA-Z0-9]+\"`)
			matches := clientID.FindAllString(string(script), -1)
			if len(matches) == 0 {
				continue
			}
			for _, match := range matches {
				s := strings.TrimPrefix(match, "client_id:")
				t := strings.Replace(s, "\"", "", -1)
				ClientID = t
				// log.Println("ClientID: " + ClientID)
				return ClientID, nil
			}
		}
	}
	return "", err
}

func fetchHTTPBody(url string) ([]byte, error) {
	res, err := http.Get(url)
	if err != nil {
		return []byte(""), err
	}
	defer res.Body.Close()
	html, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return []byte(""), err
	}
	return html, nil
}

func extractScriptURL(html []byte, prefix string) (string, error) {
	errorClientIDCouldNotBeParsed := errors.New("soundcloud: clientID could not be parsed")
	parsedHTML := soup.HTMLParse(string(html))
	scriptElements := parsedHTML.FindAll("script")
	if len(scriptElements) == 0 {
		return "", errorClientIDCouldNotBeParsed
	}
	for _, element := range scriptElements {
		if val, exists := element.Attrs()["src"]; exists {
			if strings.Contains(val, prefix+"-") {
				return val, nil
			}
		}
	}
	return "", errorClientIDCouldNotBeParsed
}

func getQualityByFileFormat(format string) string {
	if format == "mp3" {
		return "progressive"
	} else if format == "ogg" {
		return "hls"
	}
	return "progressive"
}
