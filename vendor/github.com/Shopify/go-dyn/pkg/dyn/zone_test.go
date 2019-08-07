package dyn

import (
	"fmt"
	"net/http"
	"testing"
)

func assertZone(t *testing.T, serial int, serialStyle string, zone string, zoneType string, z *Zone) {
	t.Helper()

	assertEqual(t, serial, z.Serial, "Serial")
	assertEqual(t, serialStyle, z.SerialStyle, "SerialStyle")
	assertEqual(t, zone, z.Zone, "Zone")
	assertEqual(t, zoneType, z.ZoneType, "ZoneType")
}

func TestCreateZone(t *testing.T) {
	zone := "go-dyn-test-create.go-dyn.com"

	c := mockClient("zone_service/create.json", func(w http.ResponseWriter, r *http.Request, j interface{}) {
		assertMethod(t, http.MethodPost, r)
		assertPath(t, fmt.Sprintf("/REST/Zone/%s", zone), r)

		assertUserAgent(t, "go-dyn/0.0.0", r)
		assertContentType(t, "application/json", r)
		assertAuthToken(t, "insert-token-here", r)

		assertJSON(t, "admin@example.com", "rname", j)
		assertJSON(t, nil, "serial_style", j)
		assertJSON(t, "86400", "ttl", j)

		w.Header().Set("Content-Type", "application/json")
	})

	c.token = "insert-token-here"

	if z, err := c.CreateZone(zone, "admin@example.com", 24*60*60); err != nil {
		t.Error(err)
	} else {
		assertZone(t, 0, SerialStyleIncrement, zone, ZoneTypePrimary, z)
	}
}

func TestCreateZoneWithSerialStyle(t *testing.T) {
	zone := "go-dyn-test-create.go-dyn.com"

	c := mockClient("zone_service/create_with_serial_style.json", func(w http.ResponseWriter, r *http.Request, j interface{}) {
		assertMethod(t, http.MethodPost, r)
		assertPath(t, fmt.Sprintf("/REST/Zone/%s", zone), r)

		assertUserAgent(t, "go-dyn/0.0.0", r)
		assertContentType(t, "application/json", r)
		assertAuthToken(t, "insert-token-here", r)

		assertJSON(t, "admin@example.com", "rname", j)
		assertJSON(t, SerialStyleMinute, "serial_style", j)
		assertJSON(t, "86400", "ttl", j)

		w.Header().Set("Content-Type", "application/json")
	})

	c.token = "insert-token-here"

	if z, err := c.CreateZone(zone, "admin@example.com", 24*60*60, SerialStyle(SerialStyleMinute)); err != nil {
		t.Error(err)
	} else {
		assertZone(t, 0, SerialStyleMinute, zone, ZoneTypePrimary, z)
	}
}

func TestCreateZoneError(t *testing.T) {
	zone := "go-dyn-test-create.go-dyn.com"

	c := mockClient("zone_service/create_error.json", func(w http.ResponseWriter, r *http.Request, j interface{}) {
		w.Header().Set("Content-Type", "application/json")

		w.WriteHeader(http.StatusBadRequest)
	})

	c.token = "insert-token-here"

	if _, err := c.CreateZone(zone, "admin@example.com", 24*60*60); err != nil {
		assertEqual(t, "TARGET_EXISTS: name: Name already exists", err.Error(), "error")
	} else {
		t.Error("Expected Create to fail")
	}
}

func TestGetZone(t *testing.T) {
	zone := "go-dyn.com"

	c := mockClient("zone/get.json", func(w http.ResponseWriter, r *http.Request, j interface{}) {
		assertMethod(t, http.MethodGet, r)
		assertPath(t, fmt.Sprintf("/REST/Zone/%s", zone), r)

		assertUserAgent(t, "go-dyn/0.0.0", r)
		assertContentType(t, "application/json", r)
		assertAuthToken(t, "insert-token-here", r)

		w.Header().Set("Content-Type", "application/json")
	})

	c.token = "insert-token-here"

	if z, err := c.GetZone(zone); err != nil {
		t.Error(err)
	} else {
		assertZone(t, 1902190535, SerialStyleMinute, zone, ZoneTypePrimary, z)
	}
}

