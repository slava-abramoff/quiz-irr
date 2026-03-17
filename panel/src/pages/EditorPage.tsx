import { useEffect, useMemo, useState } from 'react';
import type { ChangeEvent, FormEvent } from 'react';
import { useNavigate, useParams } from 'react-router-dom';
import { getTest, updateTest } from '../api/tests';
import {
  createQuestion,
  deleteQuestion,
  updateQuestion,
} from '../api/questions';
import {
  createOption,
  deleteOption,
  updateOption,
} from '../api/options';
import type {
  OptionResponse,
  QuestionResponse,
  TestAdminResponse,
  UpdateQuestionRequest,
  UpdateOptionRequest,
  UpdateTestRequest,
} from '../api/types';

type LocalQuestion = QuestionResponse;
type LocalOption = OptionResponse;

function parseDurationToHMS(totalSeconds: number | undefined | null) {
  if (!totalSeconds || totalSeconds <= 0) {
    return { hours: 0, minutes: 0, seconds: 0 };
  }
  const hours = Math.floor(totalSeconds / 3600);
  const minutes = Math.floor((totalSeconds % 3600) / 60);
  const seconds = totalSeconds % 60;
  return { hours, minutes, seconds };
}

function hmsToDurationSeconds(
  hours: number | undefined,
  minutes: number | undefined,
  seconds: number | undefined,
): number {
  const h = Number.isFinite(hours ?? 0) ? Number(hours ?? 0) : 0;
  const m = Number.isFinite(minutes ?? 0) ? Number(minutes ?? 0) : 0;
  const s = Number.isFinite(seconds ?? 0) ? Number(seconds ?? 0) : 0;
  return h * 3600 + m * 60 + s;
}

function backendDateTimeToInput(value: string | null | undefined): string {
  if (!value) return '';

  // Если уже ISO-формат "YYYY-MM-DDTHH:MM[:SS]Z..." — просто приводим к datetime-local
  if (value.includes('T')) {
    const iso = value;
    const date = new Date(iso);
    if (Number.isNaN(date.getTime())) return '';
    const year = date.getFullYear();
    const month = String(date.getMonth() + 1).padStart(2, '0');
    const day = String(date.getDate()).padStart(2, '0');
    const hours = String(date.getHours()).padStart(2, '0');
    const minutes = String(date.getMinutes()).padStart(2, '0');
    return `${year}-${month}-${day}T${hours}:${minutes}`;
  }

  // Старый формат: "YYYY-MM-DD HH:MM:SS"
  const [datePart, timePart] = value.split(' ');
  if (!datePart || !timePart) return '';
  const [hh, mm] = timePart.split(':');
  if (!hh || !mm) return '';
  return `${datePart}T${hh}:${mm}`;
}

function inputDateTimeToBackend(value: string | null | undefined): string | undefined {
  if (!value) return undefined;
  // value в формате "YYYY-MM-DDTHH:MM" (локальное время)
  const [datePart, timePart] = value.split('T');
  if (!datePart || !timePart) return undefined;

  // Вычисляем текущий часовой пояс клиента
  const now = new Date();
  const offsetMinutes = -now.getTimezoneOffset(); // например, +180 для UTC+3
  const sign = offsetMinutes >= 0 ? '+' : '-';
  const absMinutes = Math.abs(offsetMinutes);
  const offsetHoursPart = String(Math.floor(absMinutes / 60)).padStart(2, '0');
  const offsetMinutesPart = String(absMinutes % 60).padStart(2, '0');

  // Собираем ISO‑подобную строку с явным часовым поясом: "YYYY-MM-DDTHH:MM:00+03:00"
  return `${datePart}T${timePart}:00${sign}${offsetHoursPart}:${offsetMinutesPart}`;
}

