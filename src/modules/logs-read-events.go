package modules

import (
	"fmt"
	"github.com/tjgq/broadcast"
	"regexp"
)

type LogEvent struct {
	Name string
	Log  string
	Args []string
}

type LogMatcher struct {
	Regex string
	Name  string
}

var matchers []LogMatcher

func FireEventsFromLogs(broadcast *broadcast.Broadcaster, out *broadcast.Broadcaster) {
	logChannel := broadcast.Listen().Ch

	fmt.Println("[HATCH MODULE] Started FireEventsFromLog Module")

	for {
		select {
		case line := <-logChannel:
			msg := line.(string)

			fmt.Println("[LOG LINE]", msg)
			for _, element := range matchers {
				re := regexp.MustCompile(element.Regex)
				arr := re.FindStringSubmatch(msg)

				if len(arr) == 0 {
					continue
				}

				event := LogEvent{
					Name: element.Name,
					Log:  msg,
					Args: arr,
				}

				fmt.Println("[LOG EVENT]", event.Name, event.Args)
				out.Send(event)
			}
		}
	}
}

func AddLogMatcher(regex string, name string) {
	matchers = append(matchers, LogMatcher{
		Regex: regex,
		Name:  name,
	})
}
