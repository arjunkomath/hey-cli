---
name: hey
description: Reads HEY email, mailboxes, drafts, calendars, todos, time tracks, and journal entries through the read-only HEY CLI. Use for HEY lookup and reading tasks.
argument-hint: "[command] [args...]"
---

# HEY read-only workflow

Use this CLI only to retrieve HEY data. It cannot send email or modify mailbox, calendar, todo, habit, time-track, or journal state.

## Rules

1. Use `--json` for structured output.
2. Run `hey auth login` when authentication is required.
3. Use `--html` when raw HTML content is needed.

## Commands

| Task | Command |
|------|---------|
| List mailboxes | `hey boxes --json` |
| List emails in a mailbox | `hey box imbox --json` |
| Read an email thread | `hey threads <topic_id> --json` |
| List drafts | `hey drafts --json` |
| List calendars | `hey calendars --json` |
| List calendar recordings | `hey recordings <calendar_id> --json` |
| List todos | `hey todo list --json` |
| Show current time track | `hey timetrack current --json` |
| List time tracks | `hey timetrack list --json` |
| List journal entries | `hey journal list --json` |
| Read a journal entry | `hey journal read 2026-03-15 --json` |
| Check authentication | `hey auth status` |

`hey box` accepts a mailbox name such as `imbox`, `feedbox`, `trailbox`, `asidebox`, `laterbox`, or `bubblebox`, as well as a numeric mailbox ID. Use a posting's topic ID or the ID in its `app_url` with `hey threads`.
