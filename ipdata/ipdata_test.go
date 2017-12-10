package ipdata

import (
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

	_, err = c.GetMyIpData()
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

	for _, addr := range addrs {
		_, err := c.GetIpData(addr)
		if err != nil {
			t.Error(err)
		}
		//println(d.IP, d.CountryName, d.ContinentName, d.Organisation)
	}
}
