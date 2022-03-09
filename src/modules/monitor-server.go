package modules

import (
	"fmt"
	"github.com/qixalite/hatch/modules/services"
	"github.com/riking/homeapi/rcon"
	"github.com/tjgq/broadcast"
	"strconv"
)

func MonitorServer(broadcast *broadcast.Broadcaster) {
	logChannel := broadcast.Listen().Ch

	regex := "Executing dedicated server config file server"
	name := "ServerConfigFileExecuted"
	AddLogMatcher(regex, name)

	fmt.Println("[HATCH MODULE] Started MonitorServer Module")

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

	for _, player := range activePlayers {
		if !services.CheckPlayerBan(player.Steam) {
			continue
		}

		port, _ := strconv.Atoi(serverStatus.Port)
		conn, err := rcon.Dial(serverStatus.IP, port, serverStatus.RconPassword)
		if err != nil {
			fmt.Println("[MONITOR_PLAYERS] Failed to create rcon connection to ban the player", err)
			continue
		}

		command := fmt.Sprintf("banid 0 %s kick", player.Steam)
		if _, err := conn.Command(command); err != nil {
			fmt.Println("[MONITOR_PLAYERS] Failed to execute rcon command to ban the player", err)
			continue
		}

		if err := conn.Close(); err != nil {
			fmt.Println("[MONITOR_PLAYERS] Failed to close rcon connection", err)
			continue
		}
	}
}
