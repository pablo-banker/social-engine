// Tipos de domínio compartilhados entre server e client.
// Client-safe (sem imports de $lib/server), então componentes podem importar daqui.

export type Author = {
  name: string
  username: string
  avatarId: string
}

export type Post = {
  id: string
  author: Author
  content: string
  likes: number
  /** Se o usuário atual curtiu este post. */
  likedByMe: boolean
  comments: number
  createdAt: string
}

export type LikeResult = {
  liked: boolean
  likes: number
}

export type PostComment = {
  id: string
  author: Author
  content: string
  createdAt: string
}

export type Profile = {
  id: string
  name: string
  username: string
  bio: string
  avatarId: string
  bannerId: string
  joinedAt: string
  stats: {
    posts: number
    followers: number
    following: number
  }
}

/** Campos editáveis do perfil do usuário logado. */
export type ProfileSettings = {
  name: string
  bio: string
  avatarId: string
  bannerId: string
}

export type TrendingTopic = {
  /** Tag sem o '#'. */
  tag: string
  posts: number
}

export type SuggestedUser = {
  id: string
  name: string
  username: string
  bio: string
  avatarId: string
}
