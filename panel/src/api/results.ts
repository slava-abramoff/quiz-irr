import api from './http';
import type {
  ResultsResponse,
  SimpleMessageResponse,
} from './types';

export interface GetResultsParams {
  skip?: number;
  take?: number;
}

export async function getResultsByTest(
  testId: string,
  params: GetResultsParams = {},
): Promise<ResultsResponse> {
  const { skip = 0, take = 10 } = params;
  const { data } = await api.get<ResultsResponse>(`/results/test/${testId}`, {
    params: { skip, take },
  });
  return data;
}

export async function deleteResult(resultId: number): Promise<SimpleMessageResponse> {
  const { data } = await api.delete<SimpleMessageResponse>(`/results/result/${resultId}`);
  return data;
}

export async function deleteResultsByTest(testId: string): Promise<SimpleMessageResponse> {
  const { data } = await api.delete<SimpleMessageResponse>(`/results/test/${testId}`);
  return data;
}

