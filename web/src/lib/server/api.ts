import { env } from '$env/dynamic/private'
import type {
  LikeResult,
  Post,
  PostComment,
  Profile,
  ProfileSettings,
  SuggestedUser,
  TrendingTopic
} from '$lib/types'

/**
 * Camada HTTP tipada para a futura Go API.
 *
 * Enquanto `USE_API` for `false` (padrão), este módulo não é chamado: os
 * callers (auth/feed) usam mocks. Ao subir o backend Go, basta setar
 * `USE_API=true` e `API_BASE_URL` no ambiente — as páginas não mudam.
 */

/** Liga/desliga o backend real. Controlado por env (default: mocks). */
export const USE_API = env.USE_API === 'true'

/** Base URL da Go API, sem barra final. */
const API_BASE_URL = (env.API_BASE_URL ?? 'http://localhost:8080').replace(/\/+$/, '')

type FetchLike = typeof fetch

type RequestOptions = {
  method?: 'GET' | 'POST' | 'PATCH' | 'PUT' | 'DELETE'
  body?: unknown
  /** Bearer token (accessToken da sessão) para rotas autenticadas. */
  token?: string
  /** `event.fetch` do SvelteKit; cai para o `fetch` global se omitido. */
  fetch?: FetchLike
  signal?: AbortSignal
}

/** Erro normalizado da API — carrega o status HTTP e um código opcional. */
export class ApiError extends Error {
  readonly status: number
  readonly code?: string

  constructor(message: string, status: number, code?: string) {
    super(message)
    this.name = 'ApiError'
    this.status = status
    this.code = code
  }
}

function isRecord(value: unknown): value is Record<string, unknown> {
  return typeof value === 'object' && value !== null
}

function safeJsonParse(raw: string): unknown {
  try {
    return JSON.parse(raw)
  } catch {
    return null
  }
}

async function apiFetch<T>(path: string, options: RequestOptions = {}): Promise<T> {
  const { method = 'GET', body, token, fetch: fetchImpl = fetch, signal } = options

  const headers: Record<string, string> = {
    accept: 'application/json'
  }

  if (body !== undefined) {
    headers['content-type'] = 'application/json'
  }

  if (token) {
    headers.authorization = `Bearer ${token}`
  }

  let response: Response

  try {
    response = await fetchImpl(`${API_BASE_URL}${path}`, {
      method,
      headers,
      body: body === undefined ? undefined : JSON.stringify(body),
      signal
    })
  } catch {
    // Falha de rede/DNS/conexão recusada — a API está fora do ar.
    throw new ApiError('Não foi possível conectar à API do Social Engine.', 503)
  }

  const raw = await response.text()
  const payload = raw ? safeJsonParse(raw) : null

  if (!response.ok) {
    const message =
      (isRecord(payload) && typeof payload.message === 'string' && payload.message) ||
      'A API retornou um erro inesperado.'
    const code = isRecord(payload) && typeof payload.code === 'string' ? payload.code : undefined

    throw new ApiError(message, response.status, code)
  }

  return payload as T
}

/* ---------------------------------------------------------------- *
 * Contratos (DTOs) — o formato esperado das respostas da Go API.
 * Espelham os mocks para que as páginas não precisem mudar.
 * ---------------------------------------------------------------- */

export type ApiUser = {
  id: string
  name: string
  username: string
  email: string
  avatarId?: string
}

export type AuthResponse = {
  user: ApiUser
  accessToken: string
}

export type LoginInput = {
  email: string
  password: string
}

export type RegisterInput = {
  firstName: string
  lastName: string
  email: string
  password: string
}

export type CreatePostInput = {
  content: string
}

export type CreateCommentInput = {
  content: string
}

export type UpdateProfileInput = ProfileSettings

/* ---------------------------------------------------------------- *
 * Endpoints tipados.
 * ---------------------------------------------------------------- */

export const api = {
  auth: {
    login(input: LoginInput, fetchImpl?: FetchLike) {
      return apiFetch<AuthResponse>('/auth/login', {
        method: 'POST',
        body: input,
        fetch: fetchImpl
      })
    },

    register(input: RegisterInput, fetchImpl?: FetchLike) {
      return apiFetch<AuthResponse>('/auth/register', {
        method: 'POST',
        body: input,
        fetch: fetchImpl
      })
    }
  },

  posts: {
    list(fetchImpl?: FetchLike) {
      return apiFetch<Post[]>('/posts', {
        fetch: fetchImpl
      })
    },

    get(id: string, fetchImpl?: FetchLike) {
      return apiFetch<Post>(`/posts/${encodeURIComponent(id)}`, {
        fetch: fetchImpl
      })
    },

    create(input: CreatePostInput, token: string, fetchImpl?: FetchLike) {
      return apiFetch<Post>('/posts', {
        method: 'POST',
        body: input,
        token,
        fetch: fetchImpl
      })
    },

    toggleLike(id: string, token: string, fetchImpl?: FetchLike) {
      return apiFetch<LikeResult>(`/posts/${encodeURIComponent(id)}/like`, {
        method: 'POST',
        token,
        fetch: fetchImpl
      })
    },

    listComments(id: string, fetchImpl?: FetchLike) {
      return apiFetch<PostComment[]>(`/posts/${encodeURIComponent(id)}/comments`, {
        fetch: fetchImpl
      })
    },

    addComment(id: string, input: CreateCommentInput, token: string, fetchImpl?: FetchLike) {
      return apiFetch<PostComment>(`/posts/${encodeURIComponent(id)}/comments`, {
        method: 'POST',
        body: input,
        token,
        fetch: fetchImpl
      })
    }
  },

  users: {
    getByUsername(username: string, fetchImpl?: FetchLike) {
      return apiFetch<Profile>(`/users/${encodeURIComponent(username)}`, {
        fetch: fetchImpl
      })
    },

    listPosts(username: string, fetchImpl?: FetchLike) {
      return apiFetch<Post[]>(`/users/${encodeURIComponent(username)}/posts`, {
        fetch: fetchImpl
      })
    }
  },

  explore: {
    get(fetchImpl?: FetchLike) {
      return apiFetch<{ users: SuggestedUser[]; posts: Post[] }>('/explore', {
        fetch: fetchImpl
      })
    },

    byTag(tag: string, fetchImpl?: FetchLike) {
      return apiFetch<Post[]>(`/explore?tag=${encodeURIComponent(tag)}`, {
        fetch: fetchImpl
      })
    }
  },

  trending: {
    get(fetchImpl?: FetchLike) {
      return apiFetch<{ topics: TrendingTopic[]; posts: Post[] }>('/trending', {
        fetch: fetchImpl
      })
    }
  },

  me: {
    getSettings(token: string, fetchImpl?: FetchLike) {
      return apiFetch<ProfileSettings>('/me', {
        token,
        fetch: fetchImpl
      })
    },

    updateProfile(input: UpdateProfileInput, token: string, fetchImpl?: FetchLike) {
      return apiFetch<ApiUser>('/me', {
        method: 'PATCH',
        body: input,
        token,
        fetch: fetchImpl
      })
    }
  }
}
