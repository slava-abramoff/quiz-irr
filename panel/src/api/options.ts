import api from './http';
import type {
  CreateOptionRequest,
  OptionResponse,
  UpdateOptionRequest,
} from './types';

export async function createOption(questionId: number, payload: CreateOptionRequest): Promise<OptionResponse> {
  const { data } = await api.post<OptionResponse>(`/options/question/${questionId}`, payload);
  return data;
}

export async function updateOption(id: number, payload: UpdateOptionRequest): Promise<OptionResponse> {
  const { data } = await api.patch<OptionResponse>(`/options/${id}`, payload);
  return data;
}

export async function deleteOption(id: number): Promise<string> {
  const { data } = await api.delete<string>(`/options/${id}`);
  return data;
}

