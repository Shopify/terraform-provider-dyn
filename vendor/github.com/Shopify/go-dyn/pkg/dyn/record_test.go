package dyn

import (
	"net/http"
	"testing"
)

func assertRecord(t *testing.T, expected *Record, r *Record) {
	t.Helper()

	assertEqual(t, expected.Zone, r.Zone, "Zone")
	assertEqual(t, expected.TTL, r.TTL, "TTL")
	assertEqual(t, expected.FQDN, r.FQDN, "FQDN")
	assertEqual(t, expected.RecordType, r.RecordType, "RecordType")

	assertEqual(t, expected.RData, r.RData, "RData")
}

func TestEachRecord(t *testing.T) {
	zone := "go-dyn.com"

	c := mockClient("record/each.json", func(w http.ResponseWriter, r *http.Request, j interface{}) {
		assertMethod(t, http.MethodGet, r)
		assertPath(t, "/REST/AllRecord/go-dyn.com", r)
		assertParam(t, "Y", "detail", r)

		assertUserAgent(t, "go-dyn/0.0.0", r)
		assertContentType(t, "application/json", r)
		assertAuthToken(t, "insert-token-here", r)

		w.Header().Set("Content-Type", "application/json")
	})

	c.token = "insert-token-here"

	i := 0

	n, err := c.EachRecord(zone, func(r *Record) {
		switch i {
		case 0:
			assertRecord(t, NewARecord(zone, "a.go-dyn.com", "10.1.2.3", TTL(3600)), r)
		case 1:
			assertRecord(t, NewAAAARecord(zone, "aaaa.go-dyn.com", "fdfe:972e:46c1:6ed0:0000:0001:0002:0003", TTL(3600)), r)
		case 2:
			assertRecord(t, NewALIASRecord(zone, "alias.go-dyn.com", "go-dyn.com.", TTL(30)), r)
		case 3:
			assertRecord(t, NewCAARecord(zone, "caa.go-dyn.com", 0, "issue", "letsencrypt.org", TTL(1800)), r)
		case 4:
			assertRecord(t, NewCNAMERecord(zone, "cname.go-dyn.com", "go-dyn.example.com.", TTL(900)), r)
		case 5:
			assertRecord(t, NewMXRecord(zone, "mx.go-dyn.com", 10, "mail.example.com.", TTL(3600)), r)
		case 6:
			assertRecord(t, NewNSRecord(zone, "go-dyn.com", "ns1.p19.dynect.net.", TTL(86400)), r)
		case 7:
			assertRecord(t, NewNSRecord(zone, "go-dyn.com", "ns2.p19.dynect.net.", TTL(86400)), r)
		case 8:
			assertRecord(t, NewNSRecord(zone, "go-dyn.com", "ns3.p19.dynect.net.", TTL(86400)), r)
		case 9:
			assertRecord(t, NewNSRecord(zone, "go-dyn.com", "ns4.p19.dynect.net.", TTL(86400)), r)
		case 10:
			soa := NewRecord(zone, "go-dyn.com", "SOA", TTL(3600))
			soa.RData.RName = "admin@example.com."
			assertRecord(t, soa, r)
		case 11:
			assertRecord(t, NewSRVRecord(zone, "srv.go-dyn.com", 10, 100, 443, "www.example.com.", TTL(3600)), r)
		case 12:
			assertRecord(t, NewTXTRecord(zone, "txt.go-dyn.com", "hello world", TTL(7200)), r)
		}
		i++
	})

	if err != nil {
		t.Fatal(err)
	}

	assertEqual(t, 13, n, "count")
}
