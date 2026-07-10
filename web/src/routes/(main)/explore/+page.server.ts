import type { PageServerLoad } from './$types'
import type { Post, SuggestedUser, TrendingTopic } from '$lib/types'
import { getOptionalSession } from '$lib/server/auth'
import { api, USE_API } from '$lib/server/api'
import {
  getTrendingTopics,
  listPopularPosts,
  listPostsByTopic,
  listSuggestedUsers
} from '$lib/server/mock-data'

export const load: PageServerLoad = async ({ url, fetch, cookies }) => {
  const tag = url.searchParams.get('tag')?.toLowerCase().replace(/^#/, '') || null
  const viewer = getOptionalSession(cookies)?.user.username

  if (USE_API) {
    // Best-effort e independente do conteúdo principal (só alimenta a sidebar).
    const topicsPromise = api.trending
      .get(fetch)
      .then((trending) => trending.topics)
      .catch(() => [] as TrendingTopic[])

    try {
      const posts = tag ? await api.explore.byTag(tag, fetch) : null
      const explore = tag ? null : await api.explore.get(fetch)

      return {
        tag,
        users: explore?.users ?? ([] as SuggestedUser[]),
        posts: posts ?? explore?.posts ?? ([] as Post[]),
        topics: await topicsPromise
      }
    } catch (error) {
      console.error('Falha ao carregar o explore da API:', error)

      return {
        tag,
        users: [] as SuggestedUser[],
        posts: [] as Post[],
        topics: await topicsPromise
      }
    }
  }

  const topics = getTrendingTopics()

  if (tag) {
    return { tag, users: [] as SuggestedUser[], posts: listPostsByTopic(tag, viewer), topics }
  }

  return {
    tag: null,
    users: listSuggestedUsers(viewer),
    posts: listPopularPosts(viewer),
    topics
  }
}
