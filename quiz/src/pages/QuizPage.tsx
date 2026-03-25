import { useCallback, useEffect, useMemo, useState } from "react";
import { useParams } from "react-router-dom";
import { fetchExamInfo, saveExamAnswers, startExam } from "../api/exam";
import { ExamQuestionStep } from "../components/ExamQuestionStep";
import { ExamStartForm } from "../components/ExamStartForm";
import { QuizCollapsibleDesc } from "../components/QuizCollapsibleDesc";
import type { ExamQuestion } from "../types/examQuiz";
import type { UserAnswerPayload } from "../types/examQuiz";
import type { TestCustomerResponse } from "../types";
import {
  formatCountdownTo,
  formatDurationSeconds,
  formatLocalDateTime,
  parseMoscowDatetime,
} from "../utils/time";
import { shuffleArray } from "../utils/shuffle";
import "./QuizPage.css";

type LoadState =
  | { status: "loading" }
  | { status: "error" }
  | { status: "ok"; data: TestCustomerResponse };

type Phase = "intro" | "form" | "quiz" | "done";

type StartFormState = {
  fullName: string;
  email: string;
  org: string;
  birthYear: string;
};

type QuizSession = {
  rawId: string;
  questions: ExamQuestion[];
  currentIndex: number;
  answers: UserAnswerPayload[];
};

function prepareQuestions(questions: ExamQuestion[]): ExamQuestion[] {
  const withShuffledOptions = questions.map((q) => ({
    ...q,
    options: shuffleArray([...q.options]),
  }));
  return shuffleArray(withShuffledOptions);
}

