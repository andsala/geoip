package ipdata

import (
	"os"
	"testing"
)

func TestNewClient(t *testing.T) {
	_, err := NewClient(nil)
	if err != nil {
		t.Error("Unexpected error happened: ", err)
	}
}

func TestClient_GetMyIpData(t *testing.T) {
	c, err := NewClient(nil)
	if err != nil {
		t.Error("Unexpected error happened: ", err)
	}
	c.APIKey = os.Getenv("GEOIP_API_KEY")

	_, err = c.GetMyIPData()
	if err != nil {
		t.Error("Request error:", err)
	}
}

func TestClient_GetIpData(t *testing.T) {
	var addrs = []string{
		"1.1.1.1",
		"2.2.2.2",
		"3.3.3.3",
		"8.8.8.8",
	}

	c, err := NewClient(nil)
	if err != nil {
		t.Error("Unexpected error happened: ", err)
	}
	c.APIKey = os.Getenv("GEOIP_API_KEY")

	for _, addr := range addrs {
		_, err := c.GetIPData(addr)
		if err != nil {
			t.Error(addr, err)
		}
	}
}

func TestClient_GetIpData_Fail(t *testing.T) {
	var addrs = []string{
		"256.1.1.1",
		"10.2.2.2",
		"192.168.2.3",
	}

	c, err := NewClient(nil)
	if err != nil {
		t.Error("Unexpected error happened: ", err)
	}
	c.APIKey = os.Getenv("GEOIP_API_KEY")

	for _, addr := range addrs {
		_, err := c.GetIPData(addr)
		if err == nil {
			t.Errorf("Request with ip '%v' should fail", addr)
		}
	}
}

func TestClient_GetIpData_Fail_MissingApiKey(t *testing.T) {
	var addrs = []string{
		"256.1.1.1",
		"10.2.2.2",
		"192.168.2.3",
	}

	c, err := NewClient(nil)
	if err != nil {
		t.Error("Unexpected error happened: ", err)
	}
	c.APIKey = ""

	for _, addr := range addrs {
		_, err := c.GetIPData(addr)
		if err == nil {
			t.Errorf("Request with ip '%v' should fail", addr)
		}
	}
}
