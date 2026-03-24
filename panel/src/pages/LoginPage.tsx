import { useState } from 'react';
import type { FormEvent } from 'react';
import { useNavigate } from 'react-router-dom';
import { login } from '../api/auth';
import type { LoginRequest } from '../api/types';

export default function LoginPage() {
  const navigate = useNavigate();
  const [form, setForm] = useState<LoginRequest>({ email: '', password: '' });
  const [isSubmitting, setIsSubmitting] = useState(false);
  const [error, setError] = useState<string | null>(null);

  const handleSubmit = async (event: FormEvent) => {
    event.preventDefault();
    setError(null);
    setIsSubmitting(true);

    try {
      const data = await login(form);

      if (typeof window !== 'undefined') {
        localStorage.setItem('access_token', data.access_token);
        localStorage.setItem('refresh_token', data.refresh_token);
        localStorage.setItem('user_full_name', data.full_name);
        localStorage.setItem('user_email', data.email);
      }

      navigate('/', { replace: true });
    } catch (err: any) {
      const message = err?.response?.data?.message ?? 'Не удалось выполнить вход. Проверьте данные.';
      setError(message);
    } finally {
      setIsSubmitting(false);
    }
  };

  return (
    <div className="min-h-screen flex items-center justify-center bg-gray-50">
      <div className="w-full max-w-sm">
        <div className="bg-white border border-gray-200 rounded-lg shadow-sm p-8">
          <h1 className="text-xl font-semibold text-gray-900 text-center mb-6">
            Вход в панель
          </h1>

          {error && (
            <div className="mb-4 rounded-md bg-red-50 px-3 py-2 text-sm text-red-700 border border-red-100">
              {error}
            </div>
          )}

          <form className="space-y-4" onSubmit={handleSubmit}>
            <div>
              <label htmlFor="email" className="block text-sm font-medium text-gray-700 mb-1">
                Email
              </label>
              <input
                id="email"
                type="email"
                placeholder="you@example.com"
                value={form.email}
                onChange={(event) =>
                  setForm((prev) => ({ ...prev, email: event.target.value }))
                }
                className="w-full px-3 py-2 border border-gray-300 rounded-md text-sm text-gray-900 placeholder-gray-400 focus:outline-none focus:ring-2 focus:ring-gray-900 focus:border-transparent transition"
              />
            </div>

            <div>
              <label htmlFor="password" className="block text-sm font-medium text-gray-700 mb-1">
                Пароль
              </label>
              <input
                id="password"
                type="password"
                placeholder="••••••••"
                value={form.password}
                onChange={(event) =>
                  setForm((prev) => ({ ...prev, password: event.target.value }))
                }
                className="w-full px-3 py-2 border border-gray-300 rounded-md text-sm text-gray-900 placeholder-gray-400 focus:outline-none focus:ring-2 focus:ring-gray-900 focus:border-transparent transition"
              />
            </div>

            <button
              type="submit"
              disabled={isSubmitting}
              className="w-full py-2 px-4 bg-gray-900 text-white text-sm font-medium rounded-md hover:bg-gray-800 disabled:opacity-60 disabled:cursor-not-allowed focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-gray-900 transition"
            >
              {isSubmitting ? 'Входим...' : 'Войти'}
            </button>
          </form>
        </div>
      </div>
    </div>
  );
}
