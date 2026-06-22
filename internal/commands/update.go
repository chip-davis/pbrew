package commands

// Update wraps `brew update` with a spinner and streamed output.
//
// Suggested approach:
//   s := ui.NewSpinner("Updating Homebrew…")
//   s.Start(); defer s.Stop()
//   return brew.Stream([]string{"update"}, s.Log)
//
// On success, print a SuccessStyle summary line.
func Update(args []string) error {
	// TODO
	return nil
}

// Upgrade wraps `brew upgrade [pkg...]`. Same shape as Update.
//
// Consider: after upgrade completes, re-run Outdated() to show what's
// still pinned. Or print a tally ("upgraded 7 packages").
func Upgrade(args []string) error {
	// TODO
	return nil
}
