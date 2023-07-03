package model

// Result contains the results of a scan
type Result struct {
	Repo     string
	Webhook  string
	URL      string
	Files    []string
	Channels map[string]struct{}
}
