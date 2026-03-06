package models

type Option struct {
	ID         int
	QuestionID Question
	Text       string
	IsCorrect  bool
}

// id (UUID, PK)

// question_id (FK -> Questions.id)

// text (String) — Текст варианта.

// is_correct (Boolean) — Является ли этот вариант верным.
