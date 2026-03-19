import type { TestAdminResponse } from "../../api/types";

interface TestsListProps {
  tests: TestAdminResponse[];
  isLoading: boolean;
  error: string | null;
  currentPage: number;
  totalPages: number;
  onChangePage: (page: number) => void;
  onOpenCreate: () => void;
  onOpenTest: (id: string) => void;
}

const PAGE_MIN_HEIGHT_CLASS = "min-h-[200px]";

export default function TestsList({
  tests,
  isLoading,
  error,
  currentPage,
  totalPages,
  onChangePage,
  onOpenCreate,
  onOpenTest,
}: TestsListProps) {
  const canGoPrev = currentPage > 1;
  const canGoNext = currentPage < totalPages;

  return (
    <main className="mx-auto max-w-6xl px-4 py-6">
      <div className="mb-4 flex items-center justify-between gap-3">
        <h2 className="text-base font-medium text-gray-900">Список тестов</h2>

        <button
          type="button"
          onClick={onOpenCreate}
          className="inline-flex items-center rounded-md bg-gray-900 px-4 py-2 text-sm font-medium text-white shadow-sm hover:bg-gray-800 focus:outline-none focus:ring-2 focus:ring-gray-900 focus:ring-offset-2"
        >
          + Новый тест
        </button>
      </div>

      {error && (
        <div className="mb-4 rounded-md bg-red-50 px-3 py-2 text-sm text-red-700 border border-red-100">
          {error}
        </div>
      )}

      <div className="bg-white border border-gray-200 rounded-lg shadow-sm overflow-hidden">
        <div className={PAGE_MIN_HEIGHT_CLASS}>
          {isLoading ? (
            <div className="flex items-center justify-center py-16 text-sm text-gray-500">
              Загружаем тесты...
            </div>
          ) : tests.length === 0 ? (
            <div className="flex flex-col items-center justify-center py-16 text-center px-4">
              <p className="text-sm font-medium text-gray-900 mb-1">
                Тестов пока нет
              </p>
              <p className="text-sm text-gray-500 mb-3">
                Создайте первый тест, чтобы начать пользоваться конструктором.
              </p>
              <button
                type="button"
                onClick={onOpenCreate}
                className="inline-flex items-center rounded-md bg-gray-900 px-4 py-2 text-sm font-medium text-white shadow-sm hover:bg-gray-800"
              >
                Создать тест
              </button>
            </div>
          ) : (
            <ul className="divide-y divide-gray-100">
              {tests.map((test) => (
                <li
                  key={test.id}
                  className="px-4 py-3 flex items-center justify-between gap-4"
                >
                  <div className="min-w-0 flex-1">
                    <div className="flex items-center gap-2">
                      <p className="truncate text-sm font-medium text-gray-900">
                        {test.title}
                      </p>
                      {test.is_active && (
                        <span className="inline-flex items-center rounded-full bg-emerald-50 px-2 py-0.5 text-[11px] font-medium text-emerald-700">
                          Активен
                        </span>
                      )}
                    </div>
                    <p className="mt-0.5 line-clamp-1 text-xs text-gray-500">
                      {test.desc}
                    </p>
                    <p className="mt-1 text-[11px] text-gray-400">
                      Автор: {test.author}
                    </p>
                  </div>

                  <div className="flex items-center gap-2">
                    <button
                      type="button"
                      onClick={() => onOpenTest(test.id)}
                      className="inline-flex items-center rounded-md border border-gray-200 bg-white px-3 py-1.5 text-xs font-medium text-gray-700 hover:bg-gray-50"
                    >
                      Перейти в тест
                    </button>
                  </div>
                </li>
              ))}
            </ul>
          )}
        </div>

        <div className="flex items-center justify-between border-t border-gray-100 bg-gray-50 px-4 py-2.5">
          <p className="text-xs text-gray-500">
            Страница {currentPage} из {totalPages}
          </p>
          <div className="flex items-center gap-2">
            <button
              type="button"
              disabled={!canGoPrev}
              onClick={() => canGoPrev && onChangePage(currentPage - 1)}
              className="inline-flex items-center rounded-md border border-gray-200 bg-white px-3 py-1.5 text-xs font-medium text-gray-700 hover:bg-gray-50 disabled:opacity-50 disabled:cursor-not-allowed"
            >
              Назад
            </button>
            <button
              type="button"
              disabled={!canGoNext}
              onClick={() => canGoNext && onChangePage(currentPage + 1)}
              className="inline-flex items-center rounded-md border border-gray-200 bg-white px-3 py-1.5 text-xs font-medium text-gray-700 hover:bg-gray-50 disabled:opacity-50 disabled:cursor-not-allowed"
            >
              Вперед
            </button>
          </div>
        </div>
      </div>
    </main>
  );
}
