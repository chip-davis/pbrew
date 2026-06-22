package commands

// Search runs `brew search <term>` and pretty-prints the results.
//
// brew's search output groups by formula/cask with ==> headers. Easiest
// path is to parse those sections and render each under our own
// ui.HeaderStyle.
//
// args[0] is the search term. Return an error if missing.
func Search(args []string) error {
	// TODO
	return nil
}