function QuizContent({ testId }: { testId: string }) {
  const [state, setState] = useState<LoadState>({ status: "loading" });
  const [now, setNow] = useState(() => new Date());

  useEffect(() => {
    const shouldAllow = (target: EventTarget | null): boolean => {
      if (!(target instanceof HTMLElement)) return false;
      const tag = target.tagName;
      if (tag === "INPUT" || tag === "TEXTAREA" || tag === "SELECT") return true;
      return Boolean(target.closest("input, textarea, select, [contenteditable='true']"));
    };

    const preventIfNotAllowed = (e: Event) => {
      if (shouldAllow(e.target)) return;
      e.preventDefault();
    };

    document.addEventListener("copy", preventIfNotAllowed);
    document.addEventListener("cut", preventIfNotAllowed);
    document.addEventListener("contextmenu", preventIfNotAllowed);
    document.addEventListener("selectstart", preventIfNotAllowed);
    document.addEventListener("dragstart", preventIfNotAllowed);

    return () => {
      document.removeEventListener("copy", preventIfNotAllowed);
      document.removeEventListener("cut", preventIfNotAllowed);
      document.removeEventListener("contextmenu", preventIfNotAllowed);
      document.removeEventListener("selectstart", preventIfNotAllowed);
      document.removeEventListener("dragstart", preventIfNotAllowed);
    };
  }, []);

  const [phase, setPhase] = useState<Phase>("intro");
  const [startForm, setStartForm] = useState<StartFormState>({
    fullName: "",
    email: "",
    org: "",
    birthYear: "",
  });
  const [startError, setStartError] = useState<string | null>(null);
  const [startLoading, setStartLoading] = useState(false);

  const [session, setSession] = useState<QuizSession | null>(null);
  const [secondsLeft, setSecondsLeft] = useState(0);
  const [saveError, setSaveError] = useState<string | null>(null);
  const [saveLoading, setSaveLoading] = useState(false);
  const [pendingSave, setPendingSave] = useState<{
    rawId: string;
    answers: UserAnswerPayload[];
  } | null>(null);

  useEffect(() => {
    let cancelled = false;
    fetchExamInfo(testId)
      .then((data) => {
        if (!cancelled) setState({ status: "ok", data });
      })
      .catch(() => {
        if (!cancelled) setState({ status: "error" });
      });
    return () => {
      cancelled = true;
    };
  }, [testId]);

  useEffect(() => {
    const id = window.setInterval(() => setNow(new Date()), 1000);
    return () => window.clearInterval(id);
  }, []);

  useEffect(() => {
    if (phase !== "quiz") return;
    const id = window.setInterval(() => {
      setSecondsLeft((s) => Math.max(0, s - 1));
    }, 1000);
    return () => window.clearInterval(id);
  }, [phase]);

  const parsed = useMemo(() => {
    if (state.status !== "ok") return null;
    const start = parseMoscowDatetime(state.data.start_at);
    const end = parseMoscowDatetime(state.data.end_at);
    return {
      start,
      end,
      valid: !Number.isNaN(start.getTime()) && !Number.isNaN(end.getTime()),
    };
  }, [state]);

  const handleStartFormChange = useCallback(
    (field: keyof StartFormState, value: string) => {
      setStartForm((prev) => ({ ...prev, [field]: value }));
    },
    [],
  );

  const handleExamStart = useCallback(async () => {
    if (state.status !== "ok") return;
    setStartLoading(true);
    setStartError(null);
    try {
      const data = await startExam(testId, {
        full_name: startForm.fullName.trim(),
        email: startForm.email.trim(),
        org: startForm.org.trim(),
        birth_year: Number(startForm.birthYear),
      });
      setSession({
        rawId: data.raw_id,
        questions: prepareQuestions(data.questions),
        currentIndex: 0,
        answers: [],
      });
      setSecondsLeft(Math.max(0, Math.floor(state.data.duration)));
      setPhase("quiz");
    } catch (e) {
      setStartError(e instanceof Error ? e.message : "Не удалось начать тест");
    } finally {
      setStartLoading(false);
    }
  }, [state, startForm, testId]);

  const handleQuizAdvance = useCallback(
    (answer: UserAnswerPayload) => {
      if (!session) return;
      const isLast = session.currentIndex >= session.questions.length - 1;
      const nextAnswers = [...session.answers, answer];

      if (!isLast) {
        setSession({
          ...session,
          answers: nextAnswers,
          currentIndex: session.currentIndex + 1,
        });
        setSaveError(null);
        setPendingSave(null);
        return;
      }

      setSaveLoading(true);
      setSaveError(null);
      setPendingSave({ rawId: session.rawId, answers: nextAnswers });

      saveExamAnswers(session.rawId, { answers: nextAnswers })
        .then(() => {
          setPhase("done");
          setSession(null);
          setPendingSave(null);
        })
        .catch((e) => {
          setSaveError(
            e instanceof Error ? e.message : "Не удалось отправить ответы",
          );
        })
        .finally(() => setSaveLoading(false));
    },
    [session],
  );

  const handleRetrySave = useCallback(() => {
    if (!pendingSave) return;
    setSaveLoading(true);
    setSaveError(null);
    saveExamAnswers(pendingSave.rawId, { answers: pendingSave.answers })
      .then(() => {
        setPhase("done");
        setSession(null);
        setPendingSave(null);
      })
      .catch((e) => {
        setSaveError(
          e instanceof Error ? e.message : "Не удалось отправить ответы",
        );
      })
      .finally(() => setSaveLoading(false));
  }, [pendingSave]);

  if (state.status === "loading") {
    return (
      <div className="quiz-shell">
        <div className="quiz-card">
          <p className="quiz-loading">Загрузка…</p>
        </div>
      </div>
    );
  }

  if (state.status === "error") {
    return (
      <div className="quiz-shell">
        <div className="quiz-card quiz-card--message">
          <p className="quiz-message">
            Такого квиза нет или он пока недоступен.
          </p>
        </div>
      </div>
    );
  }

  const { data } = state;
  const durationLabel = formatDurationSeconds(data.duration);
  const beforeStart =
    parsed?.valid && now < parsed.start
      ? formatCountdownTo(parsed.start, now)
      : null;

  const inActiveWindow =
    parsed?.valid === true && now >= parsed.start && now <= parsed.end;

  if (phase === "form") {
    return (
      <div className="quiz-shell">
        <div className="quiz-card">
          <ExamStartForm
            values={startForm}
            onChange={handleStartFormChange}
            onSubmit={handleExamStart}
            loading={startLoading}
            error={startError}
          />
        </div>
      </div>
    );
  }

  if (phase === "quiz" && session) {
    const q = session.questions[session.currentIndex];
    const total = session.questions.length;
    const idx = session.currentIndex;
    const progressText = `${idx + 1} из ${total}`;
    const timerLabel = formatDurationSeconds(secondsLeft);

    return (
      <div className="quiz-shell">
        <div className="quiz-card quiz-card--quiz">
          <div className="exam-quiz-toolbar">
            <div className="exam-timer" aria-live="polite">
              <span className="exam-timer-label">Осталось</span>
              <span
                className={`exam-timer-value${secondsLeft <= 60 && secondsLeft > 0 ? " exam-timer-value--warn" : ""}${secondsLeft === 0 ? " exam-timer-value--zero" : ""}`}
              >
                {timerLabel}
              </span>
            </div>
          </div>
          {q ? (
            <ExamQuestionStep
              key={q.id}
              question={q}
              progressText={progressText}
              onAdvance={handleQuizAdvance}
              busy={saveLoading}
              isLast={idx >= total - 1}
              saveError={saveError}
              onRetrySave={handleRetrySave}
            />
          ) : null}
        </div>
      </div>
    );
  }

  if (phase === "done") {
    return (
      <div className="quiz-shell">
        <div className="quiz-card quiz-card--message">
          <p className="quiz-message">Спасибо! Ответы отправлены.</p>
        </div>
      </div>
    );
  }

  return (
    <div className="quiz-shell">
      <article className="quiz-card">
        <p className="quiz-eyebrow">Информация о тесте</p>
        <h1 className="quiz-title">{data.title}</h1>
        <QuizCollapsibleDesc text={data.desc} />

        <dl className="quiz-meta">
          <div className="quiz-meta-row">
            <dt>Время на прохождение</dt>
            <dd>{durationLabel}</dd>
          </div>
          {parsed?.valid ? (
            <>
              <div className="quiz-meta-row">
                <dt>Начало</dt>
                <dd>{formatLocalDateTime(parsed.start)}</dd>
              </div>
              <div className="quiz-meta-row">
                <dt>Окончание</dt>
                <dd>{formatLocalDateTime(parsed.end)}</dd>
              </div>
            </>
          ) : (
            <div className="quiz-meta-row">
              <dt>Период проведения</dt>
              <dd>Не удалось разобрать даты</dd>
            </div>
          )}
        </dl>

        {beforeStart !== null && (
          <section className="quiz-countdown" aria-live="polite">
            <h2 className="quiz-countdown-label">До начала</h2>
            <p className="quiz-countdown-value">{beforeStart}</p>
          </section>
        )}

        {inActiveWindow && (
          <div className="quiz-actions">
            <button
              type="button"
              className="quiz-start-btn"
              onClick={() => {
                setPhase("form");
                setStartError(null);
              }}
            >
              Приступить
            </button>
          </div>
        )}
      </article>
    </div>
  );
}

export function QuizPage() {
  const { testId } = useParams<{ testId?: string }>();

  if (!testId) {
    return (
      <div className="quiz-shell">
        <div className="quiz-card quiz-card--message">
          <p className="quiz-message">
            Укажите идентификатор квиза в адресе страницы.
          </p>
        </div>
      </div>
    );
  }

  return <QuizContent key={testId} testId={testId} />;
}
