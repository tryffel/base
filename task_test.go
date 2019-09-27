// Copyright 2019 Tero Vierimaa
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package base

import (
	"fmt"
	"sync"
	"testing"
	"time"
)

func TestTask_Start(t1 *testing.T) {
	type fields struct {
		Name        string
		lock        sync.RWMutex
		Initialized bool
		Running     bool
		ChanStop    chan bool
		Loop        func()
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{
			// Not initialized
			fields: fields{
				Name: "test-a",
				Initialized: false,
			},
			wantErr: true,
		},
		{
			// No loop
			fields: fields{
				Name: "test-b",
				Initialized: true,
			},
			wantErr: true,
		},
		{
			// Already running
			fields: fields{
				Name: "test-c",
				Initialized: true,
				Running: true,
			},
			wantErr: true,
		},
		{
			// Already running
			fields: fields{
				Name: "test-d",
				Initialized: true,
				Running: false,
				Loop: func() {
					fmt.Println("testing task loop")
					time.Sleep(time.Nanosecond)
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t1.Run(tt.name, func(t1 *testing.T) {
			t := &Task{
				Name:        tt.fields.Name,
				lock:        tt.fields.lock,
				Initialized: tt.fields.Initialized,
				Running:     tt.fields.Running,
				ChanStop:    tt.fields.ChanStop,
				Loop:        tt.fields.Loop,
			}
			if err := t.Start(); (err != nil) != tt.wantErr {
				t1.Errorf("Start() error = %v, wantErr %v", err, tt.wantErr)
			} else {
				time.Sleep(time.Nanosecond*10)
			}
		})
	}
}

func TestTask_Stop(t1 *testing.T) {

	type testLoop struct {
		Task
		testFunc func()
	}

	smallTest := testLoop{}
	eternalTest := testLoop{}

	smallTest.testFunc = func() {
		i := 0
		for i < 3 {
			fmt.Println("testing task loop")
			time.Sleep(time.Nanosecond)
			i += 1
		}
		smallTest.Stop()
	}

	eternalTest.testFunc = func() {
		<-eternalTest.ChanStop
		fmt.Println("Task stopped")
		eternalTest.Stop()
	}

	tests := []struct {
		name    string
		test    testLoop
		loop    func()
		wantErr bool
	}{
		{
			// Small test exits before calling stop
			name:    "test small loop",
			test:    smallTest,
			loop:    smallTest.testFunc,
			wantErr: false,
		},
		{
			name:    "test eternal loop",
			test:    smallTest,
			loop:    smallTest.testFunc,
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t1.Run(tt.name, func(t1 *testing.T) {
			tt.test.Initialized = true
			tt.test.Loop = tt.loop
			err := tt.test.Start()
			if err != nil {
				t1.Errorf("Start() error %v", err.Error())

			}
			time.Sleep(time.Nanosecond * 10)
			if err := tt.test.Stop(); (err != nil) != tt.wantErr {
				t1.Errorf("Stop() (%s) error = %v, wantErr %v", tt.name, err, tt.wantErr)
			}
		})
	}
}
