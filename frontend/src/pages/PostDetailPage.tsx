import { useEffect, useState } from 'react';
import { useParams } from 'react-router-dom';
import {
  Typography,
  Paper,
  IconButton,
  Stack,
  Divider,
  TextField,
  Button,
  CircularProgress,
  Box,
  Card,
  CardContent,
} from '@mui/material';
import ThumbUpIcon from '@mui/icons-material/ThumbUp';
import ThumbDownIcon from '@mui/icons-material/ThumbDown';
import type { Post } from '../types';
import { getPost, likePost, dislikePost, createComment, likeComment, dislikeComment } from '../api';

export default function PostDetailPage() {
  const { id } = useParams<{ id: string }>();
  const [post, setPost] = useState<Post | null>(null);
  const [loading, setLoading] = useState(true);
  const [commentAuthor, setCommentAuthor] = useState('');
  const [commentContent, setCommentContent] = useState('');
  const [submitting, setSubmitting] = useState(false);

  const fetchPost = () => {
    if (!id) return;
    getPost(Number(id))
      .then(setPost)
      .catch(console.error)
      .finally(() => setLoading(false));
  };

  useEffect(fetchPost, [id]);

  const handleLike = async () => {
    if (!post) return;
    setPost((prev) => (prev ? { ...prev, likes: prev.likes + 1 } : prev));
    try {
      await likePost(post.id);
    } catch {
      fetchPost();
    }
  };

  const handleDislike = async () => {
    if (!post) return;
    setPost((prev) => (prev ? { ...prev, dislikes: prev.dislikes + 1 } : prev));
    try {
      await dislikePost(post.id);
    } catch {
      fetchPost();
    }
  };

  const handleLikeComment = async (commentId: number) => {
    if (!post) return;
    setPost((prev) =>
      prev
        ? {
            ...prev,
            comments: prev.comments?.map((c) =>
              c.id === commentId ? { ...c, likes: c.likes + 1 } : c
            ),
          }
        : prev
    );
    try {
      await likeComment(commentId);
    } catch {
      fetchPost();
    }
  };

  const handleDislikeComment = async (commentId: number) => {
    if (!post) return;
    setPost((prev) =>
      prev
        ? {
            ...prev,
            comments: prev.comments?.map((c) =>
              c.id === commentId ? { ...c, dislikes: c.dislikes + 1 } : c
            ),
          }
        : prev
    );
    try {
      await dislikeComment(commentId);
    } catch {
      fetchPost();
    }
  };

  const handleAddComment = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!post || !commentAuthor.trim() || !commentContent.trim()) return;
    setSubmitting(true);
    try {
      await createComment(post.id, { author: commentAuthor, content: commentContent });
      setCommentAuthor('');
      setCommentContent('');
      fetchPost();
    } catch (err) {
      console.error(err);
    } finally {
      setSubmitting(false);
    }
  };

  if (loading) {
    return (
      <Box display="flex" justifyContent="center" mt={4}>
        <CircularProgress />
      </Box>
    );
  }

  if (!post) {
    return <Typography>Post not found.</Typography>;
  }

  return (
    <Stack spacing={3}>
      {/* Post content */}
      <Paper sx={{ p: 3 }}>
        <Typography variant="h4" gutterBottom>
          {post.title}
        </Typography>
        <Typography variant="body2" color="text.secondary" gutterBottom>
          By {post.author} · {new Date(post.createdAt).toLocaleDateString()}
        </Typography>
        <Typography variant="body1" sx={{ mt: 2, whiteSpace: 'pre-wrap' }}>
          {post.content}
        </Typography>

        <Stack direction="row" spacing={1} alignItems="center" sx={{ mt: 3 }}>
          <IconButton onClick={handleLike} color="primary">
            <ThumbUpIcon />
          </IconButton>
          <Typography>{post.likes}</Typography>
          <IconButton onClick={handleDislike} color="error">
            <ThumbDownIcon />
          </IconButton>
          <Typography>{post.dislikes}</Typography>
        </Stack>
      </Paper>

      {/* Comments */}
      <Typography variant="h5">Comments</Typography>
      <Divider />

      {post.comments && post.comments.length > 0 ? (
        post.comments.map((comment) => (
          <Card key={comment.id} variant="outlined">
            <CardContent>
              <Typography variant="subtitle2">
                {comment.author} · {new Date(comment.createdAt).toLocaleDateString()}
              </Typography>
              <Typography variant="body2" sx={{ mt: 1 }}>
                {comment.content}
              </Typography>
              <Stack direction="row" spacing={1} alignItems="center" sx={{ mt: 1 }}>
                <IconButton size="small" onClick={() => handleLikeComment(comment.id)} color="primary">
                  <ThumbUpIcon fontSize="small" />
                </IconButton>
                <Typography variant="body2">{comment.likes}</Typography>
                <IconButton size="small" onClick={() => handleDislikeComment(comment.id)} color="error">
                  <ThumbDownIcon fontSize="small" />
                </IconButton>
                <Typography variant="body2">{comment.dislikes}</Typography>
              </Stack>
            </CardContent>
          </Card>
        ))
      ) : (
        <Typography color="text.secondary">No comments yet.</Typography>
      )}

      {/* Add comment form */}
      <Paper sx={{ p: 3 }}>
        <Typography variant="h6" gutterBottom>
          Add a Comment
        </Typography>
        <Box component="form" onSubmit={handleAddComment}>
          <Stack spacing={2}>
            <TextField
              label="Author"
              value={commentAuthor}
              onChange={(e) => setCommentAuthor(e.target.value)}
              required
              fullWidth
              size="small"
            />
            <TextField
              label="Comment"
              value={commentContent}
              onChange={(e) => setCommentContent(e.target.value)}
              required
              fullWidth
              multiline
              rows={3}
            />
            <Button type="submit" variant="contained" disabled={submitting}>
              {submitting ? 'Submitting…' : 'Submit'}
            </Button>
          </Stack>
        </Box>
      </Paper>
    </Stack>
  );
}
