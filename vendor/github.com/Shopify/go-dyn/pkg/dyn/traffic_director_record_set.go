package dyn

import (
	"fmt"
	"time"
)

// TrafficDirectorRecordSet represents a Dyn Traffic Director Record Set.
type TrafficDirectorRecordSet struct {
	RecordSetID string
	Label       string
	RDataClass  string
	TTL         string
	Eligible    bool
	Automation  string
	MonitorID   string
	Records     []*TrafficDirectorRecord
}

type trafficDirectorRecordSetReference struct {
	RecordSetID string `json:"dsf_record_set_id,omitempty"`
}

type trafficDirectorRecordSetData struct {
	RecordSetID   string                      `json:"dsf_record_set_id"`
	Label         string                      `json:"label"`
	RDataClass    string                      `json:"rdata_class"`
	TTL           string                      `json:"ttl"`
	Status        string                      `json:"status"`
	LastMonitored string                      `json:"last_monitored"`
	MonitorID     string                      `json:"dsf_monitor_id"`
	PendingChange string                      `json:"pending_change"`
	Eligible      string                      `json:"eligible"`
	Automation    string                      `json:"automation"`
	Records       []trafficDirectorRecordData `json:"records"`
}

type TrafficDirectorRecordSetCURequest struct {
	ResponsePoolID string `json:"dsf_response_pool_id,omitempty"`
	Label          string `json:"label,omitempty"`
	RDataClass     string `json:"rdata_class"`
	TTL            string `json:"ttl,omitempty"`
	MonitorID      string `json:"dsf_monitor_id,omitempty"`
	Publish        string `json:"publish,omitempty"`
	Notes          string `json:"notes,omitempty"`
	Eligible       string `json:"eligible,omitempty"`
	Automation     string `json:"automation,omitempty"`
}

type trafficDirectorRecordSetDeleteRequest struct {
	Publish string `json:"publish,omitempty"`
	Notes   string `json:"notes,omitempty"`
}

type trafficDirectorRecordSetResponse struct {
	responseHeader
	trafficDirectorRecordSetData `json:"data"`
}

type trafficDirectorRecordSetAllResponse struct {
	responseHeader
	TrafficDirectorRecordSets []trafficDirectorRecordSetData `json:"data"`
}

type TrafficDirectorRecordSetOptionSetter func(*TrafficDirectorRecordSetCURequest)

func (tdrsd trafficDirectorRecordSetData) newTrafficDirectorRecordSet() *TrafficDirectorRecordSet {
	tdrs := TrafficDirectorRecordSet{
		RecordSetID: tdrsd.RecordSetID,
		Label:       tdrsd.Label,
		RDataClass:  tdrsd.RDataClass,
		TTL:         tdrsd.TTL,
		Eligible:    tdrsd.Eligible == "true",
		Automation:  tdrsd.Automation,
		MonitorID:   tdrsd.MonitorID,
		Records:     make([]*TrafficDirectorRecord, len(tdrsd.Records)),
	}

	for idx, record := range tdrsd.Records {
		tdrs.Records[idx] = record.newTrafficDirectorRecord()
	}

	return &tdrs
}

// CreateTrafficDirectorRecordSet creates a new instance of Traffic Director Record Set.
func (c *Client) CreateTrafficDirectorRecordSet(serviceID string, rDataClass string, options ...TrafficDirectorRecordSetOptionSetter) (*TrafficDirectorRecordSet, error) {
	req := TrafficDirectorRecordSetCURequest{
		RDataClass: rDataClass,
		Publish:    "Y",
	}

	for _, o := range options {
		o(&req)
	}

	var resp trafficDirectorRecordSetResponse

	// Given the issues we've found with the DynECT API, we need to allow that the first
	// request return a 'INVALID_REQUEST' for the POST method, and that either on the first
	// request or the following ones, we receive an 'OPERATION_FAILED' because 'This session
	// already has a job running'

	for try := 0; try < 10; try++ {
		err := c.post(fmt.Sprintf("DSFRecordSet/%s", serviceID), req, &resp)
		if err == nil {
			break
		}
		message, isResponseMessage := err.(responseMessage)
		if isResponseMessage {
			if try == 0 && message.ErrorCode == responseMessageInvalidRequest &&
				message.Info == "Resource does not support POST requests" {
				continue
			} else if message.ErrorCode == responseMessageOperationFailed &&
				message.Info == "token: This session already has a job running" {
				// Sleep the 5 seconds as recommended by the API specifications; we cannot
				// really use the JobID though as when we reached here during our tests, the
				// JobID would only show us that same message over and over as 'this' job failed
				time.Sleep(5 * time.Second)
				continue
			}
		}
		return nil, err
	}
	tdrs := resp.newTrafficDirectorRecordSet()

	return tdrs, nil
}

// UpdateTrafficDirectorRecordSet updates an instance of Traffic Director Record Set.
func (c *Client) UpdateTrafficDirectorRecordSet(serviceID string, recordSetID string, rDataClass string, options ...TrafficDirectorRecordSetOptionSetter) (*TrafficDirectorRecordSet, error) {
	req := TrafficDirectorRecordSetCURequest{
		RDataClass: rDataClass,
		Publish:    "Y",
	}

	for _, o := range options {
		o(&req)
	}

	var resp trafficDirectorRecordSetResponse

	if err := c.put(fmt.Sprintf("DSFRecordSet/%s/%s", serviceID, recordSetID), req, &resp); err != nil {
		return nil, err
	}

	tdrs := resp.newTrafficDirectorRecordSet()

	return tdrs, nil
}

// DeleteTrafficDirectorRecordSet deletes an instance of Traffic Director Record Set.
func (c *Client) DeleteTrafficDirectorRecordSet(serviceID string, recordSetID string) error {
	req := trafficDirectorRecordSetDeleteRequest{
		Publish: "Y",
	}

	if err := c.delete(fmt.Sprintf("DSFRecordSet/%s/%s", serviceID, recordSetID), req); err != nil {
		return err
	}

	return nil
}

// GetTrafficDirectorRecordSet returns an existing Traffic Director Record Set instance.
func (c *Client) GetTrafficDirectorRecordSet(serviceID string, recordSetID string) (*TrafficDirectorRecordSet, error) {
	var resp trafficDirectorRecordSetResponse

	if err := c.get(fmt.Sprintf("DSFRecordSet/%s/%s", serviceID, recordSetID), nil, &resp); err != nil {
		return nil, err
	}

	return resp.newTrafficDirectorRecordSet(), nil
}
