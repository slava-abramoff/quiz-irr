# Flow test для API quiz-irr

Скрипт проверяет полный сценарий: логин → создание теста → вопросы и варианты → старт экзамена → отправка ответов → анализ → результаты.

## Требования

- Node.js 18+ (нужен нативный `fetch`)
- Запущенный бэкенд на `http://localhost:8080` (или задать `BASE_URL`)

## Запуск

```bash
cd flow-test
npm run test
```

Или с кастомным URL и учёткой:

```bash
BASE_URL=http://localhost:8080 LOGIN_EMAIL=vyachik005@gmail.com LOGIN_PASSWORD=changeme node flow-test.js
```

## Логи

Каждый запрос и ответ пишутся в **`flow-test-log.txt`** в каталоге, откуда запущен скрипт:

- `>>> REQUEST:` — метод, URL и тело запроса
- `<<< STATUS:` / `<<< RESPONSE:` — код и тело ответа

Файл перезаписывается при каждом запуске.

## Детальное описание API

См. [../docs/API_ENDPOINTS.md](../docs/API_ENDPOINTS.md).
