package dto

type Pagination struct {
	CurrentPage     uint `json:"current_page"`
	TotalPages      uint `json:"total_pages"`
	TotalItems      uint `json:"total_items"`
	ItemsPerPage    uint `json:"items_perpage"`
	HasNextPage     bool `json:"has_next_page"`
	HasPreviousPage bool `json:"has_previous_page"`
}
