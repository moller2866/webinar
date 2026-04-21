# PRD: Campfire — Community Discussion Platform

## Overview

**Problem:** People want a warm, topic-based space to discuss things they care about — cooking, travel, movies, and more — but existing platforms feel cold, techy, or overwhelming. The current blog application supports basic posts and comments but lacks user accounts, community structure, and the personality needed to foster real connection.

**Solution:** Evolve the existing blog into **Campfire**, a community discussion platform organized around user-created categories. Authenticated users can post, comment in threaded conversations, and vote using a friendly "High-five / Meh" system. The experience should feel like gathering around a campfire — warm, welcoming, and human.

**Target Users:** Everyday people (non-technical) who want to share and discuss topics they're passionate about. They value simplicity, warmth, and belonging over power-user features.

---

## User Stories

### Authentication

1. **As a visitor**, I want to browse all categories, posts, and comments without signing up, so that I can get a feel for the community before committing.
   - AC: All read endpoints (list categories, list posts, get post + comments) are accessible without authentication.
   - AC: Write actions (post, comment, vote) return `401 Unauthorized` for unauthenticated users.

2. **As a new user**, I want to sign up with a display name, email, and password, so that I can participate in the community.
   - AC: Registration requires a unique email, a display name (2–30 characters), and a password (minimum 8 characters).
   - AC: Passwords are hashed before storage (bcrypt).
   - AC: On success, the user is automatically logged in and receives an auth token.
   - AC: Duplicate email returns a clear, friendly error.

3. **As a returning user**, I want to log in with my email and password, so that I can pick up where I left off.
   - AC: Valid credentials return an auth token (JWT).
   - AC: Invalid credentials return a generic "email or password is incorrect" error (no enumeration).
   - AC: The token is sent as a `Bearer` token in the `Authorization` header for subsequent requests.

4. **As a logged-in user**, I want to log out, so that my session ends.
   - AC: The frontend discards the stored token. (Stateless JWT — no server-side revocation in v1.)

### User Profiles

5. **As a logged-in user**, I want to set up my profile with a display name, avatar, and short bio, so that other community members can get to know me.
   - AC: Profile fields: display name (required, 2–30 chars), avatar (optional image upload, max 2 MB), bio (optional, max 200 chars).
   - AC: Display name is shown on all posts and comments authored by the user.
   - AC: Users can update their own profile at any time.

6. **As a visitor or user**, I want to view another user's profile, so that I can learn more about them.
   - AC: Profile page shows display name, avatar, bio, and a list of their recent posts.
   - AC: Profile is read-only to everyone except the owner.

### Categories

7. **As a logged-in user**, I want to create a new category, so that I can start a community around a topic I care about.
   - AC: A category has a name (3–50 chars, unique, case-insensitive) and a short description (max 200 chars).
   - AC: Category names are URL-safe slugs (e.g., "Cooking Tips" → `cooking-tips`).
   - AC: Duplicate names return a friendly error.

8. **As a visitor or user**, I want to browse a list of all categories, so that I can find communities that interest me.
   - AC: Categories are listed alphabetically by default.
   - AC: Each category card shows the name, description, and post count.

9. **As a visitor or user**, I want to view a single category's feed, so that I can see all posts in that topic.
   - AC: The feed shows only posts belonging to the selected category.
   - AC: Sorting can be toggled between "Newest" and "Most High-fived" (see FR-5).

### Posts

10. **As a logged-in user**, I want to create a post in a specific category, so that I can share something with that community.
    - AC: A post requires a title (3–150 chars), body text (1–10,000 chars), and a category.
    - AC: A post optionally includes one or more images (max 5 images, each max 5 MB, JPEG/PNG/WebP).
    - AC: The post is associated with the authenticated user (author is no longer a free-text field).
    - AC: The post appears at the top of the category's "Newest" feed immediately.

11. **As a visitor or user**, I want to view a single post with its full content and comments, so that I can read the discussion.
    - AC: Post detail page shows title, body, images, author (linked to profile), category, timestamp, high-five/meh counts, and threaded comments.

### Comments

