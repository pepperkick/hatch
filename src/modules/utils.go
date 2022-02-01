package modules

import "github.com/gofiber/fiber/v2"

func authPassword(password string, next fiber.Handler) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		pass := ctx.Query("password")

		if password != "" && pass != password {
			return ctx.SendStatus(403)
		}

		return next(ctx)
	}
}
