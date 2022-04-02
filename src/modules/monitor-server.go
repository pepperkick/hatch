package modules

import (
	"fmt"
	"github.com/qixalite/hatch/modules/services"
	"github.com/tjgq/broadcast"
	"time"
)

func MonitorServer(broadcast *broadcast.Broadcaster) {
	logChannel := broadcast.Listen().Ch

	regex := "Executing dedicated server config file server"
	name := "ServerConfigFileExecuted"
	AddLogMatcher(regex, name)

	fmt.Println("[HATCH MODULE] Started MonitorServer Module")

	ticker := time.NewTicker(time.Second * 30)
	quit := make(chan struct{})
	go func() {
		for {
			select {
			case <-ticker.C:
				onPlayerConnect([]string{})
			case <-quit:
				ticker.Stop()
				return
			}
		}
	}()

	for {
		select {
		case event := <-logChannel:
			msg := event.(LogEvent)

			if msg.Name == name {
				getServerInfo()
			}

			if services.InitializedVanguard && msg.Name == "PlayerConnected" {
				onPlayerConnect(msg.Args)
			}

			if services.InitializedVanguard && msg.Name == "PlayerEnteredGame" {
				onPlayerConnect(msg.Args)
			}
		}
	}
}

func getServerInfo() {

}

func onPlayerConnect(args []string) {
	fmt.Println("[MONITOR_PLAYERS] onPlayerConnect")
	ReadServerInfo()

	// Check for banned players
	for _, player := range activePlayers {
		if !services.CheckPlayerBan(player.Steam) {
			continue
		}

		command := fmt.Sprintf("banid 0 %s kick", player.Steam)
		res, err := services.ExecRconCommand(serverStatus.GetRconServer(), command)
		if err != nil {
			fmt.Println("[MONITOR_PLAYERS] Failed to execute rcon command to ban the player", err)
			continue
		}

		fmt.Println("[MONITOR_PLAYERS] Vangaurd ban response for", player.Steam, res)
	}
}
