package dyn

import (
	"fmt"
	"net/url"
	"strconv"
	"strings"
)

// TrafficDirectorMonitor represents a Dyn Traffic Director Monitor.
type TrafficDirectorMonitor struct {
	MonitorID     string
	Label         string
	Retries       int
	Protocol      string
	ResponseCount int
	ProbeInterval int
	Active        bool
	Options       TrafficDirectorMonitorOptions
	Services      []string
}

type TrafficDirectorMonitorOptions struct {
	Header   string
	Host     string
	Expected string
	Path     string
	Port     int
}

type trafficDirectorMonitorData struct {
	MonitorID     string                            `json:"dsf_monitor_id"`
	Label         string                            `json:"label"`
	Retries       string                            `json:"retries"`
	Protocol      string                            `json:"protocol"`
	ResponseCount string                            `json:"response_count"`
	ProbeInterval string                            `json:"probe_interval"`
	Active        string                            `json:"active"`
	Options       trafficDirectorMonitorOptionsData `json:"options"`
	Services      []string                          `json:"services"`
}

type trafficDirectorMonitorOptionsData struct {
	Header   string `json:"header"`
	Host     string `json:"host"`
	Expected string `json:"expected"`
	Path     string `json:"path"`
	Port     string `json:"port"`
}

type TrafficDirectorMonitorCURequest struct {
	Label         string                                 `json:"label"`
	Retries       int                                    `json:"retries"`
	Protocol      string                                 `json:"protocol"`
	ResponseCount int                                    `json:"response_count"`
	ProbeInterval int                                    `json:"probe_interval"`
	Active        string                                 `json:"active,omitempty"`
	Options       trafficDirectorMonitorCURequestOptions `json:"options,omitempty"`
	Publish       string                                 `json:"publish,omitempty"`
	Notes         string                                 `json:"notes,omitempty"`
}

type trafficDirectorMonitorCURequestOptions struct {
	Header   string `json:"header,omitempty"`
	Host     string `json:"host,omitempty"`
	Expected string `json:"expected,omitempty"`
	Path     string `json:"path,omitempty"`
	Port     int    `json:"port,omitempty"`
}

type trafficDirectorMonitorDeleteRequest struct {
	Publish string `json:"publish,omitempty"`
	Notes   string `json:"notes,omitempty"`
}

type trafficDirectorMonitorResponse struct {
	responseHeader
	trafficDirectorMonitorData `json:"data"`
}

type trafficDirectorMonitorAllResponse struct {
	responseHeader
	TrafficDirectorMonitorIDs []string `json:"data"`
}

type TrafficDirectorMonitorOptionSetter func(*TrafficDirectorMonitorCURequest)

func (tdmd trafficDirectorMonitorData) newTrafficDirectorMonitor() *TrafficDirectorMonitor {
	retries, _ := strconv.Atoi(tdmd.Retries)
	responseCount, _ := strconv.Atoi(tdmd.ResponseCount)
	probeInterval, _ := strconv.Atoi(tdmd.ProbeInterval)
	port, _ := strconv.Atoi(tdmd.Options.Port)

	tdrs := TrafficDirectorMonitor{
		MonitorID:     tdmd.MonitorID,
		Label:         tdmd.Label,
		Retries:       retries,
		Protocol:      tdmd.Protocol,
		ResponseCount: responseCount,
		ProbeInterval: probeInterval,
		Active:        tdmd.Active == "Y",
		Options: TrafficDirectorMonitorOptions{
			Header:   tdmd.Options.Header,
			Host:     tdmd.Options.Host,
			Expected: tdmd.Options.Expected,
			Path:     tdmd.Options.Path,
			Port:     port,
		},
		Services: tdmd.Services,
	}

	return &tdrs
}

// CreateTrafficDirectorMonitor creates a new instance of Traffic Director Monitor.
func (c *Client) CreateTrafficDirectorMonitor(label string, options ...TrafficDirectorMonitorOptionSetter) (*TrafficDirectorMonitor, error) {
	req := TrafficDirectorMonitorCURequest{
		Label:   label,
		Publish: "Y",
	}

	for _, o := range options {
		o(&req)
	}

	var resp trafficDirectorMonitorResponse

	if err := c.post("DSFMonitor", req, &resp); err != nil {
		return nil, err
	}

	tdrs := resp.newTrafficDirectorMonitor()

	return tdrs, nil
}

// UpdateTrafficDirectorMonitor updates an instance of Traffic Director Monitor.
func (c *Client) UpdateTrafficDirectorMonitor(monitorID string, label string, options ...TrafficDirectorMonitorOptionSetter) (*TrafficDirectorMonitor, error) {
	req := TrafficDirectorMonitorCURequest{
		Label:   label,
		Publish: "Y",
	}

	for _, o := range options {
		o(&req)
	}

	var resp trafficDirectorMonitorResponse

	if err := c.put(fmt.Sprintf("DSFMonitor/%s", monitorID), req, &resp); err != nil {
		return nil, err
	}

	tdrs := resp.newTrafficDirectorMonitor()

	return tdrs, nil
}

// DeleteTrafficDirectorMonitor deletes an instance of Traffic Director Monitor.
func (c *Client) DeleteTrafficDirectorMonitor(monitorID string) error {
	req := trafficDirectorMonitorDeleteRequest{
		Publish: "Y",
	}

	if err := c.delete(fmt.Sprintf("DSFMonitor/%s", monitorID), req); err != nil {
		return err
	}

	return nil
}

// FindTrafficDirectorMonitor returns the existing Traffic Director Monitor instance with the specified label.
func (c *Client) FindTrafficDirectorMonitor(label string) (*TrafficDirectorMonitor, error) {
	params := url.Values{}
	params.Set("label", label)

	var resp trafficDirectorMonitorAllResponse

	if err := c.get("DSFMonitor", params, &resp); err != nil {
		return nil, err
	}

	if len(resp.TrafficDirectorMonitorIDs) == 0 {
		return nil, fmt.Errorf("Unable to find a traffic director monitor for label: %s", label)
	}

	return c.GetTrafficDirectorMonitor(strings.Split(resp.TrafficDirectorMonitorIDs[0], "/")[3])
}

// GetTrafficDirectorMonitor returns an existing Traffic Director Monitor instance.
func (c *Client) GetTrafficDirectorMonitor(monitorID string) (*TrafficDirectorMonitor, error) {
	var resp trafficDirectorMonitorResponse

	if err := c.get(fmt.Sprintf("DSFMonitor/%s", monitorID), nil, &resp); err != nil {
		return nil, err
	}

	return resp.newTrafficDirectorMonitor(), nil
}