12. **As a logged-in user**, I want to comment on a post, so that I can join the conversation.
    - AC: A comment requires body text (1–2,000 chars).
    - AC: The comment is associated with the authenticated user.

13. **As a logged-in user**, I want to reply to another comment, so that conversations can branch naturally.
    - AC: A reply is a comment with a `parentId` pointing to the comment being replied to.
    - AC: Replies are displayed nested under their parent, up to 5 levels deep visually. Deeper replies are flattened to level 5.
    - AC: Each reply shows "Replying to [display name]" context.

### Voting (High-five / Meh)

14. **As a logged-in user**, I want to "High-five" a post or comment, so that I can show appreciation.
    - AC: Each user can high-five a given post/comment at most once.
    - AC: High-fiving again removes the high-five (toggle behavior).
    - AC: The total high-five count is displayed on the post/comment.

15. **As a logged-in user**, I want to "Meh" a post or comment, so that I can signal it's not useful.
    - AC: Each user can meh a given post/comment at most once.
    - AC: Meh-ing again removes the meh (toggle behavior).
    - AC: The total meh count is displayed on the post/comment.
    - AC: A user cannot have both a high-five and a meh active on the same item. Voting one way clears the other.

### Sorting & Discovery

16. **As a user browsing a category**, I want to toggle between "Newest" and "Most High-fived" sorting, so that I can find fresh or popular content.
    - AC: "Newest" sorts by `created_at DESC`.
    - AC: "Most High-fived" sorts by `high_fives DESC`, with ties broken by `created_at DESC`.
    - AC: The sort preference is reflected in the URL (query param) so it's shareable/bookmarkable.
    - AC: Default sort is "Newest".

---

## Functional Requirements

### FR-1: User Authentication System
- **Description:** Add registration, login, and JWT-based auth. Protect all write endpoints with middleware that validates the token and injects the user identity.
- **Priority:** Must-have
- **Notes:** The existing `author` free-text field on posts/comments is replaced by the authenticated user's ID. Existing like/dislike endpoints are renamed and gated behind auth.

### FR-2: User Profiles
- **Description:** Add a `users` table and profile endpoints. Support avatar image upload and short bio.
- **Priority:** Must-have
- **Notes:** Avatar images stored on disk (or a configurable path). No CDN or cloud storage in v1.

### FR-3: Categories
- **Description:** Add a `categories` table. Posts must belong to exactly one category. Any logged-in user can create a category.
- **Priority:** Must-have
- **Notes:** The existing `tags` field on posts is **removed** and replaced by the single category relationship. Categories have a unique slug derived from the name.

### FR-4: Threaded Comments
- **Description:** Add a `parent_id` (nullable, self-referencing FK) to the comments table. The API returns comments as a tree. Frontend renders nested replies.
- **Priority:** Must-have
- **Notes:** Max visual nesting depth of 5 levels in the UI. The API returns a flat list with `parentId`; the frontend assembles the tree.

### FR-5: High-five / Meh Voting
- **Description:** Replace the current like/dislike counter system with per-user voting. A `votes` table tracks `(user_id, target_type, target_id, vote_type)`. Votes are toggle-able and mutually exclusive.
- **Priority:** Must-have
- **Notes:** `target_type` is `'post'` or `'comment'`. `vote_type` is `'highfive'` or `'meh'`. Unique constraint on `(user_id, target_type, target_id)`.

### FR-6: Post Image Uploads
- **Description:** Allow up to 5 images per post. Images are uploaded as part of post creation (multipart form) and stored server-side.
- **Priority:** Should-have
- **Notes:** Accepted formats: JPEG, PNG, WebP. Max 5 MB per image. Images are served via a static file endpoint. The post model gains an `images` array of URLs.

### FR-7: Sorting Toggle
- **Description:** Category feeds support two sort modes: "Newest" and "Most High-fived", toggled via a UI control and a `?sort=` query parameter.
- **Priority:** Must-have

### FR-8: Home Page — Category Discovery
- **Description:** The home/landing page shows a browsable list of all categories with name, description, and post count.
- **Priority:** Must-have

---

## Non-Functional Requirements

