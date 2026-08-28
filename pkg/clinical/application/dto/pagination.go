package dto

import "fmt"

const defaultPageCount = 10

// PageInfo is used to add pagination information to Relay edges.
type PageInfo struct {
	// Forward pagination
	HasNextPage bool
	EndCursor   *string

	// Backward pagination
	HasPreviousPage bool
	StartCursor     *string

	BundleID *string `json:"bundleID,omitempty"`
	NextURL  *string `json:"nextURL,omitempty"`
	SelfURL  *string `json:"selfURL,omitempty"`
}

// Pagination represents paging parameters
type Pagination struct {
	// Forward pagination arguments
	First *int   `json:"first"`
	After string `json:"after"`

	// Backward pagination arguments
	Last   *int   `json:"last"`
	Before string `json:"before"`

	// A flag to indicate whether to ignore the pagination parameters
	// i.e the implementer will not perform/do pagination if this flag is true
	Skip bool `json:"skip"`
}

func (p *Pagination) Validate() error {
	if p.First != nil && p.Last != nil {
		return fmt.Errorf("cannot provide both first and last")
	}

	if p.First != nil {
		first := *p.First
		if first <= 0 {
			return fmt.Errorf("first cannot be less than 0")
		}
	}

	if p.Last != nil {
		last := *p.Last
		if last <= 0 {
			return fmt.Errorf("last cannot be less than 0")
		}
	}

	if p.First == nil && p.Last == nil {
		limit := defaultPageCount
		p.First = &limit
	}

	return nil
}
