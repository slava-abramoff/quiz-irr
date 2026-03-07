package dto

type CreateTestRequest struct {
	Title string `json:"title"`
	Desc  string `json:"desc"`
}

type UpdateTestRequest struct {
	Title    string `json:"title"`
	Desc     string `json:"desc"`
	IsActive bool   `json:"is_active"`
	StartAt  string `json:"start_at"`
	EndAt    string `json:"end_at"`
}

type TestAdminResponse struct {
	ID        string             `json:"id"`
	Title     string             `json:"title"`
	Desc      string             `json:"desc"`
	IsActive  bool               `json:"is_active"`
	StartAt   string             `json:"start_at"`
	EndAt     string             `json:"end_at"`
	Questions []QuestionResponse `json:"questions"`
}
