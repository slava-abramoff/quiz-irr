package models

import "time"

type Test struct {
	ID        int
	AuthorID  Admin
	Title     string
	Desc      string
	IsActive  bool
	Start     time.Time
	End       time.Time
	CreatedAt time.Time
	UpdatedAt time.Time
}

// id (UUID, PK)

// author_id (FK -> Admins.id)

// title (String) — Название викторины.

// description (Text) — Описание или инструкция.

// is_active (Boolean) — Можно ли сейчас проходить этот тест.

// created_at (Timestamp)
