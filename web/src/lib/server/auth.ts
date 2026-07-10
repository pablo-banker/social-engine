import { dev } from '$app/environment'
import { redirect, type Cookies } from '@sveltejs/kit'
import { DEFAULT_AVATAR_ID } from '$lib/appearance'
import { api, USE_API, type AuthResponse } from './api'
import { addMockUser, getMockLoginUser } from './mock-data'

export const SESSION_COOKIE_NAME = 'social_engine_session'

const SESSION_MAX_AGE = 60 * 60 * 24 * 7 // 7 dias

export type AuthUser = {
  id: string
  name: string
  username: string
  email: string
  avatarId: string
}

export type AuthSession = {
  user: AuthUser
  accessToken: string
}

/** Apenas os campos do usuário seguros para expor no client (sem email/token). */
export type PublicUser = Pick<AuthUser, 'id' | 'name' | 'username' | 'avatarId'>

export type PublicSession = {
  user: PublicUser
}

/**
 * Lê a sessão completa (com email/accessToken) do cookie httpOnly.
 * Use somente no server. Retorna null se não houver sessão válida.
 */
export function getOptionalSession(cookies: Cookies): AuthSession | null {
  const rawSession = cookies.get(SESSION_COOKIE_NAME)

  if (!rawSession) {
    return null
  }

  try {
    return JSON.parse(rawSession) as AuthSession
  } catch {
    // Cookie corrompido: limpa para não repetir o erro em toda requisição.
    clearSession(cookies)

    return null
  }
}

/**
 * Retorna apenas os dados públicos da sessão, prontos para ir ao client
 * via `data.session`. Nunca expõe email nem accessToken.
 */
export function getPublicSession(cookies: Cookies): PublicSession | null {
  const session = getOptionalSession(cookies)

  if (!session) {
    return null
  }

  const { id, name, username, avatarId } = session.user

  return {
    user: { id, name, username, avatarId }
  }
}

/**
 * Garante que há uma sessão. Redireciona para /login quando não houver.
 * Use em loads/actions de rotas privadas.
 */
export function requireSession(cookies: Cookies): AuthSession {
  const session = getOptionalSession(cookies)

  if (!session) {
    throw redirect(303, '/login')
  }

  return session
}

/** Grava a sessão no cookie httpOnly com as opções padrão da aplicação. */
export function createSession(cookies: Cookies, session: AuthSession) {
  cookies.set(SESSION_COOKIE_NAME, JSON.stringify(session), {
    path: '/',
    httpOnly: true,
    sameSite: 'lax',
    secure: !dev,
    maxAge: SESSION_MAX_AGE
  })
}

/** Remove a sessão (logout). */
export function clearSession(cookies: Cookies) {
  cookies.delete(SESSION_COOKIE_NAME, {
    path: '/'
  })
}

/** Mantém só os campos de sessão que guardamos no cookie. */
function toSession(response: AuthResponse): AuthSession {
  const { id, name, username, email, avatarId } = response.user

  return {
    user: { id, name, username, email, avatarId: avatarId ?? DEFAULT_AVATAR_ID },
    accessToken: response.accessToken
  }
}

/**
 * Autentica por email/senha. Usa a Go API quando USE_API=true; caso
 * contrário devolve uma sessão mockada. Lança ApiError em falha real.
 */
export async function authenticate(input: {
  email: string
  password: string
}): Promise<AuthSession> {
  if (USE_API) {
    return toSession(await api.auth.login(input))
  }

  // Mock: qualquer email/senha loga como o usuário fixo (reflete edições).
  const mockUser = getMockLoginUser()

  return {
    user: {
      id: mockUser.id,
      name: mockUser.name,
      username: mockUser.username,
      email: input.email,
      avatarId: mockUser.avatarId
    },
    accessToken: 'mock-access-token'
  }
}

/**
 * Cria uma conta. Usa a Go API quando USE_API=true; caso contrário
 * devolve uma sessão mockada. Lança ApiError em falha real.
 */
export async function registerAccount(input: {
  firstName: string
  lastName: string
  email: string
  password: string
}): Promise<AuthSession> {
  if (USE_API) {
    return toSession(await api.auth.register(input))
  }

  // Mock: registra o usuário para que ele tenha perfil editável.
  const name = `${input.firstName} ${input.lastName}`
  const username = input.email.split('@')[0]

  addMockUser({ name, username, email: input.email })

  return {
    user: {
      id: crypto.randomUUID(),
      name,
      username,
      email: input.email,
      avatarId: DEFAULT_AVATAR_ID
    },
    accessToken: 'mock-access-token'
  }
}
