package smoke_test

import "testing"

func TestTimetrackCurrent(t *testing.T) {
	heyJSON(t, "timetrack", "current")
}

func TestTimetrackList(t *testing.T) {
	resp := heyJSON(t, "timetrack", "list")
	type TimeTrack struct {
		ID int `json:"id"`
	}
	tracks := dataAs[[]TimeTrack](t, resp)
	if len(tracks) > 0 {
		html := fetchHTML(t, baseURL+"/calendar")
		if len(html) == 0 {
			t.Error("calendar page returned empty HTML")
		}
	}
}

func TestTimetrackListLimit(t *testing.T) {
	resp := heyJSON(t, "timetrack", "list", "--limit", "5")
	type TimeTrack struct {
		ID int `json:"id"`
	}
	tracks := dataAs[[]TimeTrack](t, resp)
	if len(tracks) > 5 {
		t.Errorf("expected at most 5 time tracks with --limit 5, got %d", len(tracks))
	}
}

func TestTimetrackListAll(t *testing.T) {
	heyJSON(t, "timetrack", "list", "--all")
}
