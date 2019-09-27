package repository

// QueryOpts includes options to pass to query
type QueryOpts struct {
	// Page size
	Limit           int
	// Page
	Page            int
	// Sorting
	SortField       string
	// Sorting type: ASC|DESC
	SortType        string
	// Optional transaction to use
	Transaction     interface{}
	// Whether to retrieve total records count
	GetTotalRecords bool
}

// QueryResult includes common variables for result
type QueryResult struct {
	// Possible error in query
	Error        error
	// Total records found with given parameters
	TotalRecords int
}
