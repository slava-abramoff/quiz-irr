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

type RawUpdateRequest struct {
	FullName *string `json:"full_name"`
	Email    *string `json:"email"`
	Org      *string `json:"org"`
	Status   *string `json:"status"`
	StartAt  *string `json:"start_at"`
	EndAt    *string `json:"end_at"`
}

type RawInfoResponse struct {
	ID       string `json:"id"`
	FullName string `json:"full_name"`
	Email    string `json:"email"`
	Org      string `json:"org"`
	Status   string `json:"status"`
	StartAt  string `json:"start_at"`
	EndAt    string `json:"end_at"`
}

type RawsInfoResponse struct {
	Data       []RawInfoResponse `json:"data"`
	Pagination Pagination        `json:"pagination"`
}
