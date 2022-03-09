package main

import (
	"flag"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/qixalite/hatch/modules"
	"github.com/qixalite/hatch/modules/services"
	"github.com/tjgq/broadcast"
)

var version = "2.0.5"

func main() {
	var lighthouseId, address, password, elasticHost, elasticChatIndex, elasticRconIndex, vanguardApi, vanguardSecret string
	flag.StringVar(&lighthouseId, "lighthouseId", "", "Lighthouse ID for this instance.")
	flag.StringVar(&address, "address", ":4000", "Address to listen at.")
	flag.StringVar(&password, "password", "", "Password that needs to passed as query to allow connections.")
	flag.StringVar(&elasticHost, "elasticHost", "", "URL to elastic search instance.")
	flag.StringVar(&elasticChatIndex, "elasticChatIndex", "plux-chat-development", "ElasticSearch index to store chat messages in")
	flag.StringVar(&elasticRconIndex, "elasticRconIndex", "plux-rcon-development", "ElasticSearch index to store rcon logs in")
	flag.StringVar(&vanguardApi, "vanguardApi", "http://api.qixalite.com/vanguard", "Vanguard API")
	flag.StringVar(&vanguardSecret, "vanguardSecret", "DVFDCWQ71AZRP279WTT8EYP35D0FHRSF", "Vanguard client secret")
	flag.Parse()

	app := fiber.New(fiber.Config{
		BodyLimit: 1024 * 1024 * 512,
	})

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString(fmt.Sprintf("Qixalite Hatch %s", version))
	})

	// Create thread channels
	var logBroadcast broadcast.Broadcaster
	var eventsBroadcast broadcast.Broadcaster

	// Initialize services
	services.InitializeElasticSearch(elasticHost)
	services.InitializeVanguard(vanguardApi, vanguardSecret)

	// Setup APIs
	modules.SetupFilesAPI(app, password)
	modules.SetupStatusAPI(app, password, lighthouseId)

	// Start tasks
	go modules.ReadServerLogs(&logBroadcast)
	go modules.FireEventsFromLogs(&logBroadcast, &eventsBroadcast)
	go modules.ReadChatFromLogs(&eventsBroadcast, lighthouseId, elasticChatIndex)
	go modules.ReadRconFromLogs(&eventsBroadcast, lighthouseId, elasticRconIndex)
	go modules.ReadMatchStatusFromLogs(&eventsBroadcast)
	go modules.ReadPlayerStatusFromLogs(&eventsBroadcast)
	go modules.MonitorServer(&eventsBroadcast)

	fmt.Println("Starting hatch version", version)
	err := app.Listen(address)
	if err != nil {
		fmt.Println("Server failed to start due to error", err)
		return
	}

	select {}
}
