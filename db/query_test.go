package db

import (
	"testing"
	"tryffel.net/pkg/base/repository"
)

func TestPaging(t *testing.T) {
	tests := []struct {
		name string
		opts *repository.QueryOpts
		want string
	}{
		{
			opts: &repository.QueryOpts{
				Limit:           10,
				Page:            1,
				SortField:       "",
				SortType:        "",
				Transaction:     nil,
				GetTotalRecords: false,
			},
			want: "OFFSET 0 LIMIT 10",
		},
		{
			opts: &repository.QueryOpts{
				Limit:           5,
				Page:            1,
				SortField:       "",
				SortType:        "",
				Transaction:     nil,
				GetTotalRecords: false,
			},
			want: "OFFSET 0 LIMIT 5",
		},
		{
			opts: &repository.QueryOpts{
				Limit:           10,
				Page:            4,
				SortField:       "",
				SortType:        "",
				Transaction:     nil,
				GetTotalRecords: false,
			},
			want: "OFFSET 30 LIMIT 10",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := paging(tt.opts); got != tt.want {
				t.Errorf("Paging() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSorting(t *testing.T) {
	tests := []struct {
		name string
		args *repository.QueryOpts
		want string
	}{
		{
			args: &repository.QueryOpts{
				Limit:           0,
				Page:            0,
				SortField:       "test",
				SortType:        "ASC",
				Transaction:     nil,
				GetTotalRecords: false,
			},
			want: "ORDER BY test ASC",
		},
		{
			args: &repository.QueryOpts{
				Limit:           0,
				Page:            0,
				SortField:       "test-b",
				SortType:        "DESC",
				Transaction:     nil,
				GetTotalRecords: false,
			},
			want: "ORDER BY test-b DESC",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := sorting(tt.args); got != tt.want {
				t.Errorf("Sorting() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_sortedPaging(t *testing.T) {
	tests := []struct {
		name string
		args *repository.QueryOpts
		want string
	}{
		{
			args: &repository.QueryOpts{
				Limit:           10,
				Page:            5,
				SortField:       "test",
				SortType:        "ASC",
				Transaction:     nil,
				GetTotalRecords: false,
			},
			want: "ORDER BY test ASC OFFSET 40 LIMIT 10",
		},
		{
			args: &repository.QueryOpts{
				Limit:           4,
				Page:            5,
				SortField:       "test-123",
				SortType:        "ASC",
				Transaction:     nil,
				GetTotalRecords: false,
			},
			want: "ORDER BY test-123 ASC OFFSET 16 LIMIT 4",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := sortedPaging(tt.args); got != tt.want {
				t.Errorf("sortedPaging() = %v, want %v", got, tt.want)
			}
		})
	}
}
