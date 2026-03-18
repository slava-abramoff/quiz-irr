package dto

type StartExamRequest struct {
	FullName string `json:"full_name"`
	Email    string `json:"email"`
	Org      string `json:"org"`
}

type StartExamResponse struct {
	DataID    string         `json:"raw_id"`
	Questions []ExamQuestion `json:"questions"`
}

type ExamQuestion struct {
	ID      uint         `json:"id"`
	Text    string       `json:"text"`
	Type    string       `json:"type"`
	Options []ExamOption `json:"options"`
}

type ExamOption struct {
	ID   uint   `json:"id"`
	Text string `json:"text"`
}
