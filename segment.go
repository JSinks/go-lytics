package lytics

import (
	"fmt"
	"net/url"
	"strconv"
	"strings"
	"time"
)

const (
	segmentEndpoint               = "segment/:id"
	segmentListEndpoint           = "segment"
	segmentSizeEndpoint           = "segment/:id/sizes"
	segmentSizesEndpoint          = "segment/sizes"       // ids
	segmentAttributionEndpoint    = "segment/attribution" // ids
	segmentScanEndpoint           = "segment/:id/scan"
	adHocsegmentScanEndpoint      = "segment/scan"
	segmentCollectionListEndpoint = "segmentcollection"
	segmentCollectionEndpoint     = "segmentcollection/:id"
	segmentCreateEndpoint         = segmentListEndpoint
	segmentValidateEndpoint       = "segment/validate"
)

type Segment struct {
	Id            string    `json:"id"`
	Aid           int       `json:"aid"` // Deprecated; use AccountId
	AccountId     string    `json:"account_id"`
	ShortId       string    `json:"short_id,omitempty"`
	Name          string    `json:"name"`
	IsPublic      bool      `json:"is_public"`
	PublicName    string    `json:"public_name,omitempty"`
	SlugName      string    `json:"slug_name"`
	Description   string    `json:"description,omitempty"`
	SegKind       string    `json:"kind,omitempty"`
	Table         string    `json:"table,omitempty"`
	AuthorId      string    `json:"author_id"`
	Updated       time.Time `json:"updated" bson:"updated"`
	Created       time.Time `json:"created" bson:"created"`
	SegType       string    `json:"op"`
	Negate        bool      `json:"negate"`
	Tags          []string  `json:"tags"`
	Category      string    `json:category,omitempty`
	Invalid       bool      `json:"invalid"`
	InvalidReason string    `json:"invalid_reason"`
	Deleted       bool      `json:"deleted"`
	SaveHistory   bool      `json:"save_hist"`
	FilterQL      string    `json:"segment_ql,omitempty"`
}

type SegmentSize struct {
	Id       string  `json:"id"`
	Name     string  `json:"name"`
	SlugName string  `json:"slug_name"`
	Size     float64 `json:"size"`
}

type SegmentAttribution struct {
	Id      string                      `json:"id"`
	Metrics []SegmentAttributionMetrics `json:"metrics"`
}

type SegmentAttributionMetrics struct {
	Value   int64   `json:"value"`
	Ts      string  `json:"ts"`
	Anomaly float64 `json:"anomaly"`
}

type SegmentCollection struct {
	Aid           int               `json:"aid"`
	AccountId     string            `json:"account_id"`
	Id            string            `json:"id"`
	Name          string            `json:"name"`
	Slug          string            `json:"slug_name"`
	Description   string            `json:"description,omitempty"`
	Table         string            `json:"table,omitempty"`
	AuthorId      string            `json:"author_id"`
	Updated       time.Time         `json:"updated"`
	Created       time.Time         `json:"created"`
	Internal      bool              `json:"internal"`
	Collection    []*SegColRelation `json:"collection""`
	ParentSegment string            `json:"parent_segment"`
}

type SegColRelation struct {
	Id    string `json:"id"`
	Order int    `json:"order"`
}

type SegmentScanner struct {
	SegmentID string
	SegmentQl string
	Next      string
	Previous  string
	Loader    chan []Entity
	Shutdown  chan bool
	Total     int
	Batches   []int
}

// Created is a helper method to convert the timestamp into human readable format for metrics
func (s *SegmentAttributionMetrics) Created() (time.Time, error) {
	return parseLyticsTime(s.Ts)
}

// GetSegment returns the details for a single segment based on id
// https://www.getlytics.com/developers/rest-api#segment
func (l *Client) GetSegment(id string) (Segment, error) {
	res := ApiResp{}
	data := Segment{}

	// make the request
	err := l.Get(parseLyticsURL(segmentEndpoint, map[string]string{"id": id}), nil, nil, &res, &data)
	if err != nil {
		return data, err
	}

	return data, nil
}

