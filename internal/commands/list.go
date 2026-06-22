package commands

import (
	"fmt"
	"strings"

	"github.com/chipdavis/pbrew/internal/brew"
	"github.com/chipdavis/pbrew/internal/ui"
)

type listOpts struct {
	formulae bool
	casks    bool
	versions bool
	pinned   bool
}

func List(args []string) error {
	if !shouldPrettify(args) {
		return brew.Passthrough(append([]string{"list"}, args...))
	}

	opts := parseListOpts(args)
	if opts.formulae {
		if err := renderSection("Formulae", "--formula", opts); err != nil {
			return err
		}
	}
	if opts.casks {
		if err := renderSection("Casks", "--cask", opts); err != nil {
			return err
		}
	}
	return nil
}

func parseListOpts(args []string) listOpts {
	o := listOpts{formulae: true, casks: true}
	for _, a := range args {
		switch a {
		case "--formula", "--formulae":
			o.casks = false
		case "--cask", "--casks":
			o.formulae = false
		case "--versions":
			o.versions = true
		case "--pinned":
			o.pinned = true
			// brew has no concept of pinned casks
			o.casks = false
		}
	}
	return o
}

func shouldPrettify(args []string) bool {
	for _, a := range args {
		switch a {
		case "--formula", "--formulae", "--cask", "--casks", "--versions", "--pinned":
			continue
		default:
			return false
		}
	}
	return true
}

func renderSection(heading, typeFlag string, opts listOpts) error {
	brewArgs := []string{"list", typeFlag}
	if opts.versions {
		brewArgs = append(brewArgs, "--versions")
	}
	if opts.pinned {
		brewArgs = append(brewArgs, "--pinned")
	}

	output, err := brew.Run(brewArgs...)
	if err != nil {
		return err
	}

	lines := splitLines(output)
	if len(lines) == 0 {
		return nil
	}

	fmt.Println(ui.HeaderStyle.Render(heading))

	headers := []string{"Package"}
	if opts.versions {
		headers = append(headers, "Version")
	}
	t := ui.NewTable(headers...)

	for _, line := range lines {
		parts := strings.Fields(line)
		if len(parts) == 0 {
			continue
		}
		if opts.versions {
			version := ""
			if len(parts) > 1 {
				version = strings.Join(parts[1:], " ")
			}
			t.Row(parts[0], version)
		} else {
			t.Row(parts[0])
		}
	}
	t.Render()
	return nil
}

func splitLines(b []byte) []string {
	trimmed := strings.TrimSpace(string(b))
	if trimmed == "" {
		return nil
	}
	return strings.Split(trimmed, "\n")
}
