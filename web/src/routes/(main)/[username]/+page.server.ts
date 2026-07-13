import { error } from '@sveltejs/kit'
import type { PageServerLoad } from './$types'
import { getOptionalSession } from '$lib/server/auth'
import { api, USE_API, ApiError } from '$lib/server/api'
import { getProfile, listPostsByUsername } from '$lib/server/mock-data'

export const load: PageServerLoad = async ({ params, fetch, cookies }) => {
  const username = params.username

  if (USE_API) {
    // Token opcional: faz a API marcar `likedByMe` nos posts do perfil.
    const token = getOptionalSession(cookies)?.accessToken

    try {
      const [profile, posts] = await Promise.all([
        api.users.getByUsername(username, fetch),
        api.users.listPosts(username, token, fetch)
      ])

      return { profile, posts }
    } catch (err) {
      if (err instanceof ApiError && err.status === 404) {
        throw error(404, 'Perfil não encontrado.')
      }

      throw err
    }
  }

  const profile = getProfile(username)

  if (!profile) {
    throw error(404, 'Perfil não encontrado.')
  }

  const session = getOptionalSession(cookies)

  return {
    profile,
    posts: listPostsByUsername(username, session?.user.username)
  }
}
