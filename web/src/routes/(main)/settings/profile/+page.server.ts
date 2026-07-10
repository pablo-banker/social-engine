import { fail, redirect } from '@sveltejs/kit'
import type { Actions, PageServerLoad } from './$types'
import { DEFAULT_BANNER_ID, isAvatarId, isBannerId } from '$lib/appearance'
import type { ProfileSettings } from '$lib/types'
import { createSession, requireSession } from '$lib/server/auth'
import { api, USE_API, ApiError } from '$lib/server/api'
import { getUserSettings, updateUserSettings } from '$lib/server/mock-data'

const MAX_NAME_LENGTH = 50
const MAX_BIO_LENGTH = 160

export const load: PageServerLoad = async ({ cookies, fetch }) => {
  const session = requireSession(cookies)

  let settings: ProfileSettings | null = null

  if (USE_API) {
    try {
      settings = await api.me.getSettings(session.accessToken, fetch)
    } catch (error) {
      console.error('Falha ao carregar settings da API:', error)
    }
  } else {
    settings = getUserSettings(session.user.username)
  }

  // Fallback a partir da sessão se o backend não devolver nada.
  return {
    settings: settings ?? {
      name: session.user.name,
      bio: '',
      avatarId: session.user.avatarId,
      bannerId: DEFAULT_BANNER_ID
    }
  }
}

export const actions: Actions = {
  save: async ({ request, cookies, fetch }) => {
    const session = requireSession(cookies)

    const data = await request.formData()
    const name = String(data.get('name') ?? '').trim()
    const bio = String(data.get('bio') ?? '').trim()
    const avatarId = String(data.get('avatarId') ?? '')
    const bannerId = String(data.get('bannerId') ?? '')

    const values = { name, bio, avatarId, bannerId }

    if (!name) {
      return fail(400, { error: 'Informe seu nome.', values })
    }

    if (name.length > MAX_NAME_LENGTH) {
      return fail(400, { error: `O nome deve ter no máximo ${MAX_NAME_LENGTH} caracteres.`, values })
    }

    if (bio.length > MAX_BIO_LENGTH) {
      return fail(400, { error: `A bio deve ter no máximo ${MAX_BIO_LENGTH} caracteres.`, values })
    }

    if (!isAvatarId(avatarId) || !isBannerId(bannerId)) {
      return fail(400, { error: 'Selecione um avatar e um banner válidos.', values })
    }

    const input: ProfileSettings = { name, bio, avatarId, bannerId }

    if (USE_API) {
      try {
        await api.me.updateProfile(input, session.accessToken, fetch)
      } catch (error) {
        if (error instanceof ApiError) {
          return fail(502, { error: 'Não foi possível salvar agora. Tente novamente.', values })
        }

        throw error
      }
    } else if (!updateUserSettings(session.user.username, input)) {
      return fail(404, { error: 'Perfil não encontrado.', values })
    }

    // Revalida a sessão para o header refletir o novo nome/avatar.
    createSession(cookies, {
      user: { ...session.user, name, avatarId },
      accessToken: session.accessToken
    })

    throw redirect(303, `/${session.user.username}`)
  }
}
