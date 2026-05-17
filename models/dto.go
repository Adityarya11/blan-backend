package models

type CompileRequest struct {
	SourceCode string `json:"source_code" binding:"required"`
}

type CompileResponse struct {
	Output string `json:"output"`
	Error  string `json:"error,omitempty"`
}
