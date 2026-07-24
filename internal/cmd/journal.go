package cmd

import (
	"fmt"
	"time"

	"github.com/spf13/cobra"

	"github.com/basecamp/hey-cli/internal/htmlutil"
	"github.com/basecamp/hey-cli/internal/output"
)

type journalCommand struct {
	cmd *cobra.Command
}

func newJournalCommand() *journalCommand {
	journalCommand := &journalCommand{}
	journalCommand.cmd = &cobra.Command{
		Use:   "journal",
		Short: "Read journal entries",
		Annotations: map[string]string{
			"agent_notes": "Subcommands: list, read. Read defaults to today.",
		},
	}

	journalCommand.cmd.AddCommand(newJournalListCommand().cmd)
	journalCommand.cmd.AddCommand(newJournalReadCommand().cmd)

	return journalCommand
}

// list

type journalListCommand struct {
	cmd   *cobra.Command
	limit int
	all   bool
}

func newJournalListCommand() *journalListCommand {
	journalListCommand := &journalListCommand{}
	journalListCommand.cmd = &cobra.Command{
		Use:   "list",
		Short: "List journal entries",
		Example: `  hey journal list
  hey journal list --limit 10
  hey journal list --json`,
		RunE: journalListCommand.run,
	}

	journalListCommand.cmd.Flags().IntVar(&journalListCommand.limit, "limit", 0, "Maximum number of entries to show")
	journalListCommand.cmd.Flags().BoolVar(&journalListCommand.all, "all", false, "Fetch all results (override --limit)")

	return journalListCommand
}

func (c *journalListCommand) run(cmd *cobra.Command, args []string) error {
	if err := requireAuth(); err != nil {
		return err
	}

	ctx := cmd.Context()
	resp, err := listPersonalRecordings(ctx)
	if err != nil {
		return err
	}

	entries := filterRecordingsByType(resp, "Calendar::JournalEntry")

	total := len(entries)
	if c.limit > 0 && !c.all && len(entries) > c.limit {
		entries = entries[:c.limit]
	}
	notice := output.TruncationNotice(len(entries), total)

	if writer.IsStyled() {
		if len(entries) == 0 {
			fmt.Fprintln(cmd.OutOrStdout(), "No journal entries.")
			return nil
		}

		table := newTable(cmd.OutOrStdout())
		table.addRow([]string{"ID", "Date", "Preview"})
		for _, e := range entries {
			table.addRow([]string{fmt.Sprintf("%d", e.Id), formatDate(e.StartsAt), truncate(e.Content, 60)})
		}
		table.print()
		if notice != "" {
			fmt.Fprintln(cmd.OutOrStdout(), notice)
		}
		return nil
	}

	return writeOK(entries,
		output.WithSummary(fmt.Sprintf("%d journal entries", len(entries))),
		output.WithNotice(notice),
		output.WithBreadcrumbs(output.Breadcrumb{
			Action:      "read",
			Command:     "hey journal read [date]",
			Description: "Read a journal entry",
		}),
	)
}

// read

type journalReadCommand struct {
	cmd *cobra.Command
}

func newJournalReadCommand() *journalReadCommand {
	journalReadCommand := &journalReadCommand{}
	journalReadCommand.cmd = &cobra.Command{
		Use:   "read [date]",
		Short: "Read a journal entry (default: today)",
		Example: `  hey journal read
  hey journal read 2024-01-15
  hey journal read --html
  hey journal read --json`,
		RunE: journalReadCommand.run,
		Args: cobra.MaximumNArgs(1),
	}

	return journalReadCommand
}

func (c *journalReadCommand) run(cmd *cobra.Command, args []string) error {
	if err := requireAuth(); err != nil {
		return err
	}

	date := time.Now().Format("2006-01-02")
	if len(args) > 0 {
		date = args[0]
	}

	ctx := cmd.Context()
	content, err := sdk.Journal().GetContent(ctx, date)
	if err != nil {
		return convertSDKError(err)
	}

	if content == "" {
		if writer.IsStyled() {
			fmt.Fprintf(cmd.OutOrStdout(), "Journal — %s\n\n(empty)\n", date)
			return nil
		}
		return writeOK(nil, output.WithSummary(fmt.Sprintf("No journal entry for %s", date)))
	}

	if writer.IsStyled() {
		w := cmd.OutOrStdout()
		if htmlOutput {
			fmt.Fprintln(w, content)
			return nil
		}

		fmt.Fprintf(w, "Journal — %s\n\n", date)
		fmt.Fprintln(w, htmlutil.ToText(content))
		return nil
	}

	return writeOK(map[string]string{"date": date, "content": content},
		output.WithSummary(fmt.Sprintf("Journal entry for %s", date)),
	)
}
