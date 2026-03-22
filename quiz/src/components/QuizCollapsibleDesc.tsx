import { useState } from 'react'

const DESC_PREVIEW_MAX = 170

function previewDesc(text: string): string {
  if (text.length <= DESC_PREVIEW_MAX) return text
  const slice = text.slice(0, DESC_PREVIEW_MAX)
  const lastSpace = slice.lastIndexOf(' ')
  const end =
    lastSpace > DESC_PREVIEW_MAX - 48 && lastSpace > 0 ? lastSpace : DESC_PREVIEW_MAX
  return `${text.slice(0, end).trimEnd()}…`
}

type QuizCollapsibleDescProps = {
  text: string
}

export function QuizCollapsibleDesc({ text }: QuizCollapsibleDescProps) {
  const [expanded, setExpanded] = useState(false)
  const needsCollapse = text.length > DESC_PREVIEW_MAX

  if (!needsCollapse) {
    return <p className="quiz-desc">{text}</p>
  }

  return (
    <div className="quiz-desc-block">
      <p className="quiz-desc">{expanded ? text : previewDesc(text)}</p>
      <button
        type="button"
        className="quiz-desc-toggle"
        onClick={() => setExpanded((v) => !v)}
        aria-expanded={expanded}
      >
        {expanded ? 'Свернуть' : 'Показать полностью'}
      </button>
    </div>
  )
}
