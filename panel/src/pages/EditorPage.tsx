import { useEffect, useMemo, useState } from 'react';
import type { FormEvent } from 'react';
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
import AppHeader from '../components/AppHeader';
import TestMetaForm from '../components/editor/TestMetaForm';
import QuestionsEditor from '../components/editor/QuestionsEditor';

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

  // Используем смещение именно выбранной даты/времени (важно для DST-зон).
  const localDateTime = new Date(`${datePart}T${timePart}:00`);
  if (Number.isNaN(localDateTime.getTime())) return undefined;
  const offsetMinutes = -localDateTime.getTimezoneOffset(); // например, +180 для UTC+3
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
    const testId = id;
    if (!isAuthenticated) {
      navigate('/login', { replace: true });
      return;
    }
    if (!testId) return;

    async function loadTest() {
      try {
        setIsLoading(true);
        setError(null);
        // mode "fulldata" – сервер вернёт полные данные теста
        const data = await getTest(testId as string, 'fulldata');

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

  return (
    <div className="min-h-screen bg-gray-50">
      <AppHeader subtitle="Редактор теста" />

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
            <TestMetaForm
              title={title}
              desc={desc}
              startAt={startAt}
              endAt={endAt}
              isActive={isActive}
              durationHours={durationHours}
              durationMinutes={durationMinutes}
              durationSeconds={durationSeconds}
              isSaving={isSavingMeta}
              onChangeTitle={setTitle}
              onChangeDesc={setDesc}
              onChangeStartAt={setStartAt}
              onChangeEndAt={setEndAt}
              onToggleActive={() => setIsActive((prev) => !prev)}
              onDurationChange={(field, value) => {
                if (field === 'h') setDurationHours(value);
                if (field === 'm') setDurationMinutes(value);
                if (field === 's') setDurationSeconds(value);
              }}
              onSubmit={handleMetaSubmit}
            />

            <QuestionsEditor
              questions={questions}
              onAddQuestion={handleAddQuestion}
              onUpdateQuestion={handleUpdateQuestion}
              onDeleteQuestion={handleDeleteQuestion}
              onAddOption={handleAddOption}
              onUpdateOption={handleUpdateOption}
              onDeleteOption={handleDeleteOption}
            />
          </>
        )}
      </main>
    </div>
  );
}


