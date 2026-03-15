# Анализ API бэкенда quiz-irr

## Аутентификация

- **POST /api/auth/login** (без токена)  
  - Принимает: `{ "email": string, "password": string }`  
  - Возвращает: `{ full_name, email, access_token, refresh_token }`  
  - В БД: только чтение (проверка пользователя, создание сессии/токенов).

- **POST /api/auth/refresh** (без токена)  
  - Принимает: `{ "refresh_token": string }`  
  - Возвращает: `{ access_token }`  
  - В БД: обновление/проверка refresh-токена.

---

## Админ: тесты (нужен заголовок `Authorization: Bearer <access_token>`)

- **POST /api/tests/init**  
  - Принимает: `{ "title": string, "desc": string }`  
  - Возвращает: `TestAdminResponse` (id, title, desc, is_active, start_at, end_at, author, duration, questions[])  
  - В БД: **создание** записи теста (tests).

- **GET /api/tests?skip=&take=**  
  - Query: `skip`, `take` (числа).  
  - Возвращает: `{ tests: TestAdminResponse[], pagination }`  
  - В БД: **чтение** списка тестов.

- **GET /api/tests/:id/:mode**  
  - `mode`: `preview` или что угодно (full).  
  - Возвращает: `TestAdminResponse` (preview — без вопросов/опций, full — с вопросами и опциями).  
  - В БД: **чтение** теста (и вопросов/опций при full).

- **PATCH /api/tests/:id**  
  - Принимает: `{ title?, desc?, is_active?, duration?, start_at?, end_at? }` (все опционально).  
  - Возвращает: обновлённый `TestAdminResponse`.  
  - В БД: **обновление** теста.

- **DELETE /api/tests/:id**  
  - Возвращает: `"Deleted"`.  
  - В БД: **удаление** теста (и каскадно связанных вопросов/опций и т.д.).

---

## Админ: вопросы (Bearer)

- **POST /api/questions/test/:id**  
  - Body не описан в DTO (можно пустой `{}`).  
  - Возвращает: `QuestionResponse` (id, text, type, points, options[]).  
  - В БД: **создание** вопроса для теста.

- **PATCH /api/questions/:id**  
  - Принимает: `{ "text"?: string, "type"?: string, "points"?: number }`.  
  - Возвращает: обновлённый `QuestionResponse`.  
  - В БД: **обновление** вопроса.

- **DELETE /api/questions/:id**  
  - Возвращает: `"Deleted"`.  
  - В БД: **удаление** вопроса.

---

## Админ: варианты ответов (Bearer)

- **POST /api/options/question/:id**  
  - Body не обязателен.  
  - Возвращает: `OptionResponse` (id, text, is_correct).  
  - В БД: **создание** опции у вопроса.

- **PATCH /api/options/:id**  
  - Принимает: `{ "text"?: string, "is_correct"?: boolean }`.  
  - Возвращает: обновлённый `OptionResponse`.  
  - В БД: **обновление** опции.

- **DELETE /api/options/:id**  
  - Возвращает: `"Deleted"`.  
  - В БД: **удаление** опции.

---

## Прохождение теста (экзамен, без токена)

- **GET /api/exam/info/:testId**  
  - Возвращает: `TestCustomerResponse` (title, desc, duration, start_at, end_at) — без вопросов.  
  - В БД: **чтение** теста.

- **POST /api/exam/start/:testId**  
  - Принимает: `{ "test_id": string, "full_name": string, "email": string, "org": string }`.  
  - Возвращает: `{ "raw_id": string, "questions": ExamQuestion[] }` (вопросы с options без is_correct).  
  - В БД: **создание** записи "сырой" попытки (raw) и привязка к тесту.

- **POST /api/exam/save/:rawId**  
  - Принимает: `{ "answers": [ { "answer_id": number (id вопроса), "option_ids": number[], "text_option": string } ] }`.  
  - Возвращает: `{ "message": string }`.  
  - В БД: **запись/обновление** ответов пользователя по raw_id.

---

## Сырые результаты (Bearer)

- **GET /api/raws/test/:testId?skip=&take=**  
  - Возвращает: `{ data: RawInfoResponse[], pagination }` (id, full_name, email, org, status, start_at, end_at).  
  - В БД: **чтение** списка попыток по тесту.

- **POST /api/raws/analyze/:rawId**  
  - Body не нужен.  
  - Возвращает: `{ "message": string }`.  
  - В БД: **чтение** сырых ответов, сравнение с правильными, **создание** записи результата (results), возможно обновление статуса raw.

- **DELETE /api/raws/answers/:rawId**  
  - Возвращает: `{ "message": string }`.  
  - В БД: **удаление** сырой попытки (и ответов).

- **DELETE /api/raws/test/:testId**  
  - Возвращает: `{ "message": string }`.  
  - В БД: **удаление** всех сырых попыток по тесту.

---

## Результаты и рейтинг (Bearer)

- **GET /api/results/test/:testId?skip=&take=**  
  - Возвращает: `{ data: ResultReponse[], pagination }` (id, full_name, email, org, duration, is_on_time, total_score).  
  - В БД: **чтение** результатов по тесту.

- **DELETE /api/results/result/:resultId**  
  - В БД: **удаление** одного результата.

- **DELETE /api/results/test/:testId**  
  - В БД: **удаление** всех результатов по тесту.

---

## Формат ошибок

Ответ при ошибке: `{ "status": "error", "message": string }` с соответствующим HTTP-кодом (400, 403, 422, 500 и т.д.).
