import type { LayoutServerLoad } from './$types'
import { getPublicSession } from '$lib/server/auth'

export const load: LayoutServerLoad = async ({ cookies }) => {
  return {
    session: getPublicSession(cookies)
  }
}