package modules

import (
	"fmt"
	"github.com/hpcloud/tail"
	"github.com/tjgq/broadcast"
	"os"
	timeUtils "time"
)

var logFile = "./tf/console.log"

func ReadServerLogs(broadcast *broadcast.Broadcaster) {
	// Wait until the file has been created
	retry, maxRetires := 0, 30
	for {
		if _, err := os.Stat(logFile); err == nil {
			break
		}

		timeUtils.Sleep(5 * timeUtils.Second)

		retry++
		if retry == maxRetires {
			fmt.Println("Failed to load console.log file")
			return
		}
	}

	fmt.Println("Reading console.log file...")

	// Read the log file line by line and tail it
	t, err := tail.TailFile(logFile, tail.Config{Follow: true})
	if err != nil {
		fmt.Println("Failed to read console.log", err)
		return
	}
	for line := range t.Lines {
		broadcast.Send(line.Text)
	}
}
