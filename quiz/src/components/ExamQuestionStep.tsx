import { useState } from 'react'
import type { ExamQuestion, UserAnswerPayload } from '../types/examQuiz'

type ExamQuestionStepProps = {
  question: ExamQuestion
  progressText: string
  onAdvance: (answer: UserAnswerPayload) => void
  busy: boolean
  isLast: boolean
  saveError: string | null
  onRetrySave: () => void
}

export function ExamQuestionStep({
  question,
  progressText,
  onAdvance,
  busy,
  isLast,
  saveError,
  onRetrySave,
}: ExamQuestionStepProps) {
  const [singleId, setSingleId] = useState<number | null>(null)
  const [multiIds, setMultiIds] = useState<Set<number>>(() => new Set())
  const [textValue, setTextValue] = useState('')

  function toggleMulti(id: number) {
    setMultiIds((prev) => {
      const next = new Set(prev)
      if (next.has(id)) next.delete(id)
      else next.add(id)
      return next
    })
  }

  function buildAnswer(): UserAnswerPayload | null {
    if (question.type === 'single') {
      if (singleId === null) return null
      return { answer_id: question.id, option_ids: [singleId] }
    }
    if (question.type === 'multiple') {
      if (multiIds.size === 0) return null
      return { answer_id: question.id, option_ids: [...multiIds] }
    }
    const t = textValue.trim()
    if (!t) return null
    return { answer_id: question.id, option_ids: [], text_option: t }
  }

  function handleNext() {
    const a = buildAnswer()
    if (!a || busy) return
    onAdvance(a)
  }

  const canSubmit = buildAnswer() !== null

  return (
    <div className="exam-question">
      <p className="exam-progress">{progressText}</p>
      <h2 className="exam-question-title">{question.text}</h2>

      {question.type === 'single' ? (
        <ul className="exam-options exam-options--single">
          {question.options.map((opt) => (
            <li key={opt.id}>
              <label className="exam-option-label">
                <input
                  type="radio"
                  className="exam-option-input"
                  name={`q-${question.id}`}
                  checked={singleId === opt.id}
                  onChange={() => setSingleId(opt.id)}
                  disabled={busy}
                />
                <span className="exam-option-text">{opt.text}</span>
              </label>
            </li>
          ))}
        </ul>
      ) : null}

      {question.type === 'multiple' ? (
        <ul className="exam-options exam-options--multi">
          {question.options.map((opt) => (
            <li key={opt.id}>
              <label className="exam-option-label">
                <input
                  type="checkbox"
                  className="exam-option-input"
                  checked={multiIds.has(opt.id)}
                  onChange={() => toggleMulti(opt.id)}
                  disabled={busy}
                />
                <span className="exam-option-text">{opt.text}</span>
              </label>
            </li>
          ))}
        </ul>
      ) : null}

      {question.type === 'text' ? (
        <label className="exam-field exam-field--block">
          <span className="exam-field-label">Ваш ответ</span>
          <textarea
            className="exam-textarea"
            rows={4}
            value={textValue}
            onChange={(e) => setTextValue(e.target.value)}
            disabled={busy}
          />
        </label>
      ) : null}

      {saveError ? (
        <div className="exam-save-error" role="alert">
          <p className="exam-save-error-text">{saveError}</p>
          <button
            type="button"
            className="exam-retry-btn"
            onClick={onRetrySave}
            disabled={busy}
          >
            Повторить отправку
          </button>
        </div>
      ) : null}

      <div className="exam-question-actions">
        <button
          type="button"
          className="quiz-start-btn"
          onClick={handleNext}
          disabled={!canSubmit || busy}
        >
          {busy
            ? isLast
              ? 'Отправка…'
              : 'Сохранение…'
            : isLast
              ? 'Завершить'
              : 'Далее'}
        </button>
      </div>
    </div>
  )
}
