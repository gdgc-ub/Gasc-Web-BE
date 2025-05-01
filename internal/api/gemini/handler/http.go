package searchHandler

import (
	searchServices "ProjectGolang/internal/api/gemini/service"
	"ProjectGolang/internal/middleware"
	"github.com/gofiber/fiber/v2"
	"github.com/go-playground/validator/v10"
	"github.com/sirupsen/logrus"
)

// lengkapi sebuah Handler struct dengan dependensi yang dibutuhkan
type SearchHandler struct {
	searchService searchServices.ISearchService
	validator *validator.Validate
	middleware middleware.Middleware
	log *logrus.Logger
}

// lengkapi New constructor untuk handler
func New(as searchServices.ISearchService, validate *validator.Validate, middleware middleware.Middleware, log *logrus.Logger) *SearchHandler {
	return &SearchHandler{ 
		searchService: as, 
		validator: validate, 
		middleware: middleware, 
		log: log, 
	}
}

// lengkapi Start untuk mendaftarkan route
func (h *SearchHandler) Start (srv fiber.Router) {
	search := srv.Group("/services")
	search.Post("/scan", h.AnalyzeMedicineImage)
}
