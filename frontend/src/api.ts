import type { Post, Comment, CreatePostRequest, CreateCommentRequest } from './types';

const API_BASE = 'http://localhost:8080/api';

export async function getPosts(): Promise<Post[]> {
  const response = await fetch(`${API_BASE}/posts`);
  if (!response.ok) throw new Error('Failed to fetch posts');
  return response.json();
}

export async function getPost(id: number): Promise<Post> {
  const response = await fetch(`${API_BASE}/posts/${id}`);
  if (!response.ok) throw new Error('Failed to fetch post');
  return response.json();
}

export async function createPost(data: CreatePostRequest): Promise<Post> {
  const response = await fetch(`${API_BASE}/posts`, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify(data),
  });
  if (!response.ok) throw new Error('Failed to create post');
  return response.json();
}

export async function createComment(postId: number, data: CreateCommentRequest): Promise<Comment> {
  const response = await fetch(`${API_BASE}/posts/${postId}/comments`, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify(data),
  });
  if (!response.ok) throw new Error('Failed to create comment');
  return response.json();
}

export async function likePost(id: number): Promise<void> {
  const response = await fetch(`${API_BASE}/posts/${id}/like`, { method: 'POST' });
  if (!response.ok) throw new Error('Failed to like post');
}

export async function dislikePost(id: number): Promise<void> {
  const response = await fetch(`${API_BASE}/posts/${id}/dislike`, { method: 'POST' });
  if (!response.ok) throw new Error('Failed to dislike post');
}
