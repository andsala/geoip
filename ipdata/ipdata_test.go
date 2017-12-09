package ipdata

import (
	"net"
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
	var addrs = []net.IP{
		net.IPv4(1, 1, 1, 1),
		net.IPv4(2, 2, 2, 2),
		net.IPv4(3, 3, 3, 3),
		net.IPv4(8, 8, 8, 8),
	}

	c, err := NewClient(nil)
	if err != nil {
		t.Error("Unexpected error happened: ", err)
	}

	data, errs := c.GetIpData(addrs...)
	if len(errs) > 0 {
		t.Error(errs)
	}
	if len(data) != len(addrs) {
		t.Error("Expected", len(addrs), "results, got", len(data))
	}

	//for _, d := range data {
	//	println(d.IP, d.CountryName, d.ContinentName, d.Organisation)
	//}
}
