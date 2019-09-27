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

package db

import (
	"fmt"
	"tryffel.net/pkg/base/repository"
)


func paging(opts *repository.QueryOpts) string {
	return fmt.Sprintf("OFFSET %d LIMIT %d", (opts.Page-1)*opts.Limit, opts.Limit)
}

func sorting(opts *repository.QueryOpts) string {
	return "ORDER BY " + opts.SortField + " " + opts.SortType
}

func sortedPaging(opts *repository.QueryOpts) string {
	return sorting(opts) + " " + paging(opts)
}
