package commands

// Outdated shows packages with available updates, with pinned packages in a
// separate section (per project preference).
//
// Suggested approach:
//   1. raw, err := brew.Run("outdated", "--json=v2")
//   2. result, err := brew.ParseOutdated(raw)
//   3. Render result.Active under HeaderStyle "Outdated"
//      Render result.Pinned under HeaderStyle "Pinned (held back)" with MutedStyle rows
//   4. If both are empty, print a SuccessStyle "Everything up to date" line.
//
// Columns to consider: Name | Installed | Latest | Type (formula/cask)
func Outdated(args []string) error {
	// TODO
	return nil
}
