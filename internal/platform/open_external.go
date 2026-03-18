package platform

import (
	"strings"

	tea "github.com/charmbracelet/bubbletea"
)

func OpenTargetCmd(target string) tea.Cmd {
	normalized := NormalizeTarget(target)
	if normalized == "" {
		return nil
	}

	// Over SSH we should not run open/xdg-open on the server.
	// Print a hybrid link: visible raw URL plus OSC-8 metadata for clients that support click-open.
	return tea.Printf("Open on your device: %s", TerminalLink(normalized))
}

func NormalizeTarget(target string) string {
	t := strings.TrimSpace(target)
	if t == "" {
		return ""
	}
	if strings.HasPrefix(t, "mailto:") {
		return t
	}
	if strings.Contains(t, "@") && !strings.Contains(t, "://") {
		return "mailto:" + t
	}
	if strings.HasPrefix(t, "http://") || strings.HasPrefix(t, "https://") {
		return t
	}
	return "https://" + t
}

func TerminalLink(target string) string {
	return TerminalLinkLabel(target, target)
}

func TerminalLinkLabel(label, target string) string {
	normalized := NormalizeTarget(target)
	if normalized == "" {
		return label
	}
	return "\x1b]8;;" + normalized + "\x1b\\" + label + "\x1b]8;;\x1b\\"
}
