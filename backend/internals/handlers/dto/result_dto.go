package dto

type ResultReponse struct {
	ID         uint   `json:"id"`
	FullName   string `json:"full_name"`
	Email      string `json:"email"`
	Org        string `json:"org"`
	Duration   uint   `json:"duration"`
	IsOnTime   bool   `json:"is_on_time"`
	TotalScore int    `json:"total_score"`
}

type ResultsReponse struct {
	Data       []ResultReponse `json:"data"`
	Pagination Pagination      `json:"pagination"`
}
