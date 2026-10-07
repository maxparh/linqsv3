package geoip

import (
	"net/http/httptest"
	"net/netip"
	"os"
	"testing"
)

func TestClientIPIgnoresUntrustedHeaders(t *testing.T) {
	r := httptest.NewRequest("GET", "/", nil)
	r.RemoteAddr = "[2001:4860:4860::8888]:12345"
	r.Header.Set("X-Real-IP", "1.1.1.1")
	r.Header.Set("X-Forwarded-For", "1.1.1.1")
	if got := ClientIP(r, "").String(); got != "2001:4860:4860::8888" {
		t.Fatalf("unexpected IP: %s", got)
	}
}

func TestCountryIgnoresLocalAddresses(t *testing.T) {
	r := &Resolver{}
	for _, value := range []string{"127.0.0.1", "10.0.0.1", "172.18.0.1", "192.168.1.1", "::1", "fc00::1", "::ffff:192.168.1.1"} {
		if got := r.Country(netip.MustParseAddr(value)); got != "" {
			t.Fatalf("%s: %s", value, got)
		}
	}
}

func TestDatabaseLookup(t *testing.T) {
	path := os.Getenv("GEOIP_TEST_DB")
	if path == "" {
		t.Skip("set GEOIP_TEST_DB to validate a downloaded database")
	}
	r, err := Open(path)
	if err != nil {
		t.Fatal(err)
	}
	defer r.Close()
	code := r.Country(netip.MustParseAddr("8.8.8.8"))
	if len(code) != 2 {
		t.Fatalf("expected ISO country code, got %q", code)
	}
}
