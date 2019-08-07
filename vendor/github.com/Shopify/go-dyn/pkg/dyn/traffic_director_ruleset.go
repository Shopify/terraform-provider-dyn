package dyn

import (
	"fmt"
	"strconv"
)

// TrafficDirectorRuleset represents a Dyn Traffic Director Response Pool.
type TrafficDirectorRuleset struct {
	RulesetID     string
	Label         string
	ResponsePools []*TrafficDirectorResponsePool
	CriteriaType  string
	Criteria      trafficDirectorRulesetCriteria
	Ordering      int
}

type trafficDirectorRulesetCriteriaGeolocation struct {
	Regions   []string `json:"region,omitempty"`
	Countries []string `json:"country,omitempty"`
	Provinces []string `json:"province,omitempty"`
}

type trafficDirectorRulesetCriteria struct {
	Geolocation trafficDirectorRulesetCriteriaGeolocation `json:"geoip"`
}

type trafficDirectorRulesetData struct {
	RulesetID     string                            `json:"dsf_ruleset_id"`
	Label         string                            `json:"label"`
	ResponsePools []trafficDirectorResponsePoolData `json:"response_pools"`
	CriteriaType  string                            `json:"criteria_type"`
	Criteria      trafficDirectorRulesetCriteria    `json:"criteria"`
	Ordering      string                            `json:"ordering"`
}

type TrafficDirectorRulesetCURequest struct {
	Label         string                                 `json:"label"`
	Publish       string                                 `json:"publish,omitempty"`
	ResponsePools []trafficDirectorResponsePoolReference `json:"response_pools,omitempty"`
	CriteriaType  string                                 `json:"criteria_type,omitempty"`
	Criteria      trafficDirectorRulesetCriteria         `json:"criteria,omitempty"`
	Notes         string                                 `json:"notes,omitempty"`
	Ordering      int                                    `json:"ordering,omitempty"`
}

func (tdrcq *TrafficDirectorRulesetCURequest) SetResponsePools(response_pools []string) {
	for _, rrid := range response_pools {
		tdrcq.ResponsePools = append(tdrcq.ResponsePools, trafficDirectorResponsePoolReference{
			ResponsePoolID: rrid,
		})
	}
}

func (tdrcq *TrafficDirectorRulesetCURequest) SetGeolocation(geolocations map[string][]string) {
	if len(geolocations) > 0 {
		tdrcq.CriteriaType = "geoip"
		tdrcq.Criteria = trafficDirectorRulesetCriteria{
			Geolocation: trafficDirectorRulesetCriteriaGeolocation{
				Regions:   geolocations["region"],
				Countries: geolocations["country"],
				Provinces: geolocations["province"],
			},
		}
	}
}

type trafficDirectorRulesetDeleteRequest struct {
	Publish string `json:"publish,omitempty"`
	Notes   string `json:"notes,omitempty"`
}

type trafficDirectorRulesetResponse struct {
	responseHeader
	trafficDirectorRulesetData `json:"data"`
}

type trafficDirectorRulesetAllResponse struct {
	responseHeader
	TrafficDirectorRulesets []trafficDirectorRulesetData `json:"data"`
}

type TrafficDirectorRulesetOptionSetter func(*TrafficDirectorRulesetCURequest)

func (tdrsd trafficDirectorRulesetData) newTrafficDirectorRuleset() *TrafficDirectorRuleset {
	ordering, _ := strconv.Atoi(tdrsd.Ordering)
	tdrs := TrafficDirectorRuleset{
		RulesetID:     tdrsd.RulesetID,
		Label:         tdrsd.Label,
		CriteriaType:  tdrsd.CriteriaType,
		Criteria:      tdrsd.Criteria,
		Ordering:      ordering,
		ResponsePools: make([]*TrafficDirectorResponsePool, len(tdrsd.ResponsePools)),
	}

	for idx, responsePool := range tdrsd.ResponsePools {
		tdrs.ResponsePools[idx] = responsePool.newTrafficDirectorResponsePool()
	}

	return &tdrs
}

// CreateTrafficDirectorRuleset creates a new instance of Traffic Director Response Pool.
func (c *Client) CreateTrafficDirectorRuleset(serviceID string, label string, options ...TrafficDirectorRulesetOptionSetter) (*TrafficDirectorRuleset, error) {
	req := TrafficDirectorRulesetCURequest{
		Label:        label,
		CriteriaType: "always",
		Publish:      "Y",
	}

	for _, o := range options {
		o(&req)
	}

	var resp trafficDirectorRulesetResponse

	if err := c.post(fmt.Sprintf("DSFRuleset/%s", serviceID), req, &resp); err != nil {
		return nil, err
	}

	tdrs := resp.newTrafficDirectorRuleset()

	return tdrs, nil
}

// UpdateTrafficDirectorRuleset updates an instance of Traffic Director Response Pool.
func (c *Client) UpdateTrafficDirectorRuleset(serviceID string, rulesetID string, label string, options ...TrafficDirectorRulesetOptionSetter) (*TrafficDirectorRuleset, error) {
	req := TrafficDirectorRulesetCURequest{
		Label:        label,
		CriteriaType: "always",
		Publish:      "Y",
	}

	for _, o := range options {
		o(&req)
	}

	var resp trafficDirectorRulesetResponse

	if err := c.put(fmt.Sprintf("DSFRuleset/%s/%s", serviceID, rulesetID), req, &resp); err != nil {
		return nil, err
	}

	tdrs := resp.newTrafficDirectorRuleset()

	return tdrs, nil
}

// DeleteTrafficDirectorRuleset deletes an instance of Traffic Director Response Pool.
func (c *Client) DeleteTrafficDirectorRuleset(serviceID string, rulesetID string) error {
	req := trafficDirectorRulesetDeleteRequest{
		Publish: "Y",
	}

	if err := c.delete(fmt.Sprintf("DSFRuleset/%s/%s", serviceID, rulesetID), req); err != nil {
		return err
	}

	return nil
}

// GetTrafficDirectorRuleset returns an existing Traffic Director Response Pool instance.
func (c *Client) GetTrafficDirectorRuleset(serviceID string, rulesetID string) (*TrafficDirectorRuleset, error) {
	var resp trafficDirectorRulesetResponse

	if err := c.get(fmt.Sprintf("DSFRuleset/%s/%s", serviceID, rulesetID), nil, &resp); err != nil {
		return nil, err
	}

	return resp.newTrafficDirectorRuleset(), nil
}
