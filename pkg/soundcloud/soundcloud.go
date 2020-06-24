package soundcloud

import (
	"errors"
	"io/ioutil"
	"net/http"
	"regexp"
	"strings"

	"github.com/anaskhan96/soup"
)

var soundCloudHost = "https://soundcloud.com"
var soundCloudAPIHost = "https://api-v2.soundcloud.com"

// GetClientID returns the client_id of soundcloud
func GetClientID() (string, error) {
	body, err := fetchHTTPBody(soundCloudHost + "/mt-marcy/cold-nights")
	if err != nil {
		return "", err
	}
	for _, prefix := range []string{"49", "48"} {
		scriptUrl, err := extractScriptURL(body, prefix)
		if err != nil {
			return "", err
		}
		script, err := fetchHTTPBody(scriptUrl)
		if err != nil {
			return "", err
		}
		var clientID = regexp.MustCompile(`client_id:+\"[a-zA-Z0-9]+\"`)
		matches := clientID.FindAllString(script, -1)
		if len(matches) == 0 {
			continue
		}
		for _, match := range matches {
			s := strings.TrimPrefix(match, "client_id:")
			return s, nil
		}
	}
	return "", err
}

func fetchHTTPBody(url string) (string, error) {
	res, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()
	html, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return "", err
	}
	return string(html), nil
}

func extractScriptURL(html string, prefix string) (string, error) {
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
