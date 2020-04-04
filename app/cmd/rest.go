package cmd

import (
	"fmt"
	"net/http"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"

	docmodule "github.com/LieAlbertTriAdrian/go-openapi-converter/modules/document"
)

var restCmd = &cobra.Command{
	Use:   "rest",
	Short: "Start REST server",
	Run:   startRestServer,
}

func startRestServer(cmd *cobra.Command, args []string) {
	logrus.Info("Starting http server at port 3000")

	fs := http.FileServer(http.Dir("./static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	http.HandleFunc("/openapi-conversions", func(w http.ResponseWriter, r *http.Request) {
		logrus.Info("Request is coming to openapi-conversion")
		switch r.Method {
		case "POST":
			// Call ParseForm() to parse the raw query and update r.PostForm and r.Form.
			if err := r.ParseForm(); err != nil {
				fmt.Fprintf(w, "ParseForm() err: %v", err)
				return
			}

			openAPISpecInString := r.FormValue("openapi-spec")

			doc, err := docmodule.ReadOpenAPIFromString(openAPISpecInString)
			if err != nil {
				fmt.Fprintf(w, "Error reading open api string %v", err)
				return
			}

			fmt.Println(doc)
		default:
			fmt.Fprintf(w, "Sorry, only POST methods is supported.")
		}

		return
	})

	// http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
	// 	fmt.Fprintf(w, "Welcome to my website!")
	// })

	http.ListenAndServe(":3000", nil)
	logrus.Info("Started http server at port 3000")
}
