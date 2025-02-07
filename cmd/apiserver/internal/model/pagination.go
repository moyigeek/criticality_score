package model

type PageDTO[T any] struct {
	Total   int `json:"total"`
	Current int `json:"start"`
	Take    int `json:"count"`
	Items   []T `json:"items"`
}

func NewPageDTO[T any](total int, current int, items []T) *PageDTO[T] {
	return &PageDTO[T]{
		Total:   total,
		Current: current,
		Take:    len(items),
		Items:   items,
	}
}
