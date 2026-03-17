import type { FormEvent } from 'react';
import type { CreateTestRequest } from '../../api/types';

interface CreateTestModalProps {
  isOpen: boolean;
  isCreating: boolean;
  error: string | null;
  form: CreateTestRequest;
  onClose: () => void;
  onChange: (value: CreateTestRequest) => void;
  onSubmit: (event: FormEvent) => void;
}

export default function CreateTestModal({
  isOpen,
  isCreating,
  error,
  form,
  onClose,
  onChange,
  onSubmit,
}: CreateTestModalProps) {
  if (!isOpen) return null;

  return (
    <div className="fixed inset-0 z-40 flex items-center justify-center bg-black/30 px-4">
      <div className="w-full max-w-md rounded-lg bg-white shadow-xl border border-gray-200">
        <div className="border-b border-gray-100 px-5 py-3 flex items-center justify-between">
          <h2 className="text-sm font-semibold text-gray-900">Новый тест</h2>
          <button
            type="button"
            onClick={() => !isCreating && onClose()}
            className="text-xs text-gray-400 hover:text-gray-600"
          >
            Закрыть
          </button>
        </div>

        <form onSubmit={onSubmit} className="px-5 py-4 space-y-4">
          {error && (
            <div className="rounded-md bg-red-50 px-3 py-2 text-xs text-red-700 border border-red-100">
              {error}
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
              value={form.title}
              onChange={(event) =>
                onChange({ ...form, title: event.target.value })
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
              value={form.desc}
              onChange={(event) =>
                onChange({ ...form, desc: event.target.value })
              }
              className="w-full rounded-md border border-gray-300 px-3 py-2 text-sm text-gray-900 placeholder-gray-400 focus:outline-none focus:ring-2 focus:ring-gray-900 focus:border-transparent resize-none"
              placeholder="Опишите цель теста, целевую аудиторию и формат вопросов."
            />
          </div>

          <div className="mt-2 flex items-center justify-end gap-2">
            <button
              type="button"
              disabled={isCreating}
              onClick={() => !isCreating && onClose()}
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
  );
}

