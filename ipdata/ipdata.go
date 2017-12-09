package ipdata

import (
	"bytes"
	"encoding/json"
	"io"
	"net"
	"net/http"
	"net/url"
)

type Client struct {
	baseUrl    *url.URL
	httpClient *http.Client
	UserAgent  string
	ApiKey     string
}

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
	Asn            string  `json:"asn"`
	Organisation   string  `json:"organisation"`
	Postal         string  `json:"postal"`
	Currency       string  `json:"currency"`
	CurrencySymbol string  `json:"currency_symbol"`
	CallingCode    string  `json:"calling_code"`
	Flag           string  `json:"flag"`
	TimeZone       string  `json:"time_zone"`
}

func NewClient(httpClient *http.Client) (*Client, error) {
	if httpClient == nil {
		httpClient = http.DefaultClient
	}

	baseUrl, err := url.Parse("https://api.ipdata.co")
	if err != nil {
		return nil, err
	}

	client := &Client{
		baseUrl:    baseUrl,
		httpClient: httpClient,
	}
	return client, nil
}

func (c *Client) newRequest(method, path string, body interface{}) (*http.Request, error) {
	rel := &url.URL{Path: path}
	u := c.baseUrl.ResolveReference(rel)

	var buf io.ReadWriter
	if body != nil {
		buf = new(bytes.Buffer)
		err := json.NewEncoder(buf).Encode(body)
		if err != nil {
			return nil, err
		}
	}
	req, err := http.NewRequest(method, u.String(), buf)
	if err != nil {
		return nil, err
	}
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	req.Header.Set("Accept", "application/json")
	req.Header.Set("User-Agent", c.UserAgent)
	req.Header.Set("Api-Key", c.ApiKey)

	return req, nil
}

func (c *Client) do(req *http.Request, v interface{}) (*http.Response, error) {
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	err = json.NewDecoder(resp.Body).Decode(v)
	return resp, err
}

func (c *Client) getSingleIpData(ip string) (*Data, error) {
	req, err := c.newRequest("GET", "/"+ip, nil)
	if err != nil {
		return nil, err
	}

	var ipData *Data = &Data{}
	_, err = c.do(req, ipData)
	return ipData, err
}

func (c *Client) GetMyIpData() (*Data, error) {
	return c.getSingleIpData("")
}

func (c *Client) GetIpData(addrs ...net.IP) ([]*Data, []error) {
	var ipDataList = make([]*Data, len(addrs))
	var errors []error

	for i, addr := range addrs {
		ipData, err := c.getSingleIpData(addr.String())
		if err != nil {
			errors = append(errors, err)
		}
		ipDataList[i] = ipData
	}

	return ipDataList, errors
}
