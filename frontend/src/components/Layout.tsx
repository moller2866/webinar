import { AppBar, Toolbar, Typography, Button, Container, Box } from '@mui/material';
import { Link, Outlet } from 'react-router-dom';

export default function Layout() {
  return (
    <Box>
      <AppBar position="static">
        <Toolbar>
          <Typography
            variant="h6"
            component={Link}
            to="/"
            sx={{ flexGrow: 1, textDecoration: 'none', color: 'inherit' }}
          >
            Blog
          </Typography>
          <Button color="inherit" component={Link} to="/posts/new">
            New Post
          </Button>
        </Toolbar>
      </AppBar>
      <Container maxWidth="md" sx={{ mt: 4, mb: 4 }}>
        <Outlet />
      </Container>
    </Box>
  );
}
