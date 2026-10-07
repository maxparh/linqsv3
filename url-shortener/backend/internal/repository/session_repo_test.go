package repository

import (
	"context"
	"database/sql"
	"os"
	"strings"
	"testing"
	"time"
)

func TestClickEvents(t *testing.T) {
	url := os.Getenv("ANALYTICS_TEST_DATABASE_URL")
	if url == "" {
		t.Skip("set ANALYTICS_TEST_DATABASE_URL for isolated temporary-table checks")
	}
	db, err := sql.Open("pgx", url)
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	db.SetMaxOpenConns(1)
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()
	exec := func(query string) {
		t.Helper()
		if _, err := db.ExecContext(ctx, query); err != nil {
			t.Fatal(err)
		}
	}
	exec("BEGIN")
	defer db.ExecContext(context.Background(), "ROLLBACK")
	exec("SET LOCAL TIME ZONE 'UTC'")
	exec("CREATE TEMP TABLE links (id INTEGER PRIMARY KEY, user_id INTEGER)")
	for _, filename := range []string{"004_analytics_sessions.sql", "007_analytics_click_events.sql"} {
		data, err := os.ReadFile("../../migrations/" + filename)
		if err != nil {
			t.Fatal(err)
		}
		query := strings.ReplaceAll(string(data), "CREATE TABLE", "CREATE TEMP TABLE")
		query = strings.ReplaceAll(query, "CREATE VIEW", "CREATE TEMP VIEW")
		exec(query)
	}
	exec("INSERT INTO links VALUES (1, 1), (2, 1)")
	exec(`INSERT INTO analytics_sessions (session_id, link_id, page_views, started_at)
 VALUES ('legacy', 1, 7, NOW() - INTERVAL '5 days')`)
	repo := NewSessionRepository(db)
	for _, item := range []struct {
		link    int
		session string
	}{{1, "visitor-link-1"}, {1, "visitor-link-1"}, {2, "visitor-link-2"}} {
		if err := repo.RecordSession(ctx, item.link, item.session, "hash", "NO", "desktop", "test"); err != nil {
			t.Fatal(err)
		}
	}
	// Backdate the session, not the events: today's clicks must stay today.
	exec("UPDATE analytics_sessions SET started_at = NOW() - INTERVAL '5 days' WHERE event_based")
	overview, err := repo.GetOverviewStats(ctx, 1, 1)
	if err != nil {
		t.Fatal(err)
	}
	if overview.TotalClicks != 3 || overview.UniqueClicks != 2 {
		t.Fatalf("unexpected overview: %+v", overview)
	}
	overview, err = repo.GetOverviewStats(ctx, 1, 30)
	if err != nil {
		t.Fatal(err)
	}
	if overview.TotalClicks != 10 {
		t.Fatalf("legacy counts duplicated/lost: %d", overview.TotalClicks)
	}
	var link, views int
	if err := db.QueryRowContext(ctx, "SELECT link_id, page_views FROM analytics_sessions WHERE session_id='visitor-link-1'").Scan(&link, &views); err != nil {
		t.Fatal(err)
	}
	if link != 1 || views != 2 {
		t.Fatalf("wrong link/session: %d/%d", link, views)
	}
	locations, err := repo.GetTopLocations(ctx, 1, 1, 5)
	if err != nil {
		t.Fatal(err)
	}
	if len(locations) != 1 || locations[0].CountryCode != "NO" {
		t.Fatalf("unexpected locations: %+v", locations)
	}
	devices, err := repo.GetDeviceStats(ctx, 1, 1)
	if err != nil {
		t.Fatal(err)
	}
	if len(devices) != 1 || devices[0].Value != 3 {
		t.Fatalf("unexpected devices: %+v", devices)
	}
	// Exercise chronological ordering and the Moscow midnight boundary.
	exec("UPDATE analytics_sessions SET started_at = NOW() - INTERVAL '100 days' WHERE NOT event_based")
	exec(`UPDATE analytics_click_events SET clicked_at = CASE
 WHEN id = 1 THEN date_trunc('month', NOW()) - INTERVAL '1 day' + INTERVAL '12 hours'
 ELSE date_trunc('month', NOW()) - INTERVAL '3 hours' END`)
	graph, err := repo.GetClicksOverTime(ctx, 1, 90)
	if err != nil {
		t.Fatal(err)
	}
	if len(graph.Labels) < 2 || graph.Labels[1] != "01."+time.Now().UTC().Format("01") {
		t.Fatalf("incorrect date ordering: %+v", graph)
	}
}