### Performance
- Page loads (category feed, post detail) should feel instant — target < 200ms server response for the 90th percentile.
- Category feeds should be paginated (20 posts per page) to avoid loading thousands of posts at once.
- Comment trees should be loaded in a single query per post (not N+1).

### Security
- Passwords hashed with bcrypt (cost factor ≥ 12).
- JWT tokens signed with a server-side secret (HS256). Token expiry: 7 days.
- All write endpoints require a valid JWT. Middleware validates before the handler runs.
- Image uploads validated server-side: check MIME type, enforce size limits, sanitize filenames. Never serve user-uploaded files with their original filename.
- SQL injection prevention via parameterized queries (already in place via `pgx`).
- Rate limiting on registration and login endpoints (basic: 10 requests/minute per IP) — **nice-to-have for v1**.

### Accessibility
- Semantic HTML for all pages. Proper heading hierarchy.
- All images must have `alt` text (user-provided on upload or a sensible default).
- Keyboard-navigable voting buttons and comment reply actions.
- Color contrast meeting WCAG 2.1 AA.

### UX / Visual Tone
- Warm color palette: earthy tones, campfire oranges, soft creams, forest greens.
- Friendly, rounded UI elements (buttons, cards, avatars).
- No developer jargon in the UI. Labels like "Start a conversation" instead of "Create post," "Join the discussion" instead of "Add comment."
- High-five icon: a hand / high-five emoji. Meh icon: a shrug emoji or neutral face.

---

## Technical Considerations

### Existing Codebase — What Changes

| Layer | Current State | Required Changes |
|-------|--------------|------------------|
| **Model** | `Post` (title, content, author string, tags, likes, dislikes), `Comment` (author string, likes, dislikes) | Add `User`, `Category`, `Vote` models. Replace `author string` with `userID int64`. Replace `likes/dislikes int` with computed high-five/meh counts. Add `parentId` to Comment. Replace `tags` with `categoryId`. Add `images []string` to Post. |
| **Repository** | `PostRepository`, `CommentRepository` interfaces | Add `UserRepository`, `CategoryRepository`, `VoteRepository` interfaces. Update existing interfaces for new fields, filtering by category, and sorting. |
| **Service** | `PostService` handles posts + comments | Add `UserService` (registration, login, password verification), `CategoryService`. Update `PostService` for category filtering, sorting, image handling. Add `VoteService`. |
| **Handler** | 8 endpoints, no auth | Add auth middleware. Add user, category, vote endpoints. Update existing endpoints for auth + new fields. |
| **Database** | `posts`, `comments` tables | Add `users`, `categories`, `votes` tables. Migrate existing tables (add FKs, remove deprecated columns). |
| **Frontend types** | `Post`, `Comment`, `CreatePostRequest`, `CreateCommentRequest` | Add `User`, `Category`, `Vote` types. Update existing types. |
| **Frontend API** | `api.ts` with 8 functions | Add auth headers, new API functions for users, categories, votes. Update existing functions. |
| **Frontend pages** | `PostListPage`, `PostDetailPage`, `NewPostPage` | Add `LoginPage`, `RegisterPage`, `ProfilePage`, `CategoryListPage`, `CategoryFeedPage`. Refactor `NewPostPage` for category selection + image upload. |

### New API Endpoints (Proposed)

| Method | Path | Auth | Description |
|--------|------|------|-------------|
| POST | /api/auth/register | No | Create account |
| POST | /api/auth/login | No | Log in, receive JWT |
| GET | /api/users/{id} | No | View user profile |
| PUT | /api/users/{id} | Yes (owner) | Update own profile |
| POST | /api/users/{id}/avatar | Yes (owner) | Upload avatar image |
| GET | /api/categories | No | List all categories |
| POST | /api/categories | Yes | Create a category |
| GET | /api/categories/{slug} | No | Get category details |
| GET | /api/categories/{slug}/posts | No | List posts in category (supports `?sort=newest\|highfives`) |
| POST | /api/posts | Yes | Create a post (updated — requires `categoryId`) |
| GET | /api/posts/{id} | No | Get post + threaded comments |
| POST | /api/posts/{id}/comments | Yes | Add comment (supports `parentId`) |
| POST | /api/posts/{id}/highfive | Yes | Toggle high-five on post |
| POST | /api/posts/{id}/meh | Yes | Toggle meh on post |
| POST | /api/comments/{id}/highfive | Yes | Toggle high-five on comment |
| POST | /api/comments/{id}/meh | Yes | Toggle meh on comment |

