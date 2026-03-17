import api from './http';
import type {
  SendUserAnswersRequest,
  SendUserAnswersResponse,
  StartExamRequest,
  StartExamResponse,
  TestCustomerResponse,
} from './types';

export async function getExamInfo(testId: string): Promise<TestCustomerResponse> {
  const { data } = await api.get<TestCustomerResponse>(`/exam/info/${testId}`);
  return data;
}

export async function startExam(testId: string, payload: StartExamRequest): Promise<StartExamResponse> {
  const { data } = await api.post<StartExamResponse>(`/exam/start/${testId}`, payload);
  return data;
}

export async function saveExamAnswers(rawId: string, payload: SendUserAnswersRequest): Promise<SendUserAnswersResponse> {
  const { data } = await api.post<SendUserAnswersResponse>(`/exam/save/${rawId}`, payload);
  return data;
}

