package gomusicbrainz

import (
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"net/url"
	"regexp"
	"strings"
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
	artistPath       = "artist/"
	limit            = "10"
	offset           = "0"
	aliases          = "aliases"
)

var (
	// AppName indicates the Application Name
	AppName string

	// AppVersion indicates the version of the app ex:1.2.3
	AppVersion string

	// ContactURLOrEmail should be set to either an email address or a website for MusicBraniz to reach out
	ContactURLOrEmail string

	userAgentString string
	replacer        = strings.NewReplacer("-", "", ".", "", " ", "")
)

// GetRecording returns the Recording object for a given MusicBrainzID
func GetRecording(mbid string) (*Recording, error) {
	if mbid == "" {
		return nil, errors.New("MBID is empty")
	}

	u := apiBaseURL + recordingPath + mbid

	params := make(map[string]string)
	params["inc"] = "isrcs artist-credits"
	params = addJSONParam(params)

	result, err := GET(u, params)
	if err != nil {
		return nil, err
	}

	var recording Recording
	gjson.Unmarshal(result, &recording)

	return &recording, nil
}

// GetWork returns the Music Work object per the given MusicBrainzID
func GetWork(mbid string) (*Work, error) {
	if mbid == "" {
		return nil, errors.New("MBID is empty")
	}

	u := apiBaseURL + musicalWorkPath + mbid

	params := make(map[string]string)
	params["inc"] = aliases
	params = addJSONParam(params)

	result, err := GET(u, params)
	if err != nil {
		return nil, err
	}

	var work Work
	gjson.Unmarshal(result, &work)

	return &work, nil
}

// GetRecordingsByISRC returns Recording entities for a given isrc
func GetRecordingsByISRC(isrc string) (*ISRC, error) {
	if isrc == "" {
		return nil, errors.New("ISRC is empty")
	}

	if !isISRC(isrc) {
		return nil, errors.New("Not a valid ISRC")
	}

	u := apiBaseURL + isrcPath + isrc

	params := make(map[string]string)
	params["inc"] = "isrcs artist-credits"
	params = addJSONParam(params)
	params = addResultParams(params)

	result, err := GET(u, params)
	if err != nil {
		return nil, err
	}

	var i ISRC
	gjson.Unmarshal(result, &i)

	return &i, nil
}

// GetWorksByISWC returns Work entities for a given iswc
func GetWorksByISWC(iswc string) (*ISWC, error) {
	if iswc == "" {
		return nil, errors.New("ISWC is empty")
	}

	if !isISWC(iswc) {
		return nil, errors.New("Not a valid ISWC")
	}

	u := apiBaseURL + iswcPath + iswc

	params := make(map[string]string)
	params["inc"] = aliases
	params = addJSONParam(params)
	params = addResultParams(params)

	result, err := GET(u, params)
	if err != nil {
		return nil, err
	}

	var i ISWC
	gjson.Unmarshal(result, &i)

	return &i, nil
}

// GetArtist returns the Artist entity for the given mbid
func GetArtist(mbid string) (*Artist, error) {
	if mbid == "" {
		return nil, errors.New("MBID is empty")
	}

	u := apiBaseURL + artistPath + mbid

	params := make(map[string]string)
	params["inc"] = aliases
	params = addJSONParam(params)

	result, err := GET(u, params)
	if err != nil {
		return nil, err
	}

	var artist Artist
	gjson.Unmarshal(result, &artist)

	return &artist, nil
}

// SearchArtist returns the search results of the artists given
// the artistName and the optional entry of country
func SearchArtist(artistName string, country string) (*[]Artist, error) {
	a := strings.TrimSpace(artistName)
	if a == "" {
		return nil, errors.New("artistName is empty")
	}

	u := apiBaseURL + artistPath

	params := make(map[string]string)
	params["query"] = fmt.Sprintf("artist:%s", a)

	c := strings.TrimSpace(country)
	if c != "" {
		params["query"] += fmt.Sprintf(" AND country:%s", c)
	}

	params = addJSONParam(params)
	params = addResultParams(params)

	result, err := GET(u, params)
	if err != nil {
		return nil, err
	}

	var searchArtistResult SearchArtistResult
	var artists *[]Artist
	gjson.Unmarshal(result, &searchArtistResult)

	if &searchArtistResult != nil {
		artists = &searchArtistResult.Artists
	}

	return artists, nil
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
	fmt.Println(u)

	err := validateConfig()
	if err != nil {
		return nil, err
	}

	netTransport := &http.Transport{
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

	fmt.Printf(`
STATUS: %s
RATE LIMIT: %s
RATE LIMIT REMAINING: %s

`, res.Status, res.Header.Get("X-Ratelimit-Limit"), res.Header.Get("X-Ratelimit-Remaining"))

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

	return REQUEST(http.MethodGet, u, nil)
}

func retry(attempts int, sleep time.Duration, callback func(args ...interface{}) ([]byte, error)) ([]byte, error) {
	var err error
	for i := 0; i < attempts; i++ {
		res, e := callback()
		err = e
		if err == nil {
			return res, nil
		}

		if i > attempts {
			break
		}

		time.Sleep(sleep)

		fmt.Println("Retrying after error:", err)
	}
	return nil, fmt.Errorf("After %d attempts, last error: %s", attempts, err)
}

// isISRC checks for the input being a valid ISRC
func isISRC(input string) bool {
	r1 := regexp.MustCompile("^[A-Z]{2}[A-Z0-9]{3}[0-9]{2}[0-9]{5}$")
	//r2 := regexp.MustCompile("^[A-F]{2}[A-F0-9]{3}[0-9]{2}[0-9]{5}$")

	input = replacer.Replace(input)
	upperInput := strings.ToUpper(input)

	return r1.MatchString(upperInput) //|| r2.MatchString(upperInput)
}

// isISWC checks for the input being a valid UUID
func isISWC(input string) bool {
	r := regexp.MustCompile("^T[0-9]{3}[0-9]{3}[0-9]{3}[0-9]{1}$")

	input = replacer.Replace(input)
	upperInput := strings.ToUpper(input)

	return r.MatchString(upperInput)
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

func addResultParams(params map[string]string) map[string]string {
	params["limit"] = limit
	params["offset"] = offset
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
