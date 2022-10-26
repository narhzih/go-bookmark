package models

type Filter struct {
	Limit int64
	Page  int64
}

func (f *Filter) Offset() int64 {
	return (f.Page - 1) * f.Page
}
