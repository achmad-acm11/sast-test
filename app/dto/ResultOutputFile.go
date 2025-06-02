package dto

type SemgrepResult struct {
	Results []Result `json:"results"`
}
type Result struct {
	Path  string `json:"path"`
	Start struct {
		Line int `json:"line"`
	} `json:"start"`
	Extra struct {
		Message  string `json:"message"`
		Metadata struct {
			CWE         []string `json:"cwe"`
			Impact      string   `json:"impact"`
			References  []string `json:"references"`
			Subcategory []string `json:"subcategory"`
		} `json:"metadata"`
	} `json:"extra"`
}
