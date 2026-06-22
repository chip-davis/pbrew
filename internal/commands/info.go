package commands

// Info shows details for a single package in a card layout.
//
// Suggested approach:
//   1. raw, err := brew.Run("info", "--json=v2", args[0])
//   2. info, err := brew.ParseInfo(raw)
//   3. Render as labeled rows inside a lipgloss.NewStyle().Border(...) box:
//        Name        : wget
//        Version     : 1.21.4
//        Homepage    : https://...
//        Installed   : yes / no
//        Dependencies: openssl@3, libidn2, ...
//
// Return an error if args is empty.
func Info(args []string) error {
	// TODO
	return nil
}
