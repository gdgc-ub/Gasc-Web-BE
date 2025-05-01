package gemini

import (
	"ProjectGolang/pkg/response"
)

var (
	ErrorResponse             = response.New(400, "error from user")
	ErrGeminiProcessingFailed = response.New(500, "gemini processing failed")
)
