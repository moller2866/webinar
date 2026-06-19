export interface User {
  id: number;
  email: string;
  displayName: string;
  bio: string;
  avatarUrl: string;
  createdAt: string;
}

export interface Category {
  id: number;
  name: string;
  slug: string;
  description: string;
  createdBy?: number;
  createdAt: string;
}

export interface Vote {
  id: number;
  userId: number;
  targetType: 'post' | 'comment';
  targetId: number;
  voteType: 'highfive' | 'meh';
  createdAt: string;
}

export interface Post {
  id: number;
  title: string;
  content: string;
  userId?: number;
  categoryId?: number;
  images: string[];
  createdAt: string;
  comments?: Comment[];
}

export interface Comment {
  id: number;
  postId: number;
  userId?: number;
  parentId?: number;
  content: string;
  createdAt: string;
}

export interface CreatePostRequest {
  title: string;
  content: string;
  categoryId?: number;
  images?: string[];
}

export interface CreateCommentRequest {
  content: string;
  parentId?: number;
}
