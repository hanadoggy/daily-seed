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

### Error Handling
- Handle errors explicitly. Never ignore or silently swallow errors.
- **Go:** Go does not have exceptions. Always return `error` as the last return value. Wrap errors with context using `fmt.Errorf("...: %w", err)` and check errors using `errors.Is` / `errors.As`.
- **TypeScript:** Use proper error types or result patterns when it improves clarity.

### Comments
- Write comments only when they explain **why** something is done, not what the code does.
- If you feel the need to write a long comment, consider refactoring the code first.
- Do not leave commented-out code in the final version.

### Testing & Maintainability
- Write code that is reasonably testable. Prioritize long-term maintainability over micro-optimizations.
- **Go:** Use standard `testing` package with `testify/assert` for assertions. Prefer lightweight in-memory stubs or real database test setup (e.g. SQLite in-memory / `testutil`) over complex mock frameworks.
- **TypeScript:** Use `Vitest` and `React Testing Library (RTL)`. Prioritize testing custom hooks and complex utility functions over simple UI components.
- When in doubt, choose the simpler and more readable approach.

## Language Specific Guidelines

### Go (Backend)
- **Architecture:** Follow a pragmatic 2-Layer package-by-feature architecture (`handler` and `store` per feature domain inside `internal/<domain>`). Keep components co-located with their data models.
- **Interfaces:** Follow Go's *"Accept interfaces, return structs"* idiom. Define small, consumer-side interfaces where they are used (e.g. in `handler.go`), never in central shared packages.
- **Formatting & Style:** Follow official Go formatting (`gofmt`).
- **Context:** Always pass `ctx context.Context` as the **first parameter** in functions performing I/O or cancellation-sensitive work.
- **Logging:** Use `log/slog` for structured logging with explicit attributes (e.g., `slog.String("id", id)`).
- **Validation:** Prefer explicit struct `Validate() error` methods or handler-level explicit validation over reflection-heavy tag validators.

### TypeScript / React (Frontend)
- **Architecture:** Group files by feature (`src/features`, `src/components`, `src/store`) keeping components, hooks, and types close together.
- Use TypeScript in strict mode. Avoid `any` as much as possible.
- Prefer small, focused components and custom hooks.
- Use Zustand for global state. Keep component logic understandable. Some duplication in UI components is acceptable if it improves readability.

### Zustand
- When performing optimistic updates, always implement rollback logic on API failure.
- Clearly separate UI-only state and server-synced data in the store.

## Edge Cases, Concurrency & Validation

- **Concurrency:** Handle potential race conditions gracefully. For instance, ensure a user cannot trigger the same action twice simultaneously.
- **Validation:** Always validate incoming API requests (pure TypeScript validation functions in `lib/validation.ts` for TypeScript, explicit `Validate() error` methods for Go).

## Project Conventions

- Follow the error response format: `{ code, message, details }` (`common.ErrorResponse`).
- When modifying existing code, maintain consistency with surrounding code style unless there is a clear reason to improve it.

## Final Mindset

Write code that is **easy for humans to read and modify** six months from now.  
Practical readability is more important than theoretical perfection.
