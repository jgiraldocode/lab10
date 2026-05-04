package domain

import "time"

type Source string

const (
	SourceWorkspace     Source = "workspace"
	SourceAppleScript   Source = "applescript"
	SourceAccessibility Source = "accessibility"
)

type Confidence string

const (
	ConfidenceHigh   Confidence = "high"
	ConfidenceMedium Confidence = "medium"
	ConfidenceLow    Confidence = "low"
)

type BrowserFamily string

const (
	BrowserNone     BrowserFamily = ""
	BrowserSafari   BrowserFamily = "safari"
	BrowserChromium BrowserFamily = "chromium"
	BrowserFirefox  BrowserFamily = "firefox"
	BrowserUnknown  BrowserFamily = "unknown"
)

type SampleEvent struct {
	ID            int64         `json:"id"`
	Timestamp     time.Time     `json:"timestamp"`
	BundleID      string        `json:"bundle_id"`
	AppName       string        `json:"app_name"`
	WindowTitle   string        `json:"window_title"`
	BrowserFamily BrowserFamily `json:"browser_family,omitempty"`
	TabTitle      string        `json:"tab_title,omitempty"`
	TabURLHost    string        `json:"tab_url_host,omitempty"`
	Source        Source        `json:"source"`
	Confidence    Confidence    `json:"confidence"`
}
