package dyn

import (
	"fmt"
	"net/url"
	"sort"
)

// Record represents a Dyn zone record.
type Record struct {
	recordData
}

func (r *Record) setTTL(ttl int) {
	r.TTL = ttl
}

type recordData struct {
	Zone        string `json:"zone"`
	TTL         int    `json:"ttl"`
	FQDN        string `json:"fqdn"`
	RecordType  string `json:"record_type"`
	RData       rData  `json:"rdata"`
	RecordID    int    `json:"record_id"`
	SerialStyle string `json:"serial_style,omitempty"`
}

type rData struct {
	// A, AAAA
	Address string `json:"address,omitempty"`
	// ALIAS
	Alias string `json:"alias,omitempty"`
	// CAA
	Flags int    `json:"flags,omitempty"`
	Tag   string `json:"tag,omitempty"`
	Value string `json:"value,omitempty"`
	//CNAME
	CName string `json:"cname,omitempty"`
	// MX
	Exchange   string `json:"exchange,omitempty"`
	Preference int    `json:"preference,omitempty"`
	// NS
	NSDName string `json:"nsdname,omitempty"`
	// SOA
	RName string `json:"rname,omitempty"`
	// SRV
	Priority int    `json:"priority,omitempty"`
	Weight   int    `json:"weight,omitempty"`
	Port     int    `json:"port,omitempty"`
	Target   string `json:"target,omitempty"`
	// TXT
	TXTData string `json:"txtdata,omitempty"`
}

type zoneRecordsResponse struct {
	responseHeader
	Records recordsByType `json:"data"`
}

type recordsByType map[string][]recordData

// RDataValues provides a functional way to set rData
type RDataValues func(r *rData)

// NewRecord constructs a Record
func NewRecord(zone, fqdn, recordType string, options ...TTLOption) *Record {
	r := &Record{
		recordData: recordData{
			Zone:       zone,
			FQDN:       fqdn,
			RecordType: recordType,
		},
	}

	for _, o := range options {
		o(r)
	}

	return r
}

// NewARecord constructs an A record
func NewARecord(zone, fqdn string, address string, options ...TTLOption) *Record {
	r := NewRecord(zone, fqdn, "A", options...)

	r.RData.Address = address

	return r
}

// NewAAAARecord constructs an AAAA record
func NewAAAARecord(zone, fqdn string, address string, options ...TTLOption) *Record {
	r := NewRecord(zone, fqdn, "AAAA", options...)

	r.RData.Address = address

	return r
}

// NewALIASRecord constructs an ALIAS record
func NewALIASRecord(zone, fqdn string, alias string, options ...TTLOption) *Record {
	r := NewRecord(zone, fqdn, "ALIAS", options...)

	r.RData.Alias = alias

	return r
}

// NewCAARecord constructs a CAA record
func NewCAARecord(zone, fqdn string, flags int, tag string, value string, options ...TTLOption) *Record {
	r := NewRecord(zone, fqdn, "CAA", options...)

	r.RData.Flags = flags
	r.RData.Tag = tag
	r.RData.Value = value

	return r
}

// NewCNAMERecord constructs a CNAME record
func NewCNAMERecord(zone, fqdn string, cname string, options ...TTLOption) *Record {
	r := NewRecord(zone, fqdn, "CNAME", options...)

	r.RData.CName = cname

	return r
}

// NewMXRecord constructs an MX record
func NewMXRecord(zone, fqdn string, preference int, exchange string, options ...TTLOption) *Record {
	r := NewRecord(zone, fqdn, "MX", options...)

	r.RData.Preference = preference
	r.RData.Exchange = exchange

	return r
}

// NewNSRecord constructs an NS record
func NewNSRecord(zone, fqdn string, nsdname string, options ...TTLOption) *Record {
	r := NewRecord(zone, fqdn, "NS", options...)

	r.RData.NSDName = nsdname

	return r
}

// NewSRVRecord constructs an SRV record
func NewSRVRecord(zone, fqdn string, priority, weight, port int, target string, options ...TTLOption) *Record {
	r := NewRecord(zone, fqdn, "SRV", options...)

	r.RData.Priority = priority
	r.RData.Weight = weight
	r.RData.Port = port
	r.RData.Target = target

	return r
}

// NewTXTRecord constructs a TXT record
func NewTXTRecord(zone, fqdn string, txtdata string, options ...TTLOption) *Record {
	r := NewRecord(zone, fqdn, "TXT", options...)

	r.RData.TXTData = txtdata

	return r
}

// String returns a string representation of a Record
func (r Record) String() string {
	var rdata string

	switch r.RecordType {
	case "A", "AAAA":
		rdata = r.RData.Address
	case "ALIAS":
		rdata = r.RData.Alias
	case "CAA":
		rdata = fmt.Sprintf("%d %s %q", r.RData.Flags, r.RData.Tag, r.RData.Value)
	case "CNAME":
		rdata = r.RData.CName
	case "MX":
		rdata = fmt.Sprintf("%d %s", r.RData.Preference, r.RData.Exchange)
	case "NS":
		rdata = r.RData.NSDName
	case "SOA":
		rdata = r.RData.RName
	case "SRV":
		rdata = fmt.Sprintf("%d %d %d %s", r.RData.Priority, r.RData.Weight, r.RData.Port, r.RData.Target)
	case "TXT":
		rdata = fmt.Sprintf("%q", r.RData.TXTData)
	}

	tabs := 3 - int((len(r.FQDN)+1)/8)

	if tabs < 0 {
		tabs = 0
	}

	return fmt.Sprintf("%s.%s%d\tIN\t%s\t%v", r.FQDN, "\t\t\t"[:tabs], r.TTL, r.RecordType, rdata)
}

// EachRecord calls the provided function once for each record in the zone.
func (c *Client) EachRecord(zone string, f func(r *Record)) (int, error) {
	var resp zoneRecordsResponse

	params := url.Values{}
	params.Set("detail", "Y")

	if err := c.get(fmt.Sprintf("AllRecord/%s", zone), params, &resp); err != nil {
		return 0, err
	}

	records := resp.Records.flatten()

	for _, r := range records {
		record := &Record{r}

		f(record)
	}

	return len(records), nil
}

func (records recordsByType) flatten() []recordData {
	count := 0

	for _, r := range records {
		count += len(r)
	}

	flattened := make([]recordData, 0, count)

	for _, r := range records {
		flattened = append(flattened, r...)
	}

	sort.Sort(byRecord(flattened))

	return flattened
}

type byRecord []recordData

func (r byRecord) Len() int {
	return len(r)
}

func (r byRecord) Swap(i, j int) {
	r[i], r[j] = r[j], r[i]
}

func (r byRecord) Less(i, j int) bool {
	if r[i].RecordType < r[j].RecordType {
		return true
	}

	if r[i].RecordType > r[j].RecordType {
		return false
	}

	if r[i].FQDN < r[j].FQDN {
		return true
	}

	if r[i].FQDN > r[j].FQDN {
		return false
	}

	return r[i].RecordID < r[j].RecordID
}
