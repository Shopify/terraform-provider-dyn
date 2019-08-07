package dyn

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"

	"github.com/Shopify/go-dyn/pkg/version"
)

// BaseURL is the Dyn API base URL.
const BaseURL = "https://api.dynect.net/"

// Client is used to manage a Dyn API session.
type Client struct {
	BaseURL   *url.URL
	UserAgent string
	Logger    *log.Logger

	httpClient *http.Client
	token      string
}

// NewClient creates a new API client.
func NewClient() *Client {
	baseURL, _ := url.Parse(BaseURL)

	c := &Client{
		BaseURL:   baseURL,
		UserAgent: fmt.Sprintf("go-dyn/%v", version.VERSION),

		httpClient: http.DefaultClient,
	}

	return c
}

func (c *Client) delete(resource string, requestData interface{}) error {
	return c.perform(http.MethodDelete, resource, nil, requestData, nil)
}

func (c *Client) get(resource string, params url.Values, responseData interface{}) error {
	return c.perform(http.MethodGet, resource, params, nil, responseData)
}

func (c *Client) post(resource string, requestData interface{}, responseData interface{}) error {
	return c.perform(http.MethodPost, resource, nil, requestData, responseData)
}

func (c *Client) put(resource string, requestData interface{}, responseData interface{}) error {
	return c.perform(http.MethodPut, resource, nil, requestData, responseData)
}

// perform does the actual work for the request/response cycle.
func (c *Client) perform(method, resource string, params url.Values, requestData interface{}, responseData interface{}) error {
	url := c.buildURL(resource, params)

	body, err := c.marshalJSON(requestData)
	if err != nil {
		return err
	}

	if c.Logger != nil {
		c.Logger.Println(method, url, body)
	}

	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return err
	}

	if c.token != "" {
		req.Header.Set("Auth-Token", c.token)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Add("User-Agent", c.UserAgent)

	// log.Printf("CALLING: %#v, %#v, %#v\n", method, url, requestData)
	// if jj, err := json.Marshal(requestData); err == nil {
	// log.Printf("BODY: %s", string(jj))
	// }
	// log.Printf("REQ: %#v\n", req)
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if c.Logger != nil {
		c.Logger.Println(resp.StatusCode, "RESPONSE")
	}
	// log.Printf("RESPONSE: %#v\n", resp)

	// buf := new(bytes.Buffer)
	// buf.ReadFrom(resp.Body)
	// s := buf.String()
	// log.Printf("RESPONSE BODY: %s\n", s)

	switch resp.StatusCode {
	case http.StatusOK:
		if responseData == nil {
			if c.Logger != nil {
				c.decodeError(resp)
			}

			return nil
		}

		return c.decodeJSON(resp.Body, responseData)
	}

	return c.decodeError(resp)
}

// buildURL creates a resource URL relative to the base URL.
func (c *Client) buildURL(resource string, params url.Values) string {
	path := fmt.Sprintf("/REST/%s", resource)

	rel := &url.URL{Path: path}
	rel.RawQuery = params.Encode()

	url := c.BaseURL.ResolveReference(rel)

	return url.String()
}

// marshalJSON converts a request object into JSON.
func (c *Client) marshalJSON(data interface{}) (io.Reader, error) {
	if data == nil {
		return nil, nil
	}

	b, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}

	return bytes.NewBuffer(b), nil
}

// decodeJSON converts JSON into a response object.
func (c *Client) decodeJSON(r io.Reader, data interface{}) error {
	if c.Logger != nil {
		b, _ := ioutil.ReadAll(r)
		c.Logger.Println("decodeJSON: body", string(b))
		r = bytes.NewBuffer(b)
	}

	return json.NewDecoder(r).Decode(data)
}

func (c *Client) decodeError(resp *http.Response) error {
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("%v - unable to read response body", resp.Status)
	}

	if c.Logger != nil {
		c.Logger.Println("decodeError: body", string(body))
	}

	var h responseHeader

	if err := json.Unmarshal(body, &h); err != nil {
		return fmt.Errorf("%v: %v [%v]", err, resp.Status, body)
	}

	for _, m := range h.Messages {
		if m.Level == responseMessageError {
			return m
		}
	}

	return fmt.Errorf("%v %v", resp.Status, body)
}
