package services

import (
	"fmt"
	"net/http"
)

var (
	VanguardSecret      string
	VanguardApi         string
	InitializedVanguard = false
)

func InitializeVanguard(api string, secret string) {
	fmt.Printf("[VANGUARD] %s - %s\n", api, secret)
	if api == "" || secret == "" {
		return
	}

	VanguardApi = api
	VanguardSecret = secret
	InitializedVanguard = true
}

func CheckPlayerBan(id string) bool {
	action := "LIGHTHOUSE_SERVER_CONNECT"
	url := fmt.Sprintf("%s/bans/steam/%s/%s", VanguardApi, id, action)
	fmt.Println("[VANGUARD] URL", url)

	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		fmt.Println("[VANGUARD] Failed to check ban with vanguard", err)
		return false
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", VanguardSecret))
	res, err := client.Do(req)
	if err != nil {
		fmt.Println("[VANGUARD] Failed to check ban with vanguard", err)
		return false
	}

	fmt.Println("[VANGUARD]", res)

	if res.StatusCode == 200 {
		return true
	} else {
		return false
	}
}
