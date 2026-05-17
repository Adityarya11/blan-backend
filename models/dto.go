package models

type CompileRequest struct {
	SourceCode string `json:"source_code" binding:"required"`
}

type CompileResponse struct {
	Output string `json:"output"`
	Error  string `json:"error,omitempty"`
}

type CompileAcceptedResponse struct {
	ID     string `json:"id"`
	Status string `json:"status"`
}

type JobStatusResponse struct {
	ID     string `json:"id"`
	Status string `json:"status"`
	Output string `json:"output,omitempty"`
	Error  string `json:"error,omitempty"`
}
