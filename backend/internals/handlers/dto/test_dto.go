package dto

type CreateTestRequest struct {
	Title string `json:"title"`
	Desc  string `json:"desc"`
}

type UpdateTestRequest struct {
	Title    *string `json:"title"`
	Desc     *string `json:"desc"`
	IsActive *bool   `json:"is_active"`
	Duration *uint   `json:"duration"`
	StartAt  *string `json:"start_at"`
	EndAt    *string `json:"end_at"`
}

type GetManyTestsResponse struct {
	Tests      []TestAdminResponse `json:"tests"`
	Pagination Pagination          `json:"pagination"`
}

type TestAdminResponse struct {
	ID        string             `json:"id"`
	Title     string             `json:"title"`
	Desc      string             `json:"desc"`
	IsActive  bool               `json:"is_active"`
	StartAt   string             `json:"start_at"`
	EndAt     string             `json:"end_at"`
	Author    string             `json:"author"`
	Duration  uint               `json:"duration"`
	Questions []QuestionResponse `json:"questions"`
}

type TestCustomerResponse struct {
	Title    string `json:"title"`
	Desc     string `json:"desc"`
	Duration uint   `json:"duration"`
	StartAt  string `json:"start_at"`
	EndAt    string `json:"end_at"`
}
