import api from './http';
import type {
  CreateQuestionRequest,
  QuestionResponse,
  UpdateQuestionRequest,
} from './types';

export async function createQuestion(testId: string, payload: CreateQuestionRequest): Promise<QuestionResponse> {
  const { data } = await api.post<QuestionResponse>(`/questions/test/${testId}`, payload);
  return data;
}

export async function updateQuestion(id: number, payload: UpdateQuestionRequest): Promise<QuestionResponse> {
  const { data } = await api.patch<QuestionResponse>(`/questions/${id}`, payload);
  return data;
}

export async function deleteQuestion(id: number): Promise<string> {
  const { data } = await api.delete<string>(`/questions/${id}`);
  return data;
}

