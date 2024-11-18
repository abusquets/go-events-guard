package pagination

type PaginatedResult[T any] struct {
	Items      []T   `json:"items"`
	TotalCount int64 `json:"total"`
}
