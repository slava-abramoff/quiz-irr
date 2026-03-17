import type {
  OptionResponse,
  QuestionResponse,
  UpdateOptionRequest,
  UpdateQuestionRequest,
} from '../../api/types';

interface QuestionsEditorProps {
  questions: QuestionResponse[];
  onAddQuestion: () => void;
  onUpdateQuestion: (question: QuestionResponse, partial: UpdateQuestionRequest) => void;
  onDeleteQuestion: (question: QuestionResponse) => void;
  onAddOption: (question: QuestionResponse) => void;
  onUpdateOption: (
    question: QuestionResponse,
    option: OptionResponse,
    partial: UpdateOptionRequest,
  ) => void;
  onDeleteOption: (question: QuestionResponse, option: OptionResponse) => void;
}

export default function QuestionsEditor({
  questions,
  onAddQuestion,
  onUpdateQuestion,
  onDeleteQuestion,
  onAddOption,
  onUpdateOption,
  onDeleteOption,
}: QuestionsEditorProps) {
  return (
    <section className="rounded-lg border border-gray-200 bg-white p-5 space-y-4">
      <div className="flex items-center justify-between gap-3">
        <h2 className="text-sm font-semibold text-gray-900">
          Вопросы и ответы
        </h2>
      </div>

      {questions.length === 0 ? (
        <p className="text-sm text-gray-500">
          Вопросов пока нет. Добавьте первый вопрос, чтобы начать настраивать тест.
        </p>
      ) : (
        <div className="space-y-4">
          {questions.map((question) => (
            <div
              key={question.id}
              className="rounded-md border border-gray-200 bg-gray-50/60 p-4 space-y-3"
            >
              <div className="flex items-start justify-between gap-3">
                <div className="flex-1 space-y-2">
                  <div className="flex items-center gap-2">
                    <input
                      type="text"
                      value={question.text}
                      onChange={(event) =>
                        onUpdateQuestion(question, {
                          text: event.target.value,
                        })
                      }
                      className="w-full rounded-md border border-gray-300 px-3 py-1.5 text-xs text-gray-900 placeholder-gray-400 focus:outline-none focus:ring-2 focus:ring-gray-900 focus:border-transparent"
                    />
                  </div>

                  <div className="flex flex-wrap items-center gap-3 text-xs text-gray-600">
                    <div className="flex items-center gap-1.5">
                      <span className="text-[11px] uppercase tracking-wide text-gray-400">
                        Тип:
                      </span>
                      <select
                        value={question.type}
                        onChange={(event) =>
                          onUpdateQuestion(question, {
                            type: event.target.value,
                          })
                        }
                        className="rounded-md border border-gray-300 bg-white px-2 py-1 text-xs text-gray-900 focus:outline-none focus:ring-2 focus:ring-gray-900 focus:border-transparent"
                      >
                        <option value="single">Один ответ</option>
                        <option value="multiple">Несколько ответов</option>
                        <option value="text">Текстовый ответ</option>
                      </select>
                    </div>

                    <div className="flex items-center gap-1.5">
                      <span className="text-[11px] uppercase tracking-wide text-gray-400">
                        Баллы:
                      </span>
                      <input
                        type="number"
                        min={0}
                        value={question.points}
                        onChange={(event) =>
                          onUpdateQuestion(question, {
                            points: Number.parseInt(event.target.value || '0', 10),
                          })
                        }
                        className="w-16 rounded-md border border-gray-300 px-2 py-1 text-xs text-gray-900 focus:outline-none focus:ring-2 focus:ring-gray-900 focus:border-transparent"
                      />
                    </div>
                  </div>
                </div>

                <button
                  type="button"
                  onClick={() => onDeleteQuestion(question)}
                  className="inline-flex items-center rounded-md border border-red-200 bg-red-50 px-2 py-1 text-[11px] font-medium text-red-700 hover:bg-red-100"
                >
                  Удалить
                </button>
              </div>

              <div className="space-y-2 border-t border-gray-200 pt-3">
                <div className="flex items-center justify-between gap-3">
                  <p className="text-[11px] font-medium uppercase tracking-wide text-gray-400">
                    Варианты ответов
                  </p>
                  <button
                    type="button"
                    onClick={() => onAddOption(question)}
                    className="inline-flex items-center rounded-md border border-gray-300 bg-white px-2 py-1 text-[11px] font-medium text-gray-700 hover:bg-gray-50"
                  >
                    + Вариант
                  </button>
                </div>

                {(question.options?.length ?? 0) === 0 ? (
                  <p className="text-xs text-gray-500">
                    Вариантов пока нет. Добавьте хотя бы один вариант ответа.
                  </p>
                ) : (
                  <div className="space-y-2">
                    {(question.options ?? []).map((option) => (
                      <div
                        key={option.id}
                        className="flex items-center gap-2 rounded-md bg-white px-3 py-2 border border-gray-200"
                      >
                        <button
                          type="button"
                          onClick={() =>
                            onUpdateOption(question, option, {
                              is_correct: !option.is_correct,
                            })
                          }
                          className={`inline-flex h-5 w-5 items-center justify-center rounded border text-[10px] ${
                            option.is_correct
                              ? 'border-emerald-500 bg-emerald-50 text-emerald-700'
                              : 'border-gray-300 bg-white text-gray-400'
                          }`}
                        >
                          ✓
                        </button>
                        <input
                          type="text"
                          value={option.text}
                          onChange={(event) =>
                            onUpdateOption(question, option, {
                              text: event.target.value,
                            })
                          }
                          className="flex-1 rounded-md border border-gray-300 px-2 py-1 text-xs text-gray-900 placeholder-gray-400 focus:outline-none focus:ring-2 focus:ring-gray-900 focus:border-transparent"
                        />
                        <button
                          type="button"
                          onClick={() => onDeleteOption(question, option)}
                          className="inline-flex items-center rounded-md border border-red-200 bg-red-50 px-2 py-1 text-[11px] font-medium text-red-700 hover:bg-red-100"
                        >
                          ×
                        </button>
                      </div>
                    ))}
                  </div>
                )}
              </div>
            </div>
          ))}
        </div>
      )}

      <div className="pt-2">
        <button
          type="button"
          onClick={onAddQuestion}
          className="inline-flex items-center rounded-md bg-gray-900 px-3 py-1.5 text-xs font-medium text-white shadow-sm hover:bg-gray-800"
        >
          + Добавить вопрос
        </button>
      </div>
    </section>
  );
}

