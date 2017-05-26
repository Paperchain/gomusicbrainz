package gomusicbrainz

import (
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"net/url"
	"time"

	"github.com/tidwall/gjson"
)

const (
	apiBaseURL       = "https://musicbrainz.org/ws/2/"
	releaseGroupPath = "release-group/"
	recordingPath    = "recording/"
	musicalWorkPath  = "work/"
	isrcPath         = "isrc/"
	iswcPath         = "iswc/"
)

var (
	AppName           string
	AppVersion        string
	ContactURLOrEmail string
	userAgentString   string
)

// Recording is
type Recording struct {
	Title          string   `json:"title"`
	Length         int      `json:"length"`
	ID             string   `json:"id"`
	Disambiguation string   `json:"disambiguation"`
	ISRCs          []string `json:"isrcs"`
	IsVideo        bool     `json:"video"`
}

func GetRecording(mbid string) (*Recording, error) {
	if mbid == "" {
		return nil, errors.New("MBID is empty")
	}

	u := apiBaseURL + recordingPath + mbid

	params := make(map[string]string)
	params["inc"] = "isrcs"
	params = addJSONParam(params)

	result, err := GET(u, params)
	if err != nil {
		return nil, err
	}

	var recording Recording
	gjson.Unmarshal(result, &recording)

	fmt.Printf("RECORDING: %+v\n", recording)

	return &recording, nil
}

// SetMusicBrainzConfig sets the configuration requirements
// Set these values before making any request
func SetMusicBrainzConfig(appName string, appVersion string, contactURLOrEmail string) {
	AppName = appName
	AppVersion = appVersion
	ContactURLOrEmail = contactURLOrEmail
}

// REQUEST makes a standard HTTP call
func REQUEST(method string, u string, body io.Reader) ([]byte, error) {
	err := validateConfig()
	if err != nil {
		return nil, err
	}
	var netTransport = &http.Transport{
		Dial: (&net.Dialer{
			Timeout: 5 * time.Second,
		}).Dial,
		TLSHandshakeTimeout: 5 * time.Second,
	}

	client := &http.Client{
		Timeout:   time.Second * 10,
		Transport: netTransport,
	}

	req, err := http.NewRequest(method, u, body)
	if err != nil {
		return nil, err
	}

	req.Header.Set("User-Agent", getUserAgentString())

	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	data, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	return data, nil
}

// GET makes an HTTP GET call to the specified uri
func GET(uri string, params map[string]string) ([]byte, error) {
	pu, err := url.Parse(uri)
	if err != nil {
		return nil, err
	}

	var u string
	if params != nil && len(params) > 0 {
		u = buildURLString(pu, params)
	} else {
		u = pu.String()
	}

	return REQUEST("GET", u, nil)
}

func validateConfig() error {
	if AppName == "" || AppVersion == "" || ContactURLOrEmail == "" {
		return errors.New("AppName, AppVersion or Contact parameters were not set! Make sure your app is calling SetMusicBrainzConfig() before making any requests")
	}

	return nil
}

func getUserAgentString() string {
	if userAgentString == "" {
		userAgentString = fmt.Sprintf("%s/%s (%s)", AppName, AppVersion, ContactURLOrEmail)
	}

	return userAgentString
}

func addJSONParam(params map[string]string) map[string]string {
	params["fmt"] = "json"
	return params
}

func encodeParams(u *url.URL, params map[string]string) string {
	q := u.Query()
	for k, v := range params {
		q.Add(k, v)
	}
	return q.Encode()
}

func buildURLString(u *url.URL, params map[string]string) string {
	u.RawQuery = encodeParams(u, params)
	return u.String()
}
