package lsp

type PublishDiagnosticsNotification struct {
	Notification
	Params PublishDiagnosticsParams `json:"params"`
}

type PublishDiagnosticsParams struct {
	URI         string       `json:"uri"`
	Diagnostics []Diagnostic `json:"diagnostics"`
}

type Diagnostic struct {
	Range Range `json:"range"`
	// error = 1, warning = 2, info = 3, hint = 4
	Severity int    `json:"severity"`
	Source   string `json:"source"`
	Message  string `json:"message"`
}
