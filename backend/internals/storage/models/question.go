package models

type Question struct {
	ID     int
	TestID Test
	Type   string
	Text   string
	Points int
}

// id (UUID, PK)

// test_id (FK -> Tests.id)

// type (Enum) — Тип: single (один), multiple (несколько), text (поле).

// text (Text) — Текст вопроса.

// points (Int) — Вес вопроса (сколько баллов за правильный ответ).

// order_num (Int) — Порядок отображения.
