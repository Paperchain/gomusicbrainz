package gomusicbrainz

import (
	"testing"

	"github.com/davecgh/go-spew/spew"
	"github.com/stretchr/testify/assert"
)

var (
	recordingMBID = "2cfad207-3f55-4aec-8120-86cf66e34d59"
	workMBID      = "b38119e8-260f-372c-a1ca-653d02b5577c"
	isrcID        = "USAT29900609"
	iswcID        = "T-070.080.286-3"
	artistMBID    = "678d88b2-87b0-403b-b63d-5da7465aecc3"
)

func TestGetAnythingWithoutConfigSetupWillResultInError(t *testing.T) {
	_, err := GetRecording("does_not_matter")
	t.Log(err)
	assert.True(t, err != nil, "Error should be raised!")
}

func TestWhenGettingRecordingWithoutInputWillResultInError(t *testing.T) {
	Setup()
	_, err := GetRecording("")
	t.Log(err)
	assert.True(t, err != nil, err, "Error should be raised!")
}

func TestGetRecordingWithValidMBID(t *testing.T) {
	Setup()
	result, err := GetRecording(recordingMBID)

	if err != nil {
		t.Log(err.Error())
		t.Fail()
	}

	assert.True(t, err == nil, "Error should not be raised: ")
	assert.True(t, result != nil, "Result should not be nil")
	assert.Equal(t, recordingMBID, result.ID, "Expected recording MBID is different")
	printResult(result)
}

func TestGetWorkWithValidMBID(t *testing.T) {
	Setup()
	result, err := GetWork(workMBID)

	if err != nil {
		t.Log(err.Error())
		t.Fail()
	}

	assert.True(t, err == nil, "Error should not be raised: ")
	assert.True(t, result != nil, "Result should not be nil")
	assert.Equal(t, workMBID, result.ID, "Expected recording MBID is different")
	printResult(result)
}

func TestGetRecordingsByISRCSWithValidISRC(t *testing.T) {
	Setup()
	result, err := GetRecordingsByISRC(isrcID)

	if err != nil {
		t.Log(err.Error())
		t.Fail()
	}

	assert.True(t, err == nil, "Error should not be raised: ")
	assert.True(t, result != nil, "Result should not be nil")
	assert.Equal(t, isrcID, result.ISRCID, "Expected recording ISRC is different")
	printResult(result)
}

func TestGetWorksByISWCSWithValidISWC(t *testing.T) {
	Setup()
	result, err := GetWorksByISWC(iswcID)

	if err != nil {
		t.Log(err.Error())
		t.Fail()
	}

	assert.True(t, err == nil, "Error should not be raised: ")
	assert.True(t, result != nil, "Result should not be nil")
	assert.True(t, result.WorkCount > 0, "Result.WorkCount should not be > 0")
	assert.True(t, result.Works != nil, "Result.Works should not be nil")
	printResult(result)
}

func TestSearchArtist(t *testing.T) {
	Setup()
	results, err := SearchArtist("Led Zeppelin", "GB")

	if err != nil {
		t.Log(err.Error())
		t.Fail()
	}

	assert.True(t, err == nil, "Error should not be raised: ")
	assert.True(t, results != nil, "Result should not be nil")
	printResult(results)
}

func TestGetArtistWithValidMBID(t *testing.T) {
	Setup()
	result, err := GetArtist(artistMBID)

	if err != nil {
		t.Log(err.Error())
		t.Fail()
	}

	assert.True(t, err == nil, "Error should not be raised: ")
	assert.True(t, result != nil, "Result should not be nil")
	assert.Equal(t, artistMBID, result.ID, "Expected recording MBID is different")
	printResult(result)
}

func printResult(r ...interface{}) {
	spew.Dump(r)
}

func Setup() {
	SetMusicBrainzConfig("gomusicbrainz-test", "1.0", "info@paperchain.io")
}
