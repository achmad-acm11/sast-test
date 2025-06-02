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
	Title        string `json:"title"`
	Rule         string `json:"rule"`
	Path         string `json:"path"`
	Line         string `json:"line"`
	Type         string `json:"type"`
	Description  string `json:"description"`
	Severity     string `json:"severity"`
	References   string `json:"references"`
	LastFoundAt  string `json:"last_found_at"`
	StatusResult int    `json:"status_result"`
}
