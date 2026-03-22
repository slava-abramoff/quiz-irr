/** Moscow wall time as `YYYY-MM-DD HH:mm:ss` → instant (MSK is UTC+3, no DST). */
export function parseMoscowDatetime(s: string): Date {
  const trimmed = s.trim()
  const [datePart, timePart] = trimmed.split(/\s+/)
  if (!datePart || !timePart) {
    return new Date(NaN)
  }
  return new Date(`${datePart}T${timePart}+03:00`)
}

export function formatDurationSeconds(totalSeconds: number): string {
  const s = Math.max(0, Math.floor(totalSeconds))
  const h = Math.floor(s / 3600)
  const m = Math.floor((s % 3600) / 60)
  const sec = s % 60
  return [h, m, sec].map((n) => String(n).padStart(2, '0')).join(':')
}

export function formatLocalDateTime(d: Date): string {
  return new Intl.DateTimeFormat(undefined, {
    dateStyle: 'long',
    timeStyle: 'short',
  }).format(d)
}

export function formatCountdownTo(target: Date, now: Date): string {
  const ms = target.getTime() - now.getTime()
  if (ms <= 0) return '00:00:00'
  const totalSec = Math.floor(ms / 1000)
  const days = Math.floor(totalSec / 86400)
  const h = Math.floor((totalSec % 86400) / 3600)
  const m = Math.floor((totalSec % 3600) / 60)
  const s = totalSec % 60
  const time = [h, m, s].map((n) => String(n).padStart(2, '0')).join(':')
  if (days > 0) {
    return `${days} ${pluralDays(days)} ${time}`
  }
  return time
}

function pluralDays(n: number): string {
  const mod10 = n % 10
  const mod100 = n % 100
  if (mod100 >= 11 && mod100 <= 14) return 'дней'
  if (mod10 === 1) return 'день'
  if (mod10 >= 2 && mod10 <= 4) return 'дня'
  return 'дней'
}
