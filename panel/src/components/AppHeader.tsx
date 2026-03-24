import type { ReactNode } from 'react';
import { useNavigate } from 'react-router-dom';

interface AppHeaderProps {
  title?: string;
  subtitle?: string;
  rightSlot?: ReactNode;
}

export default function AppHeader({
  title = 'Конструктор тестов',
  subtitle = 'Управляйте тестами, вопросами и результатами',
  rightSlot,
}: AppHeaderProps) {
  const navigate = useNavigate();

  return (
    <header className="border-b border-gray-200 bg-white">
      <div className="mx-auto max-w-6xl px-4 py-4 flex items-center justify-between gap-4">
        <div>
          <button
            type="button"
            onClick={() => navigate('/')}
            className="text-left"
          >
            <h1 className="text-lg font-semibold text-gray-900">{title}</h1>
            {subtitle && (
              <p className="text-sm text-gray-500">
                {subtitle}
              </p>
            )}
          </button>
        </div>

        <div className="flex items-center gap-3">
          {rightSlot}
        </div>
      </div>
    </header>
  );
}

