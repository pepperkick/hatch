package main

import (
	"flag"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"html"
	"net/http"
	"os"
	"path"
	"strings"
)

type WebFile struct {
	Name       string `json:"name"`
	URL        string `json:"url"`
	Size       int64  `json:"size"`
	ClassName  string `json:"className"`
	ModifiedAt string `json:"modifiedAt"`
}

func main() {
	var address, password string
	flag.StringVar(&address, "address", ":4000", "Address to listen at.")
	flag.StringVar(&password, "password", "", "Password that needs to passed as query to allow connections.")
	flag.Parse()

	app := fiber.New()

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Qixalite Hatch v1.0.0")
	})

	app.Get("/files/logs", AuthPassword(password, PresentFiles("./tf/logs", "")))
	app.Get("/files/logs/*", AuthPassword(password, ServePathFile("./tf/logs")))
	app.Get("/files/demos", AuthPassword(password, PresentFiles("./tf", ".dem")))
	app.Get("/files/demos/*", AuthPassword(password, ServePathFile("./tf")))

	err := app.Listen(address)
	if err != nil {
		fmt.Println("Server failed to start due to error", err)
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
