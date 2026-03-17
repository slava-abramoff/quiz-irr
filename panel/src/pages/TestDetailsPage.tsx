import { useEffect, useMemo, useState } from "react";
import { useNavigate, useParams } from "react-router-dom";
import { deleteTest, getTestPreview } from "../api/tests";
import type { TestPreviewResponse } from "../api/types";

type TabKey = "overview" | "analytics";

const MONTHS_RU = [
  "января",
  "февраля",
  "марта",
  "апреля",
  "мая",
  "июня",
  "июля",
  "августа",
  "сентября",
  "октября",
  "ноября",
  "декабря",
];

function formatDate(value: string | null | undefined): string {
  if (!value) return "—";

  // backend format: "2006-01-02 15:04:05"
  const safe = value.replace(" ", "T");
  const date = new Date(safe);

  if (Number.isNaN(date.getTime())) return value;

  const day = String(date.getDate()).padStart(2, "0");
  const monthName = MONTHS_RU[date.getMonth()];
  const year = date.getFullYear();
  const hours = String(date.getHours()).padStart(2, "0");
  const minutes = String(date.getMinutes()).padStart(2, "0");

  return `${day} ${monthName} ${year} ${hours}:${minutes}`;
}

function formatDuration(seconds: number | null | undefined): string {
  if (!seconds || seconds <= 0) return "—";

  const hrs = Math.floor(seconds / 3600);
  const mins = Math.floor((seconds % 3600) / 60);
  const secs = seconds % 60;

  if (hrs > 0) {
    return [
      String(hrs).padStart(2, "0"),
      String(mins).padStart(2, "0"),
      String(secs).padStart(2, "0"),
    ].join(":");
  }

  return [String(mins).padStart(2, "0"), String(secs).padStart(2, "0")].join(
    ":",
  );
}

