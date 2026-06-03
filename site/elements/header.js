import { rmcmsAuth } from "/auth/rmcms-auth.js";

let stylesAdded = false;

function addStyles() {
    if (stylesAdded) return;
    stylesAdded = true;

    const style = document.createElement("style");
    style.textContent = `
        .rmcms-header {
            display: flex;
            align-items: center;
            justify-content: space-between;
            gap: 12px;
            min-height: 34px;
            padding: 6px 12px;
            border-bottom: 1px solid black;
            background: #fafafa;
            font-family: monospace;
            font-size: 14px;
        }

        .rmcms-header a {
            color: inherit;
            text-decoration: none;
        }

        .rmcms-header a:hover {
            text-decoration: underline;
        }

        .rmcms-header-left,
        .rmcms-header-right {
            display: flex;
            align-items: center;
            gap: 12px;
        }

        .rmcms-auth-indicator {
            padding: 2px 8px;
            border: 1px solid black;
            border-radius: 999px;
            white-space: nowrap;
        }

        .rmcms-auth-indicator[data-state="loading"] {
            opacity: 0.7;
        }
    `;
    document.head.appendChild(style);
}

function setAuthIndicator(link, authenticated) {
    link.textContent = authenticated ? "Logged in" : "Log in";
    link.dataset.state = authenticated ? "authenticated" : "unauthenticated";
}

export function mountHeader(container, options = {}) {
    if (!container) throw new Error("mountHeader needs a container element.");

    addStyles();

    const homeHref = options.homeHref || "/convert.html";
    const loginHref = options.loginHref || "/login.html";
    const homeLabel = options.homeLabel || "Home";

    const header = document.createElement("header");
    header.className = "rmcms-header";

    const left = document.createElement("div");
    left.className = "rmcms-header-left";

    const homeLink = document.createElement("a");
    homeLink.href = homeHref;
    homeLink.textContent = homeLabel;
    left.appendChild(homeLink);

    const right = document.createElement("div");
    right.className = "rmcms-header-right";

    const authLink = document.createElement("a");
    authLink.href = loginHref;
    authLink.target = "_blank";
    authLink.rel = "noopener noreferrer";
    authLink.className = "rmcms-auth-indicator";
    authLink.dataset.state = "loading";
    authLink.textContent = "Checking...";
    right.appendChild(authLink);

    header.appendChild(left);
    header.appendChild(right);
    container.replaceChildren(header);

    (async () => {
        try {
            await rmcmsAuth.tryRestoreSession();
            setAuthIndicator(authLink, rmcmsAuth.isAuthenticated());
        } catch {
            setAuthIndicator(authLink, false);
        }
    })();
}
