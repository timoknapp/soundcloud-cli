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

var ClientID = ""

// Track represents the track from SoundCloud
type Track struct {
	Title string `json:"title"`
	Media media  `json:"media"`
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
func GetTrack(trackID string) (Track, error) {
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

// GetMediaURL returns the mediaURL of a track
func GetMediaURL(track Track, quality string) (string, error) {
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
	body, err := fetchHTTPBody(soundCloudHost + "/mt-marcy/cold-nights")
	if err != nil {
		return "", err
	}
	for _, prefix := range []string{"49", "48"} {
		scriptURL, err := extractScriptURL(body, prefix)
		if err != nil {
			return "", err
		}
		script, err := fetchHTTPBody(scriptURL)
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
			return ClientID, nil
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
