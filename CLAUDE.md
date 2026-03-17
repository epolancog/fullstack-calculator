# Claude Code Guidelines

## Communication Rules

- **NEVER assume anything.** If a requirement, design choice, or implementation detail is ambiguous, ask before proceeding.
- **Always clarify questions first.** Before starting any implementation work, summarize your understanding of the task and ask clarifying questions. Wait for confirmation before writing code.
- **Present your plan before executing.** For each session or significant task, outline what you intend to do, confirm alignment, then proceed.
- **Flag deviations from the plan.** If during implementation you discover something that requires a different approach than what's in PLAN.md, stop and discuss before changing course.

## Implementation Rules

- Follow the session plan in `PLAN.md` — do not skip steps or reorder without discussion.
- One commit per session, with the commit message specified in the plan.
- Run tests and manual test scenarios before committing.
- Update the Progress Tracker in `PLAN.md` after completing a session.
- Apply SOLID principles only where they naturally fit — do not force them.
- Do not over-engineer. Keep solutions simple and focused on the requirements.
- Do not add features, refactor code, or make "improvements" beyond what the plan specifies.

## Tech Stack

- Backend: Go 1.26+, stdlib `net/http` (no external router), standard `testing` package
- Frontend: Vite + React + TypeScript, Tailwind CSS v4, CVA (class-variance-authority)
- Testing: Go `testing` (backend), Vitest + React Testing Library (frontend)
- Monorepo: `backend/` and `frontend/` directories at repo root
