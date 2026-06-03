# Architecture & Philosophy

## Principles
- **Simplicity Above All:** Code must be short and readable. If it doesn't fit in one mental schema, it's too complex.
- **Zero Server:** The client is the engine. We avoid backend logic whenever possible.
- **Trusted Scale:** Built for ~10 users. No over-engineering for edge cases, pen-testers, or old browsers.
- **Pragmatic UX:** `alert()` is often the best error handler.

## Components
1. **Frontend:** Vanilla HTML/JS. Logic-heavy on the client.
2. **Parser (Local Go):** A simple utility to convert Markdown to HTML and optimize images.
3. **Storage:** Git-based. Static files only.
4. **Auth:** To be implemented using client-side encryption.

## Rules
- No big architectural changes without a full rethink.
- If it can be done on the client, do it there.
- Favor brevity over "completeness."
- Docs should point to working example scripts when they exist.
