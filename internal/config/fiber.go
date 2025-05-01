package config

import (
	"ProjectGolang/internal/middleware"
	"ProjectGolang/pkg/log"
	"ProjectGolang/pkg/response"
	"errors"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/utils"
	jsoniter "github.com/json-iterator/go"
	"github.com/sirupsen/logrus"
)

func NewFiber(logger *logrus.Logger) *fiber.App {
	app := fiber.New(
		fiber.Config{
			AppName:           "Radiance Backend",
			BodyLimit:         50 * 1024 * 1024,
			DisableKeepalive:  true,
			StrictRouting:     true,
			CaseSensitive:     true,
			EnablePrintRoutes: true,
			ErrorHandler:      newErrorHandler(logger),
			JSONEncoder:       jsoniter.Marshal,
			JSONDecoder:       jsoniter.Unmarshal,
		})

	// Use our custom logger middleware instead of the default one
	app.Use(middleware.LoggerConfig())

	return app
}

func newErrorHandler(logger *logrus.Logger) fiber.ErrorHandler {
	return func(ctx *fiber.Ctx, err error) error {
		// Log the error with context
		logFields := log.Fields{
			"method": ctx.Method(),
			"path":   ctx.Path(),
			"ip":     ctx.IP(),
		}

		var apiErr *response.Error
		if errors.As(err, &apiErr) {
			logFields["error_code"] = apiErr.Code
			logFields["error_message"] = apiErr.Error()

			log.Warn(logFields, "API error occurred")

			return ctx.Status(apiErr.Code).JSON(fiber.Map{
				"errors": fiber.Map{"message": apiErr.Error()},
			})
		}

		var validationErr validator.ValidationErrors
		if errors.As(err, &validationErr) {
			fieldErr := fiber.Map{}
			for _, e := range validationErr {
				fieldErr[e.Field()] = e.Error()
			}

			logFields["validation_errors"] = fieldErr
			log.Warn(logFields, "Validation error occurred")

			fieldErr["message"] = utils.StatusMessage(fiber.StatusUnprocessableEntity)
			return ctx.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{
				"errors": fieldErr,
			})
		}

		var fiberErr *fiber.Error
		if errors.As(err, &fiberErr) {
			logFields["error_code"] = fiberErr.Code
			logFields["error_message"] = fiberErr.Error()

			if fiberErr.Code >= 500 {
				log.Error(logFields, "Server error occurred")
			} else {
				log.Warn(logFields, "Client error occurred")
			}

			return ctx.Status(fiberErr.Code).JSON(fiber.Map{
				"errors": fiber.Map{"message": utils.StatusMessage(fiberErr.Code), "err": err.Error()},
			})
		}

		// For unexpected errors, generate a trace ID for easier debugging
		traceID := log.ErrorWithTraceID(logFields, "Unexpected server error")

		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"errors": fiber.Map{
				"message":  utils.StatusMessage(fiber.StatusInternalServerError),
				"err":      err.Error(),
				"trace_id": traceID,
			},
		})
	}
}
