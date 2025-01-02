package soundcloud

import (
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetTrack(t *testing.T) {
	trackID := "177318577"
	trackTitle := "Elderbrook x Andhim - How Many Times"
	quality := "mp3"
	track, _ := GetTrackById(trackID, quality)

	trackIDint, _ := strconv.Atoi(trackID)
	assert.Equal(t, trackIDint, track.ID)
	assert.Equal(t, trackTitle, track.Title)
}
