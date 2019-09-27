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
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"sync"
)

// Tasker can be run on background
type Tasker interface {
	Start() error
	Stop() error
}

// Common fields for task
type Task struct {
	// Name of the task, for logging purposes
	Name        string
	lock        sync.RWMutex
	// Initialized flag must be true in order to run the task
	Initialized bool
	Running     bool
	ChanStop    chan bool
	Loop        func()
}

func (t *Task) Start() error {
	t.lock.Lock()
	defer t.lock.Unlock()

	if t.Running {
		return errors.New("background task already running")
	}

	if !t.Initialized {
		return errors.New("task not initialized properly")
	}

	if t.Loop == nil {
		return errors.New("no loop function defined")
	}

	if t.ChanStop == nil {
		t.init()
	}


	logrus.Info("Starting task ", t.Name)
	t.Running = true
	go t.Loop()
	return nil
}

func (t *Task) Stop() error {
	t.lock.Lock()
	defer t.lock.Unlock()

	if !t.Running {
		return errors.New("background task not running")
	}

	logrus.Info("Stopping task ", t.Name)
	t.ChanStop <- true
	t.Running = false
	return nil
}

func (t *Task) init() {
	t.ChanStop = make(chan bool, 2)
}

