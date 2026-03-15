/**
 * Flow test для API quiz-irr backend.
 *
 * Сценарий:
 * 1. Логин
 * 2. Создание теста
 * 3. Вопросы и редактирование
 * 4. Варианты ответов и редактирование
 * 5a. GET exam/info при is_active = false
 * 5b. POST exam/start при is_active = false
 * 5c. PATCH теста: is_active = true
 * 6. GET exam/info при is_active = true
 * 7. POST exam/start (получаем raw_id для сохранения ответов)
 * 8. POST exam/save — отправка ответов
 * 9. GET raws/test, POST raws/analyze, GET results/test
 *
 * При ошибке флоу не останавливается; останавливаемся только когда нет данных
 * для следующего шага — об этом пишем в лог.
 * Все запросы/ответы логируются в flow-test-log.txt
 */

import fs from 'fs';
import path from 'path';

const BASE_URL = process.env.BASE_URL || 'http://localhost:8080';
const LOG_FILE = path.join(process.cwd(), 'flow-test-log.txt');

const LOGIN_EMAIL = process.env.LOGIN_EMAIL || 'vyachik005@gmail.com';
const LOGIN_PASSWORD = process.env.LOGIN_PASSWORD || 'changeme';

function log(line) {
  const ts = new Date().toISOString();
  const s = `[${ts}] ${line}\n`;
  fs.appendFileSync(LOG_FILE, s);
  console.log(line);
}

function logStep(name) {
  const sep = '='.repeat(60);
  log(`${sep}\n${name}\n${sep}`);
}

function logRequest(method, url, body = null) {
  log(`>>> REQUEST: ${method} ${url}`);
  if (body !== undefined && body !== null) {
    log('>>> BODY: ' + JSON.stringify(body, null, 2));
  }
}

function logResponse(status, data) {
  log(`<<< STATUS: ${status}`);
  log('<<< RESPONSE: ' + JSON.stringify(data, null, 2));
}

function logSkip(reason) {
  log('!!! ПРОПУСК: ' + reason);
  log('!!! Для следующего шага нужны данные, которые не были получены.');
}

/** Как request, но при любом ответе не бросает — возвращает { ok, status, data }. */
async function requestSafe(method, url, options = {}) {
  const fullUrl = url.startsWith('http') ? url : `${BASE_URL}${url}`;
  const { body, token } = options;

  const headers = { 'Content-Type': 'application/json' };
  if (token) headers['Authorization'] = `Bearer ${token}`;

  logRequest(method, fullUrl, body);

  try {
    const res = await fetch(fullUrl, {
      method,
      headers,
      body: body ? JSON.stringify(body) : undefined,
    });

    let data;
    const text = await res.text();
    try {
      data = text ? JSON.parse(text) : null;
    } catch {
      data = text;
    }

    logResponse(res.status, data);
    return { ok: res.ok, status: res.status, data };
  } catch (err) {
    logResponse('ERR', { error: err.message });
    return { ok: false, status: 0, data: { error: err.message } };
  }
}

