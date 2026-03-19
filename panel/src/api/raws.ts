import api from './http';
import type {
  AnalyzeRawResponse,
  RawsInfoResponse,
  SimpleMessageResponse,
} from './types';

export interface GetRawResultsParams {
  skip?: number;
  take?: number;
}

export async function getRawResultsByTest(
  testId: string,
  params: GetRawResultsParams = {},
): Promise<RawsInfoResponse> {
  const { skip = 0, take = 10 } = params;
  const { data } = await api.get<RawsInfoResponse>(`/raws/test/${testId}`, {
    params: { skip, take },
  });
  return data;
}

export async function analyzeRaw(rawId: string): Promise<AnalyzeRawResponse> {
  const { data } = await api.post<AnalyzeRawResponse>(`/raws/analyze/${rawId}`);
  return data;
}

export async function deleteRawAnswers(rawId: string): Promise<SimpleMessageResponse> {
  const { data } = await api.delete<SimpleMessageResponse>(`/raws/answers/${rawId}`);
  return data;
}

export async function deleteAllRawByTest(testId: string): Promise<SimpleMessageResponse> {
  const { data } = await api.delete<SimpleMessageResponse>(`/raws/test/${testId}`);
  return data;
}

export async function analyzeAllRawByTest(
  testId: string,
): Promise<SimpleMessageResponse> {
  const { data } = await api.post<SimpleMessageResponse>(
    `/raws/test/analyze/${testId}`,
  );
  return data;
}

