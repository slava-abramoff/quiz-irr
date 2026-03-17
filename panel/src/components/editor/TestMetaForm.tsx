import type { ChangeEvent, FormEvent } from 'react';

interface TestMetaFormProps {
  title: string;
  desc: string;
  startAt: string;
  endAt: string;
  isActive: boolean;
  durationHours: number;
  durationMinutes: number;
  durationSeconds: number;
  isSaving: boolean;
  onChangeTitle: (value: string) => void;
  onChangeDesc: (value: string) => void;
  onChangeStartAt: (value: string) => void;
  onChangeEndAt: (value: string) => void;
  onToggleActive: () => void;
  onDurationChange: (field: 'h' | 'm' | 's', value: number) => void;
  onSubmit: (event: FormEvent) => void;
}

export default function TestMetaForm({
  title,
  desc,
  startAt,
  endAt,
  isActive,
  durationHours,
  durationMinutes,
  durationSeconds,
  isSaving,
  onChangeTitle,
  onChangeDesc,
  onChangeStartAt,
  onChangeEndAt,
  onToggleActive,
  onDurationChange,
  onSubmit,
}: TestMetaFormProps) {
  const handleDurationInput =
    (field: 'h' | 'm' | 's') =>
    (event: ChangeEvent<HTMLInputElement>) => {
      const value = Number.parseInt(event.target.value || '0', 10);
      const safe = Number.isNaN(value) ? 0 : Math.max(0, value);
      const clamped =
        field === 'h' ? safe : Math.min(59, safe);
      onDurationChange(field, clamped);
    };

  return (
    <section className="rounded-lg border border-gray-200 bg-white p-5 space-y-4">
      <h2 className="text-sm font-semibold text-gray-900">Общие настройки теста</h2>

      <form onSubmit={onSubmit} className="space-y-4">
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
              onChange={(event) => onChangeTitle(event.target.value)}
              className="w-full rounded-md border border-gray-300 px-3 py-2 text-sm text-gray-900 placeholder-gray-400 focus:outline-none focus:ring-2 focus:ring-gray-900 focus:border-transparent"
            />
          </div>

          <div className="space-y-1.5">
            <span className="block text-xs font-medium text-gray-700">
              Статус
            </span>
            <button
              type="button"
              onClick={onToggleActive}
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
            onChange={(event) => onChangeDesc(event.target.value)}
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
              onChange={(event) => onChangeStartAt(event.target.value)}
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
              onChange={(event) => onChangeEndAt(event.target.value)}
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
                onChange={handleDurationInput('h')}
                className="w-16 rounded-md border border-gray-300 px-2 py-1 text-xs text-gray-900 focus:outline-none focus:ring-2 focus:ring-gray-900 focus:border-transparent"
              />
              <span className="text-xs text-gray-500">ч</span>
              <input
                type="number"
                min={0}
                max={59}
                value={durationMinutes}
                onChange={handleDurationInput('m')}
                className="w-16 rounded-md border border-gray-300 px-2 py-1 text-xs text-gray-900 focus:outline-none focus:ring-2 focus:ring-gray-900 focus:border-transparent"
              />
              <span className="text-xs text-gray-500">мин</span>
              <input
                type="number"
                min={0}
                max={59}
                value={durationSeconds}
                onChange={handleDurationInput('s')}
                className="w-16 rounded-md border border-gray-300 px-2 py-1 text-xs text-gray-900 focus:outline-none focus:ring-2 focus:ring-gray-900 focus:border-transparent"
              />
              <span className="text-xs text-gray-500">сек</span>
            </div>
          </div>
        </div>

        <div className="flex items-center justify-end gap-2 pt-2">
          <button
            type="submit"
            disabled={isSaving}
            className="inline-flex items-center rounded-md bg-gray-900 px-4 py-2 text-xs font-medium text-white shadow-sm hover:bg-gray-800 disabled:opacity-60 disabled:cursor-not-allowed"
          >
            {isSaving ? 'Сохраняем...' : 'Сохранить изменения'}
          </button>
        </div>
      </form>
    </section>
  );
}

