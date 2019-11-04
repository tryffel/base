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

package migrations

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"time"
)

type Schema struct {
	Level		int			`db:"level"`
	Success		bool		`db:"success"`
	Timestamp	time.Time	`db:"timestamp"`
	TookMs		int 		`db:"took_ms"`
}

// Migrator describes single migration level
type Migrator interface {
	// Get migration name
	MName() string
	// Get migration level
	MLevel() int
	// Get valid sql string to execute
	MSchema() string
}

// Migration implements migrator
type Migration struct {
	Name	string
	Level	int
	Schema	string
}

func (m *Migration) MName() string {
	return m.Name
}

func (m *Migration) MLevel() int {
	return m.Level
}

func (m *Migration) MSchema() string {
	return m.Schema
}

// Migrate runs given migrations
func Migrate(db *sqlx.DB, migrations []Migrator) error {
	current, err := CurrentVersion(db)
	if err != nil {
		return errors.Wrap(err, "failed to get schema version")
	}

	if current.Level == 0 {
		_, err := db.Exec(`
			CREATE TABLE schemas (
			    level integer,
			    success boolean not null default false,
			    timestamp timestamp with time zone default now(),
			    took_ms integer not null,
			    
			    CONSTRAINT schema_pkey PRIMARY KEY (level)
			);`)
		if err != nil {
			return errors.Wrap(err, "failed to create schema table")
		}

	} else {
		if !current.Success {
			return errors.New("previous migration has failed")
		}
	}

	if current.Level == migrations[len(migrations)-1].MLevel() {
		logrus.Debug("No new migrations to run")
		return nil
	}

	for _, v := range migrations[current.Level:len(migrations)-1] {
		err := migrateSingle(db, v)
		if err != nil {
			return errors.Wrap(err, "failed to run migrations")
		}
	}
	return nil
}

// Run single migration
func migrateSingle(db *sqlx.DB, migration Migrator) error {
	start := time.Now()
	_, merr := db.Exec(migration.MSchema())
	
	s := &Schema{
		Level:     migration.MLevel(),
		Success:   merr == nil ,
		Timestamp: time.Now(),
		TookMs: int(time.Since(start).Nanoseconds() / 1000000),
	}

	_, err := db.Exec("INSERT INTO schemas (level, success, timestamp, took_ms) " +
		"VALUES ($1, $2, $3, $4)", s.Level, s.Success, s.Timestamp, s.TookMs)

	if err != nil {
		return errors.Wrap(
			errors.Wrap(
				merr, fmt.Sprintf("migration %d failed", migration.MLevel())),
				"failed to insert new schema")
	}
	return nil
}


// CurrentVersion returns current version
func CurrentVersion(db *sqlx.DB) (Schema, error) {
	current := Schema{}
	err  := db.Get(&current, "SELECT * FROM schemas ORDER BY level DESC LIMIT 1")

	if err != nil {
		if err.Error() == "relation \"schemas\" does not exist" {
			return Schema{
				Level:     0,
				Success:   false,
				Timestamp: time.Time{},
				TookMs:    0,
			}, nil
		}

		return Schema{}, errors.Wrap(err, "failed to query schema")
	}
	return current, nil
}


