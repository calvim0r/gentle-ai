package screens

import (
	"fmt"
	"strings"

	"github.com/gentleman-programming/gentle-ai/internal/tui/styles"
)

const maxErrorLines = 15

type FailedStep struct {
	ID    string
	Error string
}

type CompletePayload struct {
	ConfiguredAgents    int
	InstalledComponents int
	FailedSteps         []FailedStep
	RollbackPerformed   bool
}

func RenderComplete(data CompletePayload) string {
	if len(data.FailedSteps) > 0 {
		return renderCompleteFailed(data)
	}
	return renderCompleteSuccess(data)
}

func renderCompleteSuccess(data CompletePayload) string {
	var b strings.Builder

	b.WriteString(styles.SuccessStyle.Render("Done! Your AI agents are ready."))
	b.WriteString("\n\n")

	b.WriteString("  " + styles.HeadingStyle.Render("Configured agents") + "  " + styles.SuccessStyle.Render(fmt.Sprintf("%d", data.ConfiguredAgents)) + "\n")
	b.WriteString("  " + styles.HeadingStyle.Render("Installed components") + "  " + styles.SuccessStyle.Render(fmt.Sprintf("%d", data.InstalledComponents)) + "\n")
	b.WriteString("\n")

	b.WriteString(styles.HeadingStyle.Render("Next steps"))
	b.WriteString("\n")
	b.WriteString(styles.UnselectedStyle.Render("  1. Set your API keys"))
	b.WriteString("\n")
	b.WriteString(styles.UnselectedStyle.Render("  2. Run your selected agent"))
	b.WriteString("\n")
	b.WriteString(styles.UnselectedStyle.Render("  3. Try /sdd-new my-feature"))
	b.WriteString("\n\n")

	b.WriteString(styles.HelpStyle.Render("Press Enter to exit."))

	return b.String()
}

func renderCompleteFailed(data CompletePayload) string {
	var b strings.Builder

	b.WriteString(styles.ErrorStyle.Render("Installation completed with errors."))
	b.WriteString("\n\n")

	b.WriteString(styles.HeadingStyle.Render("Failed steps"))
	b.WriteString("\n")
	for _, step := range data.FailedSteps {
		b.WriteString("  " + styles.ErrorStyle.Render("✗ "+step.ID))
		b.WriteString("\n")
		lines := strings.Split(step.Error, "\n")
		if len(lines) > maxErrorLines {
			lines = lines[:maxErrorLines]
			lines = append(lines, "... (truncated)")
		}
		for _, line := range lines {
			b.WriteString("    " + styles.SubtextStyle.Render(line))
			b.WriteString("\n")
		}
	}
	b.WriteString("\n")

	if data.RollbackPerformed {
		b.WriteString(styles.WarningStyle.Render("Rollback was performed — previous configuration restored."))
		b.WriteString("\n\n")
	}

	b.WriteString(styles.HeadingStyle.Render("What to do"))
	b.WriteString("\n")
	b.WriteString(styles.UnselectedStyle.Render("  1. Check the error messages above"))
	b.WriteString("\n")
	b.WriteString(styles.UnselectedStyle.Render("  2. Fix the underlying issue (missing deps, permissions, etc.)"))
	b.WriteString("\n")
	b.WriteString(styles.UnselectedStyle.Render("  3. Run gentle-ai again to retry"))
	b.WriteString("\n\n")

	b.WriteString(styles.HelpStyle.Render("Press Enter to exit."))

	return b.String()
}
