import { redirect, type RequestHandler } from '@sveltejs/kit'
import { clearSession } from '$lib/server/auth'

export const POST: RequestHandler = ({ cookies }) => {
  clearSession(cookies)

  // Volta para o feed público; o header reflete o estado deslogado.
  throw redirect(303, '/')
}
