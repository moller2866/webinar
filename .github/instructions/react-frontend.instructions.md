---
description: "Use when writing React frontend code — pages, components, API calls, or TypeScript types. Covers API layer usage, MUI conventions, state patterns, and file structure."
applyTo: "frontend/src/**/*.{ts,tsx}"
---

# React Frontend Conventions

## Project Structure

| Path | Purpose |
|------|---------|
| `pages/` | Route-level page components (`PostListPage`, `PostDetailPage`, `NewPostPage`) |
| `components/` | Reusable UI pieces shared across pages (e.g. `Layout`) |
| `api.ts` | All fetch calls — the only file that uses raw `fetch` |
| `types.ts` | All shared TypeScript interfaces |

## API Layer

All server communication goes through `api.ts`. Never call `fetch` directly in a component or page.

Add new endpoints as named async functions in `api.ts`:

```ts
export async function likePost(id: number): Promise<void> {
  const response = await fetch(`${API_BASE}/posts/${id}/like`, { method: 'POST' });
  if (!response.ok) throw new Error('Failed to like post');
}
```

## Types

Define all shared interfaces in `types.ts`. Do not co-locate type definitions with components.

## MUI v6

Use MUI components throughout — `Button`, `TextField`, `Typography`, `Stack`, `Card`, etc.
Avoid plain `<button>`, `<input>`, or `<form>` where an MUI equivalent exists.
The theme is created once in `App.tsx` via `createTheme()` and provided via `ThemeProvider`.

## State Management

Use `useState` and `useEffect` directly in page components. No state management library (no Redux, Zustand, Jotai, etc.).

## Optimistic Updates

For instant-feeling interactions (e.g. like/dislike), update local state immediately and revert on error:

```ts
setPost((prev) => prev ? { ...prev, likes: prev.likes + 1 } : prev);
try {
  await likePost(post.id);
} catch {
  fetchPost(); // revert by re-fetching
}
```

## Routing

React Router v6 — extract params with `useParams`, navigate with `useNavigate`.
Routes are defined in `App.tsx` only. New routes are added there, not in individual pages.
