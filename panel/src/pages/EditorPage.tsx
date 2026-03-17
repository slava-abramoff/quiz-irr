import { useParams } from 'react-router-dom';

export default function EditorPage() {
  const { id } = useParams<{ id: string }>();

  return (
    <div className="min-h-screen bg-gray-50">
      <header className="border-b border-gray-200 bg-white">
        <div className="mx-auto max-w-6xl px-4 py-4 flex items-center justify-between gap-4">
          <div>
            <h1 className="text-lg font-semibold text-gray-900">Редактор теста</h1>
            <p className="text-sm text-gray-500">ID теста: {id}</p>
          </div>
        </div>
      </header>

      <main className="mx-auto max-w-6xl px-4 py-6">
        <div className="rounded-lg border border-dashed border-gray-300 bg-white/60 p-8 text-center text-sm text-gray-500">
          Здесь будет интерфейс редактирования теста: общие настройки, вопросы и варианты ответов.
        </div>
      </main>
    </div>
  );
}

