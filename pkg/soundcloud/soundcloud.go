package soundcloud

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/url"
	"regexp"
	"strings"

	"github.com/anaskhan96/soup"
	"github.com/fatih/color"
)

var soundCloudHost = "https://soundcloud.com"
var soundCloudAPIHost = "https://api-v2.soundcloud.com"
var soundCloudResolveApiUrl = "https://api-widget.soundcloud.com/resolve?"

// ClientID for SoundCloud
var ClientID = ""

// Track represents the track from SoundCloud
type Track struct {
	ID           int    `json:"id"`
	Title        string `json:"title"`
	Media        media  `json:"media"`
	ArtworkURL   string `json:"artwork_url"`
	PermalinkURL string `json:"permalink_url"`
	MediaURL     string
	Artist       artist `json:"artist"`
}

// SearchResult represents the the response when a search is performed
type SearchResult struct {
	Collection   []Track `json:"collection"`
	TotalResults int     `json:"total_results"`
}

type artist struct {
	FullName  string `json:"full_name"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Image     string `json:"avatar_url"`
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

// Resolve the given url: (return info about it).
func GetTrackInfoAPIUrl(urlx string, clientId string) string {
	v := url.Values{}

	// setting all the query params
	v.Set("url", urlx)
	v.Set("format", "json")
	v.Set("client_id", clientId)

	encodedUrl := soundCloudResolveApiUrl + v.Encode()

	return encodedUrl
}

// GetTrack returns the SoundCloud track by API
func GetTrackByUrl(url string, quality string) (Track, error) {
	clientID, err := getClientID()
	if err != nil {
		return Track{}, err
	}
	if clientID == "" {
		return Track{}, errors.New("soundcloud: clientID could not be retrieved")
	}
	apiUrl := GetTrackInfoAPIUrl(url, clientID)
	body, err := fetchHTTPBody(apiUrl)
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

// GetTrackById returns the SoundCloud track by ID
func GetTrackById(trackID string, quality string) (Track, error) {
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

// GetTrackIDByURL retrieves the trackID of a SoundCloud track by its URL or path
func GetTrackIDByURL(trackURL string) (string, error) {
	if strings.HasPrefix(trackURL, soundCloudHost) {
		// fetchBody with trackPath
		return getTrackID(trackURL)
	} else if strings.HasPrefix(trackURL, "/") {
		// fetchBody with host + trackpath
		return getTrackID(soundCloudHost + trackURL)
	} else {
		// fetchBody with host + "/" + trackpath
		return getTrackID(soundCloudHost + "/" + trackURL)
	}
}

func getTrackID(trackURL string) (string, error) {
	body, err := fetchHTTPBody(trackURL)
	if err != nil {
		return "", err
	}
	parsedHTML := soup.HTMLParse(string(body))
	linkElements := parsedHTML.FindAll("link")
	if len(linkElements) == 0 {
		return "", errors.New("soundcloud: trackID could not be parsed")
	}
	for _, element := range linkElements {
		if _, exists := element.Attrs()["rel"]; exists {
			// println(val)
			if val, exists := element.Attrs()["href"]; exists {
				// println(val)
				s := strings.Split(val, ":")
				if len(s) == 3 {
					if s[0] == "android-app" || s[0] == "ios-app" {
						trackID := s[2]
						// println(trackID)
						return trackID, nil
					}
				}
			}
		}
	}
	return "", err
}

func SearchTracks(input string, limit string) ([]Track, error) {
	_, err := getClientID()
	if err != nil {
		return []Track{}, err
	}
	body, err := fetchHTTPBody(soundCloudAPIHost + "/search/tracks?q=" + input + "&client_id=" + ClientID + "&limit=" + limit + "&offset=0")
	if err != nil {
		return []Track{}, err
	}
	var searchResult SearchResult
	if err = json.Unmarshal(body, &searchResult); err != nil {
		return []Track{}, err
	}
	return searchResult.Collection, nil
}

func GetTrackByInput(input string, quality string) (Track, error) {
	color.Cyan("\nFetching meta data ...")
	trackID := input
	track := Track{}
	var err error
	// check if input is url otherwise expect trackID
	if matched, _ := regexp.MatchString(`^https?://(www\.)?soundcloud\.com/`, input); matched {
		track, err = GetTrackByUrl(input, quality)
		if err != nil {
			return Track{}, errors.New("Track could not be found")
		}
	} else {
		track, err = GetTrackById(trackID, quality)
		if err != nil {
			return Track{}, err
		}

	}
	return track, nil
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
	lastElement := scriptElements[len(scriptElements)-1]
	if val, exists := lastElement.Attrs()["src"]; exists {
		scriptBody, err := fetchHTTPBody(val)
		if err != nil {
			return "", err
		}
		var re = regexp.MustCompile(",client_id:\"([^\"]*?.[^\"]*?)\"")
		matches := re.FindAllStringSubmatch(string(scriptBody), 1)
		if len(matches) == 0 {
			return "", errors.New("soundcloud: clientID could not be parsed")
		}
		ClientID = matches[0][1]
		return matches[0][1], nil
	}
	return "", err
}

func fetchHTTPBody(url string) ([]byte, error) {
	res, err := http.Get(url)
	if err != nil {
		return []byte(""), err
	}
	defer res.Body.Close()
	html, err := io.ReadAll(res.Body)
	if err != nil {
		return []byte(""), err
	}
	return html, nil
}

func getQualityByFileFormat(format string) string {
	if format == "mp3" {
		return "progressive"
	} else if format == "ogg" {
		return "hls"
	}
	return "progressive"
}