// GetSegments returns a list of all segments for an account
// https://www.getlytics.com/developers/rest-api#segment-list
func (l *Client) GetSegments(string table) ([]Segment, error) {
	res := ApiResp{}
	data := []Segment{}
	params := url.Values{}

	params.Add("table", table)

	// make the request
	err := l.Get(segmentListEndpoint, params, &res, &data)
	if err != nil {
		return data, err
	}

	return data, nil
}

// GetSegmentSize returns the segment size information for a single segment
// https://www.getlytics.com/developers/rest-api#segment-sizes
func (l *Client) GetSegmentSize(id string) (SegmentSize, error) {
	res := ApiResp{}
	data := SegmentSize{}

	// make the request
	err := l.Get(parseLyticsURL(segmentSizeEndpoint, map[string]string{"id": id}), nil, nil, &res, &data)
	if err != nil {
		return data, err
	}

	return data, nil
}

// GetSegmentSizes returns the segment sizes for all segments on an account
// https://www.getlytics.com/developers/rest-api#segment-sizes
func (l *Client) GetSegmentSizes(segments []string) ([]SegmentSize, error) {
	params := url.Values{}
	res := ApiResp{}
	data := []SegmentSize{}

	// if we have specific segments to filter by add those to the params as comma separated string
	if len(segments) > 0 {
		params.Add("ids", strings.Join(segments, ","))
	}

	// make the request
	err := l.Get(segmentSizesEndpoint, params, nil, &res, &data)
	if err != nil {
		return data, err
	}

	return data, nil
}

// GetSegmentAttribution returns the attribution (change over time) for segments
// method accepts a string slice of 1 or more segments to query.
// NOT CURRENTLY DOCUMENTED
func (l *Client) GetSegmentAttribution(segments []string) ([]SegmentAttribution, error) {
	params := url.Values{}

	res := ApiResp{}
	data := []SegmentAttribution{}

	// if the request is for a specific set of segments add that as comma separated param
	if len(segments) > 0 {
		params.Add("ids", strings.Join(segments, ","))
	}

	// make the request
	err := l.Get(parseLyticsURL(segmentAttributionEndpoint, nil), params, nil, &res, &data)
	if err != nil {
		return data, err
	}

	return data, nil
}

// GetSegmentCollection returns a single collection of segments
// (a grouped/named lists of segments)
func (l *Client) GetSegmentCollection(id string) (SegmentCollection, error) {
	res := ApiResp{}
	data := SegmentCollection{}

	err := l.Get(parseLyticsURL(segmentCollectionEndpoint, map[string]string{"id": id}), nil, nil, &res, &data)
	if err != nil {
		return data, err
	}

	return data, nil
}

// GetSegmentCollectionList returns a list of all segment
// collections for an account
func (l *Client) GetSegmentCollectionList() ([]SegmentCollection, error) {
	res := ApiResp{}
	data := []SegmentCollection{}

	err := l.Get(parseLyticsURL(segmentCollectionListEndpoint, nil), nil, nil, &res, &data)
	if err != nil {
		return data, err
	}

	return data, nil
}

// Other Available Endpoints
// * POST    create segment
// * PUT     update segment
// * DELETE  remove segment

// **************************** START OF SEGMENT SCAN METHODS ****************************

// GetSegmentEntities returns a single page of entities (20max) for the given segment
// also returns the next value if there are more than 20 entities in the segment
// https://www.getlytics.com/developers/rest-api#segment-scan
func (l *Client) GetSegmentEntities(segment, next string, limit int) (interface{}, string, []Entity, error) {
	res := ApiResp{}
	data := []Entity{}
	params := url.Values{}

	params.Add("start", next)
	params.Add("limit", strconv.Itoa(limit))

	// make the request
	err := l.Get(parseLyticsURL(segmentScanEndpoint, map[string]string{"id": segment}), params, nil, &res, &data)
	if err != nil {
		return "", "", data, err
	}

	return res.Status, res.Next, data, nil
}

