// Catálogo de avatares (imagens em /static/avatars) e banners (gradientes).
// Client-safe: usado em componentes e no server. Guardamos só o `id`;
// a imagem/gradiente é resolvida na renderização.

export type AvatarOption = {
  id: string
  src: string
}

export type BannerOption = {
  id: string
  gradient: string
}

export const AVATARS: AvatarOption[] = [
  { id: 'a1', src: '/avatars/avatar001.jpeg' },
  { id: 'a2', src: '/avatars/avatar002.jpeg' },
  { id: 'a3', src: '/avatars/avatar003.jpeg' },
  { id: 'a4', src: '/avatars/avatar004.jpeg' },
  { id: 'a5', src: '/avatars/avatar005.jpeg' },
  { id: 'a6', src: '/avatars/avatar006.jpeg' },
  { id: 'a7', src: '/avatars/avatar007.jpeg' },
  { id: 'a8', src: '/avatars/avatar008.jpeg' },
  { id: 'a9', src: '/avatars/avatar009.png' },
  { id: 'a10', src: '/avatars/avatar010.jpeg' }
]

export const BANNERS: BannerOption[] = [
  { id: 'b1', gradient: 'linear-gradient(120deg, #0e7490, #5b21b6)' }, // teal → roxo
  { id: 'b2', gradient: 'linear-gradient(120deg, #fbbf24, #f97316)' }, // amarelo → laranja
  { id: 'b3', gradient: 'linear-gradient(120deg, #065f46, #0e7490)' }, // esmeralda → teal
  { id: 'b4', gradient: 'linear-gradient(120deg, #84cc16, #16a34a)' }, // lima → verde
  { id: 'b5', gradient: 'linear-gradient(120deg, #1e3a8a, #0369a1)' }, // azul → céu
  { id: 'b6', gradient: 'linear-gradient(120deg, #f472b6, #db2777)' }, // rosa → magenta
  { id: 'b7', gradient: 'linear-gradient(120deg, #1e293b, #0e7490)' }, // ardósia → teal
  { id: 'b8', gradient: 'linear-gradient(120deg, #38bdf8, #06b6d4)' }, // azul-claro → ciano
  { id: 'b9', gradient: 'linear-gradient(120deg, #115e59, #3730a3)' }, // teal → índigo
  { id: 'b10', gradient: 'linear-gradient(120deg, #fb923c, #f43f5e)' } // laranja → rosa (pôr do sol)
]

export const DEFAULT_AVATAR_ID = AVATARS[0].id
export const DEFAULT_BANNER_ID = BANNERS[0].id

const AVATAR_MAP = new Map(AVATARS.map((a) => [a.id, a.src]))
const BANNER_MAP = new Map(BANNERS.map((b) => [b.id, b.gradient]))

/** Caminho da imagem de um avatar (cai no padrão se o id for inválido). */
export function avatarSrc(id?: string): string {
  return (id && AVATAR_MAP.get(id)) || AVATARS[0].src
}

/** Valor CSS `background` para um avatar (imagem coberta e centralizada). */
export function avatarBackground(id?: string): string {
  return `url('${avatarSrc(id)}') center / cover no-repeat`
}

/** Gradiente CSS de um banner (cai no padrão se o id for inválido). */
export function bannerGradient(id?: string): string {
  return (id && BANNER_MAP.get(id)) || BANNERS[0].gradient
}

export function isAvatarId(id: string): boolean {
  return AVATAR_MAP.has(id)
}

export function isBannerId(id: string): boolean {
  return BANNER_MAP.has(id)
}
