# Read-only API Coverage

Mapping of HEY API endpoints used by the CLI. All data-access endpoints are read-only and use the HEY SDK (`hey-sdk/go`).

| Endpoint | Method | Client | CLI Command | Status |
|----------|--------|--------|-------------|--------|
| `/boxes.json` | GET | SDK `Boxes().List` | `hey boxes` | covered |
| `/boxes/{id}.json` | GET | SDK `Boxes().Get` | `hey box <id>` | covered |
| `/imbox.json` | GET | SDK `Boxes().GetImbox` | `hey box imbox` | covered |
| `/feedbox.json` | GET | SDK `Boxes().GetFeedbox` | `hey box feedbox` | covered |
| `/trailbox.json` | GET | SDK `Boxes().GetTrailbox` | `hey box trailbox` | covered |
| `/asidebox.json` | GET | SDK `Boxes().GetAsidebox` | `hey box asidebox` | covered |
| `/laterbox.json` | GET | SDK `Boxes().GetLaterbox` | `hey box laterbox` | covered |
| `/bubblebox.json` | GET | SDK `Boxes().GetBubblebox` | `hey box bubblebox` | covered |
| `/calendars.json` | GET | SDK `Calendars().List` | `hey calendars` | covered |
| `/calendars/{id}/recordings.json` | GET | SDK `Calendars().GetRecordings` | `hey recordings <calendar-id>`, `hey todo list`, `hey timetrack list`, `hey journal list` | covered |
| `/topics/{id}/entries` | GET (HTML) | SDK `GetHTML` | `hey threads <id>` | gap: SDK Entry lacks body |
| `/entries/drafts.json` | GET | SDK `Entries().ListDrafts` | `hey drafts` | covered |
| `/calendar/days/{date}/journal_entry.json` | GET | SDK `Journal().Get` | `hey journal read [date]` | partial: falls back to legacy |
| `/calendar/days/{date}/journal_entry/edit` | GET (HTML) | SDK `Journal().GetContent` | `hey journal read [date]` | fallback for 204 response |
| `/calendar/ongoing_time_track.json` | GET | SDK `TimeTracks().GetOngoing` | `hey timetrack current` | covered |
