package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/basecamp/hey-cli/internal/output"
)

type todoCommand struct {
	cmd *cobra.Command
}

func newTodoCommand() *todoCommand {
	todoCommand := &todoCommand{}
	todoCommand.cmd = &cobra.Command{
		Use:   "todo",
		Short: "Read todos",
		Annotations: map[string]string{
			"agent_notes": "Subcommands: list.",
		},
	}

	todoCommand.cmd.AddCommand(newTodoListCommand().cmd)

	return todoCommand
}

type todoListCommand struct {
	cmd   *cobra.Command
	limit int
	all   bool
}

func newTodoListCommand() *todoListCommand {
	todoListCommand := &todoListCommand{}
	todoListCommand.cmd = &cobra.Command{
		Use:   "list",
		Short: "List todos",
		Example: `  hey todo list
  hey todo list --limit 10
  hey todo list --json`,
		RunE: todoListCommand.run,
	}

	todoListCommand.cmd.Flags().IntVar(&todoListCommand.limit, "limit", 0, "Maximum number of todos to show")
	todoListCommand.cmd.Flags().BoolVar(&todoListCommand.all, "all", false, "Fetch all results (override --limit)")

	return todoListCommand
}

func (c *todoListCommand) run(cmd *cobra.Command, args []string) error {
	if err := requireAuth(); err != nil {
		return err
	}

	resp, err := listPersonalRecordings(cmd.Context())
	if err != nil {
		return err
	}

	todos := filterRecordingsByType(resp, "Calendar::Todo")
	total := len(todos)
	if c.limit > 0 && !c.all && len(todos) > c.limit {
		todos = todos[:c.limit]
	}
	notice := output.TruncationNotice(len(todos), total)

	if writer.IsStyled() {
		if len(todos) == 0 {
			fmt.Fprintln(cmd.OutOrStdout(), "No todos.")
			return nil
		}

		table := newTable(cmd.OutOrStdout())
		table.addRow([]string{"ID", "Title", "Date", "Done"})
		for _, todo := range todos {
			done := ""
			if !todo.CompletedAt.IsZero() {
				done = "yes"
			}
			table.addRow([]string{fmt.Sprintf("%d", todo.Id), todo.Title, formatDate(todo.StartsAt), done})
		}
		table.print()
		if notice != "" {
			fmt.Fprintln(cmd.OutOrStdout(), notice)
		}
		return nil
	}

	return writeOK(todos,
		output.WithSummary(fmt.Sprintf("%d todos", len(todos))),
		output.WithNotice(notice),
	)
}
