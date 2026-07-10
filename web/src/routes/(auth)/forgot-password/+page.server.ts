import { fail, type Actions } from '@sveltejs/kit'

export const actions: Actions = {
  forgotPassword: async ({ request }) => {
    const data = await request.formData()

    const email = String(data.get('email') ?? '').trim()

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

    return {
      sent: true,
      values: {
        email
      }
    }
  }
}