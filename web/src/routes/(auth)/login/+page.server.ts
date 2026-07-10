import { fail, redirect, type Actions } from '@sveltejs/kit'
import { authenticate, createSession } from '$lib/server/auth'
import { ApiError } from '$lib/server/api'

export const actions: Actions = {
  login: async ({ request, cookies }) => {
    const data = await request.formData()

    const email = String(data.get('email') ?? '').trim()
    const password = String(data.get('password') ?? '')

    if (!email) {
      return fail(400, {
        errors: {
          email: 'Informe seu email.'
        },
        values: {
          email
        }
      })
    }

    if (!password) {
      return fail(400, {
        errors: {
          password: 'Informe sua senha.'
        },
        values: {
          email
        }
      })
    }

    try {
      // Mock por enquanto; passa a chamar a Go API quando USE_API=true.
      const session = await authenticate({ email, password })
      createSession(cookies, session)
    } catch (error) {
      if (error instanceof ApiError) {
        return fail(error.status === 401 ? 400 : 502, {
          errors: {
            general:
              error.status === 401
                ? 'Email ou senha inválidos.'
                : 'Não foi possível entrar agora. Tente novamente em instantes.'
          },
          values: {
            email
          }
        })
      }

      throw error
    }

    throw redirect(303, '/')
  }
}