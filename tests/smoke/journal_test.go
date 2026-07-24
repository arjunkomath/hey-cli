package smoke_test

import "testing"

func TestJournalList(t *testing.T) {
	resp := heyJSON(t, "journal", "list")
	type JournalEntry struct {
		Date string `json:"date"`
	}
	entries := dataAs[[]JournalEntry](t, resp)
	if len(entries) > 0 && entries[0].Date != "" {
		html := fetchHTML(t, baseURL+"/calendar/days/"+entries[0].Date+"/journal_entry/edit")
		if len(html) == 0 {
			t.Errorf("journal entry page for %s returned empty HTML", entries[0].Date)
		}
	}
}

func TestJournalReadToday(t *testing.T) {
	_, _, code := hey(t, "journal", "read", "--json")
	if code != 0 {
		t.Error("expected journal read (today) to succeed")
	}
}

func TestJournalListLimit(t *testing.T) {
	resp := heyJSON(t, "journal", "list", "--limit", "2")
	type JournalEntry struct {
		Date string `json:"date"`
	}
	entries := dataAs[[]JournalEntry](t, resp)
	if len(entries) > 2 {
		t.Errorf("expected at most 2 entries with --limit 2, got %d", len(entries))
	}
}

func TestJournalListAll(t *testing.T) {
	heyJSON(t, "journal", "list", "--all")
}
