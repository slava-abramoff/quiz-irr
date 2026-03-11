package dto

type SendUserAnswersRequest struct {
	Answers []Answer `json:"answers"`
}

type SendUserAnswersResponse struct {
	Message string `json:"message"`
}

type Answer struct {
	ID        uint   `json:"answer_id"`
	OptionIDs []uint `json:"option_ids"`
	Text      string `json:"text_option"`
}