func TestGetZoneError(t *testing.T) {
	zone := "missing.go-dyn.com"

	c := mockClient("zone/get_error.json", func(w http.ResponseWriter, r *http.Request, j interface{}) {
		w.Header().Set("Content-Type", "application/json")

		w.WriteHeader(http.StatusNotFound)
	})

	if _, err := c.GetZone(zone); err != nil {
		assertEqual(t, "NOT_FOUND: zone: No such zone", err.Error(), "error")
	} else {
		t.Error("Expected Get to fail")
	}
}

func TestEachZone(t *testing.T) {
	c := mockClient("zone/each.json", func(w http.ResponseWriter, r *http.Request, j interface{}) {
		assertMethod(t, http.MethodGet, r)
		assertPath(t, "/REST/Zone", r)
		assertParam(t, "Y", "detail", r)

		assertUserAgent(t, "go-dyn/0.0.0", r)
		assertContentType(t, "application/json", r)
		assertAuthToken(t, "insert-token-here", r)

		w.Header().Set("Content-Type", "application/json")
	})

	c.token = "insert-token-here"

	i := 1

	n, err := c.EachZone(func(z *Zone) {
		assertZone(t, i*11, SerialStyleIncrement, fmt.Sprintf("zone-%d.co", i), ZoneTypePrimary, z)
		i++
	})

	if err != nil {
		t.Fatal(err)
	}

	assertEqual(t, 7, n, "count")
}

func TestPublishZone(t *testing.T) {
	zone := "go-dyn-test-publish.go-dyn.com"

	c := mockClient("zone/publish.json", func(w http.ResponseWriter, r *http.Request, j interface{}) {
		assertMethod(t, http.MethodPut, r)
		assertPath(t, fmt.Sprintf("/REST/Zone/%s", zone), r)

		assertUserAgent(t, "go-dyn/0.0.0", r)
		assertContentType(t, "application/json", r)
		assertAuthToken(t, "insert-token-here", r)

		assertJSON(t, true, "publish", j)
		assertJSON(t, "insert notes here", "notes", j)

		w.Header().Set("Content-Type", "application/json")
	})

	c.token = "insert-token-here"

	z, err := c.PublishZone(zone, "insert notes here")
	if err != nil {
		t.Error(err)
	} else {
		assertZone(t, 1, SerialStyleIncrement, zone, ZoneTypePrimary, z)
	}
}

func TestFreezeZone(t *testing.T) {
	zone := "go-dyn-test-freeze.go-dyn.com"

	c := mockClient("zone/freeze.json", func(w http.ResponseWriter, r *http.Request, j interface{}) {
		assertMethod(t, http.MethodPut, r)
		assertPath(t, fmt.Sprintf("/REST/Zone/%s", zone), r)

		assertUserAgent(t, "go-dyn/0.0.0", r)
		assertContentType(t, "application/json", r)
		assertAuthToken(t, "insert-token-here", r)

		w.Header().Set("Content-Type", "application/json")
	})

	c.token = "insert-token-here"

	if err := c.FreezeZone(zone); err != nil {
		t.Error(err)
	}
}

func TestThawZone(t *testing.T) {
	zone := "go-dyn-test-thaw.go-dyn.com"

	c := mockClient("zone/thaw.json", func(w http.ResponseWriter, r *http.Request, j interface{}) {
		assertMethod(t, http.MethodPut, r)
		assertPath(t, fmt.Sprintf("/REST/Zone/%s", zone), r)

		assertUserAgent(t, "go-dyn/0.0.0", r)
		assertContentType(t, "application/json", r)
		assertAuthToken(t, "insert-token-here", r)

		w.Header().Set("Content-Type", "application/json")
	})

	c.token = "insert-token-here"

	if err := c.ThawZone(zone); err != nil {
		t.Error(err)
	}
}

func TestDeleteZone(t *testing.T) {
	zone := "go-dyn-test-delete.go-dyn.com"

	c := mockClient("zone/delete.json", func(w http.ResponseWriter, r *http.Request, j interface{}) {
		assertMethod(t, http.MethodDelete, r)
		assertPath(t, fmt.Sprintf("/REST/Zone/%s", zone), r)

		assertUserAgent(t, "go-dyn/0.0.0", r)
		assertContentType(t, "application/json", r)
		assertAuthToken(t, "insert-token-here", r)

		w.Header().Set("Content-Type", "application/json")
	})

	c.token = "insert-token-here"

	if err := c.DeleteZone(zone); err != nil {
		t.Error(err)
	}
}
