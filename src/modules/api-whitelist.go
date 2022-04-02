package modules

import (
	"encoding/json"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/qixalite/hatch/modules/services"
	"strings"
)

func SetupWhitelistAPI(app *fiber.App, password string) {
	app.Get("/whitelist/enable", authPassword(password, func(ctx *fiber.Ctx) error {
		res, err := services.ExecRconCommand(serverStatus.GetRconServer(), "qix_restrict_players")
		if err != nil {
			return ctx.SendStatus(400)
		}

		fmt.Println("[WHITELIST]", res)

		if strings.Contains(res, "\"qix_restrict_players\" = \"0\"") {
			return ctx.JSON(false)
		} else if strings.Contains(res, "\"qix_restrict_players\" = \"1\"") {
			return ctx.JSON(true)
		}

		return ctx.SendStatus(500)
	}))
	app.Post("/whitelist/enable", authPassword(password, func(ctx *fiber.Ctx) error {
		_, err := services.ExecRconCommand(serverStatus.GetRconServer(), "qix_restrict_players 1")
		if err != nil {
			return ctx.SendStatus(400)
		}

		return ctx.SendStatus(200)
	}))
	app.Delete("/whitelist/enable", authPassword(password, func(ctx *fiber.Ctx) error {
		_, err := services.ExecRconCommand(serverStatus.GetRconServer(), "qix_restrict_players 0")
		if err != nil {
			return ctx.SendStatus(400)
		}

		return ctx.SendStatus(200)
	}))

	app.Post("/whitelist/player", authPassword(password, func(ctx *fiber.Ctx) error {
		body := ctx.Body()
		var player struct {
			Steam string `json:"steam"`
			Name  string `json:"name"`
			Class string `json:"class"`
			Team  string `json:"team"`
		}
		err := json.Unmarshal(body, &player)
		if err != nil {
			return ctx.SendStatus(400)
		}

		team := 0
		switch player.Team {
		case "red":
			team = 1
			break
		case "blu":
			team = 2
			break
		}

		class := 0
		switch player.Class {
		case "scout":
			class = 1
			break
		case "sniper":
			class = 2
			break
		case "soldier":
			class = 3
			break
		case "demoman":
			class = 4
			break
		case "medic":
			class = 5
			break
		case "heavy":
			class = 6
			break
		case "pyro":
			class = 7
			break
		case "spy":
			class = 8
			break
		case "engineer":
			class = 9
			break
		}

		cmd := fmt.Sprintf("qix_add_player %s %d %d %s", player.Steam, team, class, player.Name)
		fmt.Println("[WHITELIST]", cmd)
		res, err := services.ExecRconCommand(serverStatus.GetRconServer(), cmd)
		if err != nil {
			return ctx.SendStatus(400)
		}

		if strings.Contains(res, fmt.Sprintf("Added Player %s", player.Steam)) {
			return ctx.SendStatus(200)
		}

		return ctx.SendStatus(500)
	}))

	app.Delete("/whitelist/player/:steam", authPassword(password, func(ctx *fiber.Ctx) error {
		steam := ctx.Params("steam")
		res, err := services.ExecRconCommand(serverStatus.GetRconServer(), fmt.Sprintf("qix_remove_player %s", steam))
		if err != nil {
			return ctx.SendStatus(400)
		}

		if strings.Contains(res, fmt.Sprintf("Removed Player %s", steam)) {
			return ctx.SendStatus(200)
		}

		return ctx.SendStatus(500)
	}))

	fmt.Println("[HATCH MODULE] Started WhitelistAPI Module")
}
