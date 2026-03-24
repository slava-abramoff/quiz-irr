import axios, { AxiosHeaders } from 'axios';
import type { RefreshTokenRequest, RefreshTokenResponse } from './types';

const api = axios.create({
  baseURL: import.meta.env.VITE_API_BASE_URL ?? '/api',
});

api.interceptors.request.use((config) => {
  const token = typeof window !== 'undefined' ? localStorage.getItem('access_token') : null;

  if (token) {
    const headers = AxiosHeaders.from(config.headers);
    headers.set('Authorization', `Bearer ${token}`);
    // eslint-disable-next-line no-param-reassign
    config.headers = headers;
  }

  return config;
});

let isRefreshing = false;
let refreshPromise: Promise<string> | null = null;

async function refreshAccessToken(): Promise<string> {
  if (refreshPromise) {
    return refreshPromise;
  }

  const refreshToken = typeof window !== 'undefined' ? localStorage.getItem('refresh_token') : null;

  if (!refreshToken) {
    throw new Error('No refresh token');
  }

  const body: RefreshTokenRequest = {
    refresh_token: refreshToken,
  };

  refreshPromise = api
    .post<RefreshTokenResponse>('/auth/refresh', body)
    .then((res) => {
      const newAccessToken = res.data.access_token;

      if (typeof window !== 'undefined') {
        localStorage.setItem('access_token', newAccessToken);
      }

      return newAccessToken;
    })
    .finally(() => {
      isRefreshing = false;
      refreshPromise = null;
    });

  return refreshPromise;
}

api.interceptors.response.use(
  (response) => response,
  async (error) => {
    const originalRequest = error.config;

    if (!originalRequest || originalRequest._retry) {
      return Promise.reject(error);
    }

    if (error.response?.status === 401 && !isRefreshing) {
      originalRequest._retry = true;
      isRefreshing = true;

      try {
        const newAccessToken = await refreshAccessToken();

        const headers = AxiosHeaders.from(originalRequest.headers);
        headers.set('Authorization', `Bearer ${newAccessToken}`);
        // eslint-disable-next-line no-param-reassign
        originalRequest.headers = headers;

        return api(originalRequest);
      } catch (refreshError) {
        return Promise.reject(refreshError);
      }
    }

    return Promise.reject(error);
  },
);

export default api;

