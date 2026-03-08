package dto

type CreateQuestionRequest struct {
	Text   string `json:"text"`
	Type   string `json:"type"`
	Points int    `json:"points"`
}

type UpdateQuestionRequest struct {
	Text   *string `json:"text"`
	Type   *string `json:"type"`
	Points *int    `json:"points"`
}

type QuestionResponse struct {
	ID      uint             `json:"id"`
	Text    string           `json:"text"`
	Type    string           `json:"type"`
	Points  int              `json:"points"`
	Options []OptionResponse `json:"options"`
}
