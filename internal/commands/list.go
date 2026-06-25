package commands

import (
	"encoding/json"
	"slices"
	"strings"

	"github.com/chipdavis/pbrew/internal/brew"
	"github.com/chipdavis/pbrew/internal/ui"
)

type listOpts struct {
	formulae         bool
	casks            bool
	pinned           bool
	pouredFromBottle bool
}

// installedPackage describes a single installed formula or cask, joining
// data from `brew list --json --versions` with `brew list --poured-from-bottle`.
type installedPackage struct {
	Name    string
	Cask    bool
	Version string
	Pinned  bool
	Bottle  bool
}

func List(args []string) error {
	if !shouldPrettify(args) {
		return brew.Passthrough(append([]string{"list"}, args...))
	}

	opts := parseListOpts(args)

	raw, err := brew.Run("list", "--json", "--versions")
	if err != nil {
		return err
	}
	pkgs, err := parseInstalledList(raw)
	if err != nil {
		return err
	}

	bottles, err := fetchBottleSet()
	if err != nil {
		return err
	}
	for i := range pkgs {
		if !pkgs[i].Cask && bottles[pkgs[i].Name] {
			pkgs[i].Bottle = true
		}
	}

	rows := filterPackages(pkgs, opts)
	renderTable(rows, opts)
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
		case "--pinned":
			o.pinned = true
			o.casks = false
		case "--poured-from-bottle":
			o.pouredFromBottle = true
			o.casks = false
		}
	}
	return o
}

func shouldPrettify(args []string) bool {
	for _, a := range args {
		switch a {
		case "--formula", "--formulae", "--cask", "--casks", "--pinned", "--poured-from-bottle":
			continue
		default:
			return false
		}
	}
	return true
}

func filterPackages(pkgs []installedPackage, opts listOpts) []installedPackage {
	out := make([]installedPackage, 0, len(pkgs))
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
		if opts.pouredFromBottle && !p.Bottle {
			continue
		}
		out = append(out, p)
	}
	return out
}

func renderTable(pkgs []installedPackage, opts listOpts) {
	if len(pkgs) == 0 {
		return
	}

	showType := opts.formulae && opts.casks
	showPinned := showColumn(pkgs, func(p installedPackage) bool { return p.Pinned }) && !opts.pinned
	showBottle := showColumn(pkgs, func(p installedPackage) bool { return p.Bottle }) && !opts.pouredFromBottle

	headers := []string{"Name"}
	if showType {
		headers = append(headers, "Type")
	}
	headers = append(headers, "Version")
	if showPinned {
		headers = append(headers, "Pinned")
	}
	if showBottle {
		headers = append(headers, "Bottle")
	}

	t := ui.NewTable(headers...)
	for _, p := range pkgs {
		row := []string{p.Name}
		if showType {
			row = append(row, packageType(p))
		}
		row = append(row, p.Version)
		if showPinned {
			row = append(row, check(p.Pinned))
		}
		if showBottle {
			row = append(row, check(p.Bottle))
		}
		t.Row(row...)
	}
	t.Render()
}

// showColumn returns true if any row has a true value — i.e., there is
// something to display. A column where every value is the same conveys
// no information and is suppressed.
func showColumn(pkgs []installedPackage, get func(installedPackage) bool) bool {
	return slices.ContainsFunc(pkgs, get)
}

func packageType(p installedPackage) string {
	if p.Cask {
		return "cask"
	}
	return "formula"
}

func check(b bool) string {
	if b {
		return "✓"
	}
	return ""
}

// parseInstalledList parses `brew list --json --versions` into a flat
// []installedPackage. Bottle attribution is filled in separately from
// `brew list --poured-from-bottle` (the JSON payload doesn't include it).
func parseInstalledList(raw []byte) ([]installedPackage, error) {
	var payload struct {
		Formulae []struct {
			Name          string   `json:"name"`
			Versions      []string `json:"versions"`
			PinnedVersion *string  `json:"pinned_version"`
		} `json:"formulae"`
		Casks []struct {
			Token         string   `json:"token"`
			Versions      []string `json:"versions"`
			PinnedVersion *string  `json:"pinned_version"`
		} `json:"casks"`
	}
	if err := json.Unmarshal(raw, &payload); err != nil {
		return nil, err
	}

	pkgs := make([]installedPackage, 0, len(payload.Formulae)+len(payload.Casks))
	for _, f := range payload.Formulae {
		pkgs = append(pkgs, installedPackage{
			Name:    f.Name,
			Version: firstVersion(f.Versions),
			Pinned:  f.PinnedVersion != nil,
		})
	}
	for _, c := range payload.Casks {
		pkgs = append(pkgs, installedPackage{
			Name:    c.Token,
			Cask:    true,
			Version: firstVersion(c.Versions),
			Pinned:  c.PinnedVersion != nil,
		})
	}
	return pkgs, nil
}

func firstVersion(vs []string) string {
	if len(vs) == 0 {
		return ""
	}
	return vs[0]
}

func fetchBottleSet() (map[string]bool, error) {
	out, err := brew.Run("list", "--poured-from-bottle")
	if err != nil {
		return nil, err
	}
	set := map[string]bool{}
	for line := range strings.SplitSeq(strings.TrimSpace(string(out)), "\n") {
		if line == "" {
			continue
		}
		set[line] = true
	}
	return set, nil
}
