package cmd

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"

	docmodule "github.com/LieAlbertTriAdrian/go-openapi-converter/modules/document"
	guuid "github.com/google/uuid"
)

var restCmd = &cobra.Command{
	Use:   "rest",
	Short: "Start REST server",
	Run:   startRestServer,
}

func startRestServer(cmd *cobra.Command, args []string) {
	port := "80"

	logrus.Info("Starting http server at port " + port)

	fs := http.FileServer(http.Dir("./static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	http.HandleFunc("/openapi-conversions", func(w http.ResponseWriter, r *http.Request) {
		logrus.Info("Request is coming to openapi-conversion")

		currentTime := time.Now().Format("01-02-2006")
		name := currentTime + "-" + guuid.New().String()

		logrus.Info(name)

		inputName := "input-" + name + ".json"
		fileName := "output-" + name + ".docx"

		switch r.Method {
		case "POST":
			// Call ParseForm() to parse the raw query and update r.PostForm and r.Form.
			if err := r.ParseForm(); err != nil {
				fmt.Fprintf(w, "ParseForm() err: %v", err)
				return
			}

			openAPISpecInString := r.FormValue("openapi-spec")

			err := ioutil.WriteFile(inputName, []byte(openAPISpecInString), 0644)
			if err != nil {
				fmt.Fprintf(w, "Error writing input file %v", err)
				return
			}

			doc, err := docmodule.ReadDocTemplate("template/standard.docx")
			if err != nil {
				fmt.Fprintf(w, "Error reading open api string %v", err)
				return
			}

			swagger, err := docmodule.ReadOpenAPIFromString(openAPISpecInString)
			if err != nil {
				fmt.Fprintf(w, "Error reading open api string %v", err)
				return
			}

			docService := docmodule.NewDocumentHandler(doc, swagger)

			docService.BuildFrontpage()
			docService.BuildTOC()
			docService.BuildPaths()

			doc.SaveToFile(fileName)

			fileData, err := ioutil.ReadFile(fileName)
			if err != nil {
				fmt.Fprintf(w, "Error reading output data %v", err)
				return
			}

			// tell the browser the returned content should be downloaded
			w.Header().Add("Content-Disposition", "Attachment")

			http.ServeContent(w, r, fileName, time.Now(), bytes.NewReader(fileData))
		default:
			fmt.Fprintf(w, "Sorry, only POST methods is supported.")
		}

		return
	})

	// http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
	// 	fmt.Fprintf(w, "Welcome to my website!")
	// })

	http.ListenAndServe(":"+port, nil)
	logrus.Info("Started http server at port " + port)
}
