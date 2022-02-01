package modules

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"
)

var infoFile = "./tf/server.info"

func ReadServerInfo() {
	file, err := os.Open(infoFile)
	if err != nil {
		log.Fatal(err)
		return
	}
	defer file.Close()

	activePlayers = []PlayerStatus{}
	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)
	for scanner.Scan() {
		text := scanner.Text()
		parts := strings.Split(text, ": ")
		fmt.Println("[SERVER INFO]", text)

		switch parts[0] {
		case "password":
			status.Password = parts[1]
			break
		case "rcon_password":
			status.RconPassword = parts[1]
			break
		case "tv_password":
			status.TvPassword = parts[1]
			break
		case "connected_player":
			re := regexp.MustCompile("connected_player: (.*?) (.*?) \"(.*?)\" (.*)")
			arr := re.FindStringSubmatch(text)

			if len(arr) == 0 {
				break
			}

			userId, _ := strconv.Atoi(arr[1])
			activePlayers = append(activePlayers, PlayerStatus{
				UserId: userId,
				Steam:  arr[2],
				Name:   arr[3],
				IP:     arr[4],
			})

			break
		}
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
}
