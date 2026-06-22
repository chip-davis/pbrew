package brew

// OutdatedPackage describes a single outdated formula or cask as reported
// by `brew outdated --json=v2`.
type OutdatedPackage struct {
	Name             string
	InstalledVersion string
	LatestVersion    string
	Pinned           bool
	Cask             bool
}

// OutdatedResult splits outdated packages by pinned state so the UI can
// render them in separate sections (per user preference).
type OutdatedResult struct {
	Active []OutdatedPackage
	Pinned []OutdatedPackage
}

// ParseOutdated parses the JSON output of `brew outdated --json=v2`.
//
// TODO: unmarshal the JSON (it has "formulae" and "casks" arrays), flatten
// into []OutdatedPackage, then split into Active vs Pinned.
func ParseOutdated(raw []byte) (OutdatedResult, error) {
	return OutdatedResult{}, nil
}

// PackageInfo is the trimmed-down view of `brew info --json=v2 <pkg>` that
// the info command needs.
type PackageInfo struct {
	Name         string
	Description  string
	Homepage     string
	Version      string
	Installed    bool
	Dependencies []string
	// TODO: add fields as the info command needs them (size, license, etc.)
}

// ParseInfo parses `brew info --json=v2 <pkg>`.
//
// TODO: unmarshal, pull out the first formula or cask, populate PackageInfo.
func ParseInfo(raw []byte) (PackageInfo, error) {
	return PackageInfo{}, nil
}
