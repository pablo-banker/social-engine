import { DEFAULT_AVATAR_ID, DEFAULT_BANNER_ID } from '$lib/appearance'
import type {
  LikeResult,
  Post,
  PostComment,
  Profile,
  ProfileSettings,
  SuggestedUser,
  TrendingTopic
} from '$lib/types'

/**
 * Fonte única de dados mockados (usa-se só quando USE_API=false).
 * Simula o "banco": usuários + feed global em memória. Some ao reiniciar
 * o server e é compartilhado entre requisições. Será trocado por PostgreSQL.
 */

// O `avatarId` do autor é resolvido na leitura (a partir do usuário), então os
// posts/comentários são guardados só com nome + username do autor.
type StoredAuthor = {
  name: string
  username: string
}

// `likedByMe` depende de quem está vendo — preenchido na leitura via `decorate`.
type StoredPost = Omit<Post, 'likedByMe' | 'author'> & { author: StoredAuthor }

type StoredComment = Omit<PostComment, 'author'> & { author: StoredAuthor }

type MockUser = {
  id: string
  name: string
  username: string
  bio: string
  avatarId: string
  bannerId: string
  joinedAt: string
  followers: number
  following: number
}

const users: MockUser[] = [
  {
    id: 'u1',
    name: 'Pablo Banker',
    username: 'pcbanker',
    bio: 'Construindo o Social Engine com Svelte, Go e PostgreSQL. Foco em arquitetura limpa.',
    avatarId: 'a1',
    bannerId: 'b1',
    joinedAt: 'julho de 2026',
    followers: 128,
    following: 42
  },
  {
    id: 'u2',
    name: 'Social Engine',
    username: 'engine',
    bio: 'A engine social modular. Feed público, ações autenticadas.',
    avatarId: 'a5',
    bannerId: 'b5',
    joinedAt: 'julho de 2026',
    followers: 512,
    following: 3
  },
  {
    id: 'u3',
    name: 'Ada Lovelace',
    username: 'ada',
    bio: 'A primeira programadora. Curiosa por padrões e máquinas analíticas.',
    avatarId: 'a6',
    bannerId: 'b6',
    joinedAt: 'junho de 2026',
    followers: 9042,
    following: 120
  }
]

// Feed global, mais recentes primeiro.
const posts: StoredPost[] = [
  {
    id: '1',
    author: { name: 'Pablo Banker', username: 'pcbanker' },
    content:
      'Começando o Social Engine: uma plataforma social modular usando #Svelte, #Go e #PostgreSQL.',
    likes: 24,
    comments: 8,
    createdAt: 'agora'
  },
  {
    id: '2',
    author: { name: 'Social Engine', username: 'engine' },
    content:
      'O feed é público. Ações como postar, comentar, curtir e seguir exigem autenticação. #Svelte #webdev',
    likes: 12,
    comments: 3,
    createdAt: '2min'
  },
  {
    id: '3',
    author: { name: 'Ada Lovelace', username: 'ada' },
    content:
      'A máquina analítica tece padrões algébricos como o tear tece flores e folhas. #história #computação',
    likes: 340,
    comments: 21,
    createdAt: '1h'
  }
]

// Curtidas: postId -> conjunto de usernames que curtiram.
const likedBy: Record<string, Set<string>> = {}

// Comentários por postId (mais antigos primeiro).
const commentsByPost: Record<string, StoredComment[]> = {
  '1': [
    {
      id: 'c1',
      author: { name: 'Ada Lovelace', username: 'ada' },
      content: 'Ótima escolha de stack. Vou acompanhar de perto!',
      createdAt: '5min'
    },
    {
      id: 'c2',
      author: { name: 'Social Engine', username: 'engine' },
      content: 'Bem-vindo ao Social Engine 🚀',
      createdAt: '2min'
    }
  ],
  '3': [
    {
      id: 'c3',
      author: { name: 'Pablo Banker', username: 'pcbanker' },
      content: 'Clássico atemporal. Poesia e engenharia no mesmo parágrafo.',
      createdAt: '30min'
    }
  ]
}

/** Avatar atual de um autor (resolvido pelo username; padrão se não achar). */
function avatarIdOf(username: string): string {
  return users.find((user) => user.username === username)?.avatarId ?? DEFAULT_AVATAR_ID
}

/** Preenche `avatarId` do autor e `likedByMe` para o `viewer`. */
function decorate(post: StoredPost, viewer?: string): Post {
  return {
    ...post,
    author: { ...post.author, avatarId: avatarIdOf(post.author.username) },
    likedByMe: viewer ? (likedBy[post.id]?.has(viewer) ?? false) : false
  }
}

/** Feed completo, decorado para `viewer`. */
export function listPosts(viewer?: string): Post[] {
  return posts.map((post) => decorate(post, viewer))
}

/** Um post específico, decorado para `viewer`, ou null se não existir. */
export function getPost(id: string, viewer?: string): Post | null {
  const post = posts.find((candidate) => candidate.id === id)

  return post ? decorate(post, viewer) : null
}

/** Posts de um usuário específico, decorados para `viewer`. */
export function listPostsByUsername(username: string, viewer?: string): Post[] {
  return posts.filter((post) => post.author.username === username).map((post) => decorate(post, viewer))
}

/** Insere um post novo no topo do feed. */
export function addPost(post: StoredPost) {
  posts.unshift(post)
}

