package modules

import (
	"fmt"
	"github.com/tjgq/broadcast"
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
		}
	}
}

func getServerInfo() {

}
