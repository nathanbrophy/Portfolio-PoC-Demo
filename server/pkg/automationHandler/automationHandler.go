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
package automationhandler

import (
	"encoding/json"
	"time"
)

const (
	NIL     string = "<nil>"
	MESSAGE string = "Automate all the things again!"
)

// AutomateEverything is the type the server exposes as a RESTful accessor
type AutomateEverything interface {
	// SetTimeNow will set the time to the Epoch moment the function is called
	SetTimeNow() time.Time

	// String will attempt to return a JSON representation of the object
	String() ([]byte, error)
}

// AutomateEverythingV1 implements the AutomateEverything interface as a V1 version
type AutomateEverythingV1 struct {
	// Message is the message in the JSON string to report
	Message string `json:"message"`

	// Time is the time object to report in the JSON output
	//+default time.Now epoch value
	Time time.Time `json:"time"`
}

// SetTimeNow will set the time to the Epoch moment the function is called
func (a *AutomateEverythingV1) SetTimeNow() time.Time {
	if a == nil {
		a = &AutomateEverythingV1{
			Message: MESSAGE,
		}
	}
	now := time.Now()
	a.Time = now
	return now
}

// String will attempt to return a JSON representation of the object
func (a *AutomateEverythingV1) String() ([]byte, error) {
	if a == nil {
		return []byte(NIL), nil
	}

	return json.Marshal(a)
}

// Default returns a default server responder to be used in the server implementation
func Default() AutomateEverything {
	return &AutomateEverythingV1{
		Message: MESSAGE,
	}
}
