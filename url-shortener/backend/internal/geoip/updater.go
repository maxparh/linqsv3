package geoip

import (
	"compress/gzip"
	"context"
	"fmt"
	"github.com/oschwald/maxminddb-golang/v2"
	"io"
	"log"
	"net/http"
	"net/netip"
	"os"
	"path/filepath"
	"time"
)

var updateZone = time.FixedZone("Europe/Moscow", 3*60*60)

func nextUpdate(now time.Time) time.Time {
	now = now.In(updateZone)
	return time.Date(now.Year(), now.Month()+1, 1, 0, 0, 0, 0, updateZone)
}

// Stop waits for the worker before the reader can be closed.
func (r *Resolver) StartUpdater(path string) func() {
	ctx, cancel := context.WithCancel(context.Background())
	done := make(chan struct{})
	go func() {
		defer close(done)
		for {
			err := r.update(ctx, path, time.Now().In(updateZone))
			if ctx.Err() != nil {
				return
			}
			delay := time.Until(nextUpdate(time.Now()))
			if err != nil {
				log.Printf("GeoIP update failed; retaining existing database, retry in one hour: %v", err)
				if delay > time.Hour {
					delay = time.Hour
				}
			} else {
				log.Print("GeoIP database updated successfully")
			}
			timer := time.NewTimer(delay)
			select {
			case <-ctx.Done():
				timer.Stop()
				return
			case <-timer.C:
			}
		}
	}()
	return func() { cancel(); <-done }
}

func (r *Resolver) update(ctx context.Context, path string, now time.Time) error {
	ctx, cancel := context.WithTimeout(ctx, 2*time.Minute)
	defer cancel()
	url := "https://download.db-ip.com/free/dbip-country-lite-" + now.Format("2006-01") + ".mmdb.gz"
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return err
	}
	response, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer response.Body.Close()
	if response.StatusCode != http.StatusOK {
		return fmt.Errorf("DB-IP returned HTTP %d", response.StatusCode)
	}
	compressed, err := gzip.NewReader(io.LimitReader(response.Body, 64<<20))
	if err != nil {
		return err
	}
	defer compressed.Close()
	const maxSize = 64 << 20
	data, err := io.ReadAll(io.LimitReader(compressed, maxSize+1))
	if err != nil {
		return err
	}
	if len(data) > maxSize {
		return fmt.Errorf("GeoIP database exceeds size limit")
	}
	db, err := maxminddb.OpenBytes(data)
	if err != nil {
		return err
	}
	installed := false
	defer func() {
		if !installed {
			db.Close()
		}
	}()
	if err := db.Verify(); err != nil {
		return fmt.Errorf("invalid MMDB: %w", err)
	}
	var code string
	if err := db.Lookup(netip.MustParseAddr("8.8.8.8")).DecodePath(&code, "country", "iso_code"); err != nil || len(code) != 2 {
		return fmt.Errorf("MMDB does not contain expected country ISO data")
	}
	if err := ctx.Err(); err != nil {
		return err
	}
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return err
	}
	file, err := os.CreateTemp(filepath.Dir(path), ".geoip-*.mmdb")
	if err != nil {
		return err
	}
	defer os.Remove(file.Name())
	if _, err := file.Write(data); err != nil {
		file.Close()
		return err
	}
	if err := file.Sync(); err != nil {
		file.Close()
		return err
	}
	if err := file.Close(); err != nil {
		return err
	}
	if err := os.Rename(file.Name(), path); err != nil {
		return err
	}
	r.mu.Lock()
	old := r.db
	r.db = db
	installed = true
	if old != nil {
		old.Close()
	}
	r.mu.Unlock()
	return nil
}
