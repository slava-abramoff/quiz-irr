package models

type TestResult struct {
	TestID      Test
	FullName    string
	Email       string
	Org         string
	TotalScore  string
	CompletedAt string
}

// test_id (FK -> Tests.id)

// full_name (String) — ФИО.

// email (String) — Почта.

// organization (String) — Учреждение.

// total_score (Int) — Итоговый балл (считается автоматически в конце).

// completed_at (Timestamp)
