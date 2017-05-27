package gomusicbrainz

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

var (
	recordingMBID = "2cfad207-3f55-4aec-8120-86cf66e34d59"
	workMBID      = "b38119e8-260f-372c-a1ca-653d02b5577c"
	ISRCID        = "USAT29900609"
	ISWCID        = "T-070.080.286-3"
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
	result, err := GetRecording(recordingMBID)

	if err != nil {
		t.Log(err.Error())
		t.Fail()
	}

	assert.True(t, err == nil, "Error should not be raised: ")
	assert.True(t, result != nil, "Result should not be nil")
	assert.Equal(t, recordingMBID, result.ID, "Expected recording MBID is different")
	fmt.Printf("RESULT: %+v\n", result)
}

func TestGetWorkWithValidMBID(t *testing.T) {
	result, err := GetWork(workMBID)

	if err != nil {
		t.Log(err.Error())
		t.Fail()
	}

	assert.True(t, err == nil, "Error should not be raised: ")
	assert.True(t, result != nil, "Result should not be nil")
	assert.Equal(t, workMBID, result.ID, "Expected recording MBID is different")
	fmt.Printf("RESULT: %+v\n", result)
}

func TestGetRecordingsByISRCSWithValidISRC(t *testing.T) {
	result, err := GetRecordingsByISRC(ISRCID)

	if err != nil {
		t.Log(err.Error())
		t.Fail()
	}

	assert.True(t, err == nil, "Error should not be raised: ")
	assert.True(t, result != nil, "Result should not be nil")
	assert.Equal(t, ISRCID, result.ISRCID, "Expected recording ISRC is different")
	fmt.Printf("RESULT: %+v\n", result)
}

func TestGetWorksByISWCSWithValidISWC(t *testing.T) {
	result, err := GetWorksByISWC(ISWCID)

	if err != nil {
		t.Log(err.Error())
		t.Fail()
	}

	assert.True(t, err == nil, "Error should not be raised: ")
	assert.True(t, result != nil, "Result should not be nil")
	assert.True(t, result.WorkCount > 0, "Result.WorkCount should not be > 0")
	assert.True(t, result.Works != nil, "Result.Works should not be nil")
	fmt.Printf("RESULT: %+v\n", result)
}

func Setup() {
	SetMusicBrainzConfig("gomusicbrainz-test", "1.0", "info@paperchain.io")
}
