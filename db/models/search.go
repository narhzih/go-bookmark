package models

const (
	SearchTypePipes    = "pipes"
	SearchTypeTags     = "tags"
	SearchTypePlatform = "platform"
	SearchTypeAll      = "all"
)

type AllSearchResult struct {
	Type   string        `json:"type"` // either a pipe or a bookmark
	Result []interface{} `json:"result"`
}
