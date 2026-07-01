# RM-CMS Elements

`elements/` currently contains one reusable element: the site header.

The goal is to keep the element small enough that a page can mount it without
having to know any of the auth internals.

## Header

`mountHeader(container, options)` renders a two-part header:

- a left home link
- a right auth indicator link

It also tries to restore the auth session before it decides what text to show.

### Parameters

- `container`: the element that will receive the header
- `options.homeHref`: overrides the home link, default `/convert.html`
- `options.loginHref`: overrides the auth link, default `/login.html`
- `options.homeLabel`: overrides the home label, default `Home`

### Behavior

`mountHeader()`:

1. adds `elements/global.css` once
2. creates the header markup
3. replaces the container contents with the header
4. calls `rmcmsAuth.tryRestoreSession()`
5. updates the auth text to `Logged in` or `Log in`

If session restore fails, the header falls back to `Log in`.

### Example

```html
<div id="header-slot"></div>
<script type="module">
    import { mountHeader } from "/elements/header.js";

    mountHeader(document.getElementById("header-slot"), {
        homeHref: "/convert.html",
        loginHref: "/login.html",
        homeLabel: "Home",
    });
</script>
```

## Styles

`elements/global.css` contains the header styling.

The styles are intentionally minimal:

- a horizontal layout
- simple spacing
- plain text links
- a lighter state while auth is still loading

