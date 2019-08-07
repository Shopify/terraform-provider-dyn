package dyn

import (
	"fmt"
	"net/url"
	"strconv"
)

// TrafficDirector represents a Dyn Traffic Director service.
type TrafficDirector struct {
	ServiceID     string
	Label         string
	Active        bool
	TTL           int
	Nodes         []trafficDirectorNode
	Rulesets      []*TrafficDirectorRuleset
	ResponsePools []*TrafficDirectorResponsePool
}

type trafficDirectorData struct {
	ServiceID     string                       `json:"service_id"`
	Label         string                       `json:"label"`
	Active        string                       `json:"active"`
	TTL           string                       `json:"ttl"`
	Notifiers     []trafficDirectorNotifier    `json:"notifiers"`
	Rulesets      []trafficDirectorRulesetData `json:"rulesets"`
	Nodes         []trafficDirectorNode        `json:"nodes"`
	PendingChange string                       `json:"pending_change"`
}

type trafficDirectorNode struct {
	Zone string `json:"zone"`
	FQDN string `json:"fqdn"`
}

type trafficDirectorNotifier struct{}

type TrafficDirectorCURequest struct {
	Label     string                            `json:"label"`
	TTL       int                               `json:"ttl,omitempty"`
	Publish   string                            `json:"publish,omitempty"`
	Notes     string                            `json:"notes,omitempty"`
	Nodes     []trafficDirectorNode             `json:"nodes,omitempty"`
	Notifiers []trafficDirectorNotifier         `json:"notifiers,omitempty"`
	Rulesets  []TrafficDirectorRulesetCURequest `json:"rulesets,omitempty"`
}

func (tdreq *TrafficDirectorCURequest) AddNode(node map[string]string) {
	tdreq.Nodes = append(tdreq.Nodes, trafficDirectorNode{
		Zone: node["zone"],
		FQDN: node["fqdn"],
	})
}

type trafficDirectorResponse struct {
	responseHeader
	trafficDirectorData `json:"data"`
}

type trafficDirectorAllResponse struct {
	responseHeader
	TrafficDirectors []trafficDirectorData `json:"data"`
}

type TrafficDirectorOptionSetter func(*TrafficDirectorCURequest)

func (tdd trafficDirectorData) newTrafficDirector() *TrafficDirector {
	ttl, _ := strconv.Atoi(tdd.TTL)

	td := TrafficDirector{
		ServiceID: tdd.ServiceID,
		Label:     tdd.Label,
		Active:    tdd.Active == "Y",
		TTL:       ttl,
		Nodes:     tdd.Nodes,
		Rulesets:  make([]*TrafficDirectorRuleset, len(tdd.Rulesets)),
	}

	responsePools := make(map[string]*TrafficDirectorResponsePool)
	for idx, ruleset := range tdd.Rulesets {
		td.Rulesets[idx] = ruleset.newTrafficDirectorRuleset()
		for _, responsePool := range td.Rulesets[idx].ResponsePools {
			responsePools[responsePool.ResponsePoolID] = responsePool
		}
	}

	td.ResponsePools = make([]*TrafficDirectorResponsePool, 0, len(responsePools))
	for _, responsePool := range responsePools {
		td.ResponsePools = append(td.ResponsePools, responsePool)
	}

	return &td
}

// CreateTrafficDirector creates a new instance of Traffic Director.
func (c *Client) CreateTrafficDirector(label string, options ...TrafficDirectorOptionSetter) (*TrafficDirector, error) {
	req := TrafficDirectorCURequest{
		Label:   label,
		Publish: "Y",
	}

	for _, o := range options {
		o(&req)
	}

	var resp trafficDirectorResponse

	if err := c.post("DSF", req, &resp); err != nil {
		return nil, err
	}

	td := resp.newTrafficDirector()

	return td, nil
}

// UpdateTrafficDirector updates an instance of Traffic Director.
func (c *Client) UpdateTrafficDirector(serviceID string, label string, options ...TrafficDirectorOptionSetter) (*TrafficDirector, error) {
	req := TrafficDirectorCURequest{
		Label:   label,
		Publish: "Y",
	}

	for _, o := range options {
		o(&req)
	}

	var resp trafficDirectorResponse

	if err := c.put(fmt.Sprintf("DSF/%s", serviceID), req, &resp); err != nil {
		return nil, err
	}

	td := resp.newTrafficDirector()

	return td, nil
}

// DeleteTrafficDirector deletes an instance of Traffic Director.
func (c *Client) DeleteTrafficDirector(serviceID string) error {
	if err := c.delete(fmt.Sprintf("DSF/%s", serviceID), nil); err != nil {
		return err
	}

	return nil
}

// EachTrafficDirector calls the provided function for every existing Traffic Director service instance.
func (c *Client) EachTrafficDirector(f func(td *TrafficDirector)) (int, error) {
	var resp trafficDirectorAllResponse

	params := url.Values{}
	params.Set("detail", "Y")

	if err := c.get("DSF", params, &resp); err != nil {
		return 0, err
	}

	for _, td := range resp.TrafficDirectors {
		f(td.newTrafficDirector())
	}

	return len(resp.TrafficDirectors), nil
}

// FindTrafficDirector returns the service ID for the Traffic Director service instance with the specified label.
func (c *Client) FindTrafficDirector(label string) (*TrafficDirector, error) {
	params := url.Values{}
	params.Set("label", label)
	params.Set("detail", "Y")

	var resp trafficDirectorAllResponse

	if err := c.get("DSF", params, &resp); err != nil {
		return nil, err
	}

	if len(resp.TrafficDirectors) == 0 {
		return nil, fmt.Errorf("Unable to find a traffic director for label: %s", label)
	}

	return c.GetTrafficDirector(resp.TrafficDirectors[0].newTrafficDirector().ServiceID)
}

// GetTrafficDirector returns an existing Traffic Director service instance.
func (c *Client) GetTrafficDirector(serviceID string) (*TrafficDirector, error) {
	var resp trafficDirectorResponse

	if err := c.get(fmt.Sprintf("DSF/%s", serviceID), nil, &resp); err != nil {
		return nil, err
	}

	// return nil, fmt.Errorf("BLIH: %#v", resp)

	return resp.newTrafficDirector(), nil
}
