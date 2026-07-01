# RM-CMS Documentation

This document explains the project at a high level.

The goal of RM-CMS is simple: turn local content into site output, and keep the
site editing and authentication pieces separate enough to understand without
holding too much in your head at once.

## Main Areas

There are three main pieces in this repo:

- `convert.html` and `convert.js` for the content editing and publishing flow
- `auth/` for client-side login, signup, and secret encryption
- `elements/` for reusable UI pieces, currently only the header

The remaining pages are support tools:

- `site_explorer.html` and `site_explorer.js` render the site map
- `file_tree.html` and `file_tree.js` build editable file trees
- `destination_selector.html` and `destination_selector.js` choose the output target
- `examples/auth_demo.html` shows secret retrieval through the auth module

## Editing And Publishing

The main site editor lives in `convert.html`.

The flow is:

1. a user selects a local directory
2. markdown files are shown as buttons
3. a file is loaded into the editor
4. the page title and date are filled in
5. the markdown is sent to the server for preview
6. the preview HTML is written into an iframe
7. static assets and final HTML can then be pushed

This part of the app depends on the server endpoints that `convert.js` calls:

- `/upload_markdown`
- `/resources_dump`
- `/push-static-images`
- `/push-html`

## Site Tree Tools

The site tree code is separate from the editor.

`site_explorer.js` reads `site_map.json` and turns it into a nested directory
view. `file_tree.js` and `destination_selector.js` use the same tree idea to
build interactive selectors for choosing or creating paths.

The mental model is:

- `site_map.json` is the read-only shape of the site
- `site_explorer.js` displays that shape
- `file_tree.js` creates editable directory nodes
- `destination_selector.js` combines the tree with file-name entry

## Authentication

Authentication is browser-side and split across a few small scripts.

The flow is:

1. signup creates an RSA keypair for a user
2. the private key is wrapped with a password-derived AES key
3. login unwraps the private key after deriving the same AES key again
4. the unlocked private key is cached in IndexedDB
5. secret data is fetched and decrypted on demand

The core auth module is `auth/rmcms-auth.js`.

It stores two runtime values:

- the private key
- the hashed username

Those values are enough to look up the right secret data file and decrypt the
property for the current user.

## Header

`elements/header.js` exports a single reusable function, `mountHeader()`.

The header:

- renders a left link to the home page
- renders a right auth indicator
- tries to restore the auth session on mount
- changes the indicator text to show whether a session is present

The header is intentionally small. It is a reusable element, not a full component
framework.

## Current Gaps

The code documents the current behavior, including a few mismatches that are
already present in the repo:

- `login.js` and `auth/manager.js` treat `keys.json` as an array of users
- `rmcms-auth.js` treats `data.users` as an object keyed by username hash
- `signup.js` currently outputs a single user object
- `signup.html` shows a WebAuthn section, but `signup.js` does not yet wire it up
- `site_map.json` still references `auth.html`, which is not present in the repo

These are worth fixing in code, but the docs below describe the current state as
it exists now.

