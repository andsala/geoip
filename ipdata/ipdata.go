package ipdata

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strings"
)

type boolean bool

// Client represent the wrapper for ipdata.co
type Client struct {
	baseURL    string
	httpClient *http.Client
	UserAgent  string
	APIKey     string
}

// Data represent the information retrieved from ipdata.com
type Data struct {
	IP            string     `json:"ip"`
	City          string     `json:"city"`
	Region        string     `json:"region"`
	RegionCode    string     `json:"region_code"`
	CountryName   string     `json:"country_name"`
	CountryCode   string     `json:"country_code"`
	ContinentName string     `json:"continent_name"`
	ContinentCode string     `json:"continent_code"`
	Latitude      float32    `json:"latitude"`
	Longitude     float32    `json:"longitude"`
	ASN           string     `json:"asn"`
	Organisation  string     `json:"organisation"`
	Postal        string     `json:"postal"`
	CallingCode   string     `json:"calling_code"`
	Flag          string     `json:"flag"`
	EmojiFlag     string     `json:"emoji_flag"`
	EmojiUnicode  string     `json:"emoji_unicode"`
	IsEU          boolean    `json:"is_eu"`
	Languages     []Language `json:"languages"`
	Currency      Currency   `json:"currency"`
	TimeZone      TimeZone   `json:"time_zone"`
	Threat        Threat     `json:"threat"`
	JSON          *string
}

// Language information retrieved from ipdata.co
type Language struct {
	Name   string `json:"name"`
	Native string `json:"native"`
}

// Currency information retrieved from ipdata.co
type Currency struct {
	Name   string `json:"name"`
	Code   string `json:"code"`
	Symbol string `json:"symbol"`
	Native string `json:"native"`
	Plural string `json:"plural"`
}

// TimeZone information retrieved from ipdata.co
type TimeZone struct {
	Name        string  `json:"name"`
	Abbr        string  `json:"abbr"`
	Offset      string  `json:"offset"`
	IsDST       boolean `json:"is_dst"`
	CurrentTime string  `json:"current_time"`
}

// Threat information retrieved from ipdata.co
type Threat struct {
	IsTor           boolean `json:"is_tor"`
	IsProxy         boolean `json:"is_proxy"`
	IsAnonymous     boolean `json:"is_anonymous"`
	IsKnownAttacker boolean `json:"is_known_attacker"`
	IsKnownAbuser   boolean `json:"is_known_abuser"`
	IsThreat        boolean `json:"is_threat"`
	IsBogon         boolean `json:"is_bogon"`
}

// Error message received form ipdata.co
type Error struct {
	Message string `json:"message"`
}

// NewClient generates a new Client.
// If nil is passed, http.DefaultClient will be used.
func NewClient(httpClient *http.Client) (*Client, error) {
	if httpClient == nil {
		httpClient = http.DefaultClient
	}

	client := &Client{
		baseURL:    "https://api.ipdata.co",
		httpClient: httpClient,
	}
	return client, nil
}

func (c *Client) newRequest(method, path string) (*http.Request, error) {
	params := url.Values{}
	params.Add("api-key", c.APIKey)

	u, _ := url.ParseRequestURI(c.baseURL)
	u.Path = path
	u.RawQuery = params.Encode()

	req, err := http.NewRequest(method, u.String(), nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("User-Agent", c.UserAgent)

	return req, nil
}

func (c *Client) do(req *http.Request, v interface{}) (*http.Response, *string, error) {
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, nil, err
	}
	defer func() {
		err := resp.Body.Close()
		if err != nil {
			log.Fatal(err)
		}
	}()

	buf := new(bytes.Buffer)
	_, err = buf.ReadFrom(resp.Body)
	if err != nil {
		return nil, nil, err
	}
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
	if err != nil || resp.StatusCode != 200 {
		var errorResponse = &Error{}
		_ = json.Unmarshal([]byte(*body), errorResponse)

		switch resp.StatusCode {
		case 400: // Bad Request
			return nil, errors.New(errorResponse.Message)
		case 401: // Unauthorized
			return nil, errors.New(fmt.Sprintf("Unauthorized: %v", errorResponse.Message))
		case 429: // Too Many Requests
			return nil, errors.New("you have exceeded requests limit. See https://ipdata.co")
		default:
			errorString := "Unknown Error"
			if err != nil {
				errorString = fmt.Sprintf("%v: %v", errorString, err.Error())
			} else if len(errorResponse.Message) > 0 {
				errorString = fmt.Sprintf("%v: %v", errorString, errorResponse.Message)
			}
			return nil, errors.New(errorString)
		}
	}
	data.JSON = body

	return data, err
}

// GetMyIPData retrieves information about your public IP address.
func (c *Client) GetMyIPData() (*Data, error) {
	return c.GetIPData("")
}

func (bit boolean) UnmarshalJSON(data []byte) error {
	asString := strings.ToLower(string(data))
	if asString == "true" {
		bit = true
	} else {
		bit = false
	}
	return nil
}