async function run() {
  if (fs.existsSync(LOG_FILE)) fs.unlinkSync(LOG_FILE);
  log('Flow test started. BASE_URL=' + BASE_URL);
  log('');

  let accessToken = null;
  let testId = null;
  let questionIds = [];
  let optionIdsByQuestion = {};
  let rawId = null;
  const optionsSpec = [
    [
      { text: '3', is_correct: false },
      { text: '4', is_correct: true },
      { text: '5', is_correct: false },
    ],
    [
      { text: '2', is_correct: true },
      { text: '3', is_correct: false },
      { text: '4', is_correct: true },
    ],
  ];

  // --- 1. Login ---
  logStep('1. LOGIN');
  const loginRes = await requestSafe('POST', '/api/auth/login', {
    body: { email: LOGIN_EMAIL, password: LOGIN_PASSWORD },
  });
  if (loginRes.ok && loginRes.data && loginRes.data.access_token) {
    accessToken = loginRes.data.access_token;
  } else {
    logSkip('access_token не получен (логин неуспешен). Дальнейшие шаги с авторизацией будут пропущены при отсутствии данных.');
  }
  log('');

  // --- 2. Create test ---
  logStep('2. CREATE TEST');
  const createTestBody = { title: 'Flow Test Quiz', desc: 'Created by flow-test.js' };
  const testRes = await requestSafe('POST', '/api/tests/init', {
    body: createTestBody,
    token: accessToken,
  });
  if (testRes.ok && testRes.data && testRes.data.id) {
    testId = testRes.data.id;
  } else {
    logSkip('testId не получен (создание теста неуспешно). Шаги, требующие testId, будут пропущены.');
  }
  log('');

  // --- 3. Add questions ---
  logStep('3. ADD QUESTIONS');
  const questionsToAdd = [
    { text: 'Сколько будет 2+2?', type: 'single', points: 5 },
    { text: 'Выберите чётные числа', type: 'multiple', points: 10 },
  ];

  if (testId && accessToken) {
    for (const q of questionsToAdd) {
      const added = await requestSafe('POST', `/api/questions/test/${testId}`, { token: accessToken });
      if (!added.ok || !added.data || added.data.id == null) continue;
      questionIds.push(added.data.id);
      await requestSafe('PATCH', `/api/questions/${added.data.id}`, {
        body: { text: q.text, type: q.type, points: q.points },
        token: accessToken,
      });
    }
  } else {
    logSkip('нет testId или accessToken — добавление вопросов невозможно.');
  }
  log('');

  // --- 4. Add options ---
  logStep('4. ADD OPTIONS');
  if (questionIds.length && accessToken) {
    optionIdsByQuestion = {};
    for (let i = 0; i < questionIds.length; i++) {
      const qId = questionIds[i];
      optionIdsByQuestion[qId] = [];
      for (const opt of optionsSpec[i] || []) {
        const added = await requestSafe('POST', `/api/options/question/${qId}`, {
          token: accessToken,
        });
        if (!added.ok || !added.data || added.data.id == null) continue;
        await requestSafe('PATCH', `/api/options/${added.data.id}`, {
          body: { text: opt.text, is_correct: opt.is_correct },
          token: accessToken,
        });
        optionIdsByQuestion[qId].push(added.data.id);
      }
    }
  } else {
    logSkip('нет questionIds или accessToken — добавление опций невозможно.');
  }
  log('');

  // --- 5a. Exam info при is_active = false ---
  logStep('5a. EXAM INFO (is_active = false)');
  if (testId) {
    await requestSafe('GET', `/api/exam/info/${testId}`);
  } else {
    logSkip('нет testId — запрос GET exam/info невозможен.');
  }
  log('');

  // --- 5b. Start exam при is_active = false ---
  logStep('5b. START EXAM (is_active = false)');
  const startBody = {
    test_id: testId,
    full_name: 'Иван Тестов',
    email: 'ivan@test.local',
    org: 'Flow Test Org',
  };
  if (testId) {
    const startResFalse = await requestSafe('POST', `/api/exam/start/${testId}`, {
      body: startBody,
    });
    if (startResFalse.ok && startResFalse.data && startResFalse.data.raw_id) {
      rawId = startResFalse.data.raw_id;
    }
  } else {
    logSkip('нет testId — старт экзамена невозможен.');
  }
  log('');

  // --- 5c. Включаем тест: is_active = true ---
  logStep('5c. PATCH TEST is_active = true');
  if (testId && accessToken) {
    await requestSafe('PATCH', `/api/tests/${testId}`, {
      body: { is_active: true },
      token: accessToken,
    });
  } else {
    logSkip('нет testId или accessToken — включение теста невозможно.');
  }
  log('');

  // --- 6. Exam info при is_active = true ---
  logStep('6. EXAM INFO (is_active = true)');
  if (testId) {
    await requestSafe('GET', `/api/exam/info/${testId}`);
  } else {
    logSkip('нет testId.');
  }
  log('');

  // --- 7. Start exam (для raw_id, если ещё нет) ---
  logStep('7. START EXAM (is_active = true)');
  if (testId) {
    const startResTrue = await requestSafe('POST', `/api/exam/start/${testId}`, {
      body: startBody,
    });
    if (startResTrue.ok && startResTrue.data && startResTrue.data.raw_id) {
      rawId = startResTrue.data.raw_id;
    }
  } else {
    logSkip('нет testId.');
  }
  log('');

  // --- 8. Save answers ---
  logStep('8. SAVE ANSWERS');
  if (rawId && questionIds.length >= 2 && optionIdsByQuestion[questionIds[0]]?.length && optionIdsByQuestion[questionIds[1]]?.length) {
    const q1OptionIds = optionIdsByQuestion[questionIds[0]];
    const q2OptionIds = optionIdsByQuestion[questionIds[1]];
    const q1CorrectOptionId = q1OptionIds[optionsSpec[0].findIndex((o) => o.is_correct)];
    const q2CorrectOptionIds = q2OptionIds.filter((_, idx) => optionsSpec[1][idx].is_correct);
    const answersBody = {
      answers: [
        { answer_id: questionIds[0], option_ids: [q1CorrectOptionId], text_option: '' },
        { answer_id: questionIds[1], option_ids: q2CorrectOptionIds, text_option: '' },
      ],
    };
    await requestSafe('POST', `/api/exam/save/${rawId}`, { body: answersBody });
  } else {
    logSkip('нет rawId (старт экзамена не вернул raw_id) или нет questionIds/optionIds — сохранение ответов невозможно.');
  }
  log('');

  // --- 9. Raw results list ---
  logStep('9. RAW RESULTS (list by test)');
  if (testId && accessToken) {
    await requestSafe('GET', `/api/raws/test/${testId}?skip=0&take=10`, {
      token: accessToken,
    });
  } else {
    logSkip('нет testId или accessToken.');
  }
  log('');

  // --- 10. Analyze raw ---
  logStep('10. ANALYZE RAW');
  if (rawId && accessToken) {
    await requestSafe('POST', `/api/raws/analyze/${rawId}`, {
      token: accessToken,
    });
  } else {
    logSkip('нет rawId или accessToken — анализ сырых ответов невозможен.');
  }
  log('');

  // --- 11. Results by test ---
  logStep('11. RESULTS BY TEST');
  if (testId && accessToken) {
    const resultsRes = await requestSafe('GET', `/api/results/test/${testId}?skip=0&take=10`, {
      token: accessToken,
    });
    if (resultsRes.ok && resultsRes.data?.data?.length) {
      log('Result IDs: ' + resultsRes.data.data.map((r) => r.id).join(', '));
    }
  } else {
    logSkip('нет testId или accessToken.');
  }
  log('');

  logStep('DONE');
  log(`Log written to: ${LOG_FILE}`);
}

run();
