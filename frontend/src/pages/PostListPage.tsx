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
  const navigate = useNavigate();

  useEffect(() => {
    getPosts()
      .then((data) => setPosts(data ?? []))
      .catch(console.error)
      .finally(() => setLoading(false));
  }, []);

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
      {posts.map((post) => (
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
              <Stack direction="row" spacing={1}>
                <Chip icon={<ThumbUpIcon />} label={post.likes} size="small" />
                <Chip icon={<ThumbDownIcon />} label={post.dislikes} size="small" />
              </Stack>
            </CardContent>
          </CardActionArea>
        </Card>
      ))}
    </Stack>
  );
}
