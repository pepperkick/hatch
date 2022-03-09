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
	if _, err := os.Stat(infoFile); err != nil {
		fmt.Println("Server info file not found, skipping check")
		return
	}

	file, err := os.Open(infoFile)
	if err != nil {
		fmt.Println(err)
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
		case "server_ip":
			serverStatus.IP = parts[1]
			break
		case "server_port":
			serverStatus.Port = parts[1]
			break
		case "password":
			serverStatus.Password = parts[1]
			break
		case "rcon_password":
			serverStatus.RconPassword = parts[1]
			break
		case "tv_password":
			serverStatus.TvPassword = parts[1]
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
