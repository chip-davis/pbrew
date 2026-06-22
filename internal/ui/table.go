package ui

import (
	"charm.land/lipgloss/v2"
	"charm.land/lipgloss/v2/table"
)

// Table is a thin convenience over lipgloss/table for rendering aligned,
// colored tables. Commands use this rather than touching lipgloss directly
// so we can swap the implementation later (e.g. drop to tabwriter for piped
// output).
//
// TODO: implement using github.com/charmbracelet/lipgloss/table.
// Suggested API:
//
//   t := ui.NewTable("Name", "Version", "Status")
//   t.Row("wget", "1.21.4", "✓")
//   fmt.Println(t.Render())

type Table struct {
	headers []string
	rows    [][]string
}

// NewTable creates a table with the given column headers.
func NewTable(headers ...string) *Table {
	return &Table{headers: headers}
}

// Row appends a row. Cell count should match the header count.
func (t *Table) Row(cells ...string) {
	t.rows = append(t.rows, cells)
}

// Render returns the styled table as a string.
func (t *Table) Render() {
	var (
		purple    = lipgloss.Color("99")
		gray      = lipgloss.Color("245")
		lightGray = lipgloss.Color("241")

		headerStyle  = lipgloss.NewStyle().Foreground(purple).Bold(true).Align(lipgloss.Center)
		cellStyle    = lipgloss.NewStyle().Padding(0, 1)
		oddRowStyle  = cellStyle.Foreground(gray)
		evenRowStyle = cellStyle.Foreground(lightGray)
	)

	tbl := table.New().
		Border(lipgloss.NormalBorder()).
		BorderStyle(lipgloss.NewStyle().Foreground(purple)).
		StyleFunc(func(row, col int) lipgloss.Style {
			switch {
			case row == table.HeaderRow:
				return headerStyle
			case row%2 == 0:
				return evenRowStyle
			default:
				return oddRowStyle
			}
		}).
		Headers(t.headers...).
		Rows(t.rows...)

	lipgloss.Println(tbl)
}
