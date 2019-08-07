package dyn

import (
	"fmt"
	"net/http"
	"testing"
)

func assertTrafficDirector(t *testing.T, serviceID string, label string, ttl int, active bool, td *TrafficDirector) {
	t.Helper()

	assertEqual(t, serviceID, td.ServiceID, "ServiceID")
	assertEqual(t, label, td.Label, "Label")
	assertEqual(t, ttl, td.TTL, "TTL")
	assertEqual(t, active, td.Active, "Active")
}

func TestCreateTrafficDirector(t *testing.T) {
	c := mockClient("traffic_director/create.json", func(w http.ResponseWriter, r *http.Request, j interface{}) {
		assertMethod(t, http.MethodPost, r)
		assertPath(t, "/REST/DSF", r)

		assertUserAgent(t, "go-dyn/0.0.0", r)
		assertContentType(t, "application/json", r)
		assertAuthToken(t, "insert-token-here", r)

		assertJSON(t, "insert-label-here", "label", j)
		assertJSON(t, nil, "ttl", j)

		w.Header().Set("Content-Type", "application/json")
	})

	c.token = "insert-token-here"

	if td, err := c.CreateTrafficDirector("insert-label-here"); err != nil {
		t.Error(err)
	} else {
		assertTrafficDirector(t, "insert-service-id-here", "insert-label-here", 3600, true, td)
	}
}

func TestCreateTrafficDirectorWithTTL(t *testing.T) {
	c := mockClient("traffic_director/create_with_ttl.json", func(w http.ResponseWriter, r *http.Request, j interface{}) {
		assertMethod(t, http.MethodPost, r)
		assertPath(t, "/REST/DSF", r)

		assertUserAgent(t, "go-dyn/0.0.0", r)
		assertContentType(t, "application/json", r)
		assertAuthToken(t, "insert-token-here", r)

		assertJSON(t, "insert-label-here", "label", j)
		assertJSON(t, float64(60), "ttl", j)

		w.Header().Set("Content-Type", "application/json")
	})

	c.token = "insert-token-here"

	optionsSetter := func(req *TrafficDirectorCURequest) {
		req.TTL = 60
	}

	if td, err := c.CreateTrafficDirector("insert-label-here", optionsSetter); err != nil {
		t.Error(err)
	} else {
		assertTrafficDirector(t, "insert-service-id-here", "insert-label-here", 60, true, td)
	}
}

func TestEachTrafficDirector(t *testing.T) {
	c := mockClient("traffic_director/each.json", func(w http.ResponseWriter, r *http.Request, j interface{}) {
		assertMethod(t, http.MethodGet, r)
		assertPath(t, "/REST/DSF", r)
		assertParam(t, "Y", "detail", r)

		assertUserAgent(t, "go-dyn/0.0.0", r)
		assertContentType(t, "application/json", r)
		assertAuthToken(t, "insert-token-here", r)

		w.Header().Set("Content-Type", "application/json")
	})

	c.token = "insert-token-here"

	i := 1

	n, err := c.EachTrafficDirector(func(td *TrafficDirector) {
		assertTrafficDirector(t, fmt.Sprintf("service-%d", i), fmt.Sprintf("service %d", i), 15, i != 5, td)
		i++
	})

	if err != nil {
		t.Fatal(err)
	}

	assertEqual(t, 5, n, "count")
}

func TestFindTrafficDirector(t *testing.T) {
	c := mockClient("traffic_director/find.json", func(w http.ResponseWriter, r *http.Request, j interface{}) {
		assertMethod(t, http.MethodGet, r)
		assertPath(t, "/REST/DSF", r)
		assertParam(t, "service 3", "label", r)
		assertParam(t, "Y", "detail", r)

		assertUserAgent(t, "go-dyn/0.0.0", r)
		assertContentType(t, "application/json", r)
		assertAuthToken(t, "insert-token-here", r)

		w.Header().Set("Content-Type", "application/json")
	})

	c.token = "insert-token-here"

	td, err := c.FindTrafficDirector("service 3")

	if err != nil {
		t.Error(err)
	}

	assertTrafficDirector(t, "service-3", "service 3", 15, true, td)
}

func TestGetTrafficDirector(t *testing.T) {
	serviceID := "service-1"

	c := mockClient("traffic_director/get.json", func(w http.ResponseWriter, r *http.Request, j interface{}) {
		assertMethod(t, http.MethodGet, r)
		assertPath(t, fmt.Sprintf("/REST/DSF/%s", serviceID), r)

		assertUserAgent(t, "go-dyn/0.0.0", r)
		assertContentType(t, "application/json", r)
		assertAuthToken(t, "insert-token-here", r)

		w.Header().Set("Content-Type", "application/json")
	})

	c.token = "insert-token-here"

	if td, err := c.GetTrafficDirector(serviceID); err != nil {
		t.Error(err)
	} else {
		assertTrafficDirector(t, "insert-service-id-here", "insert-label-here", 60, false, td)
	}
}
