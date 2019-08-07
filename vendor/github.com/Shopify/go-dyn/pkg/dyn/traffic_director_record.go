package dyn

import (
	"fmt"
	"time"
)

// TrafficDirectorRecord represents a Dyn Traffic Director Record.
type TrafficDirectorRecord struct {
	RecordID        string
	MasterLine      string
	Label           string
	Weight          int
	Endpoints       []string
	EndpointUpCount int
	Eligible        bool
	Automation      string
}

type trafficDirectorRecordReference struct {
	RecordID string `json:"dsf_record_id,omitempty"`
}

type trafficDirectorRecordData struct {
	RecordID        string   `json:"dsf_record_id"`
	MasterLine      string   `json:"master_line"`
	Label           string   `json:"label"`
	Weight          int      `json:"weight"`
	Endpoints       []string `json:"endpoints"`
	EndpointUpCount int      `json:"endpoint_up_count"`
	Eligible        string   `json:"eligible"`
	Automation      string   `json:"automation"`
}

type TrafficDirectorRecordCURequest struct {
	MasterLine      string   `json:"master_line"`
	Label           string   `json:"label,omitempty"`
	Weight          int      `json:"weight,omitempty"`
	Endpoints       []string `json:"endpoints,omitempty"`
	EndpointUpCount int      `json:"endpoint_up_count,omitempty"`
	Eligible        string   `json:"eligible,omitempty"`
	Automation      string   `json:"automation,omitempty"`
	Publish         string   `json:"publish,omitempty"`
	Notes           string   `json:"notes,omitempty"`
}

type trafficDirectorRecordDeleteRequest struct {
	Publish string `json:"publish,omitempty"`
	Notes   string `json:"notes,omitempty"`
}

type trafficDirectorRecordResponse struct {
	responseHeader
	trafficDirectorRecordData `json:"data"`
}

type trafficDirectorRecordAllResponse struct {
	responseHeader
	TrafficDirectorRecords []trafficDirectorRecordData `json:"data"`
}

type TrafficDirectorRecordOptionSetter func(*TrafficDirectorRecordCURequest)

func (tdrd trafficDirectorRecordData) newTrafficDirectorRecord() *TrafficDirectorRecord {
	tdr := TrafficDirectorRecord{
		RecordID:        tdrd.RecordID,
		MasterLine:      tdrd.MasterLine,
		Label:           tdrd.Label,
		Weight:          tdrd.Weight,
		Endpoints:       tdrd.Endpoints,
		EndpointUpCount: tdrd.EndpointUpCount,
		Eligible:        tdrd.Eligible == "true",
		Automation:      tdrd.Automation,
	}

	return &tdr
}

// CreateTrafficDirectorRecord creates a new instance of Traffic Director Record.
func (c *Client) CreateTrafficDirectorRecord(serviceID string, recordSetID string, masterLine string, options ...TrafficDirectorRecordOptionSetter) (*TrafficDirectorRecord, error) {
	req := TrafficDirectorRecordCURequest{
		MasterLine: masterLine,
		Publish:    "Y",
	}

	for _, o := range options {
		o(&req)
	}

	var resp trafficDirectorRecordResponse

	for try := 0; try < 10; try++ {
		err := c.post(fmt.Sprintf("DSFRecord/%s/%s", serviceID, recordSetID), req, &resp)
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
	tdr := resp.newTrafficDirectorRecord()

	return tdr, nil
}

// UpdateTrafficDirectorRecord updates an instance of Traffic Director Record.
func (c *Client) UpdateTrafficDirectorRecord(serviceID string, recordID string, masterLine string, options ...TrafficDirectorRecordOptionSetter) (*TrafficDirectorRecord, error) {
	req := TrafficDirectorRecordCURequest{
		MasterLine: masterLine,
		Publish:    "Y",
	}

	for _, o := range options {
		o(&req)
	}

	var resp trafficDirectorRecordResponse

	if err := c.put(fmt.Sprintf("DSFRecord/%s/%s", serviceID, recordID), req, &resp); err != nil {
		return nil, err
	}

	tdr := resp.newTrafficDirectorRecord()

	return tdr, nil
}

// DeleteTrafficDirectorRecord deletes an instance of Traffic Director Record.
func (c *Client) DeleteTrafficDirectorRecord(serviceID string, recordID string) error {
	req := trafficDirectorRecordDeleteRequest{
		Publish: "Y",
	}

	if err := c.delete(fmt.Sprintf("DSFRecord/%s/%s", serviceID, recordID), req); err != nil {
		return err
	}

	return nil
}

// GetTrafficDirectorRecord returns an existing Traffic Director Record instance.
func (c *Client) GetTrafficDirectorRecord(serviceID string, recordID string) (*TrafficDirectorRecord, error) {
	var resp trafficDirectorRecordResponse

	if err := c.get(fmt.Sprintf("DSFRecord/%s/%s", serviceID, recordID), nil, &resp); err != nil {
		return nil, err
	}

	return resp.newTrafficDirectorRecord(), nil
}
