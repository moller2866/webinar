import { useEffect, useState } from 'react';
import { useNavigate, Link as RouterLink } from 'react-router-dom';
import {
  Card,
  CardContent,
  CardActionArea,
  Typography,
  Stack,
  CircularProgress,
  Box,
  Button,
  Chip,
} from '@mui/material';
import ThumbUpIcon from '@mui/icons-material/ThumbUp';
import ThumbDownIcon from '@mui/icons-material/ThumbDown';
import type { Post } from '../types';
import { getPosts } from '../api';

export default function PostListPage() {
  const [posts, setPosts] = useState<Post[]>([]);
  const [loading, setLoading] = useState(true);
  const [selectedTags, setSelectedTags] = useState<string[]>([]);
  const navigate = useNavigate();

  useEffect(() => {
    getPosts()
      .then((data) => setPosts(data ?? []))
      .catch(console.error)
      .finally(() => setLoading(false));
  }, []);

  const allTags = Array.from(new Set(posts.flatMap((p) => p.tags ?? []))).sort();

  const displayedPosts =
    selectedTags.length === 0
      ? posts
      : posts.filter((p) => selectedTags.some((t) => (p.tags ?? []).includes(t)));

  const toggleTag = (tag: string) => {
    setSelectedTags((prev) =>
      prev.includes(tag) ? prev.filter((t) => t !== tag) : [...prev, tag],
    );
  };

  if (loading) {
    return (
      <Box display="flex" justifyContent="center" mt={4}>
        <CircularProgress />
      </Box>
    );
  }

  if (posts.length === 0) {
    return (
      <Box textAlign="center" mt={4}>
        <Typography variant="h6" gutterBottom>
          No posts yet
        </Typography>
        <Button variant="contained" component={RouterLink} to="/posts/new">
          Create the first post
        </Button>
      </Box>
    );
  }

  return (
    <Stack spacing={2}>
      {allTags.length > 0 && (
        <Box>
          <Typography variant="caption" color="text.secondary" sx={{ mb: 0.5, display: 'block' }}>
            Filter by tag
          </Typography>
          <Stack direction="row" spacing={1} flexWrap="wrap" useFlexGap>
            {allTags.map((tag) => (
              <Chip
                key={tag}
                label={tag}
                onClick={() => toggleTag(tag)}
                color={selectedTags.includes(tag) ? 'primary' : 'default'}
                variant={selectedTags.includes(tag) ? 'filled' : 'outlined'}
                size="small"
              />
            ))}
          </Stack>
        </Box>
      )}
      {displayedPosts.map((post) => (
        <Card key={post.id}>
          <CardActionArea onClick={() => navigate(`/posts/${post.id}`)}>
            <CardContent>
              <Typography variant="h5" gutterBottom>
                {post.title}
              </Typography>
              <Typography variant="body2" color="text.secondary" gutterBottom>
                By {post.author} · {new Date(post.createdAt).toLocaleDateString()}
              </Typography>
              <Typography variant="body1" sx={{ mb: 1 }}>
                {post.content.length > 150
                  ? post.content.slice(0, 150) + '…'
                  : post.content}
              </Typography>
              <Stack direction="row" spacing={1} flexWrap="wrap" useFlexGap>
                <Chip icon={<ThumbUpIcon />} label={post.likes} size="small" />
                <Chip icon={<ThumbDownIcon />} label={post.dislikes} size="small" />
                {(post.tags ?? []).map((tag) => (
                  <Chip key={tag} label={tag} size="small" variant="outlined" color="primary" />
                ))}
              </Stack>
            </CardContent>
          </CardActionArea>
        </Card>
      ))}
    </Stack>
  );
}
