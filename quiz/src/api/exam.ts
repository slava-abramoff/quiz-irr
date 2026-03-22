import type {
  ExamSavePayload,
  ExamStartPayload,
  ExamStartedResponse,
  QuestionType,
} from '../types/examQuiz'
import type { TestCustomerResponse } from '../types'

function apiBase(): string {
  return import.meta.env.VITE_API_BASE ?? ''
}

async function readErrorMessage(res: Response): Promise<string> {
  const fallback = `Запрос не выполнен (${res.status})`
  try {
    const text = await res.text()
    if (!text) return fallback
    try {
      const j: unknown = JSON.parse(text)
      if (j && typeof j === 'object') {
        const o = j as Record<string, unknown>
        if (typeof o.message === 'string' && o.message.trim()) return o.message
        if (typeof o.error === 'string' && o.error.trim()) return o.error
        if (typeof o.detail === 'string' && o.detail.trim()) return o.detail
      }
    } catch {
      if (text.trim().length < 200) return text.trim()
    }
  } catch {
    /* ignore */
  }
  return fallback
}

function apiUrl(path: string): string {
  const base = apiBase().replace(/\/$/, '')
  return `${base}${path.startsWith('/') ? path : `/${path}`}`
}

export async function fetchExamInfo(testId: string): Promise<TestCustomerResponse> {
  const url = apiUrl(`/api/exam/info/${encodeURIComponent(testId)}`)
  const res = await fetch(url)
  if (!res.ok) {
    throw new Error('unavailable')
  }
  const data: unknown = await res.json()
  if (!data || typeof data !== 'object') {
    throw new Error('invalid')
  }
  const o = data as Record<string, unknown>
  if (
    typeof o.title !== 'string' ||
    typeof o.desc !== 'string' ||
    typeof o.duration !== 'number' ||
    typeof o.start_at !== 'string' ||
    typeof o.end_at !== 'string'
  ) {
    throw new Error('invalid')
  }
  return {
    title: o.title,
    desc: o.desc,
    duration: o.duration,
    start_at: o.start_at,
    end_at: o.end_at,
  }
}

function isQuestionType(v: unknown): v is QuestionType {
  return v === 'single' || v === 'multiple' || v === 'text'
}

function parseExamStarted(data: unknown): ExamStartedResponse {
  if (!data || typeof data !== 'object') {
    throw new Error('invalid')
  }
  const o = data as Record<string, unknown>
  if (typeof o.raw_id !== 'string' || !Array.isArray(o.questions)) {
    throw new Error('invalid')
  }
  const questions: ExamStartedResponse['questions'] = []
  for (const q of o.questions) {
    if (!q || typeof q !== 'object') throw new Error('invalid')
    const qo = q as Record<string, unknown>
    if (
      typeof qo.id !== 'number' ||
      typeof qo.text !== 'string' ||
      !isQuestionType(qo.type) ||
      !Array.isArray(qo.options)
    ) {
      throw new Error('invalid')
    }
    const options: ExamStartedResponse['questions'][0]['options'] = []
    for (const opt of qo.options) {
      if (!opt || typeof opt !== 'object') throw new Error('invalid')
      const oo = opt as Record<string, unknown>
      if (typeof oo.id !== 'number' || typeof oo.text !== 'string') {
        throw new Error('invalid')
      }
      options.push({ id: oo.id, text: oo.text })
    }
    questions.push({
      id: qo.id,
      text: qo.text,
      type: qo.type,
      options,
    })
  }
  return { raw_id: o.raw_id, questions }
}

export async function startExam(
  testId: string,
  payload: ExamStartPayload,
): Promise<ExamStartedResponse> {
  const url = apiUrl(`/api/exam/start/${encodeURIComponent(testId)}`)
  const res = await fetch(url, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify(payload),
  })
  if (!res.ok) {
    const msg = await readErrorMessage(res)
    throw new Error(msg)
  }
  const data: unknown = await res.json()
  return parseExamStarted(data)
}

export async function saveExamAnswers(rawId: string, payload: ExamSavePayload): Promise<void> {
  const url = apiUrl(`/api/exam/save/${encodeURIComponent(rawId)}`)
  const res = await fetch(url, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify(payload),
  })
  if (!res.ok) {
    const msg = await readErrorMessage(res)
    throw new Error(msg)
  }
}
