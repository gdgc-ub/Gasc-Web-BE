package searchServices

import (
	"ProjectGolang/internal/api/gemini"
	"context"
	"mime/multipart"

	"github.com/sirupsen/logrus"
)

type searchService struct {
	log *logrus.Logger
}

type ISearchService interface {
	AnalyzeMedicineImage(ctx context.Context, imageFile *multipart.FileHeader) (gemini.GeminiResponse, error)
}

func New(log *logrus.Logger) ISearchService {
	return &searchService{
		log: log,
	}
}
