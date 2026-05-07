import { useEffect, useState } from 'react';
import { useParams } from 'react-router-dom';
import {
  Typography,
  Paper,
  Stack,
  Divider,
  TextField,
  Button,
  CircularProgress,
  Box,
  Card,
  CardContent,
} from '@mui/material';
import type { Post } from '../types';
import { getPost, createComment } from '../api';

export default function PostDetailPage() {
  const { id } = useParams<{ id: string }>();
  const [post, setPost] = useState<Post | null>(null);
  const [loading, setLoading] = useState(true);
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

  const handleAddComment = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!post || !commentContent.trim()) return;
    setSubmitting(true);
    try {
      await createComment(post.id, { content: commentContent });
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
          {new Date(post.createdAt).toLocaleDateString()}
        </Typography>
        <Typography variant="body1" sx={{ mt: 2, whiteSpace: 'pre-wrap' }}>
          {post.content}
        </Typography>
      </Paper>

      {/* Comments */}
      <Typography variant="h5">Comments</Typography>
      <Divider />

      {post.comments && post.comments.length > 0 ? (
        post.comments.map((comment) => (
          <Card key={comment.id} variant="outlined">
            <CardContent>
              <Typography variant="body2" color="text.secondary">
                {new Date(comment.createdAt).toLocaleDateString()}
              </Typography>
              <Typography variant="body2" sx={{ mt: 1 }}>
                {comment.content}
              </Typography>
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