// GetAdHocSegmentEntities returns a single page of entities for the given Ad Hoc segment
// also returns the next value if there are more than 20 entities in the segment
// https://www.getlytics.com/developers/rest-api#segment-scan
func (l *Client) GetAdHocSegmentEntities(ql, next string, limit int) (interface{}, string, []Entity, error) {
	res := ApiResp{}
	data := []Entity{}
	params := url.Values{}

	params.Add("start", next)
	params.Add("limit", strconv.Itoa(limit))

	err := l.Get(adHocsegmentScanEndpoint, params, ql, params, &res, &data)
	if err != nil {
		return "", "", data, err
	}

	return res.Status, res.Next, data, nil
}

// CreateScanner generates a segment scanner so that we can process entities as they are loaded
// reduces the load as some segments can have hundreds of thousands of users and it is impractical
// to return a single result
func (l *Client) CreateScanner() error {
	scanner := SegmentScanner{}

	// create loader and shutdown
	loader := make(chan []Entity)
	shutdown := make(chan bool)

	// save to scanner
	scanner.Loader = loader
	scanner.Shutdown = shutdown

	// add the scanner to the client
	l.Scan = &scanner

	return nil
}

// LoadEntity does the heavy lifting when it comes to paging. It loops through all available pages
// and emits the batch of entities to be processed along the way. Also maintains a slice of batch
// counts and total count to help with debugging and reporting.
func (l *Client) LoadEntity(segType string) {
	var (
		entities []Entity
		err      error
		fails    int
		maxTries int
	)

	maxTries = 10

	// make calls for next batch of segments until we run out of next params
	for {
		switch segType {
		case "ql":
			_, l.Scan.Next, entities, err = l.GetAdHocSegmentEntities(l.Scan.SegmentQl, l.Scan.Next, 100)
			break

		default:
			_, l.Scan.Next, entities, err = l.GetSegmentEntities(l.Scan.SegmentID, l.Scan.Next, 100)
			break
		}

		if err != nil {
			fails++

			if fails > maxTries {
				panic(fmt.Sprintf("Failed to get entities, exceeded max tries(%d): %v", maxTries, err))
			}

			// if we fail and have not exceeded the limit, try again
			continue
		}

		// for logging add the batch details to the scanner
		l.Scan.Batches = append(l.Scan.Batches, len(entities))

		// for logging add the total entites returned to the scanner
		l.Scan.Total = l.Scan.Total + len(entities)

		// emit the entities to the loader for processing
		l.Scan.Loader <- entities

		// if there are no more pages we will have a blank next, just break and return
		if l.Scan.Next == "" {
			break
		}
	}

	// everything worked, shutdown safely
	l.Scan.Shutdown <- true
}

// PageMembers sets the segment on the master scanner and initiates the main go routine
// for paging
func (l *Client) PageMembers(segment string) error {
	// save the target segment on the scanner
	l.Scan.SegmentID = segment

	// fire up the go routine for paging entities
	go l.LoadEntity("id")
	return nil
}

// PageAdHocSegment sets the ad-hoc segment ql on the master scanner and initiates
// the main go routine for paging
func (l *Client) PageAdHocSegment(ql string) error {
	l.Scan.SegmentQl = ql

	go l.LoadEntity("ql")
	return nil
}

// CreateSegment creates a new segment from a Segment QL logic expression
// https://www.getlytics.com/developers/rest-api#segment
func (l *Client) CreateSegment(name, ql, slug string) (Segment, error) {
	res := ApiResp{}
	data := Segment{}

	payload := Segment{
		Name:     name,
		FilterQL: ql,
		Slug:     slug,
	}

	// make the request
	err := l.Post(segmentCreateEndpoint, nil, payload, &res, &data)
	if err != nil {
		return data, err
	}

	return res.Status, res.Next, data, nil
}

// ValidateSegment validates a single segment QL statement
// https://www.getlytics.com/developers/rest-api#segment-validate
func (l *Client) ValidateSegment(ql string) (bool, error) {
	res := ApiResp{}

	err := l.Post(segmentValidateEndpoint, nil, ql, &res, nil)
	if err != nil {
		return false, err
	}

	if res.Message == "success" {
		return true, nil
	}

	return false, nil
}
