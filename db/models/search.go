package models

const (
	SearchTypePipes = "pipes"
	SearchTypeTags  = "tags"
	SearchTypeAll   = "all"
)

type AllSearchResult struct {
	Type   string        `json:"type"` // either a pipe or a bookmark
	Result []interface{} `json:"result"`
}
