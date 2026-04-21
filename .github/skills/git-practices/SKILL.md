---
name: git-practices
description: Enforce Git best practices — branch naming, commit message conventions, and workflow rules. Use when working on any code change that should be committed.
---

# Git Practices

## Branch Rules

- **Never commit directly to `main`**. Always create a branch first.
- Use the appropriate branch prefix:

| Prefix | When to use |
|--------|-------------|
| `feature/` | New functionality |
| `bugfix/` | Fixing a bug |
| `refactor/` | Code restructuring without behavior change |
| `docs/` | Documentation-only changes |
| `chore/` | Tooling, config, dependency updates |

- Branch names should be lowercase, hyphen-separated, and descriptive.  
  Example: `feature/add-comment-likes`, `bugfix/fix-null-post-id`

## Commit Message Conventions

Use [Conventional Commits](https://www.conventionalcommits.org/):

```
<type>: <short imperative description>
```

**Types:**

| Type | When to use |
|------|-------------|
| `feat:` | New feature |
| `fix:` | Bug fix |
| `docs:` | Documentation change |
| `refactor:` | Code change that neither fixes a bug nor adds a feature |
| `test:` | Adding or updating tests |
| `chore:` | Tooling, config, CI, dependency updates |
| `style:` | Formatting, whitespace (no logic change) |

**Rules:**
- Keep the subject line under 72 characters
- Use imperative mood: "add feature" not "added feature" or "adds feature"
- No period at the end of the subject line
- Do not include issue numbers unless the project convention requires it

**Examples:**
```
feat: add like button to comment cards
fix: return 404 when post not found
refactor: extract comment query into helper
docs: update API endpoint table in README
chore: bump pgx to v5.7
```

## Commit Frequency

- Commit after each logical, self-contained unit of work
- Do not bundle unrelated changes in one commit
- Prefer many small commits over one large commit at the end
- A commit should leave the codebase in a working state

## Workflow Checklist

Before starting any work:
1. Confirm you are not on `main` — if you are, create a branch first
2. Choose the right prefix for the branch
3. Name the branch after the change being made

While working:
4. Commit frequently as you complete logical steps
5. Each commit message must start with the correct type prefix

Before finishing:
6. Review staged changes — only include what belongs in this commit
7. Verify the branch has not diverged unexpectedly from `main`