package modules

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/qixalite/hatch/modules/services"
)

func SetupCommonAPI(app *fiber.App, password string) {
	app.Get("/common/kickall", authPassword(password, func(ctx *fiber.Ctx) error {
		_, err := services.ExecRconCommand(serverStatus.GetRconServer(), "kickall")
		if err != nil {
			return ctx.SendStatus(400)
		}
		return ctx.SendStatus(200)
	}))
	fmt.Println("[HATCH MODULE] Started CommonAPI Module")
}
