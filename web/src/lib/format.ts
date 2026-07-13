// Formatação de datas para exibição.
//
// A Go API envia datas em ISO 8601 (RFC3339). Os mocks já mandam rótulos
// prontos ("agora", "julho de 2026"). Estes helpers formatam apenas o que
// parece ISO e devolvem qualquer outra string inalterada — assim os mesmos
// componentes funcionam com a API e com os mocks.

const MONTHS_PT = [
  'janeiro',
  'fevereiro',
  'março',
  'abril',
  'maio',
  'junho',
  'julho',
  'agosto',
  'setembro',
  'outubro',
  'novembro',
  'dezembro'
]

function parseISO(value: string | null | undefined): Date | null {
  if (!value || !/^\d{4}-\d{2}-\d{2}T/.test(value)) {
    return null
  }
  const date = new Date(value)
  return Number.isNaN(date.getTime()) ? null : date
}

/** Tempo relativo curto para posts/comentários (ex.: "agora", "2 h", "3 d"). */
export function relativeTime(value: string | null | undefined, now: Date = new Date()): string {
  const date = parseISO(value)
  if (!date) {
    return value ?? ''
  }

  const seconds = Math.max(0, Math.floor((now.getTime() - date.getTime()) / 1000))
  if (seconds < 60) return 'agora'

  const minutes = Math.floor(seconds / 60)
  if (minutes < 60) return `${minutes} min`

  const hours = Math.floor(minutes / 60)
  if (hours < 24) return `${hours} h`

  const days = Math.floor(hours / 24)
  if (days < 7) return `${days} d`

  return `${date.getDate()} de ${MONTHS_PT[date.getMonth()]}`
}

/** Mês e ano para a data de entrada no perfil (ex.: "julho de 2026"). */
export function joinedLabel(value: string | null | undefined): string {
  const date = parseISO(value)
  if (!date) {
    return value ?? ''
  }
  return `${MONTHS_PT[date.getMonth()]} de ${date.getFullYear()}`
}
