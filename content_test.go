package lytics

import (
	"github.com/bmizerany/assert"
	"github.com/jarcoal/httpmock"
	"github.com/lytics/go-lytics/mock"
	"testing"
)

func TestGetUserContentRecommendation(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()
	mock.RegisterContentMocks()

	client := NewLytics(mock.MockApiKey, nil, nil)
	recs, err := client.GetUserContentRecommendation("user_id", mock.MockUserID, "", 1, false)

	assert.Equal(t, err, nil)
	assert.Equal(t, len(recs), 1)
	assert.Equal(t, recs[0].Url, "www.testwebsite.com/some/url")
	assert.Equal(t, recs[0].Confidence, 0.8328074038765564)
}

func TestGetSegmentContentRecommendation(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()
	mock.RegisterContentMocks()

	client := NewLytics(mock.MockApiKey, nil, nil)
	recs, err := client.GetSegmentContentRecommendation(mock.MockSegmentID1, "", 1, false)

	assert.Equal(t, err, nil)
	assert.Equal(t, len(recs), 1)
	assert.Equal(t, recs[0].Url, "www.testwebsite.com/some/url")
	assert.Equal(t, recs[0].Confidence, 0.8328074038765564)
}
