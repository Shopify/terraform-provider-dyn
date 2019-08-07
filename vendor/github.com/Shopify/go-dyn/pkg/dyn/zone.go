package dyn

import (
	"fmt"
	"net/url"
	"strconv"
)

// SerialStyle values
const (
	SerialStyleDefault   = ""
	SerialStyleIncrement = "increment" // Serial incremented by 1 on every change. Default setting.
	SerialStyleEpoch     = "epoch"     // Serial is UNIX timestamp at the time of the publish.
	SerialStyleDay       = "day"       // Serial form of YYYYMMDDxx where xx is incremented for each change in a day.
	SerialStyleMinute    = "minute"    // Serial form of YYMMDDHHMM.
)

// ZoneType values
const (
	ZoneTypePrimary   = "Primary"
	ZoneTypeSecondary = "Secondary"
)

// Zone represents a Dyn zone.
type Zone struct {
	Serial      int    `json:"serial"`
	SerialStyle string `json:"serial_style"`
	Zone        string `json:"zone"`
	ZoneType    string `json:"zone_type"`
}

type zoneCreateRequest struct {
	RName       string `json:"rname"`
	SerialStyle string `json:"serial_style,omitempty"`
	TTL         string `json:"ttl"`
}

type zoneUpdateRequest struct {
	Freeze  bool   `json:"freeze,omitempty"`
	Thaw    bool   `json:"thaw,omitempty"`
	Publish bool   `json:"publish,omitempty"`
	Notes   string `json:"notes,omitempty"`
}

type zoneResponseData struct {
	TaskID string `json:"task_id"`
	Zone
}

type zoneResponse struct {
	responseHeader
	zoneResponseData `json:"data"`
}

type zoneGetResponse struct {
	responseHeader
	Zone `json:"data"`
}

type zoneAllResponse struct {
	responseHeader
	Zones []Zone `json:"data"`
}

// CreateZone creates a new Zone.
func (c *Client) CreateZone(zone string, rName string, ttl int, options ...ZoneOption) (*Zone, error) {
	req := zoneCreateRequest{
		RName: rName,
		TTL:   strconv.Itoa(ttl),
	}

	for _, o := range options {
		o(&req)
	}

	var resp zoneResponse

	if err := c.post(fmt.Sprintf("Zone/%s", zone), req, &resp); err != nil {
		return nil, err
	}

	return &resp.Zone, nil
}

// GetZone returns an existing Zone.
func (c *Client) GetZone(zone string) (*Zone, error) {
	var resp zoneGetResponse

	if err := c.get(fmt.Sprintf("Zone/%s", zone), nil, &resp); err != nil {
		return nil, err
	}

	return &resp.Zone, nil
}

// EachZone calls the provided function for every existing Dyn Managed DNS zone.
func (c *Client) EachZone(f func(z *Zone)) (int, error) {
	var resp zoneAllResponse

	params := url.Values{}
	params.Set("detail", "Y")

	if err := c.get("Zone", params, &resp); err != nil {
		return 0, err
	}

	for _, z := range resp.Zones {
		f(&z)
	}

	return len(resp.Zones), nil
}

// PublishZone causes pending changes to become part of the Zone.
func (c *Client) PublishZone(zone, notes string) (*Zone, error) {
	req := zoneUpdateRequest{
		Publish: true,
		Notes:   notes,
	}

	var resp zoneResponse

	if err := c.put(fmt.Sprintf("Zone/%s", zone), req, &resp); err != nil {
		return nil, err
	}

	return &resp.Zone, nil
}

// FreezeZone prevents changes to the Zone.
func (c *Client) FreezeZone(zone string) error {
	req := zoneUpdateRequest{
		Freeze: true,
	}

	return c.put(fmt.Sprintf("Zone/%s", zone), req, nil)
}

// ThawZone allows changes to again be made to the Zone.
func (c *Client) ThawZone(zone string) error {
	req := zoneUpdateRequest{
		Thaw: true,
	}

	return c.put(fmt.Sprintf("Zone/%s", zone), req, nil)
}

// DeleteZone removes the Zone.
func (c *Client) DeleteZone(zone string) error {
	return c.delete(fmt.Sprintf("Zone/%s", zone), nil)
}
