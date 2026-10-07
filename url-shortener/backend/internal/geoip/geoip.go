package geoip

import (
	"context"
	"net"
	"net/http"
	"net/netip"
	"strings"
	"sync"
	"time"

	"github.com/oschwald/maxminddb-golang/v2"
)

type Resolver struct {
	mu sync.RWMutex
	db *maxminddb.Reader
}

func Open(path string) (*Resolver, error) {
	db, err := maxminddb.Open(path)
	if err != nil {
		return nil, err
	}
	return &Resolver{db: db}, nil
}

func (r *Resolver) Country(ip netip.Addr) string {
	if !ip.IsValid() {
		return ""
	}
	ip = ip.Unmap()
	if !ip.IsGlobalUnicast() || ip.IsPrivate() {
		return ""
	}
	r.mu.RLock()
	defer r.mu.RUnlock()
	if r.db == nil {
		return ""
	}
	var record struct {
		Country struct {
			ISOCode string `maxminddb:"iso_code"`
		} `maxminddb:"country"`
	}
	if err := r.db.Lookup(ip).Decode(&record); err != nil {
		return ""
	}
	return record.Country.ISOCode
}

func (r *Resolver) Close() error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if r.db == nil {
		return nil
	}
	err := r.db.Close()
	r.db = nil
	return err
}

// Only the configured reverse proxy may supply the client address.
func ClientIP(r *http.Request, proxyHost string) netip.Addr {
	host, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		host = r.RemoteAddr
	}
	peer, err := netip.ParseAddr(host)
	if err != nil {
		return netip.Addr{}
	}
	peer = peer.Unmap()
	if proxyHost == "" || r.Header.Get("X-Real-IP") == "" {
		return peer
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	proxies, err := net.DefaultResolver.LookupNetIP(ctx, "ip", proxyHost)
	if err != nil {
		return peer
	}
	for _, proxy := range proxies {
		if proxy.Unmap() == peer {
			if client, err := netip.ParseAddr(strings.TrimSpace(r.Header.Get("X-Real-IP"))); err == nil {
				return client.Unmap()
			}
		}
	}
	return peer
}
