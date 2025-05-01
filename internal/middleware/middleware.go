package middleware

import (
	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
)

type middleware struct {
	token             *tokenMiddleware
	rateLimitter      *rateLimiter
	loggingMiddleware *loggingMiddleware
	log               *logrus.Logger
}

type Middleware interface {
	NewRateLimiter(ctx *fiber.Ctx) error
	NewTokenMiddleware(ctx *fiber.Ctx) error
	NewLoggingMiddleware(ctx *fiber.Ctx) error
}

func New(logger *logrus.Logger) Middleware {
	rateLimit := newRateLimiter(50, 100)
	token := newTokenMiddleware()
	logging := newLoggingMiddleware(logger)

	return &middleware{
		token:             token,
		rateLimitter:      rateLimit,
		loggingMiddleware: logging,
		log:               logger,
	}
}
