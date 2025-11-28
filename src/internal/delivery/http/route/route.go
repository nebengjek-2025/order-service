package route

import (
	"order-service/src/internal/delivery/http"
	"order-service/src/internal/delivery/http/middleware"

	"github.com/gofiber/fiber/v2"
)

type RouteConfig struct {
	App              *fiber.App
	UserController   *http.UserController
	DriverController *http.DriverController
	AuthMiddleware   fiber.Handler
}

func (c *RouteConfig) Setup() {
	c.App.Use(middleware.NewLogger())
	c.App.Get("/health", func(ctx *fiber.Ctx) error {
		return ctx.SendString("OK")
	})
	c.SetupAuthRoute()

}

func (c *RouteConfig) SetupAuthRoute() {
	c.App.Use(c.AuthMiddleware)
	// passanger routes
	c.App.Get("/users/v1/profile", c.UserController.GetProfile)
	c.App.Post("/order/v1/location", c.UserController.PostLocation)
	c.App.Post("/order/v1/find-driver", c.UserController.FindDriver)
	c.App.Post("/order/v1/confirm", c.UserController.ConfirmOrder)
	c.App.Post("/order/v1/cancel", c.UserController.CancelOrder)
	c.App.Get("/order/v1/driver-pickup/:orderId", c.UserController.GetDriverPickupRequest)
	c.App.Get("/users/v1/order-status/:orderId", c.UserController.GetOrderStatus)

	// driver routes
	c.App.Post("/drivers/v1/pickup-passanger", c.DriverController.PickupPassanger)
	// c.App.Post("/drivers/v1/complete-trip", c.UserController.CompletedTrip)
	// c.App.Get("/drivers/v1/detail-trip/:orderId", c.UserController.DetailTrip)
}
