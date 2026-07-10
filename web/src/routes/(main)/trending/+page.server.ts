import type { PageServerLoad } from './$types'
import type { Post, TrendingTopic } from '$lib/types'
import { getOptionalSession } from '$lib/server/auth'
import { api, USE_API } from '$lib/server/api'
import { getTrendingTopics, listPopularPosts } from '$lib/server/mock-data'

export const load: PageServerLoad = async ({ fetch, cookies }) => {
  if (USE_API) {
    try {
      const { topics, posts } = await api.trending.get(fetch)

      return { topics, posts }
    } catch (error) {
      console.error('Falha ao carregar o trending da API:', error)

      return {
        topics: [] as TrendingTopic[],
        posts: [] as Post[]
      }
    }
  }

  const viewer = getOptionalSession(cookies)?.user.username

  return {
    topics: getTrendingTopics(),
    posts: listPopularPosts(viewer)
  }
}
