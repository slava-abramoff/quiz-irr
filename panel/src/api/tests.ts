import api from './http';
import type {
  CreateTestRequest,
  GetManyTestsResponse,
  TestAdminResponse,
  TestPreviewResponse,
  UpdateTestRequest,
} from './types';

export interface GetTestsParams {
  skip?: number;
  take?: number;
}

export async function createTest(payload: CreateTestRequest): Promise<TestAdminResponse> {
  const { data } = await api.post<TestAdminResponse>('/tests/init', payload);
  return data;
}

export async function getTests(params: GetTestsParams = {}): Promise<GetManyTestsResponse> {
  const { skip = 0, take = 10 } = params;
  const { data } = await api.get<GetManyTestsResponse>('/tests', {
    params: { skip, take },
  });
  return data;
}

export async function getTest(id: string, mode: 'preview' | 'full' | string = 'full'): Promise<TestAdminResponse> {
  const { data } = await api.get<TestAdminResponse>(`/tests/${id}/${mode}`);
  return data;
}

export async function getTestPreview(id: string): Promise<TestPreviewResponse> {
  const { data } = await api.get<TestPreviewResponse>(`/tests/${id}/preview`);
  return data;
}

export async function updateTest(id: string, payload: UpdateTestRequest): Promise<TestAdminResponse> {
  const { data } = await api.patch<TestAdminResponse>(`/tests/${id}`, payload);
  return data;
}

export async function deleteTest(id: string): Promise<string> {
  const { data } = await api.delete<string>(`/tests/${id}`);
  return data;
}

