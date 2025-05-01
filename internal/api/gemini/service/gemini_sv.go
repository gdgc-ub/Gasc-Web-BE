package searchServices

import (
	"ProjectGolang/internal/api/gemini"
	"context"
	"mime/multipart"
)


func (s *searchService) AnalyzeMedicineImage(ctx context.Context, imageFile *multipart.FileHeader) (gemini.GeminiResponse, error) {
	return gemini.GeminiResponse{}, nil 
}
