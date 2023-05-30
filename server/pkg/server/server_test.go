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
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"sync"
	"testing"

	automationhandler "github.com/nathanbrophy/portfolio-demo/server/pkg/automationHandler"
)

func TestServerStartupAndShutdown(t *testing.T) {
	wg := &sync.WaitGroup{}
	wg.Add(1)
	s := Server(wg)

	if err := Shutdown(s); err != nil {
		t.Fatalf("error shutting down server %v", err)
	}
}

// This serves as an e2e test
func TestHandler(t *testing.T) {
	localAddr := "http://:8081"

	if ok, _ := strconv.ParseBool(os.Getenv("RUN_SERVER_TESTS")); !ok {
		t.Skipf("TestHandler is skipped due to env settings.")
	}

	wg := &sync.WaitGroup{}
	wg.Add(1)
	s := Server(wg)

	type args struct {
		addr   string
		method string
	}
	tests := []struct {
		name        string
		args        args
		wantErr     bool
		wantStatus  int
		wantMarshal bool
	}{
		{
			name: "deafult scenario",
			args: args{
				fmt.Sprintf("%s/example", localAddr),
				http.MethodGet,
			},
			wantErr:     false,
			wantStatus:  http.StatusOK,
			wantMarshal: true,
		},
		{
			name: "404 scenario",
			args: args{
				fmt.Sprintf("%s/example-not-a-route", localAddr),
				http.MethodGet,
			},
			wantErr:     false,
			wantStatus:  http.StatusNotFound,
			wantMarshal: false,
		},
		{
			name: "405 scenario",
			args: args{
				fmt.Sprintf("%s/example", localAddr),
				http.MethodPost,
			},
			wantErr:     false,
			wantStatus:  http.StatusMethodNotAllowed,
			wantMarshal: false,
		},
	}

	for _, tt := range tests {
		fmt.Printf("test case: %v\n", tt.name)
		u, _ := url.Parse(tt.args.addr)
		req := &http.Request{
			Method: tt.args.method,
			URL:    u,
		}
		c := &http.Client{}
		resp, err := c.Do(req)

		if (err != nil) != tt.wantErr {
			t.Fatalf("expected no error to occur %v", err)
			return
		}

		if resp.StatusCode != tt.wantStatus {
			t.Fatalf("mismatching reponse codes want: %d, got: %d", tt.wantStatus, resp.StatusCode)
		}

		d, _ := io.ReadAll(resp.Body)
		container := &automationhandler.AutomateEverythingV1{}

		if err := json.Unmarshal(d, container); err != nil && tt.wantMarshal {
			t.Fatalf("could not reverse the marshalling of the response to a response object: %v", err)
		}
	}

	if err := Shutdown(s); err != nil {
		t.Fatalf("error shutting down server %v", err)
	}
}
