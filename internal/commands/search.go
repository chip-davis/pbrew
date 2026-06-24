package commands

import (
	"fmt"
	"strings"

	"charm.land/lipgloss/v2"
	"github.com/chipdavis/pbrew/internal/brew"
	"github.com/chipdavis/pbrew/internal/ui"
)

type searchOpts struct {
	formulae bool
	casks    bool
}

func Search(args []string) error {
	if !shouldPrettifySearch(args) {
		return brew.Passthrough(append([]string{"search"}, args...))
	}

	opts, term, err := parseSearchOpts(args)
	if err != nil {
		return err
	}

	s := ui.NewSpinner(fmt.Sprintf("Searching for %s…", term))
	s.Start()

	raw, err := brew.Run("search", term)
	if err != nil {
		s.Stop()
		return err
	}
	s.Stop()

	formulae, casks := parseSearchResults(raw, opts)
	renderSearchResults(formulae, casks, opts)
	return nil
}

func shouldPrettifySearch(args []string) bool {
	for _, a := range args {
		if !strings.HasPrefix(a, "-") {
			continue
		}
		switch a {
		case "--formula", "--formulae", "--cask", "--casks":
			continue
		default:
			return false
		}
	}
	return true
}

func parseSearchOpts(args []string) (searchOpts, string, error) {
	opts := searchOpts{formulae: true, casks: true}
	var term string
	for _, a := range args {
		switch a {
		case "--formula", "--formulae":
			opts.casks = false
		case "--cask", "--casks":
			opts.formulae = false
		default:
			if term != "" {
				return opts, "", fmt.Errorf("search requires a single term")
			}
			term = a
		}
	}
	if term == "" {
		return opts, "", fmt.Errorf("search requires a term")
	}
	return opts, term, nil
}

// parseSearchResults splits brew's output into formulae and cask name slices.
// brew separates the two blocks with a blank line; when a filter is active
// only a single block is emitted, and we attribute it based on opts.
func parseSearchResults(raw []byte, opts searchOpts) (formulae, casks []string) {
	lines := strings.Split(strings.TrimSpace(string(raw)), "\n")

	var blocks [][]string
	current := []string{}
	for _, line := range lines {
		if line == "" {
			if len(current) > 0 {
				blocks = append(blocks, current)
				current = nil
			}
			continue
		}
		current = append(current, line)
	}
	if len(current) > 0 {
		blocks = append(blocks, current)
	}

	switch len(blocks) {
	case 0:
		return nil, nil
	case 1:
		if opts.casks && !opts.formulae {
			return nil, blocks[0]
		}
		return blocks[0], nil
	default:
		return blocks[0], blocks[1]
	}
}

func renderSearchResults(formulae, casks []string, opts searchOpts) {
	var b strings.Builder

	if opts.formulae && len(formulae) > 0 {
		b.WriteString(ui.HeaderStyle.Render("Formulae"))
		b.WriteString("\n")
		for _, name := range formulae {
			b.WriteString("  " + name + "\n")
		}
	}
	if opts.casks && len(casks) > 0 {
		if b.Len() > 0 {
			b.WriteString("\n")
		}
		b.WriteString(ui.HeaderStyle.Render("Casks"))
		b.WriteString("\n")
		for _, name := range casks {
			b.WriteString("  " + name + "\n")
		}
	}
	box := lipgloss.NewStyle().Border(lipgloss.RoundedBorder()).BorderForeground(lipgloss.Color("#FBB040")).Padding(1, 2).Render(b.String())
	fmt.Println(box)
}
