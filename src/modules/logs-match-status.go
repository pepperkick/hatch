package modules

import (
	"fmt"
	"github.com/tjgq/broadcast"
	"strconv"
	"strings"
	"time"
)

type MatchStatus struct {
	Status       string         `json:"status"`
	EndedMidGame string         `json:"endedMidGame"`
	TeamScore    map[string]int `json:"teamScore"`
	StartTime    time.Time      `json:"startTime"`
	EndTime      time.Time      `json:"endTime"`
	LogFile      string         `json:"logFile"`
	LogstfID     string         `json:"logstfId"`
	LogstfURL    string         `json:"logstfUrl"`
	DemoFile     string         `json:"demoFile"`
	DemostfID    string         `json:"demostfId"`
	DemostfURL   string         `json:"demostfUrl"`
}

var matches []MatchStatus
var currentMatchIndex = 0

func ReadMatchStatusFromLogs(broadcast *broadcast.Broadcaster) {
	logChannel := broadcast.Listen().Ch

	matchStartedRegex := "\\[logs\\.tf\\] Match started"
	matchStartedName := "MatchStarted"
	AddLogMatcher(matchStartedRegex, matchStartedName)

	matchResetRegex := "\\[logs\\.tf\\] Match reset"
	matchResetName := "MatchReset"
	AddLogMatcher(matchResetRegex, matchResetName)

	matchEndedRegex := "\\[logs\\.tf\\] Match ended \\(midgame (.*?)\\)"
	matchEndedName := "MatchEnded"
	AddLogMatcher(matchEndedRegex, matchEndedName)

	teamCurrentScoreRegex := "L (.*?) - (.*?): Team \\\"(.*?)\\\" current score \\\"(.*?)\\\" with \\\"(.*?)\\\" players"
	teamCurrentScoreName := "TeamCurrentScore"
	AddLogMatcher(teamCurrentScoreRegex, teamCurrentScoreName)

	teamFinalScoreRegex := "L (.*?) - (.*?): Team \\\"(.*?)\\\" final score \\\"(.*?)\\\" with \\\"(.*?)\\\" players"
	teamFinalScoreName := "TeamFinalScore"
	AddLogMatcher(teamFinalScoreRegex, teamFinalScoreName)

	logstfFileRegex := "L (.*?) - (.*?): \\[logs\\.tf\\] Uploading (.*)"
	logstfFileName := "LogstfFile"
	AddLogMatcher(logstfFileRegex, logstfFileName)

	logstfUploadedRegex := "L (.*?) - (.*?): \\[logs\\.tf\\] Uploaded (.*?) \\((.*?)\\)"
	logstfUploadedName := "LogstfUploaded"
	AddLogMatcher(logstfUploadedRegex, logstfUploadedName)

	demostfFileRegex := "L (.*?) - (.*?): \\[demos\\.tf\\] Uploading (.*)"
	demostfFileName := "DemostfFile"
	AddLogMatcher(demostfFileRegex, demostfFileName)

	demostfUploadedRegex := "L (.*?) - (.*?): \\[demos\\.tf\\] Response STV available at: (.*)"
	demostfUploadedName := "DemostfUploaded"
	AddLogMatcher(demostfUploadedRegex, demostfUploadedName)

	fmt.Println("[HATCH MODULE] Started ReadMatchStatusFromLogs Module")

	for {
		select {
		case event := <-logChannel:
			msg := event.(LogEvent)

			switch msg.Name {
			case matchStartedName:
				onMatchStart()
				break
			case matchResetName:
				onMatchReset()
				break
			case matchEndedName:
				onMatchEnded(msg.Args)
				break
			case teamCurrentScoreName:
				onTeamScore(msg.Args)
				break
			case teamFinalScoreName:
				onTeamScore(msg.Args)
				break
			case logstfFileName:
				onLogFile(msg.Args)
				break
			case logstfUploadedName:
				onLogUploaded(msg.Args)
				break
			case demostfFileName:
				onDemoFile(msg.Args)
				break
			case demostfUploadedName:
				onDemoUploaded(msg.Args)
				break
			}
		}
	}
}

func onMatchStart() {
	match := MatchStatus{
		Status:    "live",
		StartTime: time.Now(),
		TeamScore: map[string]int{},
	}
	matches = append(matches, match)
	currentMatchIndex = len(matches) - 1
}

func onMatchReset() {
	matches[currentMatchIndex].Status = "reset"
}

func onMatchEnded(args []string) {
	matches[currentMatchIndex].Status = "ended"
	matches[currentMatchIndex].EndedMidGame = args[1]
	matches[currentMatchIndex].EndTime = time.Now()
}

func onTeamScore(args []string) {
	team, score := args[3], args[4]
	matches[currentMatchIndex].TeamScore[team], _ = strconv.Atoi(score)
}

func onLogFile(args []string) {
	file := args[3]
	matches[currentMatchIndex].LogFile = file
}

func onLogUploaded(args []string) {
	id, url := args[3], args[4]
	matches[currentMatchIndex].LogstfID = id
	matches[currentMatchIndex].LogstfURL = url
}

func onDemoFile(args []string) {
	file := args[3]
	matches[currentMatchIndex].DemoFile = file
}

func onDemoUploaded(args []string) {
	url := args[3]
	parts := strings.Split(url, "/")
	id := parts[len(parts)-1]

	matches[currentMatchIndex].DemostfID = id
	matches[currentMatchIndex].DemostfURL = url
}
