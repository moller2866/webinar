export interface Post {
  id: number;
  title: string;
  content: string;
  author: string;
  likes: number;
  dislikes: number;
  createdAt: string;
  comments?: Comment[];
}

export interface Comment {
  id: number;
  postId: number;
  author: string;
  content: string;
  createdAt: string;
}

export interface CreatePostRequest {
  title: string;
  content: string;
  author: string;
}

export interface CreateCommentRequest {
  author: string;
  content: string;
}
