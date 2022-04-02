package services

import (
	"fmt"
	"github.com/forewing/csgo-rcon"
	"time"
)

type RconServer struct {
	IP           string `json:"serverIp"`
	Port         string `json:"serverPort"`
	RconPassword string `json:"rconPassword"`
}

func ExecRconCommand(serverStatus RconServer, command string) (string, error) {
	conn := rcon.New(fmt.Sprintf("%s:%s", serverStatus.IP, serverStatus.Port), serverStatus.RconPassword, time.Second*1)

	res, err := conn.Execute(command)
	if err != nil {
		return "", err
	}

	return res, nil
}
