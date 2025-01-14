package schemas

type PaginationRequest struct {
	Page  int   `json:"page"`
	Limit int64 `json:"limit"`
}

type PaginationResponse struct {
	Page      int   `json:"current_page"`
	Limit     int64 `json:"limit_per_page"`
	Total     int   `json:"total_data"`
	TotalPage int   `json:"total_page"`
}

type PaginationData[T any] struct {
	Data []T `json:"-"`
}
