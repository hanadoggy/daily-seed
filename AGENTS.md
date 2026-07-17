# AGENTS.md — Coding Guidelines

You are an expert software engineer. Your primary goal is to write **readable, maintainable, and practical** code rather than overly idealistic code.

## Agent Behavior & Output Rules

- **Do not hallucinate packages or external tools.** If a required utility or library is missing, ask for clarification before inventing a solution.
- **Provide concise explanations.** Omit pleasantries and skip basic setup instructions unless explicitly asked.
- **Code modification:** When proposing changes, clearly indicate where the change occurs. Do not output the entire file unless necessary; use clear comments like `// ... existing code ...` to represent unchanged parts.
- If user requirements are ambiguous or conflict with these guidelines, **stop and ask clarifying questions** before writing code.

## Core Principles

### Meaningful Names
- Use clear, intention-revealing names.
- Avoid single-letter variable names except for simple loop counters.
- A name should explain what the variable or function does without needing a comment.

### Functions & Abstraction
- Extract a function when **the behavior has clear meaning** and **repetition exists**.
- Do not extract functions just because they are short (e.g., one-liners). 
- Some repetition is acceptable if extracting it would reduce readability or add unnecessary abstraction.
- Focus on making the **overall flow easy to understand** rather than forcing every small piece into its own function.
- Good example: Extract `calculateCompletionRate()` if the same logic appears in multiple places.
- Bad example: Extracting `getTrue()` or `formatDateString()` for a one-time simple operation.

### Error Handling
- Handle errors explicitly. Never ignore or silently swallow errors.
- In Go, always check and propagate errors appropriately.
- In TypeScript, use proper error types or result patterns when it improves clarity.
- Prefer returning errors over throwing exceptions when reasonable (especially in Go).

### Comments
- Write comments only when they explain **why** something is done, not what the code does.
- If you feel the need to write a long comment, consider refactoring the code first.
- Do not leave commented-out code in the final version.

### Testing & Maintainability
- Write code that is reasonably testable. Prioritize long-term maintainability over micro-optimizations.
- **Go:** Use the standard `testing` package along with `stretchr/testify` for assertions and mocking. Focus on unit testing business logic in the service layer using interface mocks.
- **TypeScript:** Use `Vitest` and `React Testing Library (RTL)`. Prioritize testing custom hooks and complex utility functions over simple UI components.
- When in doubt, choose the simpler and more readable approach.

## Language Specific Guidelines

### Go (Backend)
- **Architecture:** Follow a layered architecture: `handler` (HTTP), `service` (business logic), `repository` (DB operations). Place domain types and interfaces in a shared location to avoid circular dependencies.
- Follow official Go formatting (`gofmt`).
- Use `slog` for structured logging with request IDs.
- Keep HTTP handlers relatively thin. Move complex logic to the service layer.
- Use interfaces for dependency injection when it improves testability.
- Always pass `context.Context` in public functions that may involve I/O or cancellation.

### TypeScript / React (Frontend)
- **Architecture:** Group files by feature (Feature-sliced design) rather than strictly by file type (e.g., keep a feature's components, hooks, and types close together).
- Use TypeScript in strict mode. Avoid `any` as much as possible.
- Prefer small, focused components and custom hooks.
- Use Zustand for global state. Keep component logic understandable. Some duplication in UI components is acceptable if it improves readability.

### Zustand
- When performing optimistic updates, always implement rollback logic on API failure.
- Clearly separate UI-only state and server-synced data in the store.

## Edge Cases, Concurrency & Validation

- **Concurrency:** Handle potential race conditions gracefully. For instance, ensure a user cannot trigger the same action twice simultaneously (e.g., using optimistic locking or state checks before updating).
- **Validation:** Always validate incoming API requests. Use `Zod` for TypeScript and `go-playground/validator` for Go to enforce strict input boundaries.

## Project Conventions

- Follow the error response format: `{ code, message, details }`.
- When modifying existing code, try to maintain consistency with surrounding code style unless there is a clear reason to improve it.


## Final Mindset

Write code that is **easy for humans to read and modify** six months from now.  
Practical readability is more important than theoretical perfection.
