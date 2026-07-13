import { error, fail } from '@sveltejs/kit'
import type { Actions, PageServerLoad } from './$types'
import { getOptionalSession } from '$lib/server/auth'
import { api, USE_API, ApiError } from '$lib/server/api'
import { addComment, getPost, listComments, toggleLike } from '$lib/server/mock-data'

const MAX_COMMENT_LENGTH = 300

export const load: PageServerLoad = async ({ params, fetch, cookies }) => {
  const id = params.id

  if (USE_API) {
    // Token opcional: faz a API marcar `likedByMe` para o usuário logado.
    const token = getOptionalSession(cookies)?.accessToken

    try {
      const [post, comments] = await Promise.all([
        api.posts.get(id, token, fetch),
        api.posts.listComments(id, fetch)
      ])

      return { post, comments }
    } catch (err) {
      if (err instanceof ApiError && err.status === 404) {
        throw error(404, 'Post não encontrado.')
      }

      throw err
    }
  }

  const session = getOptionalSession(cookies)
  const post = getPost(id, session?.user.username)

  if (!post) {
    throw error(404, 'Post não encontrado.')
  }

  return {
    post,
    comments: listComments(id)
  }
}

export const actions: Actions = {
  toggleLike: async ({ params, cookies, fetch }) => {
    const session = getOptionalSession(cookies)

    if (!session) {
      return fail(401, {
        error: 'Você precisa estar logado para curtir.'
      })
    }

    if (USE_API) {
      try {
        return await api.posts.toggleLike(params.id, session.accessToken, fetch)
      } catch (err) {
        if (err instanceof ApiError) {
          return fail(err.status === 401 ? 401 : 502, {
            error: 'Não foi possível curtir agora. Tente novamente em instantes.'
          })
        }

        throw err
      }
    }

    const result = toggleLike(params.id, session.user.username)

    if (!result) {
      throw error(404, 'Post não encontrado.')
    }

    return result
  },

  addComment: async ({ request, params, cookies, fetch }) => {
    const session = getOptionalSession(cookies)

    if (!session) {
      return fail(401, {
        error: 'Você precisa estar logado para comentar.'
      })
    }

    const data = await request.formData()
    const content = String(data.get('content') ?? '').trim()

    if (!content) {
      return fail(400, {
        error: 'Escreva um comentário.'
      })
    }

    if (content.length > MAX_COMMENT_LENGTH) {
      return fail(400, {
        error: `O comentário deve ter no máximo ${MAX_COMMENT_LENGTH} caracteres.`
      })
    }

    if (USE_API) {
      try {
        await api.posts.addComment(params.id, { content }, session.accessToken, fetch)
      } catch (err) {
        if (err instanceof ApiError) {
          return fail(err.status === 401 ? 401 : 502, {
            error: 'Não foi possível comentar agora. Tente novamente em instantes.'
          })
        }

        throw err
      }

      return {
        success: true
      }
    }

    // Mock: injeta o comentário e incrementa o contador do post.
    addComment(params.id, {
      id: crypto.randomUUID(),
      author: {
        name: session.user.name,
        username: session.user.username
      },
      content,
      createdAt: 'agora'
    })

    return {
      success: true
    }
  }
}
