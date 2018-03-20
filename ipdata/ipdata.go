package ipdata

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/url"
)

// Client represent the wrapper for ipdata.co
type Client struct {
	baseURL    *url.URL
	httpClient *http.Client
	UserAgent  string
	APIKey     string
}

// Data represent the information retrieved from ipdata.com
type Data struct {
	IP             string  `json:"ip"`
	City           string  `json:"city"`
	Region         string  `json:"region"`
	CountryName    string  `json:"country_name"`
	CountryCode    string  `json:"country_code"`
	ContinentName  string  `json:"continent_name"`
	ContinentCode  string  `json:"continent_code"`
	Latitude       float32 `json:"latitude"`
	Longitude      float32 `json:"longitude"`
	ASN            string  `json:"asn"`
	Organisation   string  `json:"organisation"`
	Postal         string  `json:"postal"`
	Currency       string  `json:"currency"`
	CurrencySymbol string  `json:"currency_symbol"`
	CallingCode    string  `json:"calling_code"`
	Flag           string  `json:"flag"`
	TimeZone       string  `json:"time_zone"`
	JSON           *string
}

// NewClient generates a new Client.
// If nil is passed, http.DefaultClient will be used.
func NewClient(httpClient *http.Client) (*Client, error) {
	if httpClient == nil {
		httpClient = http.DefaultClient
	}

	baseURL, err := url.Parse("https://api.ipdata.co")
	if err != nil {
		return nil, err
	}

	client := &Client{
		baseURL:    baseURL,
		httpClient: httpClient,
	}
	return client, nil
}

func (c *Client) newRequest(method, path string) (*http.Request, error) {
	rel := &url.URL{Path: path}
	u := c.baseURL.ResolveReference(rel)

	req, err := http.NewRequest(method, u.String(), nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("User-Agent", c.UserAgent)
	req.Header.Set("Api-Key", c.APIKey)

	return req, nil
}

func (c *Client) do(req *http.Request, v interface{}) (*http.Response, *string, error) {
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, nil, err
	}
	defer resp.Body.Close()

	buf := new(bytes.Buffer)
	buf.ReadFrom(resp.Body)
	body := buf.String()

	err = json.Unmarshal(buf.Bytes(), v)

	return resp, &body, err
}

// GetIPData retrieves information about the ip from ipdata.co and
// returns a valid Data if no error occurs.
func (c *Client) GetIPData(ip string) (*Data, error) {
	req, err := c.newRequest("GET", "/"+ip)
	if err != nil {
		return nil, err
	}

	var data = &Data{}
	resp, body, err := c.do(req, data)
	if err != nil {
		switch resp.StatusCode {
		case 400: // Bad Request
			return nil, errors.New(*body)
		case 429: // Too Many Requests
			return nil, errors.New("you have exceeded requests limit. See https://ipdata.co")
		default:
			return nil, err
		}
	}
	data.JSON = body

	return data, err
}

// GetMyIPData retrieves information about your public IP address.
func (c *Client) GetMyIPData() (*Data, error) {
	return c.GetIPData("")
}
