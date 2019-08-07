package dyn

// values for zoneNote type
const (
	zoneNoteTypePublish = "publish"
	zoneNoteTypeTask    = "task"
	zoneNoteTypeRemote  = "remove"
)

type zoneNotesRequest struct {
	Zone   string `json:"zone"`
	Limit  int    `json:"limit,omitempty"`
	Offset int    `json:"offset,omitempty"`
}

func (req *zoneNotesRequest) setLimit(limit int) {
	req.Limit = limit
}

func (req *zoneNotesRequest) setOffset(offset int) {
	req.Offset = offset
}

// ZoneNote is a note for a Dyn Managed DNS zone
type ZoneNote struct {
	Zone      string `json:"zone"`
	Serial    int    `json:"serial"`
	Type      string `json:"type"`
	Note      string `json:"note"`
	Timestamp string `json:"timestamp"`
	UserName  string `json:"user_name"`
}

type zoneNotesResponse struct {
	responseHeader
	Notes []ZoneNote `json:"data"`
}

// GetZoneNotes generates a report containing the Zone Notes for the Zone.
func (c *Client) GetZoneNotes(zone string, options ...PaginationOption) ([]ZoneNote, error) {
	req := zoneNotesRequest{
		Zone: zone,
	}

	for _, o := range options {
		o(&req)
	}

	var resp zoneNotesResponse

	if err := c.post("ZoneNoteReport", req, &resp); err != nil {
		return nil, err
	}

	return resp.Notes, nil
}
