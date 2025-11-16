package middleware

import (
	"time"

	"order-service/src/pkg/log"
	"order-service/src/pkg/utils"

	"github.com/gofiber/fiber/v2"
)

func NewLogger() fiber.Handler {
	logger := log.GetLogger()

	return func(c *fiber.Ctx) error {
		ctx := c.Context()
		start := time.Now()
		err := c.Next()
		latency := time.Since(start)

		status := c.Response().StatusCode()
		method := c.Method()
		path := c.Path()

		meta := "latency=" + latency.String() + " status=" + utils.ConvertString(status)

		if err != nil || status >= 500 {
			logger.Error(ctx.String(), "HTTP Request-"+method+" "+path, "middleware", meta)
		} else if status >= 400 {
			logger.Info(ctx.String(), "HTTP Request Warning"+method+" "+path, "middleware", meta)
		} else {
			logger.Info(ctx.String(), "HTTP Request "+method+" "+path, "middleware", meta)
		}

		return err
	}
}
