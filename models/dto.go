package models

type CompileRequest struct {
	SourceCode string `json:"source_code" binding:"required"`
}

type CompileResponse struct {
	Output string `json:"output"`
	Error  string `json:"error,omitempty"`
	Cached bool   `json:"cached"`
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

type SignUpRequest struct {
	Username string `json:"username" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type AuthResponse struct {
	Token string `json:"token"`
}
