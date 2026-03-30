import { useState } from 'react';
import { useNavigate } from 'react-router-dom';
import { Typography, Paper, TextField, Button, Stack, Box } from '@mui/material';
import { createPost } from '../api';

export default function NewPostPage() {
  const navigate = useNavigate();
  const [title, setTitle] = useState('');
  const [author, setAuthor] = useState('');
  const [content, setContent] = useState('');
  const [tagsInput, setTagsInput] = useState('');
  const [submitting, setSubmitting] = useState(false);

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!title.trim() || !author.trim() || !content.trim()) return;
    setSubmitting(true);
    const tags = tagsInput
      .split(',')
      .map((t) => t.trim())
      .filter(Boolean);
    try {
      const post = await createPost({ title, author, content, tags });
      navigate(`/posts/${post.id}`);
    } catch (err) {
      console.error(err);
      setSubmitting(false);
    }
  };

  return (
    <Paper sx={{ p: 3 }}>
      <Typography variant="h5" gutterBottom>
        Create New Post
      </Typography>
      <Box component="form" onSubmit={handleSubmit}>
        <Stack spacing={2}>
          <TextField
            label="Title"
            value={title}
            onChange={(e) => setTitle(e.target.value)}
            required
            fullWidth
          />
          <TextField
            label="Author"
            value={author}
            onChange={(e) => setAuthor(e.target.value)}
            required
            fullWidth
          />
          <TextField
            label="Content"
            value={content}
            onChange={(e) => setContent(e.target.value)}
            required
            fullWidth
            multiline
            rows={6}
          />
          <TextField
            label="Tags (comma-separated, optional)"
            value={tagsInput}
            onChange={(e) => setTagsInput(e.target.value)}
            fullWidth
            placeholder="e.g. go, react, tutorial"
          />
          <Button type="submit" variant="contained" size="large" disabled={submitting}>
            {submitting ? 'Publishing…' : 'Publish'}
          </Button>
        </Stack>
      </Box>
    </Paper>
  );
}
