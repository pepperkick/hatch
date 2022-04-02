package modules

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/qixalite/hatch/modules/services"
)

type ServerStatus struct {
	IP           string         `json:"serverIp"`
	Port         string         `json:"serverPort"`
	Password     string         `json:"password"`
	RconPassword string         `json:"rconPassword"`
	TvPassword   string         `json:"TvPassword"`
	LighthouseID string         `json:"lighthouseId"`
	Players      []PlayerStatus `json:"players"`
	Matches      []MatchStatus  `json:"matches"`
}

var serverStatus ServerStatus

func SetupStatusAPI(app *fiber.App, password string, lighthouseId string) {
	serverStatus.LighthouseID = lighthouseId
	serverStatus.Matches = []MatchStatus{}
	serverStatus.Players = []PlayerStatus{}

	app.Get("/status", authPassword(password, func(c *fiber.Ctx) error {
		ReadServerInfo()

		serverStatus.Matches = matches
		serverStatus.Players = activePlayers

		return c.JSON(serverStatus)
	}))

	fmt.Println("[HATCH MODULE] Started StatusAPI Module")
}

func (s *ServerStatus) GetRconServer() services.RconServer {
	return services.RconServer{
		IP:           s.IP,
		Port:         s.Port,
		RconPassword: s.RconPassword,
	}
}
