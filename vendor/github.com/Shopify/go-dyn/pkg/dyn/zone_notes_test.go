package dyn

import (
	"fmt"
	"net/http"
	"testing"
)

func TestGetZoneNotes(t *testing.T) {
	zone := "go-dyn.com"

	c := mockClient("zone_notes/get_zone_notes.json", func(w http.ResponseWriter, r *http.Request, j interface{}) {
		assertMethod(t, http.MethodPost, r)
		assertPath(t, "/REST/ZoneNoteReport", r)

		assertUserAgent(t, "go-dyn/0.0.0", r)
		assertContentType(t, "application/json", r)
		assertAuthToken(t, "insert-token-here", r)

		assertJSON(t, "go-dyn.com", "zone", j)

		w.Header().Set("Content-Type", "application/json")
	})

	c.token = "insert-token-here"

	if notes, err := c.GetZoneNotes(zone); err != nil {
		t.Error(err)
	} else {
		assertEqual(t, 4, len(notes), "length")

		for i, n := range notes {
			assertEqual(t, fmt.Sprintf("note %d\n", 4-i), n.Note, "note")
		}
	}
}

func TestGetZoneNotesWithPagination(t *testing.T) {
	zone := "go-dyn.com"

	c := mockClient("zone_notes/get_zone_notes_with_pagination.json", func(w http.ResponseWriter, r *http.Request, j interface{}) {
		assertMethod(t, http.MethodPost, r)
		assertPath(t, "/REST/ZoneNoteReport", r)

		assertUserAgent(t, "go-dyn/0.0.0", r)
		assertContentType(t, "application/json", r)
		assertAuthToken(t, "insert-token-here", r)

		assertJSON(t, "go-dyn.com", "zone", j)
		assertJSON(t, float64(2), "limit", j)
		assertJSON(t, float64(1), "offset", j)

		w.Header().Set("Content-Type", "application/json")
	})

	c.token = "insert-token-here"

	if notes, err := c.GetZoneNotes(zone, Limit(2), Offset(1)); err != nil {
		t.Error(err)
	} else {
		assertEqual(t, 2, len(notes), "length")

		for i, n := range notes {
			assertEqual(t, fmt.Sprintf("note %d\n", 3-i), n.Note, "note")
		}
	}
}
