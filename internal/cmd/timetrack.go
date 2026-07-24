package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/basecamp/hey-cli/internal/output"
)

type timetrackCommand struct {
	cmd *cobra.Command
}

func newTimetrackCommand() *timetrackCommand {
	timetrackCommand := &timetrackCommand{}
	timetrackCommand.cmd = &cobra.Command{
		Use:   "timetrack",
		Short: "Read time tracking",
		Annotations: map[string]string{
			"agent_notes": "Subcommands: current, list.",
		},
	}

	timetrackCommand.cmd.AddCommand(newTimetrackCurrentCommand().cmd)
	timetrackCommand.cmd.AddCommand(newTimetrackListCommand().cmd)

	return timetrackCommand
}

// current

type timetrackCurrentCommand struct {
	cmd *cobra.Command
}

func newTimetrackCurrentCommand() *timetrackCurrentCommand {
	timetrackCurrentCommand := &timetrackCurrentCommand{}
	timetrackCurrentCommand.cmd = &cobra.Command{
		Use:   "current",
		Short: "Show current time tracking status",
		Example: `  hey timetrack current
  hey timetrack current --json`,
		RunE: timetrackCurrentCommand.run,
	}

	return timetrackCurrentCommand
}

func (c *timetrackCurrentCommand) run(cmd *cobra.Command, args []string) error {
	if err := requireAuth(); err != nil {
		return err
	}

	ctx := cmd.Context()
	track, err := sdk.TimeTracks().GetOngoing(ctx)
	if err != nil {
		return convertSDKError(err)
	}

	if writer.IsStyled() {
		w := cmd.OutOrStdout()
		if track == nil {
			fmt.Fprintln(w, "No active time track.")
			return nil
		}

		fmt.Fprintf(w, "Active time track #%d\n", track.Id)
		fmt.Fprintf(w, "Started: %s\n", formatTimestamp(track.StartsAt))
		if track.Title != "" {
			fmt.Fprintf(w, "Title:   %s\n", track.Title)
		}
		return nil
	}

	if track == nil {
		return writeOK(nil, output.WithSummary("No active time track"))
	}

	return writeOK(track,
		output.WithSummary(fmt.Sprintf("Active time track #%d", track.Id)),
	)
}

// list

type timetrackListCommand struct {
	cmd   *cobra.Command
	limit int
	all   bool
}

func newTimetrackListCommand() *timetrackListCommand {
	timetrackListCommand := &timetrackListCommand{}
	timetrackListCommand.cmd = &cobra.Command{
		Use:   "list",
		Short: "List time tracks",
		Example: `  hey timetrack list
  hey timetrack list --limit 10
  hey timetrack list --json`,
		RunE: timetrackListCommand.run,
	}

	timetrackListCommand.cmd.Flags().IntVar(&timetrackListCommand.limit, "limit", 0, "Maximum number of time tracks to show")
	timetrackListCommand.cmd.Flags().BoolVar(&timetrackListCommand.all, "all", false, "Fetch all results (override --limit)")

	return timetrackListCommand
}

func (c *timetrackListCommand) run(cmd *cobra.Command, args []string) error {
	if err := requireAuth(); err != nil {
		return err
	}

	ctx := cmd.Context()
	resp, err := listPersonalRecordings(ctx)
	if err != nil {
		return err
	}

	tracks := filterRecordingsByType(resp, "Calendar::TimeTrack")

	total := len(tracks)
	if c.limit > 0 && !c.all && len(tracks) > c.limit {
		tracks = tracks[:c.limit]
	}
	notice := output.TruncationNotice(len(tracks), total)

	if writer.IsStyled() {
		if len(tracks) == 0 {
			fmt.Fprintln(cmd.OutOrStdout(), "No time tracks.")
			return nil
		}

		table := newTable(cmd.OutOrStdout())
		table.addRow([]string{"ID", "Title", "Start", "End"})
		for _, t := range tracks {
			table.addRow([]string{fmt.Sprintf("%d", t.Id), t.Title, formatTimestamp(t.StartsAt), formatTimestamp(t.EndsAt)})
		}
		table.print()
		if notice != "" {
			fmt.Fprintln(cmd.OutOrStdout(), notice)
		}
		return nil
	}

	return writeOK(tracks,
		output.WithSummary(fmt.Sprintf("%d time tracks", len(tracks))),
		output.WithNotice(notice),
	)
}
