package commands

import (
	"encoding/json"
	"fmt"

	"github.com/chipdavis/pbrew/internal/brew"
	"github.com/chipdavis/pbrew/internal/ui"
)

type outdatedOpts struct {
	formulae bool
	casks    bool
	pinned   bool
}

type outdatedPackage struct {
	Name             string
	Cask             bool
	InstalledVersion string
	CurrentVersion   string
	Pinned           bool
	PinnedVersion    string
}

func Outdated(args []string) error {
	if !shouldPrettifyOutdated(args) {
		return brew.Passthrough(append([]string{"outdated"}, args...))
	}

	opts := parseOutdatedOpts(args)

	raw, err := brew.Run("outdated", "--json=v2")
	if err != nil {
		return err
	}

	pkgs, err := parseOutdated(raw)
	if err != nil {
		return err
	}

	pkgs = filterOutdated(pkgs, opts)
	if len(pkgs) == 0 {
		fmt.Println(ui.SuccessStyle.Render("✓ Everything up to date"))
		return nil
	}

	renderOutdatedTable(pkgs, opts)
	return nil
}

func shouldPrettifyOutdated(args []string) bool {
	for _, a := range args {
		switch a {
		case "--formula", "--formulae", "--cask", "--casks", "--pinned":
			continue
		default:
			return false
		}
	}
	return true
}

func parseOutdatedOpts(args []string) outdatedOpts {
	o := outdatedOpts{formulae: true, casks: true}
	for _, a := range args {
		switch a {
		case "--formula", "--formulae":
			o.casks = false
		case "--cask", "--casks":
			o.formulae = false
		case "--pinned":
			o.pinned = true
			o.casks = false
		}
	}
	return o
}

func parseOutdated(raw []byte) ([]outdatedPackage, error) {
	var payload struct {
		Formulae []struct {
			Name              string   `json:"name"`
			InstalledVersions []string `json:"installed_versions"`
			CurrentVersion    string   `json:"current_version"`
			Pinned            bool     `json:"pinned"`
			PinnedVersion     string   `json:"pinned_version"`
		} `json:"formulae"`
		Casks []struct {
			Name              string   `json:"name"`
			InstalledVersions []string `json:"installed_versions"`
			CurrentVersion    string   `json:"current_version"`
			Pinned            bool     `json:"pinned"`
			PinnedVersion     string   `json:"pinned_version"`
		} `json:"casks"`
	}
	if err := json.Unmarshal(raw, &payload); err != nil {
		return nil, err
	}

	pkgs := make([]outdatedPackage, 0, len(payload.Formulae)+len(payload.Casks))
	for _, f := range payload.Formulae {
		pkgs = append(pkgs, outdatedPackage{
			Name:             f.Name,
			InstalledVersion: firstVersion(f.InstalledVersions),
			CurrentVersion:   f.CurrentVersion,
			Pinned:           f.Pinned,
			PinnedVersion:    f.PinnedVersion,
		})
	}
	for _, c := range payload.Casks {
		pkgs = append(pkgs, outdatedPackage{
			Name:             c.Name,
			Cask:             true,
			InstalledVersion: firstVersion(c.InstalledVersions),
			CurrentVersion:   c.CurrentVersion,
			Pinned:           c.Pinned,
			PinnedVersion:    c.PinnedVersion,
		})
	}
	return pkgs, nil
}

func filterOutdated(pkgs []outdatedPackage, opts outdatedOpts) []outdatedPackage {
	out := make([]outdatedPackage, 0, len(pkgs))
	for _, p := range pkgs {
		if p.Cask && !opts.casks {
			continue
		}
		if !p.Cask && !opts.formulae {
			continue
		}
		if opts.pinned && !p.Pinned {
			continue
		}
		out = append(out, p)
	}
	return out
}

func renderOutdatedTable(pkgs []outdatedPackage, opts outdatedOpts) {
	showType := opts.formulae && opts.casks
	showPinned := anyPinned(pkgs) && !opts.pinned

	headers := []string{"Name"}
	if showType {
		headers = append(headers, "Type")
	}
	headers = append(headers, "Installed", "Latest")
	if showPinned {
		headers = append(headers, "Pinned")
	}

	t := ui.NewTable(headers...)
	for _, p := range pkgs {
		row := []string{p.Name}
		if showType {
			row = append(row, packageType(installedPackage{Cask: p.Cask}))
		}
		row = append(row, p.InstalledVersion, p.CurrentVersion)
		if showPinned {
			row = append(row, check(p.Pinned))
		}
		t.Row(row...)
	}
	t.Render()
}

func anyPinned(pkgs []outdatedPackage) bool {
	for _, p := range pkgs {
		if p.Pinned {
			return true
		}
	}
	return false
}
