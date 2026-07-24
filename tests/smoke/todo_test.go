package smoke_test

import "testing"

func TestTodoListAll(t *testing.T) {
	heyJSON(t, "todo", "list", "--all")
}

func TestTodoListLimit(t *testing.T) {
	resp := heyJSON(t, "todo", "list", "--limit", "2")
	type Todo struct {
		ID int `json:"id"`
	}
	todos := dataAs[[]Todo](t, resp)
	if len(todos) > 2 {
		t.Errorf("expected at most 2 todos with --limit 2, got %d", len(todos))
	}
}
