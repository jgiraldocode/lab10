package tracker

import (
	"context"
	"fmt"
	"net/url"
	"strings"

	"github.com/jgiraldo/focustrack/internal/domain"
)

type BrowserInfo struct {
	Family    domain.BrowserFamily
	ScriptApp string
	Type      string // "safari" or "chromium"
}

type BrowserDetector struct {
	browsers map[string]BrowserInfo
}

func NewBrowserDetector() *BrowserDetector {
	return &BrowserDetector{
		browsers: map[string]BrowserInfo{
			"com.apple.Safari":              {Family: domain.BrowserSafari, ScriptApp: "Safari", Type: "safari"},
			"com.google.Chrome":             {Family: domain.BrowserChromium, ScriptApp: "Google Chrome", Type: "chromium"},
			"com.brave.Browser":             {Family: domain.BrowserChromium, ScriptApp: "Brave Browser", Type: "chromium"},
			"com.microsoft.edgemac":         {Family: domain.BrowserChromium, ScriptApp: "Microsoft Edge", Type: "chromium"},
			"com.microsoft.edgemac.Dev":     {Family: domain.BrowserChromium, ScriptApp: "Microsoft Edge Dev", Type: "chromium"},
			"company.thebrowser.Browser":    {Family: domain.BrowserChromium, ScriptApp: "Arc", Type: "chromium"},
		},
	}
}

func (d *BrowserDetector) Lookup(bundleID string) (BrowserInfo, bool) {
	info, ok := d.browsers[bundleID]
	return info, ok
}

func (d *BrowserDetector) GetActiveTab(ctx context.Context, info BrowserInfo) (title string, host string, err error) {
	var script string
	switch info.Type {
	case "safari":
		script = fmt.Sprintf(`tell application "%s"
	set tabURL to URL of current tab of front window
	set tabName to name of current tab of front window
	return tabURL & "|||" & tabName
end tell`, info.ScriptApp)
	case "chromium":
		script = fmt.Sprintf(`tell application "%s"
	set tabURL to URL of active tab of front window
	set tabName to title of active tab of front window
	return tabURL & "|||" & tabName
end tell`, info.ScriptApp)
	default:
		return "", "", fmt.Errorf("unsupported browser type: %s", info.Type)
	}

	out, err := runAppleScript(ctx, script)
	if err != nil {
		return "", "", fmt.Errorf("browser tab script: %w", err)
	}

	parts := strings.SplitN(out, "|||", 2)
	if len(parts) < 2 {
		return "", "", fmt.Errorf("unexpected output: %s", out)
	}

	tabURL := strings.TrimSpace(parts[0])
	title = strings.TrimSpace(parts[1])

	if parsed, err := url.Parse(tabURL); err == nil {
		host = parsed.Hostname()
	}

	return title, host, nil
}
