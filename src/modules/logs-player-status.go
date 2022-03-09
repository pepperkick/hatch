package modules

import (
	"fmt"
	"github.com/tjgq/broadcast"
)

type PlayerStatus struct {
	UserId int    `json:"userId"`
	Steam  string `json:"steam"`
	Name   string `json:"name"`
	IP     string `json:"ip"`
}

var activePlayers []PlayerStatus

func ReadPlayerStatusFromLogs(broadcast *broadcast.Broadcaster) {
	logChannel := broadcast.Listen().Ch

	connectRegex := "Client \"(.*?)\" connected \\((.*?)\\)."
	connectName := "PlayerConnected"
	AddLogMatcher(connectRegex, connectName)

	disconnectRegex := "Dropped (.*?) from server \\((.*?)\\)"
	disconnectName := "PlayerDisconnected"
	AddLogMatcher(disconnectRegex, disconnectName)

	enteredGameRegex := "L (.*?) - (.*?): \"(.*?)<(.*?)><(.*?)><>\" entered the game"
	enteredGameName := "PlayerEnteredGame"
	AddLogMatcher(enteredGameRegex, enteredGameName)

	steamVerifiedRegex := "L (.*?) - (.*?): \"(.*?)<(.*?)><(.*?)><>\" STEAM USERID validated"
	steamVerifiedName := "PlayerSteamVerified"
	AddLogMatcher(steamVerifiedRegex, steamVerifiedName)

	fmt.Println("[HATCH MODULE] Started ReadPlayerStatusFromLogs Module")

	for {
		select {
		case event := <-logChannel:
			msg := event.(LogEvent)

			if msg.Name == connectName {

			}
			if msg.Name == disconnectName {

			}
		}
	}
}