### New Database Tables (Proposed)

```sql
CREATE TABLE users (
    id          BIGSERIAL PRIMARY KEY,
    email       TEXT UNIQUE NOT NULL,
    display_name TEXT NOT NULL,
    password_hash TEXT NOT NULL,
    bio         TEXT DEFAULT '',
    avatar_url  TEXT DEFAULT '',
    created_at  TIMESTAMPTZ DEFAULT NOW()
);

CREATE TABLE categories (
    id          BIGSERIAL PRIMARY KEY,
    name        TEXT UNIQUE NOT NULL,
    slug        TEXT UNIQUE NOT NULL,
    description TEXT DEFAULT '',
    created_by  BIGINT REFERENCES users(id),
    created_at  TIMESTAMPTZ DEFAULT NOW()
);

CREATE TABLE votes (
    id          BIGSERIAL PRIMARY KEY,
    user_id     BIGINT NOT NULL REFERENCES users(id),
    target_type TEXT NOT NULL CHECK (target_type IN ('post', 'comment')),
    target_id   BIGINT NOT NULL,
    vote_type   TEXT NOT NULL CHECK (vote_type IN ('highfive', 'meh')),
    created_at  TIMESTAMPTZ DEFAULT NOW(),
    UNIQUE (user_id, target_type, target_id)
);
```

Existing tables modified:
- `posts`: add `category_id BIGINT REFERENCES categories(id)`, add `user_id BIGINT REFERENCES users(id)`, add `images TEXT[]`, remove `tags`, remove `likes`/`dislikes` (computed from `votes`).
- `comments`: add `user_id BIGINT REFERENCES users(id)`, add `parent_id BIGINT REFERENCES comments(id)`, remove `likes`/`dislikes`.

### Dependencies & Constraints
- **No external router** — continue using Go 1.22+ `http.ServeMux` method+pattern routing.
- **No frontend state library** — auth state managed via `useState` + React Context for the current user/token.
- **No axios** — continue using `fetch` via `api.ts`.
- **PostgreSQL only** — all new tables live in the same database.
- **JWT library** — will need a Go JWT library (e.g., `golang-jwt/jwt/v5`).
- **bcrypt** — use `golang.org/x/crypto/bcrypt`.
- **Image storage** — local filesystem in v1. Images served via a `/uploads/` static file handler.

---

## Out of Scope

- **Moderation tools** (reporting, content removal, bans) — deferred to v2.
- **Search** (full-text search across posts) — deferred.
- **Notifications** (new replies, new posts in followed categories) — deferred.
- **Email verification** — registration works immediately without confirming email.
- **Password reset / forgot password** — deferred.
- **OAuth / social login** — deferred.
- **Rich text editor** — posts and comments are plain text (with line breaks) in v1.
- **Direct messages** between users — deferred.
- **Category moderation** (category-specific rules, mods) — deferred.
- **Dark mode** — deferred.
- **Rate limiting** — nice-to-have, not required for v1.
- **Cloud image storage** (S3, etc.) — local disk only in v1.

---

## Open Questions

1. **Category limits** — Should there be a cap on how many categories a single user can create (to prevent spam)?
2. **Post editing / deletion** — Can a user edit or delete their own post after publishing? (Not mentioned in requirements — needs a decision.)
3. **Comment editing / deletion** — Same question for comments.
4. **Home page vs. category page** — Should the home page show a global "all posts" feed in addition to the category list, or only categories?
5. **Image upload flow** — Should images be uploaded inline during post composition (with preview), or selected all at once on submit? This affects the UX significantly.
6. **Display name uniqueness** — Should display names be unique across the platform, or can multiple users share a name (distinguished by their profile)?
