package dyn

import (
	"fmt"
	"time"
)

// TrafficDirectorResponsePool represents a Dyn Traffic Director Response Pool.
type TrafficDirectorResponsePool struct {
	ResponsePoolID string
	Label          string
	Eligible       bool
	Automation     string
	RecordSets     []*TrafficDirectorRecordSet
}

type trafficDirectorResponsePoolReference struct {
	ResponsePoolID string `json:"dsf_response_pool_id,omitempty"`
}

type trafficDirectorResponsePoolData struct {
	ResponsePoolID  string                              `json:"dsf_response_pool_id"`
	Label           string                              `json:"label"`
	Rulesets        []trafficDirectorRulesetData        `json:"rulesets"`
	RecordSetChains []trafficDirectorRecordSetChainData `json:"rs_chains"`
	Status          string                              `json:"status"`
	LastMonitored   string                              `json:"last_monitored"`
	PendingChange   string                              `json:"pending_change"`
	Eligible        string                              `json:"eligible"`
	Automation      string                              `json:"automation"`
}

type trafficDirectorRecordSetChainData struct {
	RecordSets []trafficDirectorRecordSetData `json:"record_sets"`
}

type trafficDirectorResponsePoolCreateRequest struct {
	Label      string `json:"label"`
	Publish    string `json:"publish,omitempty"`
	Notes      string `json:"notes,omitempty"`
	Eligible   string `json:"eligible,omitempty"`
	Automation string `json:"automation,omitempty"`
}

type trafficDirectorResponsePoolDeleteRequest struct {
	Publish string `json:"publish,omitempty"`
	Notes   string `json:"notes,omitempty"`
}

type trafficDirectorResponsePoolResponse struct {
	responseHeader
	trafficDirectorResponsePoolData `json:"data"`
}

type trafficDirectorResponsePoolAllResponse struct {
	responseHeader
	TrafficDirectorResponsePools []trafficDirectorResponsePoolData `json:"data"`
}

func (tdrpd trafficDirectorResponsePoolData) newTrafficDirectorResponsePool() *TrafficDirectorResponsePool {
	tdrp := TrafficDirectorResponsePool{
		ResponsePoolID: tdrpd.ResponsePoolID,
		Label:          tdrpd.Label,
		Eligible:       tdrpd.Eligible == "true",
		Automation:     tdrpd.Automation,
		RecordSets:     make([]*TrafficDirectorRecordSet, 0),
	}

	for _, recordSetChain := range tdrpd.RecordSetChains {
		for _, recordSet := range recordSetChain.RecordSets {
			tdrp.RecordSets = append(tdrp.RecordSets, recordSet.newTrafficDirectorRecordSet())
		}
	}

	return &tdrp
}

// CreateTrafficDirectorResponsePool creates a new instance of Traffic Director Response Pool.
func (c *Client) CreateTrafficDirectorResponsePool(serviceID string, label string) (*TrafficDirectorResponsePool, error) {
	req := trafficDirectorResponsePoolCreateRequest{
		Label:   label,
		Publish: "Y",
	}

	var resp trafficDirectorResponsePoolResponse

	// Given the issues we've found with the DynECT API, we need to allow that the first
	// request return a 'INVALID_REQUEST' for the POST method, and that either on the first
	// request or the following ones, we receive an 'OPERATION_FAILED' because 'This session
	// already has a job running'

	for try := 0; try < 10; try++ {
		err := c.post(fmt.Sprintf("DSFResponsePool/%s", serviceID), req, &resp)
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
	tdrp := resp.newTrafficDirectorResponsePool()

	return tdrp, nil
}

// UpdateTrafficDirectorResponsePool updates an instance of Traffic Director Response Pool.
func (c *Client) UpdateTrafficDirectorResponsePool(serviceID string, responsePoolID string, label string) (*TrafficDirectorResponsePool, error) {
	req := trafficDirectorResponsePoolCreateRequest{
		Label:   label,
		Publish: "Y",
	}

	var resp trafficDirectorResponsePoolResponse

	if err := c.put(fmt.Sprintf("DSFResponsePool/%s/%s", serviceID, responsePoolID), req, &resp); err != nil {
		return nil, err
	}

	tdrp := resp.newTrafficDirectorResponsePool()

	return tdrp, nil
}

// DeleteTrafficDirectorResponsePool deletes an instance of Traffic Director Response Pool.
func (c *Client) DeleteTrafficDirectorResponsePool(serviceID string, responsePoolID string) error {
	req := trafficDirectorResponsePoolDeleteRequest{
		Publish: "Y",
	}

	if err := c.delete(fmt.Sprintf("DSFResponsePool/%s/%s", serviceID, responsePoolID), req); err != nil {
		return err
	}

	return nil
}

// GetTrafficDirectorResponsePool returns an existing Traffic Director Response Pool instance.
func (c *Client) GetTrafficDirectorResponsePool(serviceID string, responsePoolID string) (*TrafficDirectorResponsePool, error) {
	var resp trafficDirectorResponsePoolResponse

	if err := c.get(fmt.Sprintf("DSFResponsePool/%s/%s", serviceID, responsePoolID), nil, &resp); err != nil {
		return nil, err
	}

	return resp.newTrafficDirectorResponsePool(), nil
}
