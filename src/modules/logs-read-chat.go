package modules

import (
	"encoding/json"
	"fmt"
	"github.com/qixalite/hatch/modules/services"
	"github.com/tjgq/broadcast"
	timeUtils "time"
)

type PlayerMessage struct {
	CreatedAt    timeUtils.Time `json:"createdAt"`
	Name         string         `json:"name"`
	UserID       string         `json:"userid"`
	Steam        string         `json:"steam"`
	Team         string         `json:"team"`
	Command      string         `json:"command"`
	Content      string         `json:"content"`
	LighthouseID string         `json:"lighthouseId"`
}

func ReadChatFromLogs(broadcast *broadcast.Broadcaster, lighthouseId string, elasticIndex string) {
	if lighthouseId == "" || elasticIndex == "" || services.ElasticSearchClient == nil {
		return
	}

	logChannel := broadcast.Listen().Ch

	regex := "L (.*?) - (.*?): \"(.*?)<(.*?)><(.*?)><(.*?)>\" (say|say_team) \"(.*?)\""
	name := "ChatMessage"
	AddLogMatcher(regex, name)

	fmt.Println("[HATCH MODULE] Started ReadChatFromLogs Module")

	for {
		select {
		case event := <-logChannel:
			msg := event.(LogEvent)

			if msg.Name == name {
				readChatLog(msg.Args, lighthouseId, elasticIndex)
			}
		}
	}
}

func readChatLog(arr []string, lighthouseId string, elasticIndex string) {
	// Read the data from the msg line and create a structure out of it
	date, time, name, userid, steam, team, command, content :=
		arr[1], arr[2], arr[3], arr[4], arr[5], arr[6], arr[7], arr[8]

	// Create a timestamp from msg line date and time
	layout := "01/02/2006 - 15:04:05"
	str := fmt.Sprintf("%s - %s", date, time)
	timestamp, err := timeUtils.Parse(layout, str)

	if err != nil {
		fmt.Println("Failed to parse timestamp", err)
		return
	}

	// Covert msg line data to struct to pass in request later
	msg := PlayerMessage{
		CreatedAt:    timestamp,
		Name:         name,
		UserID:       userid,
		Steam:        steam,
		Team:         team,
		Command:      command,
		Content:      content,
		LighthouseID: lighthouseId,
	}
	fmt.Println("[CHAT EVENT]", msg)
	msgBytes, _ := json.Marshal(msg)

	// Index the message
	services.IndexMessageInElasticSearch(elasticIndex, msgBytes)
}
