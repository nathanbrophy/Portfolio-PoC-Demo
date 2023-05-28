/*
Copyright 2023 Nathan Brophy

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
package server

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"sync"

	"github.com/sirupsen/logrus"

	automationhandler "github.com/nathanbrophy/portfolio-demo/server/pkg/automationHandler"
)

const (
	PORT                      int    = 8081
	METHOD_NOT_SUPPORTED      string = "Method Not Supported"
	INTERNAL_SERVER_ERROR_MSG string = "Encountered an unexpected error"
)

var (
	defaultHandler automationhandler.AutomateEverything = automationhandler.Default()
	logger         *logrus.Logger
)

// setupLogger sets the log instance up the server will use
func setupLogger() {
	logger = logrus.New()
	logger.Formatter = new(logrus.JSONFormatter)
	logger.Out = os.Stdout
}

// handleExample is a handler func that registers the /example route to return the GET method
func handleExample(w http.ResponseWriter, req *http.Request) {
	logger.Info(fmt.Sprintf("connection handled for client %s", req.RemoteAddr))

	// Mutates the in memory instance for the handler to
	// set the time to when the request was received.
	defaultHandler.SetTimeNow()

	var (
		clientMessage    []byte
		clientStatusCode int = http.StatusOK
	)

	// Only GET methods are supported at this time
	// all other REST requests will receive the HTTP Status Not Allowed
	// status code with an appropriate error message.
	switch req.Method {
	case http.MethodGet:
		_clientMessage, marshalErr := defaultHandler.String()
		if marshalErr != nil {
			clientStatusCode = http.StatusInternalServerError
			clientMessage = []byte(INTERNAL_SERVER_ERROR_MSG)

			logger.Error("encountered error in handler for route /example attempting to marshal JSON object", marshalErr)
		}
		clientMessage = _clientMessage

	default:
		clientMessage = []byte(METHOD_NOT_SUPPORTED)
		clientStatusCode = http.StatusMethodNotAllowed
	}

	w.WriteHeader(clientStatusCode)
	_, err := w.Write(clientMessage)

	if err != nil {
		logger.Error("encountered error in handler for route /example", err)
	}
}

// server will setup and define the server for handling the RESTful requests
func Server(wg *sync.WaitGroup) *http.Server {
	addr := fmt.Sprintf(":%d", PORT)
	mux := http.NewServeMux()

	logger.Info("port and address", addr)

	// Define the GET / LIST routes for the REST server
	// Mutating HTTP methods are not allowed
	mux.HandleFunc("/example", handleExample)

	server := &http.Server{
		Addr:    addr,
		Handler: mux,
	}

	// The goroutine is required to prevent the server startup
	// from becoming a blocking operation.
	go func() {
		defer wg.Done()

		server.ListenAndServe()
	}()

	return server
}

// Shutdown handles tearing down the active listener
func Shutdown(s *http.Server) error {
	if s == nil {
		return nil
	}

	return s.Shutdown(context.Background())
}

func init() {
	setupLogger()
}
