import { fail, type Actions } from '@sveltejs/kit'
import type { PageServerLoad } from './$types'
import type { Post, TrendingTopic } from '$lib/types'
import { getOptionalSession } from '$lib/server/auth'
import { api, USE_API, ApiError } from '$lib/server/api'
import { addPost, getTrendingTopics, listPosts } from '$lib/server/mock-data'

const MAX_POST_LENGTH = 500

export const load: PageServerLoad = async ({ fetch, cookies }) => {
  if (USE_API) {
    // Token opcional: faz a API marcar `likedByMe` para o usuário logado.
    const token = getOptionalSession(cookies)?.accessToken

    // Independentes: uma falha no trending (sidebar) não pode esvaziar o feed.
    const [postsResult, trendingResult] = await Promise.allSettled([
      api.posts.list(token, fetch),
      api.trending.get(token, fetch)
    ])

    if (postsResult.status === 'rejected') {
      console.error('Falha ao carregar posts da API:', postsResult.reason)
    }

    return {
      posts: postsResult.status === 'fulfilled' ? postsResult.value : ([] as Post[]),
      topics: trendingResult.status === 'fulfilled' ? trendingResult.value.topics : ([] as TrendingTopic[])
    }
  }

  const session = getOptionalSession(cookies)

  return {
    posts: listPosts(session?.user.username),
    topics: getTrendingTopics()
  }
}

export const actions: Actions = {
  createPost: async ({ request, cookies, fetch }) => {
    const session = getOptionalSession(cookies)

    if (!session) {
      return fail(401, {
        error: 'Você precisa estar logado para publicar.'
      })
    }

    const data = await request.formData()
    const content = String(data.get('content') ?? '').trim()

    if (!content) {
      return fail(400, {
        error: 'Escreva algo antes de publicar.'
      })
    }

    if (content.length > MAX_POST_LENGTH) {
      return fail(400, {
        error: `O post deve ter no máximo ${MAX_POST_LENGTH} caracteres.`
      })
    }

    if (USE_API) {
      try {
        await api.posts.create({ content }, session.accessToken, fetch)
      } catch (error) {
        if (error instanceof ApiError) {
          return fail(error.status === 401 ? 401 : 502, {
            error: 'Não foi possível publicar agora. Tente novamente em instantes.'
          })
        }

        throw error
      }

      return {
        success: true
      }
    }

    // Mock: injeta o post no topo do feed em memória.
    addPost({
      id: crypto.randomUUID(),
      author: {
        name: session.user.name,
        username: session.user.username
      },
      content,
      likes: 0,
      comments: 0,
      createdAt: 'agora'
    })

    return {
      success: true
    }
  }
}