export default function TestDetailsPage() {
  const { id } = useParams<{ id: string }>();
  const navigate = useNavigate();
  const [test, setTest] = useState<TestPreviewResponse | null>(null);
  const [isLoading, setIsLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [activeTab, setActiveTab] = useState<TabKey>("overview");
  const [isDescExpanded, setIsDescExpanded] = useState(false);
  const [isDeleting, setIsDeleting] = useState(false);

  const isAuthenticated =
    typeof window !== "undefined" && !!localStorage.getItem("access_token");

  useEffect(() => {
    if (!isAuthenticated) {
      navigate("/login", { replace: true });
      return;
    }
    if (!id) return;

    async function loadTest() {
      try {
        setIsLoading(true);
        setError(null);
        const data = await getTestPreview(id);
        setTest(data);
      } catch (err: any) {
        const message =
          err?.response?.data?.message ?? "Не удалось загрузить тест.";
        setError(message);
      } finally {
        setIsLoading(false);
      }
    }

    void loadTest();
  }, [id, isAuthenticated, navigate]);

  const formattedStartAt = useMemo(
    () => formatDate(test?.start_at),
    [test?.start_at],
  );
  const formattedEndAt = useMemo(
    () => formatDate(test?.end_at),
    [test?.end_at],
  );
  const formattedDuration = useMemo(
    () => formatDuration(test?.duration),
    [test?.duration],
  );

  const shouldClampDesc = (test?.desc?.length ?? 0) > 320;
  const visibleDesc = useMemo(() => {
    if (!test?.desc) return "";
    if (!shouldClampDesc || isDescExpanded) return test.desc;
    return `${test.desc.slice(0, 320)}…`;
  }, [isDescExpanded, shouldClampDesc, test?.desc]);

  const handleDeleteTest = async () => {
    if (!id) return;
    if (
      !window.confirm(
        "Удалить тест и все связанные вопросы, результаты и ответы? Это действие необратимо.",
      )
    ) {
      return;
    }

    try {
      setIsDeleting(true);
      await deleteTest(id);
      navigate("/", { replace: true });
    } catch (err: any) {
      const message =
        err?.response?.data?.message ?? "Не удалось удалить тест.";
      setError(message);
    } finally {
      setIsDeleting(false);
    }
  };

  return (
    <div className="min-h-screen bg-gray-50">
      <header className="border-b border-gray-200 bg-white">
        <div className="mx-auto max-w-6xl px-4 py-4 flex items-center justify-between gap-4">
          <div>
            <button
              type="button"
              onClick={() => navigate("/")}
              className="text-left"
            >
              <h1 className="text-lg font-semibold text-gray-900">
                Конструктор тестов
              </h1>
              <p className="text-sm text-gray-500">
                Управляйте тестами, вопросами и результатами
              </p>
            </button>
          </div>

          <div className="flex items-center gap-3">
            {/* сюда позже можно добавить профиль пользователя или кнопку выхода */}
          </div>
        </div>
      </header>

      <main className="mx-auto max-w-6xl px-4 py-6 space-y-6">
        {error && (
          <div className="rounded-md bg-red-50 px-3 py-2 text-sm text-red-700 border border-red-100">
            {error}
          </div>
        )}

        {/* Блок информации о тесте */}
        <section className="bg-white border border-gray-200 rounded-lg shadow-sm p-5">
          {isLoading && !test ? (
            <div className="flex items-center justify-center py-10 text-sm text-gray-500">
              Загружаем данные теста...
            </div>
          ) : !test ? (
            <div className="flex items-center justify-center py-10 text-sm text-gray-500">
              Тест не найден.
            </div>
          ) : (
            <div className="space-y-4">
              <div className="flex items-start justify-between gap-4">
                <div className="min-w-0 flex-1">
                  <h2 className="text-base font-semibold text-gray-900 underline underline-offset-4 decoration-gray-300">
                    {test.title}
                  </h2>
                  {test.desc && (
                    <div className="mt-2 text-sm text-gray-700">
                      <p className="whitespace-pre-line break-words">
                        {visibleDesc}
                      </p>
                      {shouldClampDesc && (
                        <button
                          type="button"
                          onClick={() => setIsDescExpanded((prev) => !prev)}
                          className="mt-1 text-xs font-medium text-gray-500 hover:text-gray-700 underline underline-offset-2"
                        >
                          {isDescExpanded ? "Свернуть" : "Показать полностью"}
                        </button>
                      )}
                    </div>
                  )}
                </div>

                <div className="flex flex-col items-end gap-2 text-xs">
                  <span
                    className={`inline-flex items-center rounded-full px-2 py-0.5 font-medium ${
                      test.is_active
                        ? "bg-emerald-50 text-emerald-700"
                        : "bg-gray-100 text-gray-500"
                    }`}
                  >
                    {test.is_active ? "Активен" : "Неактивен"}
                  </span>
                  <p className="text-gray-500">
                    Автор:{" "}
                    <span className="font-medium text-gray-800">
                      {test.author || "—"}
                    </span>
                  </p>
                  <button
                    type="button"
                    onClick={() => id && navigate(`/editor/${id}`)}
                    className="mt-1 inline-flex items-center rounded-md bg-gray-900 px-3 py-1.5 text-[11px] font-medium text-white shadow-sm hover:bg-gray-800"
                  >
                    Редактировать
                  </button>
                  <button
                    type="button"
                    onClick={handleDeleteTest}
                    disabled={isDeleting}
                    className="inline-flex items-center rounded-md border border-red-200 bg-red-50 px-3 py-1.5 text-[11px] font-medium text-red-700 hover:bg-red-100 disabled:opacity-60 disabled:cursor-not-allowed"
                  >
                    {isDeleting ? "Удаляем..." : "Удалить тест"}
                  </button>
                </div>
              </div>

              <dl className="grid grid-cols-1 gap-4 border-t border-gray-100 pt-4 text-sm text-gray-600 sm:grid-cols-3">
                <div>
                  <dt className="text-xs uppercase tracking-wide text-gray-400">
                    Начало тестирования
                  </dt>
                  <dd className="mt-0.5 font-medium text-gray-800">
                    {formattedStartAt}
                  </dd>
                </div>
                <div>
                  <dt className="text-xs uppercase tracking-wide text-gray-400">
                    Окончание тестирования
                  </dt>
                  <dd className="mt-0.5 font-medium text-gray-800">
                    {formattedEndAt}
                  </dd>
                </div>
                <div>
                  <dt className="text-xs uppercase tracking-wide text-gray-400">
                    Время на прохождение
                  </dt>
                  <dd className="mt-0.5 font-medium text-gray-800">
                    {formattedDuration}
                  </dd>
                </div>
              </dl>
            </div>
          )}
        </section>

        {/* Нижний блок с вкладками */}
        <section className="bg-white border border-gray-200 rounded-lg shadow-sm">
          <div className="border-b border-gray-100 px-4 pt-3">
            <div className="flex items-center gap-4">
              <button
                type="button"
                onClick={() => setActiveTab("overview")}
                className={`border-b-2 px-1.5 pb-2 text-sm font-medium ${
                  activeTab === "overview"
                    ? "border-gray-900 text-gray-900"
                    : "border-transparent text-gray-500 hover:text-gray-800"
                }`}
              >
                Необработанные ответы
              </button>
              <button
                type="button"
                onClick={() => setActiveTab("analytics")}
                className={`border-b-2 px-1.5 pb-2 text-sm font-medium ${
                  activeTab === "analytics"
                    ? "border-gray-900 text-gray-900"
                    : "border-transparent text-gray-500 hover:text-gray-800"
                }`}
              >
                Результаты теста
              </button>
            </div>
          </div>

          <div className="px-4 py-6 text-sm text-gray-500">
            {activeTab === "overview" ? (
              <p>
                Здесь позже появится содержимое первой вкладки (обзор теста,
                сводка и быстрые действия).
              </p>
            ) : (
              <p>
                Здесь позже появится содержимое второй вкладки (аналитика,
                результаты, отчёты).
              </p>
            )}
          </div>
        </section>
      </main>
    </div>
  );
}