/** Alterna a curtida de `username` no post e devolve o novo estado. */
export function toggleLike(postId: string, username: string): LikeResult | null {
  const post = posts.find((candidate) => candidate.id === postId)

  if (!post) {
    return null
  }

  const set = (likedBy[postId] ??= new Set())

  if (set.has(username)) {
    set.delete(username)
    post.likes -= 1

    return { liked: false, likes: post.likes }
  }

  set.add(username)
  post.likes += 1

  return { liked: true, likes: post.likes }
}

/** Comentários de um post (mais antigos primeiro), com avatar resolvido. */
export function listComments(postId: string): PostComment[] {
  return (commentsByPost[postId] ?? []).map((comment) => ({
    ...comment,
    author: { ...comment.author, avatarId: avatarIdOf(comment.author.username) }
  }))
}

/** Adiciona um comentário e incrementa o contador do post. */
export function addComment(postId: string, comment: StoredComment) {
  ;(commentsByPost[postId] ??= []).push(comment)

  const post = posts.find((candidate) => candidate.id === postId)

  if (post) {
    post.comments += 1
  }
}

/** Extrai as tags (#tag) únicas e minúsculas de um conteúdo. */
function extractTags(content: string): string[] {
  const matches = content.match(/#[\p{L}\p{N}_]+/gu) ?? []

  return [...new Set(matches.map((match) => match.slice(1).toLowerCase()))]
}

/** Tópicos em alta: contagem de posts por tag, mais populares primeiro. */
export function getTrendingTopics(): TrendingTopic[] {
  const counts = new Map<string, number>()

  for (const post of posts) {
    for (const tag of extractTags(post.content)) {
      counts.set(tag, (counts.get(tag) ?? 0) + 1)
    }
  }

  return [...counts.entries()]
    .map(([tag, count]) => ({ tag, posts: count }))
    .sort((a, b) => b.posts - a.posts || a.tag.localeCompare(b.tag))
}

/** Posts que contêm uma tag específica, decorados para `viewer`. */
export function listPostsByTopic(tag: string, viewer?: string): Post[] {
  const target = tag.toLowerCase()

  return posts
    .filter((post) => extractTags(post.content).includes(target))
    .map((post) => decorate(post, viewer))
}

/** Todos os posts ordenados por curtidas (mais curtidos primeiro). */
export function listPopularPosts(viewer?: string): Post[] {
  return [...posts].sort((a, b) => b.likes - a.likes).map((post) => decorate(post, viewer))
}

/** Sugestões de quem seguir (todos os usuários, menos `exclude`). */
export function listSuggestedUsers(exclude?: string): SuggestedUser[] {
  return users
    .filter((user) => user.username !== exclude)
    .map(({ id, name, username, bio, avatarId }) => ({ id, name, username, bio, avatarId }))
}

/** Perfil público por username, ou null se não existir. */
export function getProfile(username: string): Profile | null {
  const user = users.find((candidate) => candidate.username === username)

  if (!user) {
    return null
  }

  return {
    id: user.id,
    name: user.name,
    username: user.username,
    bio: user.bio,
    avatarId: user.avatarId,
    bannerId: user.bannerId,
    joinedAt: user.joinedAt,
    stats: {
      posts: listPostsByUsername(username).length,
      followers: user.followers,
      following: user.following
    }
  }
}

/** Campos editáveis do perfil, ou null se o usuário não existir. */
export function getUserSettings(username: string): ProfileSettings | null {
  const user = users.find((candidate) => candidate.username === username)

  if (!user) {
    return null
  }

  return {
    name: user.name,
    bio: user.bio,
    avatarId: user.avatarId,
    bannerId: user.bannerId
  }
}

/**
 * Atualiza o perfil e propaga o novo nome para posts/comentários já
 * existentes (o avatar é resolvido na leitura, então acompanha sozinho).
 * Retorna false se o usuário não existir.
 */
export function updateUserSettings(username: string, input: ProfileSettings): boolean {
  const user = users.find((candidate) => candidate.username === username)

  if (!user) {
    return false
  }

  user.name = input.name
  user.bio = input.bio
  user.avatarId = input.avatarId
  user.bannerId = input.bannerId

  for (const post of posts) {
    if (post.author.username === username) {
      post.author.name = input.name
    }
  }

  for (const list of Object.values(commentsByPost)) {
    for (const comment of list) {
      if (comment.author.username === username) {
        comment.author.name = input.name
      }
    }
  }

  return true
}

/** Usuário canônico do login mockado (reflete edições no avatar/nome). */
export function getMockLoginUser() {
  const user = users.find((candidate) => candidate.username === 'pcbanker') ?? users[0]

  return {
    id: user.id,
    name: user.name,
    username: user.username,
    avatarId: user.avatarId
  }
}

/** Registra um novo usuário no mock (para register criar perfil editável). */
export function addMockUser(input: { name: string; username: string; email: string }): void {
  if (users.some((user) => user.username === input.username)) {
    return
  }

  users.push({
    id: crypto.randomUUID(),
    name: input.name,
    username: input.username,
    bio: '',
    avatarId: DEFAULT_AVATAR_ID,
    bannerId: DEFAULT_BANNER_ID,
    joinedAt: 'agora',
    followers: 0,
    following: 0
  })
}
