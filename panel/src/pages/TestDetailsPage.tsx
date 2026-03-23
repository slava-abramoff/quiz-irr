import { useEffect, useMemo, useState } from "react";
import { useNavigate, useParams } from "react-router-dom";
import { deleteTest, getTestPreview } from "../api/tests";
import {
  analyzeAllRawByTest,
  analyzeRaw,
  deleteAllRawByTest,
  deleteRawAnswers,
  getRawResultsByTest,
} from "../api/raws";
import {
  deleteResult,
  deleteResultsByTest,
  getResultsByTest,
} from "../api/results";
import type {
  RawsInfoResponse,
  ResultsResponse,
  TestPreviewResponse,
} from "../api/types";

type TabKey = "unprocessed" | "results";

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

function formatRawStatus(status: string | null | undefined): string {
  if (status === "end") return "окончен";
  if (status === "started") return "начат";
  if (status === "handled") return "обработан";
  return status ?? "—";
}

export default function TestDetailsPage() {
  const { id } = useParams<{ id: string }>();
  const navigate = useNavigate();
  const [test, setTest] = useState<TestPreviewResponse | null>(null);
  const [isLoading, setIsLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [activeTab, setActiveTab] = useState<TabKey>("unprocessed");
  const [isDescExpanded, setIsDescExpanded] = useState(false);
  const [isDeleting, setIsDeleting] = useState(false);

  const PAGE_SIZE = 10;

  const [raws, setRaws] = useState<RawsInfoResponse | null>(null);
  const [rawsPage, setRawsPage] = useState(1);
  const [isLoadingRaws, setIsLoadingRaws] = useState(false);
  const [rawsError, setRawsError] = useState<string | null>(null);

  const [results, setResults] = useState<ResultsResponse | null>(null);
  const [resultsPage, setResultsPage] = useState(1);
  const [isLoadingResults, setIsLoadingResults] = useState(false);
  const [resultsError, setResultsError] = useState<string | null>(null);

  const [analyzingRawIds, setAnalyzingRawIds] = useState<
    Record<string, boolean>
  >({});
  const [deletingRawIds, setDeletingRawIds] = useState<Record<string, boolean>>(
    {},
  );
  const [isDeletingAllRaws, setIsDeletingAllRaws] = useState(false);
  const [isAnalyzingAllRaws, setIsAnalyzingAllRaws] = useState(false);

  const [deletingResultIds, setDeletingResultIds] = useState<
    Record<number, boolean>
  >({});
  const [isDeletingAllResults, setIsDeletingAllResults] = useState(false);

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

  const loadRaws = async (page: number) => {
    if (!id) return;
    setIsLoadingRaws(true);
    setRawsError(null);
    try {
      const skip = (page - 1) * PAGE_SIZE;
      const data = await getRawResultsByTest(id, { skip, take: PAGE_SIZE });
      setRaws(data);
    } catch (err: any) {
      const message =
        err?.response?.data?.message ??
        "Не удалось загрузить необработанные ответы.";
      setRawsError(message);
    } finally {
      setIsLoadingRaws(false);
    }
  };

  const loadResults = async (page: number) => {
    if (!id) return;
    setIsLoadingResults(true);
    setResultsError(null);
    try {
      const skip = (page - 1) * PAGE_SIZE;
      const data = await getResultsByTest(id, { skip, take: PAGE_SIZE });
      setResults(data);
    } catch (err: any) {
      const message =
        err?.response?.data?.message ??
        "Не удалось загрузить результаты теста.";
      setResultsError(message);
    } finally {
      setIsLoadingResults(false);
    }
  };

  useEffect(() => {
    if (!id) return;
    if (activeTab === "unprocessed") {
      void loadRaws(rawsPage);
    }
    if (activeTab === "results") {
      void loadResults(resultsPage);
    }
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [activeTab, rawsPage, resultsPage, id]);

  const handleAnalyzeRaw = async (rawId: string) => {
    setAnalyzingRawIds((prev) => ({ ...prev, [rawId]: true }));
    try {
      await analyzeRaw(rawId);
      await loadRaws(rawsPage);
    } catch (err: any) {
      const message =
        err?.response?.data?.message ?? "Не удалось выполнить анализ.";
      setRawsError(message);
    } finally {
      setAnalyzingRawIds((prev) => ({ ...prev, [rawId]: false }));
    }
  };

  const handleAnalyzeRawToResults = async (rawId: string) => {
    setAnalyzingRawIds((prev) => ({ ...prev, [rawId]: true }));
    try {
      await analyzeRaw(rawId);
      setActiveTab("results");
      setResultsPage(1);
    } catch (err: any) {
      const message =
        err?.response?.data?.message ?? "Не удалось выполнить анализ.";
      setRawsError(message);
    } finally {
      setAnalyzingRawIds((prev) => ({ ...prev, [rawId]: false }));
    }
  };

  const handleDeleteRaw = async (rawId: string) => {
    setDeletingRawIds((prev) => ({ ...prev, [rawId]: true }));
    try {
      await deleteRawAnswers(rawId);
      await loadRaws(rawsPage);
    } catch (err: any) {
      const message =
        err?.response?.data?.message ?? "Не удалось удалить запись.";
      setRawsError(message);
    } finally {
      setDeletingRawIds((prev) => ({ ...prev, [rawId]: false }));
    }
  };

  const handleAnalyzeAllRaws = async () => {
    if (!id || !raws?.data?.length) return;
    setIsAnalyzingAllRaws(true);
    setRawsError(null);
    try {
      await analyzeAllRawByTest(id);
      await loadRaws(rawsPage);
    } catch (err: any) {
      const message =
        err?.response?.data?.message ?? "Не удалось выполнить анализ для всех.";
      setRawsError(message);
    } finally {
      setIsAnalyzingAllRaws(false);
    }
  };

  const handleDeleteAllRaws = async () => {
    if (!id) return;
    if (!window.confirm("Удалить все необработанные ответы для этого теста?")) {
      return;
    }

    setIsDeletingAllRaws(true);
    setRawsError(null);
    try {
      await deleteAllRawByTest(id);
      setRaws(null);
      setRawsPage(1);
    } catch (err: any) {
      const message = err?.response?.data?.message ?? "Не удалось удалить все.";
      setRawsError(message);
    } finally {
      setIsDeletingAllRaws(false);
    }
  };

  const handleDeleteResultRow = async (resultId: number) => {
    setDeletingResultIds((prev) => ({ ...prev, [resultId]: true }));
    try {
      await deleteResult(resultId);
      await loadResults(resultsPage);
    } catch (err: any) {
      const message =
        err?.response?.data?.message ?? "Не удалось удалить результат.";
      setResultsError(message);
    } finally {
      setDeletingResultIds((prev) => ({ ...prev, [resultId]: false }));
    }
  };

  const handleDeleteAllResults = async () => {
    if (!id) return;
    if (!window.confirm("Удалить все результаты теста?")) return;

    setIsDeletingAllResults(true);
    setResultsError(null);
    try {
      await deleteResultsByTest(id);
      setResults(null);
      setResultsPage(1);
      await loadResults(1);
    } catch (err: any) {
      const message =
        err?.response?.data?.message ?? "Не удалось удалить все результаты.";
      setResultsError(message);
    } finally {
      setIsDeletingAllResults(false);
    }
  };

  const rawRows = raws?.data ?? [];
  const rawPagination = raws?.pagination;
  const resultsRows = results?.data ?? [];
  const resultsPagination = results?.pagination;

  const TrashIcon = ({ className }: { className?: string }) => (
    <svg
      className={className ?? "h-4 w-4"}
      viewBox="0 0 24 24"
      fill="none"
      xmlns="http://www.w3.org/2000/svg"
    >
      <path
        d="M9 3H15"
        stroke="currentColor"
        strokeWidth="2"
        strokeLinecap="round"
      />
      <path
        d="M4 6H20"
        stroke="currentColor"
        strokeWidth="2"
        strokeLinecap="round"
      />
      <path
        d="M6 6L7.2 20.2C7.27 21.03 7.97 21.65 8.8 21.65H15.2C16.03 21.65 16.73 21.03 16.8 20.2L18 6"
        stroke="currentColor"
        strokeWidth="2"
        strokeLinecap="round"
        strokeLinejoin="round"
      />
      <path
        d="M10 11V17"
        stroke="currentColor"
        strokeWidth="2"
        strokeLinecap="round"
      />
      <path
        d="M14 11V17"
        stroke="currentColor"
        strokeWidth="2"
        strokeLinecap="round"
      />
    </svg>
  );

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
                onClick={() => setActiveTab("unprocessed")}
                className={`border-b-2 px-1.5 pb-2 text-sm font-medium ${
                  activeTab === "unprocessed"
                    ? "border-gray-900 text-gray-900"
                    : "border-transparent text-gray-500 hover:text-gray-800"
                }`}
              >
                Необработанные ответы
              </button>
              <button
                type="button"
                onClick={() => setActiveTab("results")}
                className={`border-b-2 px-1.5 pb-2 text-sm font-medium ${
                  activeTab === "results"
                    ? "border-gray-900 text-gray-900"
                    : "border-transparent text-gray-500 hover:text-gray-800"
                }`}
              >
                Результаты теста
              </button>
            </div>
          </div>

          <div className="px-4 py-6">
            {activeTab === "unprocessed" ? (
              <div className="space-y-4">
                {rawsError && (
                  <div className="rounded-md bg-red-50 px-3 py-2 text-sm text-red-700 border border-red-100">
                    {rawsError}
                  </div>
                )}

                <div className="flex items-center justify-between gap-3">
                  <div className="flex items-center gap-2">
                    <button
                      type="button"
                      onClick={handleAnalyzeAllRaws}
                      disabled={
                        isAnalyzingAllRaws || isLoadingRaws || !rawRows.length
                      }
                      className="inline-flex items-center rounded-md bg-gray-900 px-3 py-1.5 text-xs font-medium text-white shadow-sm hover:bg-gray-800 disabled:opacity-60 disabled:cursor-not-allowed"
                    >
                      {isAnalyzingAllRaws ? "Проверяем..." : "Проверить все"}
                    </button>
                    <button
                      type="button"
                      onClick={handleDeleteAllRaws}
                      disabled={
                        isDeletingAllRaws || isLoadingRaws || !rawRows.length
                      }
                      className="inline-flex items-center rounded-md border border-red-200 bg-red-50 px-3 py-1.5 text-xs font-medium text-red-700 hover:bg-red-100 disabled:opacity-60 disabled:cursor-not-allowed"
                    >
                      Удалить все
                    </button>
                  </div>

                  {rawPagination && (
                    <p className="text-xs text-gray-500">
                      Страница {rawPagination.current_page} из{" "}
                      {rawPagination.total_pages}
                    </p>
                  )}
                </div>

                {isLoadingRaws ? (
                  <div className="flex items-center justify-center py-10 text-sm text-gray-500">
                    Загружаем...
                  </div>
                ) : (
                  <div className="overflow-x-auto">
                    <table className="min-w-full divide-y divide-gray-200 text-sm">
                      <thead>
                        <tr className="bg-gray-50 text-left text-xs font-semibold uppercase tracking-wide text-gray-500">
                          <th className="px-3 py-2">ФИО</th>
                          <th className="px-3 py-2">Почта</th>
                          <th className="px-3 py-2">Учреждение</th>
                          <th className="px-3 py-2">Статус</th>
                          <th className="px-3 py-2">Начат</th>
                          <th className="px-3 py-2">Закончен</th>
                          <th className="px-3 py-2 text-right">Действие</th>
                        </tr>
                      </thead>
                      <tbody className="divide-y divide-gray-100 bg-white">
                        {rawRows.length === 0 ? (
                          <tr>
                            <td
                              colSpan={7}
                              className="px-3 py-6 text-center text-sm text-gray-500"
                            >
                              Не найдено необработанных ответов.
                            </td>
                          </tr>
                        ) : (
                          rawRows.map((row) => {
                            const isStarted = row.status === "started";
                            const isAnalyzing = !!analyzingRawIds[row.id];
                            const isDeleting = !!deletingRawIds[row.id];
                            return (
                              <tr key={row.id} className="align-top">
                                <td className="px-3 py-3 font-medium text-gray-900">
                                  {row.full_name || "—"}
                                </td>
                                <td className="px-3 py-3 text-gray-700">
                                  {row.email || "—"}
                                </td>
                                <td className="px-3 py-3 text-gray-700">
                                  {row.org || "—"}
                                </td>
                                <td className="px-3 py-3 text-gray-700">
                                  {formatRawStatus(row.status)}
                                </td>
                                <td className="px-3 py-3 text-gray-700">
                                  {row.start_at
                                    ? formatDate(row.start_at)
                                    : "—"}
                                </td>
                                <td className="px-3 py-3 text-gray-700">
                                  {row.end_at ? formatDate(row.end_at) : "—"}
                                </td>
                                <td className="px-3 py-3">
                                  <div className="flex items-center justify-end gap-2">
                                    <button
                                      type="button"
                                      onClick={() => handleAnalyzeRaw(row.id)}
                                      disabled={
                                        isStarted || isAnalyzing || isDeleting
                                      }
                                      className="inline-flex items-center rounded-md bg-gray-900 px-2.5 py-1 text-[11px] font-medium text-white shadow-sm hover:bg-gray-800 disabled:opacity-60 disabled:cursor-not-allowed"
                                    >
                                      Проверить
                                    </button>
                                    <button
                                      type="button"
                                      onClick={() =>
                                        handleAnalyzeRawToResults(row.id)
                                      }
                                      disabled={
                                        isStarted || isAnalyzing || isDeleting
                                      }
                                      className="inline-flex items-center rounded-md border border-gray-200 bg-white px-2.5 py-1 text-[11px] font-medium text-gray-700 hover:bg-gray-50 disabled:opacity-60 disabled:cursor-not-allowed"
                                    >
                                      Анализ
                                    </button>
                                    <button
                                      type="button"
                                      onClick={() => handleDeleteRaw(row.id)}
                                      disabled={isDeleting || isAnalyzing}
                                      className="inline-flex items-center justify-center rounded-md border border-red-200 bg-red-50 p-2 text-red-700 hover:bg-red-100 disabled:opacity-60 disabled:cursor-not-allowed"
                                      aria-label="Удалить"
                                    >
                                      <TrashIcon className="h-4 w-4" />
                                    </button>
                                  </div>
                                </td>
                              </tr>
                            );
                          })
                        )}
                      </tbody>
                    </table>
                  </div>
                )}

                {rawPagination && rawPagination.total_pages > 1 && (
                  <div className="flex items-center justify-between border-t border-gray-100 pt-3">
                    <div className="text-xs text-gray-500">
                      Всего: {rawPagination.total_items}
                    </div>
                    <div className="flex items-center gap-2">
                      <button
                        type="button"
                        disabled={!rawPagination.has_previous_page}
                        onClick={() =>
                          rawPagination.has_previous_page &&
                          setRawsPage((p) => p - 1)
                        }
                        className="inline-flex items-center rounded-md border border-gray-200 bg-white px-3 py-1.5 text-xs font-medium text-gray-700 hover:bg-gray-50 disabled:opacity-50 disabled:cursor-not-allowed"
                      >
                        Назад
                      </button>
                      <button
                        type="button"
                        disabled={!rawPagination.has_next_page}
                        onClick={() =>
                          rawPagination.has_next_page &&
                          setRawsPage((p) => p + 1)
                        }
                        className="inline-flex items-center rounded-md border border-gray-200 bg-white px-3 py-1.5 text-xs font-medium text-gray-700 hover:bg-gray-50 disabled:opacity-50 disabled:cursor-not-allowed"
                      >
                        Вперед
                      </button>
                    </div>
                  </div>
                )}
              </div>
            ) : (
              <div className="space-y-4">
                {resultsError && (
                  <div className="rounded-md bg-red-50 px-3 py-2 text-sm text-red-700 border border-red-100">
                    {resultsError}
                  </div>
                )}

                {isLoadingResults ? (
                  <div className="flex items-center justify-center py-10 text-sm text-gray-500">
                    Загружаем...
                  </div>
                ) : (
                  <>
                    <div className="flex items-center justify-between gap-3">
                      <div className="flex items-center gap-2">
                        <button
                          type="button"
                          onClick={() =>
                            window.alert(
                              "Функция «Уведомить участников» пока в разработке.",
                            )
                          }
                          disabled={isLoadingResults || isDeletingAllResults}
                          className="inline-flex items-center rounded-md border border-gray-200 bg-white px-3 py-1.5 text-xs font-medium text-gray-700 hover:bg-gray-50 disabled:opacity-60 disabled:cursor-not-allowed"
                        >
                          Уведомить участников
                        </button>
                        <button
                          type="button"
                          onClick={() =>
                            window.alert(
                              "Функция «Импорт Excel» пока в разработке.",
                            )
                          }
                          disabled={isLoadingResults || isDeletingAllResults}
                          className="inline-flex items-center rounded-md border border-gray-200 bg-white px-3 py-1.5 text-xs font-medium text-gray-700 hover:bg-gray-50 disabled:opacity-60 disabled:cursor-not-allowed"
                        >
                          Импорт Excel
                        </button>
                      </div>

                      <button
                        type="button"
                        onClick={handleDeleteAllResults}
                        disabled={
                          isDeletingAllResults ||
                          isLoadingResults ||
                          !resultsRows.length
                        }
                        className="inline-flex items-center rounded-md border border-red-200 bg-red-50 px-3 py-1.5 text-xs font-medium text-red-700 hover:bg-red-100 disabled:opacity-60 disabled:cursor-not-allowed"
                      >
                        Удалить все
                      </button>
                    </div>

                    <div className="overflow-x-auto">
                      <table className="min-w-full divide-y divide-gray-200 text-sm">
                        <thead>
                          <tr className="bg-gray-50 text-left text-xs font-semibold uppercase tracking-wide text-gray-500">
                            <th className="px-3 py-2">ФИО</th>
                            <th className="px-3 py-2">Почта</th>
                            <th className="px-3 py-2">Учреждение</th>
                            <th className="px-3 py-2">Время</th>
                            <th className="px-3 py-2">В срок</th>
                            <th className="px-3 py-2">Баллы</th>
                            <th className="px-3 py-2 text-right">Действие</th>
                          </tr>
                        </thead>
                        <tbody className="divide-y divide-gray-100 bg-white">
                          {resultsRows.length === 0 ? (
                            <tr>
                              <td
                                colSpan={7}
                                className="px-3 py-6 text-center text-sm text-gray-500"
                              >
                                Пока нет результатов.
                              </td>
                            </tr>
                          ) : (
                            resultsRows.map((row) => (
                              <tr key={row.id} className="align-top">
                                <td className="px-3 py-3 font-medium text-gray-900">
                                  {row.full_name || "—"}
                                </td>
                                <td className="px-3 py-3 text-gray-700">
                                  {row.email || "—"}
                                </td>
                                <td className="px-3 py-3 text-gray-700">
                                  {row.org || "—"}
                                </td>
                                <td className="px-3 py-3 text-gray-700">
                                  {formatDuration(row.duration)}
                                </td>
                                <td className="px-3 py-3 text-gray-700">
                                  {row.is_on_time ? "да" : "нет"}
                                </td>
                                <td className="px-3 py-3 text-gray-700">
                                  {row.total_score}
                                </td>
                                <td className="px-3 py-3 text-right">
                                  <button
                                    type="button"
                                    onClick={() =>
                                      handleDeleteResultRow(row.id)
                                    }
                                    disabled={
                                      isDeletingAllResults ||
                                      !!deletingResultIds[row.id]
                                    }
                                    className="inline-flex items-center justify-center rounded-md border border-red-200 bg-red-50 p-2 text-red-700 hover:bg-red-100 disabled:opacity-60 disabled:cursor-not-allowed"
                                    aria-label="Удалить результат"
                                  >
                                    <TrashIcon className="h-4 w-4" />
                                  </button>
                                </td>
                              </tr>
                            ))
                          )}
                        </tbody>
                      </table>
                    </div>
                  </>
                )}

                {resultsPagination && resultsPagination.total_pages > 1 && (
                  <div className="flex items-center justify-between border-t border-gray-100 pt-3">
                    <div className="text-xs text-gray-500">
                      Всего: {resultsPagination.total_items}
                    </div>
                    <div className="flex items-center gap-2">
                      <button
                        type="button"
                        disabled={!resultsPagination.has_previous_page}
                        onClick={() =>
                          resultsPagination.has_previous_page &&
                          setResultsPage((p) => p - 1)
                        }
                        className="inline-flex items-center rounded-md border border-gray-200 bg-white px-3 py-1.5 text-xs font-medium text-gray-700 hover:bg-gray-50 disabled:opacity-50 disabled:cursor-not-allowed"
                      >
                        Назад
                      </button>
                      <button
                        type="button"
                        disabled={!resultsPagination.has_next_page}
                        onClick={() =>
                          resultsPagination.has_next_page &&
                          setResultsPage((p) => p + 1)
                        }
                        className="inline-flex items-center rounded-md border border-gray-200 bg-white px-3 py-1.5 text-xs font-medium text-gray-700 hover:bg-gray-50 disabled:opacity-50 disabled:cursor-not-allowed"
                      >
                        Вперед
                      </button>
                    </div>
                  </div>
                )}
              </div>
            )}
          </div>
        </section>
      </main>
    </div>
  );
}
