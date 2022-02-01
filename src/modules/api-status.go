package modules

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
)

type ServerStatus struct {
	Password     string         `json:"password"`
	RconPassword string         `json:"rconPassword"`
	TvPassword   string         `json:"TvPassword"`
	LighthouseID string         `json:"lighthouseId"`
	Players      []PlayerStatus `json:"players"`
	Matches      []MatchStatus  `json:"matches"`
}

var status ServerStatus

func SetupStatusAPI(app *fiber.App, password string, lighthouseId string) {
	status.LighthouseID = lighthouseId
	status.Matches = []MatchStatus{}
	status.Players = []PlayerStatus{}

	app.Get("/status", authPassword(password, func(c *fiber.Ctx) error {
		ReadServerInfo()

		status.Matches = matches
		status.Players = activePlayers

		return c.JSON(status)
	}))

	fmt.Println("[HATCH MODULE] Started StatusAPI Module")
}
