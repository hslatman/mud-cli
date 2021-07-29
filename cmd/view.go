/*
Copyright © 2021 NAME HERE <EMAIL ADDRESS>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"log"
	"net/http"
	"time"

	"github.com/hslatman/mud-cli/internal"
	"github.com/hslatman/mud-cli/web"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

// viewCmd represents the view command
var viewCmd = &cobra.Command{
	Use:   "view",
	Short: "Provides a graphical view of a MUD file",
	Args:  cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {

		filepath := args[0]
		mudfile, err := internal.ReadMUDFileFrom(filepath)
		if err != nil {
			return errors.Wrap(err, "could not get contents")
		}

		json, err := internal.JSON(mudfile)
		if err != nil {
			return errors.Wrap(err, "getting JSON representation of MUD file failed")
		}

		// TODO: open browser; show MUD Visualizer with the chosen MUD file
		// TODO: provide option to show it in terminal?

		mudHandler := newMUDHandler(json)

		// Strip / and prepend build, so that a file `a/b.js` would be
		// found in web/build/a/b.js, but served from localhost:8080/a/b.js.
		webHandler := web.AssetHandler("/", "build")

		mux := http.NewServeMux()
		mux.Handle("/mud", mudHandler)
		mux.Handle("/", webHandler)
		mux.Handle("/*filepath", webHandler)

		s := &http.Server{
			Addr:           ":8080",
			Handler:        mux,
			ReadTimeout:    10 * time.Second,
			WriteTimeout:   10 * time.Second,
			MaxHeaderBytes: 1 << 20,
		}

		log.Fatal(s.ListenAndServe()) // TODO: add some logging?

		return nil
	},
}

type mudHandler struct {
	json string
}

func (m *mudHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte(m.json))
}

func newMUDHandler(json string) http.Handler {
	return &mudHandler{
		json: json,
	}
}

func init() {
	rootCmd.AddCommand(viewCmd)
}
