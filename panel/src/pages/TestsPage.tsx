import { useEffect, useState } from 'react';
import type { FormEvent } from 'react';
import { useNavigate } from 'react-router-dom';
import { getTests, createTest } from '../api/tests';
import type { CreateTestRequest, TestAdminResponse } from '../api/types';
import AppHeader from '../components/AppHeader';
import TestsList from '../components/tests/TestsList';
import CreateTestModal from '../components/tests/CreateTestModal';

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

  return (
    <div className="min-h-screen bg-gray-50">
      <AppHeader
        rightSlot={(
          <button
            type="button"
            onClick={handleLogout}
            className="inline-flex items-center rounded-md border border-gray-300 bg-white px-3 py-1.5 text-xs font-medium text-gray-700 shadow-sm hover:bg-gray-50"
          >
            Выйти
          </button>
        )}
      />

      <TestsList
        tests={tests}
        isLoading={isLoading}
        error={error}
        currentPage={currentPage}
        totalPages={totalPages}
        onChangePage={setCurrentPage}
        onOpenCreate={handleCreateTest}
        onOpenTest={(id) => navigate(`/test/${id}`)}
      />

      <CreateTestModal
        isOpen={isCreateOpen}
        isCreating={isCreating}
        error={createError}
        form={createForm}
        onClose={() => setIsCreateOpen(false)}
        onChange={setCreateForm}
        onSubmit={handleSubmitCreate}
      />
    </div>
  );
}

