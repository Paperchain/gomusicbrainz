package gomusicbrainz

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

var (
	recordingMBID = "2cfad207-3f55-4aec-8120-86cf66e34d59"
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
	recording, err := GetRecording(recordingMBID)

	if err != nil {
		t.Log(err.Error())
		t.Fail()
	}

	assert.True(t, err == nil, "Error should not be raised: ")
	assert.True(t, recording != nil, "Recording should not be nil")
	assert.Equal(t, recordingMBID, recording.ID, "Expected recording MBID is different")
}

func Setup() {
	SetMusicBrainzConfig("gomusicbrainz-test", "1.0", "info@paperchain.io")
}
