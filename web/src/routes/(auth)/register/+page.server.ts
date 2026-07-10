import { fail, redirect, type Actions } from '@sveltejs/kit'
import { createSession, registerAccount } from '$lib/server/auth'
import { ApiError } from '$lib/server/api'

export const actions: Actions = {
  register: async ({ request, cookies }) => {
    const data = await request.formData()

    const firstName = String(data.get('firstName') ?? '').trim()
    const lastName = String(data.get('lastName') ?? '').trim()
    const email = String(data.get('email') ?? '').trim()
    const password = String(data.get('password') ?? '')
    const confirmPassword = String(data.get('confirmPassword') ?? '')

    const values = {
      firstName,
      lastName,
      email
    }

    if (!firstName) {
      return fail(400, {
        errors: {
          firstName: 'Informe seu nome.'
        },
        values
      })
    }

    if (!lastName) {
      return fail(400, {
        errors: {
          lastName: 'Informe seu sobrenome.'
        },
        values
      })
    }

    if (!email) {
      return fail(400, {
        errors: {
          email: 'Informe seu email.'
        },
        values
      })
    }

    if (!password) {
      return fail(400, {
        errors: {
          password: 'Informe sua senha.'
        },
        values
      })
    }

    if (password.length < 6) {
      return fail(400, {
        errors: {
          password: 'A senha deve ter pelo menos 6 caracteres.'
        },
        values
      })
    }

    if (password !== confirmPassword) {
      return fail(400, {
        errors: {
          confirmPassword: 'As senhas não conferem.'
        },
        values
      })
    }

    try {
      // Mock por enquanto; passa a chamar a Go API quando USE_API=true.
      const session = await registerAccount({ firstName, lastName, email, password })
      createSession(cookies, session)
    } catch (error) {
      if (error instanceof ApiError) {
        return fail(error.status === 409 ? 400 : 502, {
          errors: {
            general:
              error.status === 409
                ? 'Este email já está em uso.'
                : 'Não foi possível criar sua conta agora. Tente novamente em instantes.'
          },
          values
        })
      }

      throw error
    }

    throw redirect(303, '/')
  }
}