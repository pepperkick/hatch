package main

import (
	"bytes"
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"github.com/elastic/go-elasticsearch/v7"
	"github.com/elastic/go-elasticsearch/v7/esapi"
	"github.com/gofiber/fiber/v2"
	"github.com/hpcloud/tail"
	"html"
	"log"
	"net/http"
	"os"
	"path"
	"regexp"
	"strings"
	timeUtils "time"
)

var version = "1.1.0"

type WebFile struct {
	Name       string `json:"name"`
	URL        string `json:"url"`
	Size       int64  `json:"size"`
	ClassName  string `json:"className"`
	ModifiedAt string `json:"modifiedAt"`
}

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

func main() {
	var lighthouseId, address, password, elasticHost, elasticChatIndex string
	flag.StringVar(&lighthouseId, "lighthouseId", "", "Lighthouse ID for this instance.")
	flag.StringVar(&address, "address", ":4000", "Address to listen at.")
	flag.StringVar(&password, "password", "", "Password that needs to passed as query to allow connections.")
	flag.StringVar(&elasticHost, "elasticHost", "", "URL to elastic search instance.")
	flag.StringVar(&elasticChatIndex, "elasticChatIndex", "plux-chat-development", "Index to store chat in")
	flag.Parse()

	go setupServer(address, password)
	go readChatFromLogs(lighthouseId, elasticHost, elasticChatIndex)

	select {}
}

func setupServer(address string, password string) {
	app := fiber.New(fiber.Config{
		BodyLimit: 1024 * 1024 * 512,
	})

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString(fmt.Sprintf("Qixalite Hatch %s", version))
	})

	app.Get("/files/logs", AuthPassword(password, PresentFiles("./tf/logs", "")))
	app.Get("/files/logs/*", AuthPassword(password, ServePathFile("./tf/logs")))
	app.Get("/files/demos", AuthPassword(password, PresentFiles("./tf", ".dem")))
	app.Get("/files/demos/*", AuthPassword(password, ServePathFile("./tf")))
	app.Get("/maps", AuthPassword(password, PresentFiles("./tf/maps", "")))
	app.Get("/maps/*", AuthPassword(password, ServePathFile("./tf/maps")))
	app.Post("/maps", AuthPassword(password, func(ctx *fiber.Ctx) error {
		file, err := ctx.FormFile("map")

		if err != nil {
			return err
		}

		if ctx.SaveFile(file, fmt.Sprintf("./tf/maps/%s", file.Filename)) != nil {
			return ctx.SendStatus(500)
		}

		return ctx.SendStatus(201)
	}))

	err := app.Listen(address)
	if err != nil {
		fmt.Println("Server failed to start due to error", err)
		return
	}
}

func readChatFromLogs(lighthouseId string, elasticHost string, elasticIndex string) {
	logFile := "./tf/console.log"
	es, err := elasticsearch.NewClient(elasticsearch.Config{
		Addresses: []string{elasticHost},
	})
	if err != nil {
		fmt.Println("Failed to connect with elasticsearch", err)
		return
	}

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
		// Regex match log line for say and say_team logs
		re := regexp.MustCompile("L (.*?) - (.*?): \"(.*?)<(.*?)><(.*?)><(.*?)>\" (say|say_team) \"(.*?)\"")
		arr := re.FindStringSubmatch(line.Text)
		if len(arr) > 0 {
			// Read the data from the log line and create a structure out of it
			date, time, name, userid, steam, team, command, content :=
				arr[1], arr[2], arr[3], arr[4], arr[5], arr[6], arr[7], arr[8]

			// Create a timestamp from log line date and time
			layout := "01/02/2006 - 15:04:05"
			str := fmt.Sprintf("%s - %s", date, time)
			timestamp, err := timeUtils.Parse(layout, str)

			// Covert log line data to struct to pass in request later
			playerMsg := PlayerMessage{
				CreatedAt:    timestamp,
				Name:         name,
				UserID:       userid,
				Steam:        steam,
				Team:         team,
				Command:      command,
				Content:      content,
				LighthouseID: lighthouseId,
			}
			playerMsgBytes, _ := json.Marshal(playerMsg)

			// Index the player chat in elastic search
			req := esapi.IndexRequest{
				Index:   elasticIndex,
				Body:    bytes.NewReader(playerMsgBytes),
				Refresh: "true",
			}

			res, err := req.Do(context.Background(), es)
			if err != nil {
				log.Fatalf("Error getting response: %s", err)
			}

			if res.IsError() {
				log.Printf("[%s] Error indexing document", res.Status())
			} else {
				// Deserialize the response into a map.
				var r map[string]interface{}
				if err := json.NewDecoder(res.Body).Decode(&r); err != nil {
					log.Printf("Error parsing the response body: %s", err)
				} else {
					// Print the response status and indexed document version.
					log.Printf("[%s] %s; version=%d", res.Status(), r["result"], int(r["_version"].(float64)))
				}
			}
		}
	}
}

func ServePathFile(path string) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		fileName := ctx.Params("*")
		filePath := path + "/" + fileName
		stat, err := os.Stat(filePath)
		if err != nil {
			return err
		}
		if stat.IsDir() {
			return ctx.SendStatus(403)
		}
		return ctx.Download(filePath, fileName)
	}
}

func AuthPassword(password string, next fiber.Handler) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		pass := ctx.Query("password")

		if password != "" && pass != password {
			return ctx.SendStatus(403)
		}

		return next(ctx)
	}
}

func PresentFiles(p string, filter string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		fs := http.Dir(p)

		if _, err := os.Stat(p); os.IsNotExist(err) {
			err := os.Mkdir(p, 0777)
			if err != nil {
				return err
			}
		}

		file, err := fs.Open("")
		if err != nil {
			return err
		}

		fileInfo, err := file.Readdir(-1)
		if err != nil {
			return err
		}

		fm := make(map[string]os.FileInfo, len(fileInfo))
		filenames := make([]string, 0, len(fileInfo))
		webFiles := make([]WebFile, 0)
		for _, fi := range fileInfo {
			name := fi.Name()
			fm[name] = fi
			if strings.Contains(name, filter) {
				fi = fm[name]
				webFile := WebFile{
					Name:       name,
					URL:        html.EscapeString(path.Join(c.Path() + "/" + name)),
					Size:       0,
					ClassName:  "dir",
					ModifiedAt: fi.ModTime().String(),
				}

				if !fi.IsDir() {
					webFile.Size = fi.Size()
					webFile.ClassName = "file"
				}

				filenames = append(filenames, name)
				webFiles = append(webFiles, webFile)
			}
		}

		return c.JSON(webFiles)
	}
}
