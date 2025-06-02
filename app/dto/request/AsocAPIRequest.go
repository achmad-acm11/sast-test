package request

type AsocProjectCreateRequest struct {
	Key         string `json:"key"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

type AsocSendResultRequest struct {
	ProjectKey  string         `json:"project_key"`
	ScanVersion int            `json:"scan_version"`
	CreatedAt   string         `json:"created_at"`
	Issues      []IssueRequest `json:"issues"`
}

type IssueRequest struct {
	Type        string `json:"type"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Severity    string `json:"severity"`
}
