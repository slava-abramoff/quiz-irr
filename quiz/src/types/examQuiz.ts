export type QuestionType = 'single' | 'multiple' | 'text'

export type ExamOption = {
  id: number
  text: string
}

export type ExamQuestion = {
  id: number
  text: string
  type: QuestionType
  options: ExamOption[]
}

export type ExamStartPayload = {
  full_name: string
  email: string
  org: string
}

export type ExamStartedResponse = {
  raw_id: string
  questions: ExamQuestion[]
}

export type UserAnswerPayload = {
  answer_id: number
  option_ids: number[]
  text_option?: string
}

export type ExamSavePayload = {
  answers: UserAnswerPayload[]
}
