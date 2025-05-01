package config

import (
	searchHandler "ProjectGolang/internal/api/gemini/handler"
	searchServices "ProjectGolang/internal/api/gemini/service"
	"ProjectGolang/internal/middleware"
	"fmt"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/sirupsen/logrus"
	"os"
)

type Server struct {
	engine     *fiber.App
	log        *logrus.Logger
	middleware middleware.Middleware
	validator  *validator.Validate
	handlers   []handler
}

type handler interface {
	Start(srv fiber.Router)
}

func NewServer(fiberApp *fiber.App, log *logrus.Logger, validator *validator.Validate) (*Server, error) {
	bootstrap := &Server{
		engine:     fiberApp,
		log:        log,
		validator:  validator,
		middleware: middleware.New(log),
	}

	return bootstrap, nil
}

func (s *Server) RegisterHandler() {
	//Another Domain
	searchService := searchServices.New(s.log)
	searchHandlers := searchHandler.New(searchService, s.validator, s.middleware, s.log)

	s.checkHealth()
	s.handlers = append(s.handlers, searchHandlers)
}

func (s *Server) Run() error {
	s.engine.Use(cors.New(cors.Config{
		AllowOrigins:     "http://localhost:5173, http://localhost:3000",
		AllowHeaders:     "Origin, Content-Type, Accept",
		AllowMethods:     "GET, POST, PUT, DELETE, OPTIONS",
		AllowCredentials: true,
	}))
	s.engine.Use(s.middleware.NewLoggingMiddleware)
	router := s.engine.Group("/api/v1")

	for _, h := range s.handlers {
		h.Start(router)
	}

	port := os.Getenv("APP_PORT")
	if port == "" {
		port = "3000"
	}

	s.log.Infof("Starting server on port %s", port)

	if err := s.engine.Listen(fmt.Sprintf(":%s", port)); err != nil {
		return err
	}
	return nil
}

func (s *Server) checkHealth() {
	s.engine.Get("/", func(ctx *fiber.Ctx) error {
		s.log.Info("Health check endpoint called")
		return ctx.JSON(fiber.Map{
			"message": "Server is Healthy!",
		})
	})
}