export default function EditorPage() {
  const { id } = useParams<{ id: string }>();
  const navigate = useNavigate();
  const [test, setTest] = useState<TestAdminResponse | null>(null);
  const [isLoading, setIsLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [isSavingMeta, setIsSavingMeta] = useState(false);

  const [title, setTitle] = useState('');
  const [desc, setDesc] = useState('');
  const [startAt, setStartAt] = useState('');
  const [endAt, setEndAt] = useState('');
  const [isActive, setIsActive] = useState(false);
  const [durationHours, setDurationHours] = useState(0);
  const [durationMinutes, setDurationMinutes] = useState(0);
  const [durationSeconds, setDurationSeconds] = useState(0);

  const isAuthenticated =
    typeof window !== 'undefined' && !!localStorage.getItem('access_token');

  useEffect(() => {
    if (!isAuthenticated) {
      navigate('/login', { replace: true });
      return;
    }
    if (!id) return;

    async function loadTest() {
      try {
        setIsLoading(true);
        setError(null);
        // mode "fulldata" – сервер вернёт полные данные теста
        const data = await getTest(id, 'fulldata');

        const normalized: TestAdminResponse = {
          ...data,
          questions: (data.questions ?? []).map((q) => ({
            ...q,
            options: q.options ?? [],
          })),
        };

        setTest(normalized);

        setTitle(data.title);
        setDesc(data.desc);
        setStartAt(backendDateTimeToInput(data.start_at));
        setEndAt(backendDateTimeToInput(data.end_at));
        setIsActive(data.is_active);
        const { hours, minutes, seconds } = parseDurationToHMS(data.duration);
        setDurationHours(hours);
        setDurationMinutes(minutes);
        setDurationSeconds(seconds);
      } catch (err: any) {
        const message = err?.response?.data?.message ?? 'Не удалось загрузить тест.';
        setError(message);
      } finally {
        setIsLoading(false);
      }
    }

    void loadTest();
  }, [id, isAuthenticated, navigate]);

  const questions: LocalQuestion[] = useMemo(
    () => test?.questions ?? [],
    [test?.questions],
  );

  const handleMetaSubmit = async (event: FormEvent) => {
    event.preventDefault();
    if (!id) return;

    const payload: UpdateTestRequest = {
      title: title.trim() || undefined,
      desc: desc.trim() || undefined,
      is_active: isActive,
      start_at: inputDateTimeToBackend(startAt),
      end_at: inputDateTimeToBackend(endAt),
      duration: hmsToDurationSeconds(durationHours, durationMinutes, durationSeconds),
    };

    try {
      setIsSavingMeta(true);
      const updated = await updateTest(id, payload);
      setTest((prev) =>
        prev
          ? {
              ...prev,
              // обновляем только мета-поля, не трогаем вопросы и ответы
              title: updated.title,
              desc: updated.desc,
              is_active: updated.is_active,
              start_at: updated.start_at,
              end_at: updated.end_at,
              duration: updated.duration,
            }
          : prev,
      );
    } catch (err: any) {
      const message = err?.response?.data?.message ?? 'Не удалось сохранить изменения теста.';
      setError(message);
    } finally {
      setIsSavingMeta(false);
    }
  };

  const handleAddQuestion = async () => {
    if (!id) return;
    try {
      const created = await createQuestion(id, {
        text: 'Новый вопрос',
        type: 'single',
        points: 1,
      });
      setTest((prev) =>
        prev
          ? {
              ...prev,
              questions: [...prev.questions, created],
            }
          : prev,
      );
    } catch (err: any) {
      const message = err?.response?.data?.message ?? 'Не удалось добавить вопрос.';
      setError(message);
    }
  };

  const handleUpdateQuestion = async (
    question: LocalQuestion,
    partial: UpdateQuestionRequest,
  ) => {
    try {
      const updated = await updateQuestion(question.id, partial);
      const normalized: LocalQuestion = {
        ...updated,
        options: updated.options ?? [],
      };
      setTest((prev) =>
        prev
          ? {
              ...prev,
              questions: prev.questions.map((q) => (q.id === question.id ? normalized : q)),
            }
          : prev,
      );
    } catch (err: any) {
      const message = err?.response?.data?.message ?? 'Не удалось обновить вопрос.';
      setError(message);
    }
  };

  const handleDeleteQuestion = async (question: LocalQuestion) => {
    if (!window.confirm('Удалить этот вопрос и все его варианты ответов?')) {
      return;
    }
    try {
      await deleteQuestion(question.id);
      setTest((prev) =>
        prev
          ? {
              ...prev,
              questions: prev.questions.filter((q) => q.id !== question.id),
            }
          : prev,
      );
    } catch (err: any) {
      const message = err?.response?.data?.message ?? 'Не удалось удалить вопрос.';
      setError(message);
    }
  };

  const handleAddOption = async (question: LocalQuestion) => {
    try {
      // Ограничения по типу: у text только один ответ
      if (question.type === 'text' && (question.options?.length ?? 0) >= 1) {
        setError('Для текстового вопроса можно добавить только один вариант ответа.');
        return;
      }

      const created = await createOption(question.id, {
        text: 'Новый вариант',
        is_correct: question.type === 'text',
      });
      setTest((prev) =>
        prev
          ? {
              ...prev,
              questions: prev.questions.map((q) =>
                q.id === question.id
                  ? { ...q, options: [...(q.options ?? []), created] }
                  : q,
              ),
            }
          : prev,
      );
    } catch (err: any) {
      const message = err?.response?.data?.message ?? 'Не удалось добавить вариант.';
      setError(message);
    }
  };

  const handleUpdateOption = async (
    question: LocalQuestion,
    option: LocalOption,
    partial: UpdateOptionRequest,
  ) => {
    try {
      // Валидация: у single не может быть нескольких правильных
      if (
        question.type === 'single' &&
        partial.is_correct === true &&
        (question.options ?? []).some((o) => o.id !== option.id && o.is_correct)
      ) {
        setError('У вопроса с типом single может быть только один правильный ответ.');
        return;
      }

      // Валидация: у text максимум один вариант
      if (question.type === 'text' && partial.is_correct === false) {
        setError('Для текстового вопроса вариант должен быть единственным и правильным.');
        return;
      }

      const updated = await updateOption(option.id, partial);
      setTest((prev) =>
        prev
          ? {
              ...prev,
              questions: prev.questions.map((q) =>
                q.id === question.id
                  ? {
                      ...q,
                      options: (q.options ?? []).map((o) =>
                        o.id === option.id ? updated : o,
                      ),
                    }
                  : q,
              ),
            }
          : prev,
      );
    } catch (err: any) {
      const message = err?.response?.data?.message ?? 'Не удалось обновить вариант.';
      setError(message);
    }
  };

  const handleDeleteOption = async (question: LocalQuestion, option: LocalOption) => {
    if (!window.confirm('Удалить этот вариант ответа?')) {
      return;
    }
    try {
      await deleteOption(option.id);
      setTest((prev) =>
        prev
          ? {
              ...prev,
              questions: prev.questions.map((q) =>
                q.id === question.id
                  ? {
                      ...q,
                      options: (q.options ?? []).filter((o) => o.id !== option.id),
                    }
                  : q,
              ),
            }
          : prev,
      );
    } catch (err: any) {
      const message = err?.response?.data?.message ?? 'Не удалось удалить вариант.';
      setError(message);
    }
  };

  const handleDurationChange =
    (field: 'h' | 'm' | 's') =>
    (event: ChangeEvent<HTMLInputElement>) => {
      const value = Number.parseInt(event.target.value || '0', 10);
      const safe = Number.isNaN(value) ? 0 : Math.max(0, value);
      if (field === 'h') setDurationHours(safe);
      if (field === 'm') setDurationMinutes(Math.min(59, safe));
      if (field === 's') setDurationSeconds(Math.min(59, safe));
    };

  return (
    <div className="min-h-screen bg-gray-50">
      <header className="border-b border-gray-200 bg-white">
        <div className="mx-auto max-w-6xl px-4 py-4 flex items-center justify-between gap-4">
          <div>
            <button
              type="button"
              onClick={() => navigate('/')}
              className="text-left"
            >
              <h1 className="text-lg font-semibold text-gray-900">Конструктор тестов</h1>
              <p className="text-sm text-gray-500">Редактор теста</p>
            </button>
          </div>

          <div className="flex items-center gap-3">
            {/* здесь позже можно добавить элементы управления пользователем */}
          </div>
        </div>
      </header>

      <main className="mx-auto max-w-6xl px-4 py-6 space-y-6">
        {error && (
          <div className="rounded-md bg-red-50 px-3 py-2 text-sm text-red-700 border border-red-100">
            {error}
          </div>
        )}

        {isLoading && !test ? (
          <div className="rounded-lg border border-gray-200 bg-white p-8 text-center text-sm text-gray-500">
            Загружаем данные теста...
          </div>
        ) : !test ? (
          <div className="rounded-lg border border-gray-200 bg-white p-8 text-center text-sm text-gray-500">
            Тест не найден.
          </div>
        ) : (
          <>
            {/* Блок 1: информация о тесте */}
            <section className="rounded-lg border border-gray-200 bg-white p-5 space-y-4">
              <h2 className="text-sm font-semibold text-gray-900">Общие настройки теста</h2>

              <form onSubmit={handleMetaSubmit} className="space-y-4">
                <div className="grid grid-cols-1 gap-4 md:grid-cols-2">
                  <div className="space-y-1.5">
                    <label
                      htmlFor="test-title"
                      className="block text-xs font-medium text-gray-700"
                    >
                      Название теста
                    </label>
                    <input
                      id="test-title"
                      type="text"
                      value={title}
                      onChange={(event) => setTitle(event.target.value)}
                      className="w-full rounded-md border border-gray-300 px-3 py-2 text-sm text-gray-900 placeholder-gray-400 focus:outline-none focus:ring-2 focus:ring-gray-900 focus:border-transparent"
                    />
                  </div>

                  <div className="space-y-1.5">
                    <span className="block text-xs font-medium text-gray-700">
                      Статус
                    </span>
                    <button
                      type="button"
                      onClick={() => setIsActive((prev) => !prev)}
                      className={`inline-flex items-center rounded-full px-3 py-1 text-xs font-medium ${
                        isActive
                          ? 'bg-emerald-50 text-emerald-700'
                          : 'bg-gray-100 text-gray-500'
                      }`}
                    >
                      {isActive ? 'Активен' : 'Неактивен'}
                    </button>
                  </div>
                </div>

                <div className="space-y-1.5">
                  <label
                    htmlFor="test-desc"
                    className="block text-xs font-medium text-gray-700"
                  >
                    Описание
                  </label>
                  <textarea
                    id="test-desc"
                    rows={4}
                    value={desc}
                    onChange={(event) => setDesc(event.target.value)}
                    className="w-full rounded-md border border-gray-300 px-3 py-2 text-sm text-gray-900 placeholder-gray-400 focus:outline-none focus:ring-2 focus:ring-gray-900 focus:border-transparent resize-none"
                  />
                </div>

                <div className="grid grid-cols-1 gap-4 md:grid-cols-3">
                  <div className="space-y-1.5">
                    <label
                      htmlFor="test-start-at"
                      className="block text-xs font-medium text-gray-700"
                    >
                      Начало тестирования
                    </label>
                    <input
                      id="test-start-at"
                      type="datetime-local"
                      value={startAt}
                      onChange={(event) => setStartAt(event.target.value)}
                      className="w-full rounded-md border border-gray-300 px-3 py-2 text-sm text-gray-900 placeholder-gray-400 focus:outline-none focus:ring-2 focus:ring-gray-900 focus:border-transparent"
                    />
                  </div>

                  <div className="space-y-1.5">
                    <label
                      htmlFor="test-end-at"
                      className="block text-xs font-medium text-gray-700"
                    >
                      Окончание тестирования
                    </label>
                    <input
                      id="test-end-at"
                      type="datetime-local"
                      value={endAt}
                      onChange={(event) => setEndAt(event.target.value)}
                      className="w-full rounded-md border border-gray-300 px-3 py-2 text-sm text-gray-900 placeholder-gray-400 focus:outline-none focus:ring-2 focus:ring-gray-900 focus:border-transparent"
                    />
                  </div>

                  <div className="space-y-1.5">
                    <span className="block text-xs font-medium text-gray-700">
                      Длительность теста
                    </span>
                    <div className="flex items-center gap-2">
                      <input
                        type="number"
                        min={0}
                        value={durationHours}
                        onChange={handleDurationChange('h')}
                        className="w-16 rounded-md border border-gray-300 px-2 py-1 text-xs text-gray-900 focus:outline-none focus:ring-2 focus:ring-gray-900 focus:border-transparent"
                      />
                      <span className="text-xs text-gray-500">ч</span>
                      <input
                        type="number"
                        min={0}
                        max={59}
                        value={durationMinutes}
                        onChange={handleDurationChange('m')}
                        className="w-16 rounded-md border border-gray-300 px-2 py-1 text-xs text-gray-900 focus:outline-none focus:ring-2 focus:ring-gray-900 focus:border-transparent"
                      />
                      <span className="text-xs text-gray-500">мин</span>
                      <input
                        type="number"
                        min={0}
                        max={59}
                        value={durationSeconds}
                        onChange={handleDurationChange('s')}
                        className="w-16 rounded-md border border-gray-300 px-2 py-1 text-xs text-gray-900 focus:outline-none focus:ring-2 focus:ring-gray-900 focus:border-transparent"
                      />
                      <span className="text-xs text-gray-500">сек</span>
                    </div>
                  </div>
                </div>

                <div className="flex items-center justify-end gap-2 pt-2">
                  <button
                    type="submit"
                    disabled={isSavingMeta}
                    className="inline-flex items-center rounded-md bg-gray-900 px-4 py-2 text-xs font-medium text-white shadow-sm hover:bg-gray-800 disabled:opacity-60 disabled:cursor-not-allowed"
                  >
                    {isSavingMeta ? 'Сохраняем...' : 'Сохранить изменения'}
                  </button>
                </div>
              </form>
            </section>

            {/* Блок 2: вопросы и ответы */}
            <section className="rounded-lg border border-gray-200 bg-white p-5 space-y-4">
              <div className="flex items-center justify-between gap-3">
                <h2 className="text-sm font-semibold text-gray-900">
                  Вопросы и ответы
                </h2>
              </div>

              {questions.length === 0 ? (
                <p className="text-sm text-gray-500">
                  Вопросов пока нет. Добавьте первый вопрос, чтобы начать настраивать тест.
                </p>
              ) : (
                <div className="space-y-4">
                  {questions.map((question) => (
                    <div
                      key={question.id}
                      className="rounded-md border border-gray-200 bg-gray-50/60 p-4 space-y-3"
                    >
                      <div className="flex items-start justify-between gap-3">
                        <div className="flex-1 space-y-2">
                          <div className="flex items-center gap-2">
                            <input
                              type="text"
                              value={question.text}
                              onChange={(event) =>
                                handleUpdateQuestion(question, {
                                  text: event.target.value,
                                })
                              }
                              className="w-full rounded-md border border-gray-300 px-3 py-1.5 text-xs text-gray-900 placeholder-gray-400 focus:outline-none focus:ring-2 focus:ring-gray-900 focus:border-transparent"
                            />
                          </div>

                          <div className="flex flex-wrap items-center gap-3 text-xs text-gray-600">
                            <div className="flex items-center gap-1.5">
                              <span className="text-[11px] uppercase tracking-wide text-gray-400">
                                Тип:
                              </span>
                              <select
                                value={question.type}
                                onChange={(event) =>
                                  handleUpdateQuestion(question, {
                                    type: event.target.value,
                                  })
                                }
                                className="rounded-md border border-gray-300 bg-white px-2 py-1 text-xs text-gray-900 focus:outline-none focus:ring-2 focus:ring-gray-900 focus:border-transparent"
                              >
                                <option value="single">Один ответ</option>
                                <option value="multiple">Несколько ответов</option>
                                <option value="text">Текстовый ответ</option>
                              </select>
                            </div>

                            <div className="flex items-center gap-1.5">
                              <span className="text-[11px] uppercase tracking-wide text-gray-400">
                                Баллы:
                              </span>
                              <input
                                type="number"
                                min={0}
                                value={question.points}
                                onChange={(event) =>
                                  handleUpdateQuestion(question, {
                                    points: Number.parseInt(event.target.value || '0', 10),
                                  })
                                }
                                className="w-16 rounded-md border border-gray-300 px-2 py-1 text-xs text-gray-900 focus:outline-none focus:ring-2 focus:ring-gray-900 focus:border-transparent"
                              />
                            </div>
                          </div>
                        </div>

                        <button
                          type="button"
                          onClick={() => handleDeleteQuestion(question)}
                          className="inline-flex items-center rounded-md border border-red-200 bg-red-50 px-2 py-1 text-[11px] font-medium text-red-700 hover:bg-red-100"
                        >
                          Удалить
                        </button>
                      </div>

                      <div className="space-y-2 border-t border-gray-200 pt-3">
                        <div className="flex items-center justify-between gap-3">
                          <p className="text-[11px] font-medium uppercase tracking-wide text-gray-400">
                            Варианты ответов
                          </p>
                          <button
                            type="button"
                            onClick={() => handleAddOption(question)}
                            className="inline-flex items-center rounded-md border border-gray-300 bg-white px-2 py-1 text-[11px] font-medium text-gray-700 hover:bg-gray-50"
                          >
                            + Вариант
                          </button>
                        </div>

                        {(question.options?.length ?? 0) === 0 ? (
                          <p className="text-xs text-gray-500">
                            Вариантов пока нет. Добавьте хотя бы один вариант ответа.
                          </p>
                        ) : (
                          <div className="space-y-2">
                            {(question.options ?? []).map((option) => (
                              <div
                                key={option.id}
                                className="flex items-center gap-2 rounded-md bg-white px-3 py-2 border border-gray-200"
                              >
                                <button
                                  type="button"
                                  onClick={() =>
                                    handleUpdateOption(question, option, {
                                      is_correct: !option.is_correct,
                                    })
                                  }
                                  className={`inline-flex h-5 w-5 items-center justify-center rounded border text-[10px] ${
                                    option.is_correct
                                      ? 'border-emerald-500 bg-emerald-50 text-emerald-700'
                                      : 'border-gray-300 bg-white text-gray-400'
                                  }`}
                                >
                                  ✓
                                </button>
                                <input
                                  type="text"
                                  value={option.text}
                                  onChange={(event) =>
                                    handleUpdateOption(question, option, {
                                      text: event.target.value,
                                    })
                                  }
                                  className="flex-1 rounded-md border border-gray-300 px-2 py-1 text-xs text-gray-900 placeholder-gray-400 focus:outline-none focus:ring-2 focus:ring-gray-900 focus:border-transparent"
                                />
                                <button
                                  type="button"
                                  onClick={() => handleDeleteOption(question, option)}
                                  className="inline-flex items-center rounded-md border border-red-200 bg-red-50 px-2 py-1 text-[11px] font-medium text-red-700 hover:bg-red-100"
                                >
                                  ×
                                </button>
                              </div>
                            ))}
                          </div>
                        )}
                      </div>
                    </div>
                  ))}
                </div>
              )}

              <div className="pt-2">
                <button
                  type="button"
                  onClick={handleAddQuestion}
                  className="inline-flex items-center rounded-md bg-gray-900 px-3 py-1.5 text-xs font-medium text-white shadow-sm hover:bg-gray-800"
                >
                  + Добавить вопрос
                </button>
              </div>
            </section>
          </>
        )}
      </main>
    </div>
  );
}


