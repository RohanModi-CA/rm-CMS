# Development Guide

## Local Pipeline
1. **The Server:** `cd src && go run server.go`. 
2. **The Site:** Served at `http://localhost:8080`.
3. **The Parser:** Custom regex-based Markdown engine in `src/parsers/`.

## Core Workflow
1. **Upload:** Select a local folder in the UI.
2. **Edit:** Modify Markdown in the browser.
3. **Optimize:** Images are auto-converted to WebP with SHA-256 hashes.
4. **Publish:** Pushes to GitHub via sparse checkout.

## Rules
- **Stay Client-Side:** If it can be done in JS, do it there.
- **Keep it Simple:** Avoid large frameworks or complex state management.
