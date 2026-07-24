package cmd

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/mattn/go-runewidth"
	"golang.org/x/term"
)

var colorDisabled bool

func init() {
	_, noColor := os.LookupEnv("NO_COLOR")
	colorDisabled = noColor || !term.IsTerminal(int(os.Stdout.Fd())) //nolint:gosec // G115: fd fits in int on all supported platforms
}

type table struct {
	w            io.Writer
	columnWidths map[int]int
	rows         [][]string
}

func newTable(w io.Writer) *table {
	return &table{
		w:            w,
		columnWidths: map[int]int{},
		rows:         [][]string{},
	}
}

func (t *table) addRow(row []string) {
	t.updateColumnWidths(row)
	t.rows = append(t.rows, row)
}

func (t *table) print() {
	for rownum, row := range t.rows {
		for i, cell := range row {
			cellStyle := plain
			if rownum == 0 {
				cellStyle = italic
			}
			if rownum > 0 && i == 0 {
				cellStyle = bold
			}

			pad := max(t.columnWidths[i]-runewidth.StringWidth(cell), 0)
			fmt.Fprintf(t.w, "%s%s  ", cellStyle.format(cell), strings.Repeat(" ", pad))
		}
		fmt.Fprintln(t.w)
	}
}

func (t *table) updateColumnWidths(row []string) {
	for i, cell := range row {
		w := runewidth.StringWidth(cell)
		if w > t.columnWidths[i] {
			t.columnWidths[i] = w
		}
	}
}

type style string

const (
	plain  style = ""
	bold   style = "1;34"
	italic style = "3;94"
)

func (s style) format(value string) string {
	if s == plain || colorDisabled {
		return value
	}
	return "\033[" + string(s) + "m" + value + "\033[0m"
}

func truncate(s string, maxWidth int) string {
	if runewidth.StringWidth(s) <= maxWidth {
		return s
	}
	return runewidth.Truncate(s, maxWidth, "...")
}

func stdinIsTerminal() bool {
	return term.IsTerminal(int(os.Stdin.Fd())) //nolint:gosec // G115: fd fits in int
}

func stdoutIsTerminal() bool {
	return term.IsTerminal(int(os.Stdout.Fd())) //nolint:gosec // G115: fd fits in int
}
