import { useEffect, useState } from 'react';
import type { FormEvent } from 'react';
import { useNavigate } from 'react-router-dom';
import { getTests, deleteTest, createTest } from '../api/tests';
import type { CreateTestRequest, TestAdminResponse } from '../api/types';

const PAGE_SIZE = 10;

export default function TestsPage() {
  const navigate = useNavigate();
  const [tests, setTests] = useState<TestAdminResponse[]>([]);
  const [currentPage, setCurrentPage] = useState(1);
  const [totalPages, setTotalPages] = useState(1);
  const [isLoading, setIsLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [isCreateOpen, setIsCreateOpen] = useState(false);
  const [isCreating, setIsCreating] = useState(false);
  const [createError, setCreateError] = useState<string | null>(null);
  const [createForm, setCreateForm] = useState<CreateTestRequest>({
    title: '',
    desc: '',
  });

  const isAuthenticated =
    typeof window !== 'undefined' && !!localStorage.getItem('access_token');

  useEffect(() => {
    if (!isAuthenticated) {
      navigate('/login', { replace: true });
      return;
    }

    const controller = new AbortController();

    async function loadTests() {
      try {
        setIsLoading(true);
        setError(null);

        const skip = (currentPage - 1) * PAGE_SIZE;
        const data = await getTests({ skip, take: PAGE_SIZE });

        setTests(data.tests);
        setTotalPages(data.pagination.total_pages || 1);
      } catch (err: any) {
        const message = err?.response?.data?.message ?? 'Не удалось загрузить список тестов.';
        setError(message);
      } finally {
        setIsLoading(false);
      }
    }

    void loadTests();

    return () => {
      controller.abort();
    };
  }, [currentPage, isAuthenticated, navigate]);

  const handleLogout = () => {
    if (typeof window !== 'undefined') {
      localStorage.removeItem('access_token');
      localStorage.removeItem('refresh_token');
      localStorage.removeItem('user_full_name');
      localStorage.removeItem('user_email');
    }

    navigate('/login', { replace: true });
  };

  const handleCreateTest = () => {
    setCreateError(null);
    setCreateForm({ title: '', desc: '' });
    setIsCreateOpen(true);
  };

  const handleSubmitCreate = async (event: FormEvent) => {
    event.preventDefault();
    if (!createForm.title.trim()) {
      setCreateError('Введите название теста.');
      return;
    }

    try {
      setIsCreating(true);
      setCreateError(null);
      const created = await createTest({
        title: createForm.title.trim(),
        desc: createForm.desc.trim(),
      });

      setIsCreateOpen(false);
      navigate(`/test/${created.id}`);
    } catch (err: any) {
      const message = err?.response?.data?.message ?? 'Не удалось создать тест.';
      setCreateError(message);
    } finally {
      setIsCreating(false);
    }
  };

  const handleDelete = async (id: string) => {
    if (!window.confirm('Удалить тест и все связанные вопросы/результаты?')) {
      return;
    }

    try {
      await deleteTest(id);
      setTests((prev) => prev.filter((t) => t.id !== id));
    } catch (err: any) {
      const message = err?.response?.data?.message ?? 'Не удалось удалить тест.';
      setError(message);
    }
  };

  const canGoPrev = currentPage > 1;
  const canGoNext = currentPage < totalPages;

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
              <p className="text-sm text-gray-500">Управляйте тестами, вопросами и результатами</p>
            </button>
          </div>

          <div className="flex items-center gap-3">
            <button
              type="button"
              onClick={handleLogout}
              className="inline-flex items-center rounded-md border border-gray-300 bg-white px-3 py-1.5 text-xs font-medium text-gray-700 shadow-sm hover:bg-gray-50"
            >
              Выйти
            </button>
          </div>
        </div>
      </header>

      <main className="mx-auto max-w-6xl px-4 py-6">
        <div className="mb-4 flex items-center justify-between gap-3">
          <h2 className="text-base font-medium text-gray-900">
            Список тестов
          </h2>

          <button
            type="button"
            onClick={handleCreateTest}
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
          <div className="min-h-[200px]">
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
                  onClick={handleCreateTest}
                  className="inline-flex items-center rounded-md bg-gray-900 px-4 py-2 text-sm font-medium text-white shadow-sm hover:bg-gray-800"
                >
                  Создать тест
                </button>
              </div>
            ) : (
              <ul className="divide-y divide-gray-100">
                {tests.map((test) => (
                  <li key={test.id} className="px-4 py-3 flex items-center justify-between gap-4">
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
                        Вопросов: {test.questions?.length ?? 0} · Автор: {test.author}
                      </p>
                    </div>

                    <div className="flex items-center gap-2">
                      <button
                        type="button"
                        onClick={() => navigate(`/test/${test.id}`)}
                        className="inline-flex items-center rounded-md border border-gray-200 bg-white px-3 py-1.5 text-xs font-medium text-gray-700 hover:bg-gray-50"
                      >
                        Перейти в тест
                      </button>
                      <button
                        type="button"
                        onClick={() => handleDelete(test.id)}
                        className="inline-flex items-center rounded-md border border-red-200 bg-red-50 px-3 py-1.5 text-xs font-medium text-red-700 hover:bg-red-100"
                      >
                        Удалить
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
                onClick={() => canGoPrev && setCurrentPage((prev) => prev - 1)}
                className="inline-flex items-center rounded-md border border-gray-200 bg-white px-3 py-1.5 text-xs font-medium text-gray-700 hover:bg-gray-50 disabled:opacity-50 disabled:cursor-not-allowed"
              >
                Назад
              </button>
              <button
                type="button"
                disabled={!canGoNext}
                onClick={() => canGoNext && setCurrentPage((prev) => prev + 1)}
                className="inline-flex items-center rounded-md border border-gray-200 bg-white px-3 py-1.5 text-xs font-medium text-gray-700 hover:bg-gray-50 disabled:opacity-50 disabled:cursor-not-allowed"
              >
                Вперед
              </button>
            </div>
          </div>
        </div>
      </main>

      {isCreateOpen && (
        <div className="fixed inset-0 z-40 flex items-center justify-center bg-black/30 px-4">
          <div className="w-full max-w-md rounded-lg bg-white shadow-xl border border-gray-200">
            <div className="border-b border-gray-100 px-5 py-3 flex items-center justify-between">
              <h2 className="text-sm font-semibold text-gray-900">Новый тест</h2>
              <button
                type="button"
                onClick={() => !isCreating && setIsCreateOpen(false)}
                className="text-xs text-gray-400 hover:text-gray-600"
              >
                Закрыть
              </button>
            </div>

            <form onSubmit={handleSubmitCreate} className="px-5 py-4 space-y-4">
              {createError && (
                <div className="rounded-md bg-red-50 px-3 py-2 text-xs text-red-700 border border-red-100">
                  {createError}
                </div>
              )}

              <div className="space-y-1.5">
                <label
                  htmlFor="new-test-title"
                  className="block text-xs font-medium text-gray-700"
                >
                  Название теста
                </label>
                <input
                  id="new-test-title"
                  type="text"
                  required
                  maxLength={200}
                  value={createForm.title}
                  onChange={(event) =>
                    setCreateForm((prev) => ({ ...prev, title: event.target.value }))
                  }
                  className="w-full rounded-md border border-gray-300 px-3 py-2 text-sm text-gray-900 placeholder-gray-400 focus:outline-none focus:ring-2 focus:ring-gray-900 focus:border-transparent"
                  placeholder="Например, «Тест по продукту для новых сотрудников»"
                />
              </div>

              <div className="space-y-1.5">
                <label
                  htmlFor="new-test-desc"
                  className="block text-xs font-medium text-gray-700"
                >
                  Описание
                </label>
                <textarea
                  id="new-test-desc"
                  rows={4}
                  value={createForm.desc}
                  onChange={(event) =>
                    setCreateForm((prev) => ({ ...prev, desc: event.target.value }))
                  }
                  className="w-full rounded-md border border-gray-300 px-3 py-2 text-sm text-gray-900 placeholder-gray-400 focus:outline-none focus:ring-2 focus:ring-gray-900 focus:border-transparent resize-none"
                  placeholder="Опишите цель теста, целевую аудиторию и формат вопросов."
                />
              </div>

              <div className="mt-2 flex items-center justify-end gap-2">
                <button
                  type="button"
                  disabled={isCreating}
                  onClick={() => !isCreating && setIsCreateOpen(false)}
                  className="inline-flex items-center rounded-md border border-gray-200 bg-white px-4 py-2 text-xs font-medium text-gray-700 hover:bg-gray-50 disabled:opacity-60 disabled:cursor-not-allowed"
                >
                  Отменить
                </button>
                <button
                  type="submit"
                  disabled={isCreating}
                  className="inline-flex items-center rounded-md bg-gray-900 px-4 py-2 text-xs font-medium text-white shadow-sm hover:bg-gray-800 disabled:opacity-60 disabled:cursor-not-allowed"
                >
                  {isCreating ? 'Создаём...' : 'Создать'}
                </button>
              </div>
            </form>
          </div>
        </div>
      )}
    </div>
  );
}

