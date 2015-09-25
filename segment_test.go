package lytics

import (
	"fmt"
	"github.com/bmizerany/assert"
	"github.com/jarcoal/httpmock"
	"github.com/lytics/go-lytics/mock"
	"testing"
)

func TestGetSegments(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()
	mock.RegisterSegmentMocks()

	client := NewLytics(mock.MockApiKey, nil)
	segs, err := client.GetSegments()
	assert.Equal(t, err, nil)
	assert.T(t, len(segs) > 1)
}

func TestGetSegment(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()
	mock.RegisterSegmentMocks()

	client := NewLytics(mock.MockApiKey, nil)
	seg, err := client.GetSegment(mock.MockSegmentID1)
	assert.Equal(t, err, nil)
	assert.T(t, seg.Id == mock.MockSegmentID1)
}

func TestGetSegmentSize(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()
	mock.RegisterSegmentMocks()

	client := NewLytics(mock.MockApiKey, nil)
	seg, err := client.GetSegmentSize(mock.MockSegmentID1)

	assert.Equal(t, err, nil)
	assert.T(t, seg.Id == mock.MockSegmentID1)
}

func TestGetSegmentSizes(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()
	mock.RegisterSegmentMocks()

	var segments []string

	client := NewLytics(mock.MockApiKey, nil)

	segments = []string{
		mock.MockSegmentID1,
	}

	seg, err := client.GetSegmentSizes(segments)
	assert.Equal(t, err, nil)
	assert.T(t, seg[0].Id == segments[0])

	segments = []string{
		mock.MockSegmentID1,
		mock.MockSegmentID2,
	}

	// params
	seg, err = client.GetSegmentSizes(segments)
	assert.Equal(t, err, nil)
	assert.T(t, len(seg) == 2)
}

func TestGetSegmentAttribution(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()
	mock.RegisterSegmentMocks()

	var (
		segments []string
		limit    int
	)

	segments = []string{
		mock.MockSegmentID1,
		mock.MockSegmentID2,
	}

	limit = 5

	client := NewLytics(mock.MockApiKey, nil)
	attr, err := client.GetSegmentAttribution(segments, limit)

	assert.Equal(t, err, nil)
	assert.T(t, len(attr[0].Metrics) == 5)
	assert.T(t, len(attr[1].Metrics) == 5)
}

func TestSegmentPager(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()
	mock.RegisterSegmentMocks()

	var (
		completed     bool
		countCalls    int
		countEntities int
	)

	client := NewLytics(mock.MockApiKey, nil)

	// create the segment scanner
	err := client.CreateScanner()
	assert.Equal(t, err, nil)

	// start the paging routine
	err = client.PageMembers(mock.MockSegmentID1)
	assert.Equal(t, err, nil)

	// handle processing the entities
PagingComplete:
	for {
		select {
		case entities := <-client.Scan.Loader:
			countCalls++

			for _, v := range entities {
				countEntities++
				assert.Equal(t, v["email"], fmt.Sprintf("email%d@email.com", countEntities))
			}

		case shutdown := <-client.Scan.Shutdown:
			if shutdown {
				completed = true
				break PagingComplete
			}
		}
	}
	assert.Equal(t, countCalls, 3)
	assert.Equal(t, completed, true)
	fmt.Printf("*** COMPLETED SCAN: Loaded %d batches and %d total entities", len(client.Scan.Batches), client.Scan.Total)
}
