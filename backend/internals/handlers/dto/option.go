package dto

type CreateOptionRequest struct {
	Text      string `json:"text"`
	IsCorrect bool   `json:"is_correct"`
}

type UpdateOptionRequest struct {
	Text      *string `json:"text"`
	IsCorrect *bool   `json:"is_correct"`
}

type OptionResponse struct {
	ID        uint   `json:"id"`
	Text      string `json:"text"`
	IsCorrect bool   `json:"is_correct"`
}
