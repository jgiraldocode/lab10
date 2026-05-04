package tracker

import (
	"context"
	"fmt"
	"os/exec"
	"strings"
	"time"

	"github.com/jgiraldo/focustrack/internal/domain"
)

const frontmostScript = `
tell application "System Events"
	set frontApp to first application process whose frontmost is true
	set appName to name of frontApp
	set bundleID to bundle identifier of frontApp
	set winTitle to ""
	try
		set winTitle to name of front window of frontApp
	end try
	return appName & "|||" & bundleID & "|||" & winTitle
end tell
`

type MacTracker struct {
	browsers *BrowserDetector
}

func NewMacTracker() *MacTracker {
	return &MacTracker{
		browsers: NewBrowserDetector(),
	}
}

func (t *MacTracker) Poll(ctx context.Context) (domain.SampleEvent, error) {
	event := domain.SampleEvent{
		Timestamp:  time.Now().UTC(),
		Source:     domain.SourceAppleScript,
		Confidence: domain.ConfidenceHigh,
	}

	out, err := runAppleScript(ctx, frontmostScript)
	if err != nil {
		return event, fmt.Errorf("frontmost app: %w", err)
	}

	parts := strings.SplitN(out, "|||", 3)
	if len(parts) >= 2 {
		event.AppName = strings.TrimSpace(parts[0])
		event.BundleID = strings.TrimSpace(parts[1])
	}
	if len(parts) >= 3 {
		event.WindowTitle = strings.TrimSpace(parts[2])
	}

	if info, ok := t.browsers.Lookup(event.BundleID); ok {
		event.BrowserFamily = info.Family
		tabTitle, tabHost, err := t.browsers.GetActiveTab(ctx, info)
		if err == nil {
			event.TabTitle = tabTitle
			event.TabURLHost = tabHost
		} else {
			event.Confidence = domain.ConfidenceMedium
		}
	}

	return event, nil
}

func runAppleScript(ctx context.Context, script string) (string, error) {
	cmd := exec.CommandContext(ctx, "osascript", "-e", script)
	out, err := cmd.Output()
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(out)), nil
}
