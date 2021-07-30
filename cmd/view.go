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
	"context"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/hslatman/mud-cli/internal"
	"github.com/hslatman/mud-cli/web"
	"github.com/pkg/browser"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"

	"github.com/antage/eventsource"
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

		// TODO: provide option to show it in terminal with some ASCII art?
		// TODO: provide option to show it in a Wails UI (instead of default browser)?

		closeChan := make(chan struct{}, 1)
		stopChan := make(chan error, 1)

		mudHandler := newMUDHandler(json)
		streamHandler := newStreamHandler(closeChan)

		// Strip / and prepend build, so that a file `a/b.js` would be
		// found in web/build/a/b.js, but served from localhost:8080/a/b.js.
		webHandler := web.AssetHandler("/", "build")

		mux := http.NewServeMux()
		mux.Handle("/mud", mudHandler)
		mux.Handle("/heartbeat", streamHandler)
		mux.Handle("/", webHandler)
		mux.Handle("/*filepath", webHandler)

		s := &http.Server{
			Addr:           "localhost:8080",
			Handler:        mux,
			ReadTimeout:    10 * time.Second,
			WriteTimeout:   10 * time.Second,
			MaxHeaderBytes: 1 << 20,
		}

		s.RegisterOnShutdown(streamHandler.Close)

		go func() {
			<-closeChan
			log.Println("closing server ...")
			ctx := context.Background()
			ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
			defer cancel()
			stopChan <- s.Shutdown(ctx)
		}()

		// TODO: add some more logging?

		log.Println("serving at ...")
		go s.ListenAndServe()

		log.Println("go to ...")
		go browser.OpenURL("http://localhost:8080/")

		err = <-stopChan

		return err
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

type streamHandler struct {
	es eventsource.EventSource
	c  chan struct{}
}

func newStreamHandler(c chan struct{}) *streamHandler {
	es := eventsource.New(
		&eventsource.Settings{
			Timeout:        2 * time.Second,
			CloseOnTimeout: true,
			IdleTimeout:    2 * time.Second,
			Gzip:           true,
		},
		func(req *http.Request) [][]byte {
			return [][]byte{
				[]byte("X-Accel-Buffering: no"),
				[]byte("Access-Control-Allow-Origin: *"),
			}
		},
	)
	return &streamHandler{
		es: es,
		c:  c,
	}
}

func (s *streamHandler) Close() {
	s.es.Close()
}

func (s *streamHandler) ServeHTTP(resp http.ResponseWriter, req *http.Request) {

	numberOfConsumers := s.es.ConsumersCount()

	// serve sse to the (new) connection
	s.es.ServeHTTP(resp, req)

	// the routine for pings should only be started once
	if numberOfConsumers > 0 {
		return
	}

	// send a ping to all consumers every second
	go func() {
		var id int
		for {
			id++
			time.Sleep(1 * time.Second)
			s.es.SendEventMessage("ping", "message", strconv.Itoa(id))
			// break out of the loop when 0 consumers are reached again
			if s.es.ConsumersCount() == 0 {
				break
			}
		}
		s.c <- struct{}{}
	}()

}

func init() {
	rootCmd.AddCommand(viewCmd)
}
