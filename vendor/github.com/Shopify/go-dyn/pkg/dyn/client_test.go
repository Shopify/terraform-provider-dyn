package dyn

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"testing"
)

func writeFixture(w http.ResponseWriter, fixture string) error {
	path := fmt.Sprintf("testdata/%s", fixture)

	f, err := os.Open(path)
	if err != nil {
		return err
	}
	defer f.Close()

	_, err = io.Copy(w, f)

	return err
}

func mockClient(fixture string, handlerFunc func(w http.ResponseWriter, r *http.Request, j interface{})) *Client {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var j interface{}

		if r.Body != nil {
			b, _ := ioutil.ReadAll(r.Body)
			json.Unmarshal(b, &j)
		}

		handlerFunc(w, r, j)

		if err := writeFixture(w, fixture); err != nil {
			panic(err)
		}
	}))

	c := NewClient()
	c.BaseURL, _ = url.Parse(ts.URL)

	return c
}

func assertEqual(t *testing.T, expected interface{}, actual interface{}, desc string) {
	t.Helper()

	if actual != expected {
		t.Errorf("Expected %v of `%v`, got `%v`.", desc, expected, actual)
	}
}

func assertHeader(t *testing.T, expected string, header string, r *http.Request) {
	t.Helper()

	assertEqual(t, expected, r.Header.Get(header), header)
}

func assertAuthToken(t *testing.T, expected string, r *http.Request) {
	t.Helper()

	assertHeader(t, expected, "Auth-Token", r)
}

func assertContentType(t *testing.T, expected string, r *http.Request) {
	t.Helper()

	assertHeader(t, expected, "Content-Type", r)
}

func assertUserAgent(t *testing.T, expected string, r *http.Request) {
	t.Helper()

	assertHeader(t, expected, "User-Agent", r)
}

func assertMethod(t *testing.T, expected string, r *http.Request) {
	t.Helper()

	assertEqual(t, expected, r.Method, "request method")
}

func assertPath(t *testing.T, expected string, r *http.Request) {
	t.Helper()

	assertEqual(t, expected, r.URL.EscapedPath(), "request path")
}

func assertParam(t *testing.T, expected string, name string, r *http.Request) {
	t.Helper()

	assertEqual(t, expected, r.URL.Query().Get(name), name)
}

func assertJSON(t *testing.T, expected interface{}, path string, j interface{}) {
	t.Helper()

	var m map[string]interface{}
	m = j.(map[string]interface{})

	assertEqual(t, expected, m[path], path)
}
