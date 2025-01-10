package schemas

type PaginationRequest struct {
	Page  int   `json:"page"`
	Limit int64 `json:"limit"`
}

type PaginationResponse[T any] struct {
	Page      int   `json:"current_page"`
	Limit     int64 `json:"limit_per_page"`
	Total     int   `json:"total_data"`
	TotalPage int   `json:"total_page"`
	Data      []T   `json:"data"`
}
