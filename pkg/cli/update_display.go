package cli

import (
	"fmt"
	"io"
	"os"

	"github.com/github/gh-aw/pkg/console"
)

// showUpdateSummary displays a summary of workflow updates using console helpers
func showUpdateSummary(successfulUpdates []string, failedUpdates []updateFailure) {
	showUpdateSummaryWithWriter(os.Stderr, successfulUpdates, failedUpdates)
}

// showUpdateSummaryWithWriter displays a summary of workflow updates using console helpers,
// writing output to the provided writer
func showUpdateSummaryWithWriter(w io.Writer, successfulUpdates []string, failedUpdates []updateFailure) {
	fmt.Fprintln(w, "")

	// Show successful updates
	if len(successfulUpdates) > 0 {
		fmt.Fprintln(w, console.FormatSuccessMessage(fmt.Sprintf("Successfully updated and compiled %d workflow(s):", len(successfulUpdates))))
		for _, name := range successfulUpdates {
			fmt.Fprintln(w, console.FormatListItem(name))
		}
		fmt.Fprintln(w, "")
	}

	// Show failed updates
	if len(failedUpdates) > 0 {
		fmt.Fprintln(w, console.FormatErrorMessage(fmt.Sprintf("Failed to update %d workflow(s):", len(failedUpdates))))
		for _, failure := range failedUpdates {
			fmt.Fprintf(w, "  %s: %s\n", failure.Name, failure.Error)
		}
		fmt.Fprintln(w, "")
	}
}
