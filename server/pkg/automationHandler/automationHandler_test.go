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
	"reflect"
	"testing"
	"time"
)

func TestAutomateEverythingV1_SetTimeNow(t *testing.T) {
	type fields struct {
		Message string
		Time    time.Time
	}
	tests := []struct {
		name   string
		fields fields
	}{
		{
			name: "default scenario",
			fields: fields{
				Message: "test 1",
			},
		},
		{
			name: "override old time scenario",
			fields: fields{
				Message: "test 2",
				Time:    time.Now().Add(1000 * time.Minute),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := &AutomateEverythingV1{
				Message: tt.fields.Message,
				Time:    tt.fields.Time,
			}
			set := a.SetTimeNow()
			if !set.Equal(a.Time) {
				t.Fatalf("the time was not set to now")
			}
		})
	}
}

func TestAutomateEverythingV1_String(t *testing.T) {
	timeTest, _ := time.Parse("2006-01-02T15:04:05.000Z", "2023-05-27T11:00:00.000Z")

	type fields struct {
		Message string
		Time    time.Time
	}
	tests := []struct {
		name    string
		fields  fields
		want    []byte
		wantErr bool
	}{
		{
			name: "default scenario",
			fields: fields{
				"test 1",
				timeTest,
			},
			want:    []byte(`{"message":"test 1","time":"2023-05-27T11:00:00Z"}`),
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := &AutomateEverythingV1{
				Message: tt.fields.Message,
				Time:    tt.fields.Time,
			}
			got, err := a.String()
			if (err != nil) != tt.wantErr {
				t.Errorf("AutomateEverythingV1.String() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("AutomateEverythingV1.String() = %s, want %s", got, tt.want)
			}
		})
	}
}

func TestDefault(t *testing.T) {
	tests := []struct {
		name string
		want AutomateEverything
	}{
		{
			name: "deafult scenario",
			want: &AutomateEverythingV1{
				Message: MESSAGE,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Default(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Default() = %v, want %v", got, tt.want)
			}
		})
	}
}
