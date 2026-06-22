package ui

// Spinner wraps a Bubble Tea spinner program for long-running brew commands
// (update/upgrade). It runs in its own goroutine and prints streamed output
// lines below the spinner.
//
// TODO: implement with github.com/charmbracelet/bubbletea + bubbles/spinner.
// Suggested API:
//
//   s := ui.NewSpinner("Updating Homebrew…")
//   s.Start()
//   defer s.Stop()
//   s.Log("==> Fetching formula/x") // appears below the spinner
//
// Alternative if Bubble Tea feels heavy for this: just print a static
// message and stream brew's output verbatim. Worth trying both.

type Spinner struct {
	// TODO
}

func NewSpinner(label string) *Spinner {
	return &Spinner{}
}

func (s *Spinner) Start()           {}
func (s *Spinner) Stop()            {}
func (s *Spinner) Log(line string)  {}
func (s *Spinner) SetLabel(string)  {}
