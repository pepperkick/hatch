package modules

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"html"
	"net/http"
	"os"
	"path"
	"strings"
)

type webFile struct {
	Name       string `json:"name"`
	URL        string `json:"url"`
	Size       int64  `json:"size"`
	ClassName  string `json:"className"`
	ModifiedAt string `json:"modifiedAt"`
}

func SetupFilesAPI(app *fiber.App, password string) {
	app.Get("/files/logs", authPassword(password, presentFiles("./tf/logs", "")))
	app.Get("/files/logs/*", authPassword(password, servePathFile("./tf/logs")))
	app.Get("/files/demos", authPassword(password, presentFiles("./tf", ".dem")))
	app.Get("/files/demos/*", authPassword(password, servePathFile("./tf")))
	app.Get("/maps", authPassword(password, presentFiles("./tf/maps", "")))
	app.Get("/maps/*", authPassword(password, servePathFile("./tf/maps")))
	app.Post("/maps", authPassword(password, func(ctx *fiber.Ctx) error {
		file, err := ctx.FormFile("map")

		if err != nil {
			return err
		}

		if ctx.SaveFile(file, fmt.Sprintf("./tf/maps/%s", file.Filename)) != nil {
			return ctx.SendStatus(500)
		}

		return ctx.SendStatus(201)
	}))

	fmt.Println("[HATCH MODULE] Started FilesAPI Module")
}

func servePathFile(path string) fiber.Handler {
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

func presentFiles(p string, filter string) fiber.Handler {
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
		webFiles := make([]webFile, 0)
		for _, fi := range fileInfo {
			name := fi.Name()
			fm[name] = fi
			if strings.Contains(name, filter) {
				fi = fm[name]
				webFile := webFile{
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
