package dto

type Pagination struct {
	Pagination bool  `json:"pagination"`
	PageNumber int64 `json:"page"`
	PageSize   int64 `json:"page_size"`
}

type PaginationResponse struct {
	Count      int64 `json:"count"`
	PageNumber int64 `json:"page"`
}
