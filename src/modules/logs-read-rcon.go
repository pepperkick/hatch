package modules

import (
	"encoding/json"
	"fmt"
	"github.com/qixalite/hatch/modules/services"
	"github.com/tjgq/broadcast"
	timeUtils "time"
)

type RconCommandLog struct {
	CreatedAt    timeUtils.Time `json:"createdAt"`
	Origin       string         `json:"origin"`
	Command      string         `json:"command"`
	LighthouseID string         `json:"lighthouseId"`
}

func ReadRconFromLogs(broadcast *broadcast.Broadcaster, lighthouseId string, elasticIndex string) {
	if lighthouseId == "" || elasticIndex == "" || services.ElasticSearchClient == nil {
		return
	}

	logChannel := broadcast.Listen().Ch

	regex := "rcon from \"(.*?)\": command \"(.*?)\""
	name := "RconCommand"
	AddLogMatcher(regex, name)

	fmt.Println("[HATCH MODULE] Started ReadRconFromLogs Module")

	for {
		select {
		case event := <-logChannel:
			msg := event.(LogEvent)

			if msg.Name == name {
				readRconLog(msg.Args, lighthouseId, elasticIndex)
			}
		}
	}
}

func readRconLog(arr []string, lighthouseId string, elasticIndex string) {
	// Read the data from the msg line and create a structure out of it
	origin, command := arr[1], arr[2]

	// Covert msg line data to struct to pass in request later
	msg := RconCommandLog{
		CreatedAt:    timeUtils.Now(),
		Origin:       origin,
		Command:      command,
		LighthouseID: lighthouseId,
	}
	fmt.Println("[RCON EVENT]", msg)
	msgBytes, _ := json.Marshal(msg)

	// Index the message
	services.IndexMessageInElasticSearch(elasticIndex, msgBytes)
}
