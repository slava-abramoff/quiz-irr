package dto

type SendUserAnswersRequest struct {
	Answers []Answer `json:"answers"`
}

type Answer struct {
	ID       uint   `json:"answer_id"`
	OptionID uint   `json:"option_id"`
	Text     string `json:"text_option"`
}
