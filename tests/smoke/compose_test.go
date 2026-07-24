package smoke_test

import (
	"fmt"
	"testing"
)

func TestThreads(t *testing.T) {
	resp := heyJSON(t, "box", "imbox")
	type Posting struct {
		AppURL string `json:"app_url"`
	}
	type BoxResp struct {
		Postings []Posting `json:"postings"`
	}
	data := dataAs[BoxResp](t, resp)
	if len(data.Postings) == 0 {
		t.Fatal("no postings in imbox to test threads")
	}

	topicID := extractTopicID(data.Postings[0].AppURL)
	if topicID == "" {
		t.Fatalf("could not extract topic ID from app_url: %s", data.Postings[0].AppURL)
	}

	threadsResp := heyJSON(t, "threads", topicID)
	type Entry struct {
		Summary string `json:"summary"`
	}
	entries := dataAs[[]Entry](t, threadsResp)
	if len(entries) == 0 {
		t.Error("expected at least one entry in thread")
	}

	html := fetchHTML(t, fmt.Sprintf("%s/topics/%s", baseURL, topicID))
	if len(entries) > 0 && entries[0].Summary != "" {
		assertContains(t, html, entries[0].Summary)
	}
}

func TestDrafts(t *testing.T) {
	resp := heyJSON(t, "drafts")
	if resp.Data == nil {
		return
	}
	type Draft struct {
		ID int `json:"id"`
	}
	_ = dataAs[[]Draft](t, resp)
}

func TestDraftsLimit(t *testing.T) {
	resp := heyJSON(t, "drafts", "--limit", "2")
	if resp.Data == nil {
		return
	}
	type Draft struct {
		ID int `json:"id"`
	}
	drafts := dataAs[[]Draft](t, resp)
	if len(drafts) > 2 {
		t.Errorf("expected at most 2 drafts with --limit 2, got %d", len(drafts))
	}
}

func TestDraftsAll(t *testing.T) {
	heyJSON(t, "drafts", "--all")
}

func TestThreadsNoArgument(t *testing.T) {
	heyFail(t, "threads", "--json")
}
